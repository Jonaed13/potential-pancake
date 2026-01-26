package blockchain

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllTokenAccounts(t *testing.T) {
	// Mock server to handle RPC calls
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req RPCRequest
		json.Unmarshal(body, &req)

		if req.Method == "getTokenAccountsByOwner" {
			params := req.Params
			// params[0] is owner
			// params[1] is filter map

			// Handle map[string]interface{} (from json unmarshal) vs map[string]string (if locally typed)
			// Since we unmarshal into interface{}, it will be map[string]interface{}
			filter, ok := params[1].(map[string]interface{})
			if !ok {
				// It might be because params is []interface{} and unmarshaling makes it map[string]interface{}
				// In the test, request is marshaled then unmarshaled.
			}

			var accounts []interface{}

            createAccount := func(mint, amount string) interface{} {
                return map[string]interface{}{
                    "pubkey": "somepubkey",
                    "account": map[string]interface{}{
                        "data": map[string]interface{}{
                            "parsed": map[string]interface{}{
                                "info": map[string]interface{}{
                                    "mint": mint,
                                    "tokenAmount": map[string]interface{}{
                                        "amount": amount,
                                        "decimals": float64(6),
                                    },
                                },
                            },
                        },
                    },
                }
            }

			if filter["programId"] == TokenProgramID {
                accounts = append(accounts, createAccount("mint1", "100"))
			} else if filter["programId"] == Token2022ProgramID {
                accounts = append(accounts, createAccount("mint1", "50"))
                accounts = append(accounts, createAccount("mint2", "200"))
			}

			resp := RPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
			}
            resultBytes, _ := json.Marshal(map[string]interface{}{"value": accounts})
            resp.Result = resultBytes
			json.NewEncoder(w).Encode(resp)
			return
		}

        // Default error
        http.Error(w, "not implemented", http.StatusNotImplemented)
	}))
	defer server.Close()

	client := NewRPCClient(server.URL, server.URL, "test-key")

    ctx := context.Background()
    balances, err := client.GetAllTokenAccounts(ctx, "owner123")
    if err != nil {
        t.Fatalf("GetAllTokenAccounts failed: %v", err)
    }

    // mint1 should be 100 (legacy) + 50 (token22) = 150
    if balances["mint1"] != 150 {
        t.Errorf("Expected mint1 balance 150, got %d", balances["mint1"])
    }

    // mint2 should be 200 (token22)
    if balances["mint2"] != 200 {
        t.Errorf("Expected mint2 balance 200, got %d", balances["mint2"])
    }
}
