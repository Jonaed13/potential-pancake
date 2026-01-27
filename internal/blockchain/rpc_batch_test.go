package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// MockRoundTripper for capturing requests and returning mock responses
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestGetAllTokenAccounts_Batch(t *testing.T) {
	// Setup mock responses for Legacy and Token-2022 calls
	mockResponseLegacy := `
	{
		"jsonrpc": "2.0",
		"id": 1,
		"result": {
			"value": [
				{
					"pubkey": "LegacyAccount1",
					"account": {
						"data": {
							"parsed": {
								"info": {
									"mint": "LegacyMint1",
									"tokenAmount": {
										"amount": "1000",
										"decimals": 9
									}
								}
							}
						}
					}
				}
			]
		}
	}`

	mockResponseToken2022 := `
	{
		"jsonrpc": "2.0",
		"id": 1,
		"result": {
			"value": [
				{
					"pubkey": "Token2022Account1",
					"account": {
						"data": {
							"parsed": {
								"info": {
									"mint": "Token2022Mint1",
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
		}
	}`

	// Create a mock client
	client := NewRPCClient("http://mock-primary", "http://mock-fallback", "apikey")
	client.httpClient.Transport = &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Read body to determine which request this is
			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body

			var rpcReq RPCRequest
			json.Unmarshal(bodyBytes, &rpcReq)

			// Check programId param to decide response
			params := rpcReq.Params
			if len(params) > 1 {
				config := params[1].(map[string]interface{})
				programID, ok := config["programId"].(string)
				if ok {
					if programID == TokenProgramID {
						return &http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(mockResponseLegacy)),
						}, nil
					} else if programID == Token2022ProgramID {
						return &http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(mockResponseToken2022)),
						}, nil
					}
				}
			}

			// Fallback (should not happen in this test)
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": "unknown request"}`)),
			}, nil
		},
	}

	// Run the method
	accounts, err := client.GetAllTokenAccounts(context.Background(), "WalletOwner")
	if err != nil {
		t.Fatalf("GetAllTokenAccounts failed: %v", err)
	}

	// Verification
	// Should have 2 accounts (1 from Legacy, 1 from Token2022)
	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}

	// Verify details
	legacyFound := false
	token2022Found := false

	for _, acc := range accounts {
		if acc.Mint == "LegacyMint1" && acc.Amount == 1000 {
			legacyFound = true
		}
		if acc.Mint == "Token2022Mint1" && acc.Amount == 2000 {
			token2022Found = true
		}
	}

	if !legacyFound {
		t.Error("Legacy account not found or incorrect")
	}
	if !token2022Found {
		t.Error("Token-2022 account not found or incorrect")
	}
}

func TestGetAllTokenAccounts_PartialFailure(t *testing.T) {
	// Create a mock client that fails for Token-2022
	client := NewRPCClient("http://mock-primary", "http://mock-fallback", "apikey")
	client.httpClient.Transport = &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			var rpcReq RPCRequest
			json.Unmarshal(bodyBytes, &rpcReq)

			// Extract programId
			if len(rpcReq.Params) > 1 {
				config := rpcReq.Params[1].(map[string]interface{})
				programID := config["programId"].(string)

				if programID == TokenProgramID {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString(`{"jsonrpc":"2.0","result":{"value":[]}}`)),
					}, nil
				}
				if programID == Token2022ProgramID {
					return &http.Response{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString("fail")),
					}, nil
				}
			}
			return nil, nil
		},
	}

	_, err := client.GetAllTokenAccounts(context.Background(), "WalletOwner")
	if err == nil {
		t.Error("Expected error on partial failure, got nil")
	}
}
