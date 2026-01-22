package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// RPCClient handles Solana RPC calls
type RPCClient struct {
	primaryURL   string
	fallbackURL  string
	apiKey       string
	httpClient   *http.Client
	
	// Circuit breaker state
	mu           sync.RWMutex
	failures     int
	lastFailure  time.Time
	circuitOpen  bool
}

// RPCRequest is the JSON-RPC 2.0 request format
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}

// RPCResponse is the JSON-RPC 2.0 response format
type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError is the JSON-RPC 2.0 error format
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

// BlockhashResult is the result of getLatestBlockhash
type BlockhashResult struct {
	Value struct {
		Blockhash            string `json:"blockhash"`
		LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`
	} `json:"value"`
}

// BalanceResult is the result of getBalance
type BalanceResult struct {
	Value uint64 `json:"value"`
}

// SendTxResult is the result of sendTransaction
type SendTxResult string

// NewRPCClient creates a new RPC client
func NewRPCClient(primaryURL, fallbackURL, apiKey string) *RPCClient {
	// Configure HTTP transport for keep-alives and connection pooling
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	return &RPCClient{
		primaryURL:  primaryURL,
		fallbackURL: fallbackURL,
		apiKey:      apiKey,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

// GetLatestBlockhash fetches the latest blockhash
func (c *RPCClient) GetLatestBlockhash(ctx context.Context) (*BlockhashResult, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getLatestBlockhash",
		Params:  []interface{}{map[string]string{"commitment": "confirmed"}},
	}

	var result BlockhashResult
	if err := c.call(ctx, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBalance fetches the SOL balance for a public key
func (c *RPCClient) GetBalance(ctx context.Context, pubkey string) (uint64, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getBalance",
		Params:  []interface{}{pubkey, map[string]string{"commitment": "confirmed"}},
	}

	var result BalanceResult
	if err := c.call(ctx, req, &result); err != nil {
		return 0, err
	}

	return result.Value, nil
}

// SendTransaction sends a signed transaction
func (c *RPCClient) SendTransaction(ctx context.Context, signedTx string, skipPreflight bool) (string, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "sendTransaction",
		Params: []interface{}{
			signedTx,
			map[string]interface{}{
				"encoding":       "base64",
				"skipPreflight":  skipPreflight,
				"preflightCommitment": "processed",
				"maxRetries":     3,
			},
		},
	}

	var result SendTxResult
	if err := c.call(ctx, req, &result); err != nil {
		return "", err
	}

	return string(result), nil
}

// GetTokenAccountBalance fetches SPL token balance
func (c *RPCClient) GetTokenAccountBalance(ctx context.Context, tokenAccount string) (uint64, uint8, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getTokenAccountBalance",
		Params:  []interface{}{tokenAccount},
	}

	var result struct {
		Value struct {
			Amount   string `json:"amount"`
			Decimals uint8  `json:"decimals"`
		} `json:"value"`
	}

	if err := c.call(ctx, req, &result); err != nil {
		return 0, 0, err
	}

	var amount uint64
	fmt.Sscanf(result.Value.Amount, "%d", &amount)
	return amount, result.Value.Decimals, nil
}

func (c *RPCClient) call(ctx context.Context, req RPCRequest, result interface{}) error {
	// Check circuit breaker
	if c.isCircuitOpen() {
		// Try fallback
		return c.callURL(ctx, c.fallbackURL, req, result)
	}

	err := c.callURL(ctx, c.primaryURL, req, result)
	if err != nil {
		c.recordFailure()
		// Try fallback
		log.Warn().Err(err).Msg("primary RPC failed, trying fallback")
		return c.callURL(ctx, c.fallbackURL, req, result)
	}

	c.recordSuccess()
	return nil
}

func (c *RPCClient) callURL(ctx context.Context, url string, rpcReq RPCRequest, result interface{}) error {
	body, err := json.Marshal(rpcReq)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http status %d: %s", resp.StatusCode, string(respBody))
	}

	var rpcResp RPCResponse
	// Optimized: Use Decoder to stream response instead of ReadAll+Unmarshal
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if rpcResp.Error != nil {
		return rpcResp.Error
	}

	if err := json.Unmarshal(rpcResp.Result, result); err != nil {
		return fmt.Errorf("unmarshal result: %w", err)
	}

	return nil
}

// Circuit breaker methods
func (c *RPCClient) isCircuitOpen() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.circuitOpen {
		return false
	}

	// Check if circuit should reset (30 seconds)
	if time.Since(c.lastFailure) > 30*time.Second {
		return false
	}

	return true
}

func (c *RPCClient) recordFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.failures++
	c.lastFailure = time.Now()

	// Open circuit after 5 consecutive failures
	if c.failures >= 5 {
		c.circuitOpen = true
		log.Warn().Msg("RPC circuit breaker opened")
	}
}

func (c *RPCClient) recordSuccess() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.failures = 0
	c.circuitOpen = false
}

// LatencyMs returns estimated latency to RPC (for display)
func (c *RPCClient) LatencyMs() int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	_, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		return -1
	}
	return time.Since(start).Milliseconds()
}

// SignatureStatus represents the status of a transaction signature
type SignatureStatus struct {
	Slot               uint64  `json:"slot"`
	Confirmations      *uint64 `json:"confirmations"` // nil = finalized
	Err                interface{} `json:"err"`       // nil = success, object = error details
	ConfirmationStatus string  `json:"confirmationStatus"` // "processed", "confirmed", "finalized"
}

// GetSignatureStatuses checks the status of transaction signatures
func (c *RPCClient) GetSignatureStatuses(ctx context.Context, signatures []string) ([]*SignatureStatus, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getSignatureStatuses",
		Params: []interface{}{
			signatures,
			map[string]bool{"searchTransactionHistory": true},
		},
	}

	var result struct {
		Value []*SignatureStatus `json:"value"`
	}

	if err := c.call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result.Value, nil
}

// CheckTransaction checks a single transaction and returns status details
func (c *RPCClient) CheckTransaction(ctx context.Context, signature string) (*TxCheckResult, error) {
	statuses, err := c.GetSignatureStatuses(ctx, []string{signature})
	if err != nil {
		return nil, err
	}

	result := &TxCheckResult{
		Signature: signature,
	}

	if len(statuses) == 0 || statuses[0] == nil {
		result.Status = "NOT_FOUND"
		result.Message = "Transaction not found (may still be processing)"
		return result, nil
	}

	status := statuses[0]
	result.Slot = status.Slot
	result.ConfirmationStatus = status.ConfirmationStatus

	if status.Confirmations != nil {
		result.Confirmations = *status.Confirmations
	} else {
		result.Confirmations = 0 // Finalized
	}

	if status.Err == nil {
		result.Status = "SUCCESS"
		result.Message = fmt.Sprintf("Transaction confirmed (%s)", status.ConfirmationStatus)
	} else {
		result.Status = "FAILED"
		// Extract error details
		errBytes, _ := json.Marshal(status.Err)
		result.Message = string(errBytes)
		result.ErrorDetails = status.Err
	}

	return result, nil
}

// TxCheckResult is a human-readable transaction check result
type TxCheckResult struct {
	Signature          string
	Status             string      // "SUCCESS", "FAILED", "NOT_FOUND", "PENDING"
	Message            string
	Slot               uint64
	Confirmations      uint64
	ConfirmationStatus string
	ErrorDetails       interface{}
}

// String returns a formatted string of the result
func (r *TxCheckResult) String() string {
	if r.Status == "SUCCESS" {
		return fmt.Sprintf("✅ %s | Slot: %d | Status: %s", r.Status, r.Slot, r.ConfirmationStatus)
	} else if r.Status == "FAILED" {
		return fmt.Sprintf("❌ %s | Slot: %d | Error: %s", r.Status, r.Slot, r.Message)
	}
	return fmt.Sprintf("⏳ %s | %s", r.Status, r.Message)
}

// TokenAccountInfo holds token account data
type TokenAccountInfo struct {
	Address  string
	Mint     string
	Amount   uint64
	Decimals uint8
}

const (
	TokenProgramID     = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	Token2022ProgramID = "TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"
)

// GetTokenAccountsByOwner fetches all token accounts for an owner.
// If mint is non-empty, filters by mint.
// If mint is empty, checks both Token Program and Token-2022 Program.
func (c *RPCClient) GetTokenAccountsByOwner(ctx context.Context, owner, mint string) ([]TokenAccountInfo, error) {
	// If mint is specified, we just query for that mint (which works for both programs implicitly via RPC usually,
	// but standard getProgramAccounts or getTokenAccountsByOwner requires programId or mint).
	// When mint is provided, Solana RPC `getTokenAccountsByOwner` handles it if we pass the mint filter.
	// But wait, `getTokenAccountsByOwner` requires a programId argument in some versions or `mint` in config?
	// The standard JSON RPC takes (pubkey, config). Config can have mint.
	// Actually, the method signature is `getTokenAccountsByOwner(pubkey, {mint: ...} | {programId: ...}, ...)`
	// So we can only filter by ONE.

	if mint != "" {
		// If mint is provided, we try to find it.
		// NOTE: We might need to try both programs if the mint could be on either and we don't know.
		// But usually filtering by mint is sufficient if the RPC supports scanning both?
		// No, `getTokenAccountsByOwner` requires specifying which program we are querying accounts FOR.
		// Wait, the params are `(string, object, object)`. The first object MUST contain `mint` OR `programId`.
		// It does NOT specify the program ID of the token account itself if `mint` is used?
		// Actually, `getTokenAccountsByOwner` is a method OF a program? No, it's a general RPC.
		// Documentation says: "Required: One of the following: mint, programId".
		// But it implies searching ALL token accounts owned by `pubkey`.
		// However, does it search across ALL Token programs?
		// Usually it defaults to the original Token Program if not specified?
		// Actually, looking at docs, you usually have to specify the program ID in the filter if you want all accounts.
		// If you filter by mint, it should find it regardless of program?
		// Let's stick to the previous implementation for `mint != ""` case which used `mint` filter.
		// But for the bulk fetch, we need to be explicit.

		return c.fetchTokenAccounts(ctx, owner, map[string]string{"mint": mint})
	}

	// If no mint, we want ALL accounts. We must query both programs.
	// 1. Token Program
	accounts, err := c.fetchTokenAccounts(ctx, owner, map[string]string{"programId": TokenProgramID})
	if err != nil {
		return nil, err
	}

	// 2. Token-2022 Program
	accounts2022, err := c.fetchTokenAccounts(ctx, owner, map[string]string{"programId": Token2022ProgramID})
	if err != nil {
		// Bolt Safety: If we fail to fetch Token-2022 accounts, we must fail the whole batch.
		// Returning partial data would cause the executor to think Token-2022 positions have 0 balance,
		// leading to false "sold/failed" state and incorrect PnL (-100%).
		return nil, fmt.Errorf("failed to fetch Token-2022 accounts: %w", err)
	}
	accounts = append(accounts, accounts2022...)

	return accounts, nil
}

func (c *RPCClient) fetchTokenAccounts(ctx context.Context, owner string, filter map[string]string) ([]TokenAccountInfo, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getTokenAccountsByOwner",
		Params: []interface{}{
			owner,
			filter,
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
		var amount uint64
		fmt.Sscanf(v.Account.Data.Parsed.Info.TokenAmount.Amount, "%d", &amount)
		accounts = append(accounts, TokenAccountInfo{
			Address:  v.Pubkey,
			Mint:     v.Account.Data.Parsed.Info.Mint,
			Amount:   amount,
			Decimals: v.Account.Data.Parsed.Info.TokenAmount.Decimals,
		})
	}

	return accounts, nil
}
