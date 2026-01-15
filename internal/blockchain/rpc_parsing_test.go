package blockchain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTokenAccountsByOwner_Parsing(t *testing.T) {
	// Mock response from Solana RPC
	mockResponse := `{
		"jsonrpc": "2.0",
		"result": {
			"context": {
				"apiVersion": "2.0.15",
				"slot": 240000000
			},
			"value": [
				{
					"pubkey": "Account1Address",
					"account": {
						"lamports": 2039280,
						"owner": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
						"data": {
							"program": "spl-token",
							"parsed": {
								"info": {
									"isNative": false,
									"mint": "MintAddress1",
									"owner": "OwnerAddress",
									"state": "initialized",
									"tokenAmount": {
										"amount": "123456789",
										"decimals": 6,
										"uiAmount": 123.456789,
										"uiAmountString": "123.456789"
									}
								},
								"type": "account"
							},
							"space": 165
						},
						"executable": false,
						"rentEpoch": 0
					}
				}
			]
		},
		"id": 1
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewRPCClient(server.URL, server.URL, "test-api-key")
	// Since we are creating a new client, it uses the default transport.
	// But we passed server.URL so it will hit our mock server.

	ctx := context.Background()
	accounts, err := client.GetTokenAccountsByOwner(ctx, "OwnerAddress", "MintAddress1")
	if err != nil {
		t.Fatalf("GetTokenAccountsByOwner failed: %v", err)
	}

	if len(accounts) != 1 {
		t.Fatalf("Expected 1 account, got %d", len(accounts))
	}

	expectedAmount := uint64(123456789)
	if accounts[0].Amount != expectedAmount {
		t.Errorf("Expected amount %d, got %d", expectedAmount, accounts[0].Amount)
	}
}

func TestGetTokenAccountBalance_Parsing(t *testing.T) {
	mockResponse := `{
		"jsonrpc": "2.0",
		"result": {
			"context": {
				"slot": 240000000
			},
			"value": {
				"amount": "9876543210",
				"decimals": 9,
				"uiAmount": 9.87654321,
				"uiAmountString": "9.87654321"
			}
		},
		"id": 1
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewRPCClient(server.URL, server.URL, "test-api-key")

	ctx := context.Background()
	amount, decimals, err := client.GetTokenAccountBalance(ctx, "AccountAddress")
	if err != nil {
		t.Fatalf("GetTokenAccountBalance failed: %v", err)
	}

	expectedAmount := uint64(9876543210)
	if amount != expectedAmount {
		t.Errorf("Expected amount %d, got %d", expectedAmount, amount)
	}
	if decimals != 9 {
		t.Errorf("Expected decimals 9, got %d", decimals)
	}
}

func TestGetTokenAccountBalance_InvalidParsing(t *testing.T) {
	mockResponse := `{
		"jsonrpc": "2.0",
		"result": {
			"context": { "slot": 240000000 },
			"value": {
				"amount": "invalid-number",
				"decimals": 9
			}
		},
		"id": 1
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewRPCClient(server.URL, server.URL, "test-api-key")

	ctx := context.Background()
	amount, _, err := client.GetTokenAccountBalance(ctx, "AccountAddress")

	// Should not error on RPC call success, but parsing fails so amount should be 0
	if err != nil {
		t.Fatalf("GetTokenAccountBalance failed: %v", err)
	}

	if amount != 0 {
		t.Errorf("Expected amount 0 for invalid input, got %d", amount)
	}
}
