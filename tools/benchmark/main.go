package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"
)

// Test wallet for getBalance (random known wallet)
const testWallet = "vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg"

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}

type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func main() {
	fmt.Println("🔬 RPC Endpoint Benchmark")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Printf("Time: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// RPC endpoints
	shyftKey := os.Getenv("SHYFT_API_KEY")
	heliusKey := os.Getenv("HELIUS_API_KEY")

	endpoints := make(map[string]string)

	if shyftKey != "" {
		endpoints["Shyft"] = "https://rpc.shyft.to?api_key=" + shyftKey
	} else {
		fmt.Println("Warning: SHYFT_API_KEY not set (skipping Shyft)")
	}

	if heliusKey != "" {
		endpoints["Helius"] = "https://mainnet.helius-rpc.com/?api-key=" + heliusKey
	} else {
		fmt.Println("Warning: HELIUS_API_KEY not set (skipping Helius)")
	}

	if len(endpoints) == 0 {
		fmt.Println("No API keys provided! Set SHYFT_API_KEY or HELIUS_API_KEY.")
		os.Exit(1)
	}

	// Methods to test
	allMethods := []struct {
		name   string
		method string
		params []interface{}
	}{
		{"getLatestBlockhash", "getLatestBlockhash", []interface{}{map[string]string{"commitment": "confirmed"}}},
		{"getBalance", "getBalance", []interface{}{testWallet}},
		{"getSlot", "getSlot", nil},
		{"getHealth", "getHealth", nil},
		{"getVersion", "getVersion", nil},
		{"getBlockHeight", "getBlockHeight", nil},
		{"getRecentPrioritizationFees", "getRecentPrioritizationFees", nil},
	}

	// Bot-critical methods
	botMethods := []struct {
		name   string
		method string
		params []interface{}
	}{
		{"getLatestBlockhash", "getLatestBlockhash", []interface{}{map[string]string{"commitment": "confirmed"}}},
		{"getBalance", "getBalance", []interface{}{testWallet}},
		{"getSlot", "getSlot", nil},
	}

	iterations := 20

	fmt.Println("📊 GENERAL RPC BENCHMARK (all methods)")
	fmt.Println("-" + string(make([]byte, 60)))
	fmt.Printf("Iterations per method: %d\n\n", iterations)

	for name, url := range endpoints {
		fmt.Printf("\n🔗 %s\n", name)
		displayURL := url
		if len(url) > 50 {
			displayURL = url[:50] + "..."
		}
		fmt.Printf("   URL: %s\n\n", displayURL)

		for _, m := range allMethods {
			latencies := benchmark(url, m.method, m.params, iterations)
			if len(latencies) == 0 {
				fmt.Printf("   %-30s ❌ FAILED\n", m.name)
				continue
			}
			p50, p95, p99, avg := calcStats(latencies)
			fmt.Printf("   %-30s p50: %4dms  p95: %4dms  p99: %4dms  avg: %4dms\n",
				m.name, p50, p95, p99, avg)
		}
	}

	fmt.Println()
	fmt.Println("🤖 BOT-CRITICAL METHODS (higher iterations)")
	fmt.Println("-" + string(make([]byte, 60)))
	iterations = 50
	fmt.Printf("Iterations per method: %d\n", iterations)

	for name, url := range endpoints {
		fmt.Printf("\n🔗 %s\n", name)

		var allLatencies []int64
		for _, m := range botMethods {
			latencies := benchmark(url, m.method, m.params, iterations)
			if len(latencies) == 0 {
				fmt.Printf("   %-30s ❌ FAILED\n", m.name)
				continue
			}
			allLatencies = append(allLatencies, latencies...)
			p50, p95, p99, avg := calcStats(latencies)
			fmt.Printf("   %-30s p50: %4dms  p95: %4dms  p99: %4dms  avg: %4dms\n",
				m.name, p50, p95, p99, avg)
		}

		// Overall stats for bot methods
		if len(allLatencies) > 0 {
			p50, p95, p99, avg := calcStats(allLatencies)
			fmt.Printf("\n   %-30s p50: %4dms  p95: %4dms  p99: %4dms  avg: %4dms\n",
				"📈 OVERALL (bot methods)", p50, p95, p99, avg)
		}
	}

	// Summary comparison
	fmt.Println()
	fmt.Println("📋 SUMMARY COMPARISON")
	fmt.Println("-" + string(make([]byte, 60)))

	for _, m := range botMethods {
		fmt.Printf("\n%s:\n", m.name)
		for name, url := range endpoints {
			latencies := benchmark(url, m.method, m.params, 30)
			if len(latencies) > 0 {
				_, _, p99, avg := calcStats(latencies)
				fmt.Printf("   %-10s → avg: %4dms, p99: %4dms\n", name, avg, p99)
			}
		}
	}

	fmt.Println("\n✅ Benchmark complete")
}

func benchmark(url, method string, params []interface{}, iterations int) []int64 {
	client := &http.Client{Timeout: 10 * time.Second}
	latencies := make([]int64, 0, iterations)

	for i := 0; i < iterations; i++ {
		start := time.Now()

		req := RPCRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  method,
			Params:  params,
		}

		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(httpReq)
		if err != nil {
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var rpcResp RPCResponse
		json.Unmarshal(respBody, &rpcResp)

		if rpcResp.Error != nil {
			continue
		}

		elapsed := time.Since(start).Milliseconds()
		latencies = append(latencies, elapsed)

		// Small delay between requests
		time.Sleep(50 * time.Millisecond)
	}

	return latencies
}

func calcStats(latencies []int64) (p50, p95, p99, avg int64) {
	if len(latencies) == 0 {
		return 0, 0, 0, 0
	}

	sorted := make([]int64, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	n := len(sorted)
	p50 = sorted[n*50/100]
	p95 = sorted[n*95/100]
	if n > 1 {
		p99 = sorted[n*99/100]
	} else {
		p99 = sorted[n-1]
	}

	var sum int64
	for _, l := range sorted {
		sum += l
	}
	avg = sum / int64(n)

	return
}
