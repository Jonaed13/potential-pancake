package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllTokenAccounts(t *testing.T) {
	// Mock response
	mockResponse := `{
		"jsonrpc": "2.0",
		"result": {
			"value": [
				{
					"pubkey": "Account1",
					"account": {
						"data": {
							"parsed": {
								"info": {
									"mint": "Mint1",
									"tokenAmount": {
										"amount": "1000",
										"decimals": 6
									}
								}
							}
						}
					}
				},
				{
					"pubkey": "Account2",
					"account": {
						"data": {
							"parsed": {
								"info": {
									"mint": "Mint2",
									"tokenAmount": {
										"amount": "2000",
										"decimals": 9
									}
								}
							}
						}
					}
				}
			]
		},
		"id": 1
	}`

	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Verify request body
		var req RPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Method != "getTokenAccountsByOwner" {
			t.Errorf("expected method getTokenAccountsByOwner, got %s", req.Method)
		}

		// Verify params
		if len(req.Params) < 3 {
			t.Fatalf("expected at least 3 params, got %d", len(req.Params))
		}

		// Param 1: Owner
		if req.Params[0] != "OwnerAddress" {
			t.Errorf("expected owner 'OwnerAddress', got %v", req.Params[0])
		}

		// Param 2: Filter (should contain programId)
		filter, ok := req.Params[1].(map[string]interface{})
		if !ok {
			// It might be unmarshaled as map[string]string if we typed it strictly,
			// but RPCRequest uses interface{}. JSON unmarshal produces map[string]interface{}.
			t.Errorf("expected filter to be a map, got %T", req.Params[1])
		}

		if filter["programId"] != TokenProgramID {
			t.Errorf("expected programId %s, got %v", TokenProgramID, filter["programId"])
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mockResponse)
	}))
	defer ts.Close()

	// Create RPC client
	client := NewRPCClient(ts.URL, ts.URL, "test-api-key")

	// Call GetAllTokenAccounts
	accounts, err := client.GetAllTokenAccounts(context.Background(), "OwnerAddress")
	if err != nil {
		t.Fatalf("GetAllTokenAccounts failed: %v", err)
	}

	// Verify results
	if len(accounts) != 2 {
		t.Errorf("expected 2 accounts, got %d", len(accounts))
	}

	// Account 1
	if accounts[0].Mint != "Mint1" {
		t.Errorf("expected account 0 mint 'Mint1', got %s", accounts[0].Mint)
	}
	if accounts[0].Amount != 1000 {
		t.Errorf("expected account 0 amount 1000, got %d", accounts[0].Amount)
	}
	if accounts[0].Decimals != 6 {
		t.Errorf("expected account 0 decimals 6, got %d", accounts[0].Decimals)
	}

	// Account 2
	if accounts[1].Mint != "Mint2" {
		t.Errorf("expected account 1 mint 'Mint2', got %s", accounts[1].Mint)
	}
	if accounts[1].Amount != 2000 {
		t.Errorf("expected account 1 amount 2000, got %d", accounts[1].Amount)
	}
}
