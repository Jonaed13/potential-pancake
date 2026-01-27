package signal

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

func TestRateLimiting(t *testing.T) {
	// Create a minimal handler
	handler := &Handler{
		parser: NewParser(),
		signalChan: make(chan *Signal, 10),
		minEntry: func() float64 { return 1.0 },
		takeProfit: func() float64 { return 2.0 },
		resolveMint: func(s string) (string, error) { return "mint", nil },
	}

	server := NewServer("localhost", 0, handler)

	// Valid payload
	payload := []byte(`{"text": "📈 TOKEN is up 50% 📈", "msg_id": 12345}`)

	// Send 5 allowed requests
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("POST", "/signal", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.app.Test(req)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}

		// Should be 200 OK (or whatever success code)
		// Or 400 if validation fails, but NOT 429
		if resp.StatusCode == 429 {
			t.Fatalf("Request %d was rate limited prematurely", i)
		}
	}

	// Send 6th request - should be blocked
	req := httptest.NewRequest("POST", "/signal", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)
	if err != nil {
		t.Fatalf("Request 6 failed: %v", err)
	}

	if resp.StatusCode != 429 {
		t.Errorf("Expected status 429 for 6th request, got %d", resp.StatusCode)
	}
}
