package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// MockRoundTripper for capturing requests and sending mock responses
type MockRoundTripper struct {
	Func func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Func(req)
}

func TestGetAllTokenAccounts(t *testing.T) {
	mockTransport := &MockRoundTripper{
		Func: func(req *http.Request) (*http.Response, error) {
			// Read request body
			body, _ := io.ReadAll(req.Body)
			bodyStr := string(body)
            // request body is consumed, usually we should restore it if needed, but here we just mock response

			// Prepare mock response
			var accounts []interface{}

			if strings.Contains(bodyStr, TokenProgramID) {
				// Return legacy tokens
				accounts = []interface{}{
					map[string]interface{}{
						"pubkey": "LegacyAcc1",
						"account": map[string]interface{}{
							"data": map[string]interface{}{
								"parsed": map[string]interface{}{
									"info": map[string]interface{}{
										"mint": "MintA",
										"tokenAmount": map[string]interface{}{
											"amount":   "1000",
											"decimals": 6,
										},
									},
								},
							},
						},
					},
				}
			} else if strings.Contains(bodyStr, Token2022ProgramID) {
				// Return Token-2022 tokens
				accounts = []interface{}{
					map[string]interface{}{
						"pubkey": "Token2022Acc1",
						"account": map[string]interface{}{
							"data": map[string]interface{}{
								"parsed": map[string]interface{}{
									"info": map[string]interface{}{
										"mint": "MintB",
										"tokenAmount": map[string]interface{}{
											"amount":   "2000",
											"decimals": 9,
										},
									},
								},
							},
						},
					},
				}
			} else {
                // Unexpected request
                return &http.Response{
                    StatusCode: 400,
                    Body:       io.NopCloser(bytes.NewBufferString("Bad Request")),
                }, nil
            }

			respObj := RPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Result:  mustMarshal(map[string]interface{}{"value": accounts}),
			}
			respBytes, _ := json.Marshal(respObj)

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer(respBytes)),
			}, nil
		},
	}

	client := &RPCClient{
		primaryURL: "http://mock",
		httpClient: &http.Client{Transport: mockTransport},
	}

	accounts, err := client.GetAllTokenAccounts(context.Background(), "Owner1")
	if err != nil {
		t.Fatalf("GetAllTokenAccounts failed: %v", err)
	}

	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}

    // Verify contents
    foundA := false
    foundB := false
    for _, acc := range accounts {
        if acc.Mint == "MintA" && acc.Amount == 1000 {
            foundA = true
        }
        if acc.Mint == "MintB" && acc.Amount == 2000 {
            foundB = true
        }
    }

    if !foundA {
        t.Error("MintA not found or incorrect")
    }
    if !foundB {
        t.Error("MintB not found or incorrect")
    }
}

func mustMarshal(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
