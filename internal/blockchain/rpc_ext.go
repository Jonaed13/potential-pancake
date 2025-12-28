package blockchain

import (
	"context"
	"strconv"
)

// GetAllTokenAccounts fetches all SPL token accounts for an owner
func (c *RPCClient) GetAllTokenAccounts(ctx context.Context, owner string) ([]TokenAccountInfo, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getTokenAccountsByOwner",
		Params: []interface{}{
			owner,
			map[string]string{"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"},
			map[string]string{
				"encoding": "jsonParsed",
			},
		},
	}

	var result struct {
		Value []struct {
			Pubkey  string `json:"pubkey"`
			Account struct {
				Data struct {
					Parsed struct {
						Info struct {
							Mint        string `json:"mint"`
							TokenAmount struct {
								Amount   string `json:"amount"`
								Decimals uint8  `json:"decimals"`
							} `json:"tokenAmount"`
						} `json:"info"`
					} `json:"parsed"`
				} `json:"data"`
			} `json:"account"`
		} `json:"value"`
	}

	if err := c.call(ctx, req, &result); err != nil {
		return nil, err
	}

	accounts := make([]TokenAccountInfo, 0, len(result.Value))
	for _, v := range result.Value {
		amount, _ := strconv.ParseUint(v.Account.Data.Parsed.Info.TokenAmount.Amount, 10, 64)
		accounts = append(accounts, TokenAccountInfo{
			Address:  v.Pubkey,
			Mint:     v.Account.Data.Parsed.Info.Mint,
			Amount:   amount,
			Decimals: v.Account.Data.Parsed.Info.TokenAmount.Decimals,
		})
	}

	return accounts, nil
}
