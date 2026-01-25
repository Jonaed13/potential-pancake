package signal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestServer_RateLimit(t *testing.T) {
	// Mock handler dependencies
	signalChan := make(chan *Signal, 100)
	minEntry := func() float64 { return 50.0 }
	takeProfit := func() float64 { return 2.0 }
	resolveMint := func(name string) (string, error) { return "MOCK_MINT", nil }

	handler := NewHandler(signalChan, minEntry, takeProfit, resolveMint)
	// Port 0 is fine for testing as we don't Listen()
	server := NewServer("0.0.0.0", 0, handler)

	// Create a payload
	payload := ParsedSignal{
		Text:      "📈 TEST is up 100% 📈",
		MsgID:     123,
		Timestamp: time.Now().Unix(),
	}
	body, _ := json.Marshal(payload)

	// Send rapid requests
	// We expect a rate limit around 5-10 req/s.
	// Sending 50 requests should definitely hit it.

	limitHit := false
	for i := 0; i < 50; i++ {
		req, _ := http.NewRequest("POST", "/signal", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Use Fiber's Test method
		resp, err := server.app.Test(req, 1000) // 1s timeout
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}

		if resp.StatusCode == 429 {
			limitHit = true
			break
		}
	}

	if !limitHit {
		t.Error("Security Vulnerability: Rate limit was NOT hit after 50 requests")
	}
}
