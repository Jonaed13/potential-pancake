package config

import (
	"os"
	"testing"
)

func TestGetShyftRPCURL(t *testing.T) {
	// Setup temporary config
	os.Setenv("SHYFT_API_KEY", "test-api-key")
	defer os.Unsetenv("SHYFT_API_KEY")

	// Create a dummy manager with config
	cfg := &Config{
		RPC: RPCConfig{
			ShyftURL:       "https://rpc.shyft.to",
			ShyftAPIKeyEnv: "SHYFT_API_KEY",
		},
	}
	m := &Manager{config: cfg}

	// Test case 1: Basic URL
	url := m.GetShyftRPCURL()
	expected := "https://rpc.shyft.to?api_key=test-api-key"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}

	// Test case 2: URL with existing query param
	m.config.RPC.ShyftURL = "https://rpc.shyft.to?foo=bar"
	url = m.GetShyftRPCURL()
	expected = "https://rpc.shyft.to?foo=bar&api_key=test-api-key"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}

	// Test case 3: URL with API key already present
	m.config.RPC.ShyftURL = "https://rpc.shyft.to?api_key=existing-key"
	url = m.GetShyftRPCURL()
	expected = "https://rpc.shyft.to?api_key=existing-key"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}

	// Test case 4: API key env var missing
	os.Unsetenv("SHYFT_API_KEY")
	m.config.RPC.ShyftURL = "https://rpc.shyft.to"
	url = m.GetShyftRPCURL()
	expected = "https://rpc.shyft.to"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestGetFallbackRPCURL(t *testing.T) {
	os.Setenv("HELIUS_API_KEY", "test-helius-key")
	defer os.Unsetenv("HELIUS_API_KEY")

	cfg := &Config{
		RPC: RPCConfig{
			FallbackURL:       "https://mainnet.helius-rpc.com",
			FallbackAPIKeyEnv: "HELIUS_API_KEY",
		},
	}
	m := &Manager{config: cfg}

	// Test case 1: Basic URL (Helius uses api-key)
	url := m.GetFallbackRPCURL()
	expected := "https://mainnet.helius-rpc.com?api-key=test-helius-key"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestGetShyftWSURL(t *testing.T) {
	os.Setenv("SHYFT_API_KEY", "test-ws-key")
	defer os.Unsetenv("SHYFT_API_KEY")

	cfg := &Config{
		WebSocket: WebSocketConfig{
			ShyftURL: "wss://rpc.shyft.to",
		},
		RPC: RPCConfig{
			ShyftAPIKeyEnv: "SHYFT_API_KEY",
		},
	}
	m := &Manager{config: cfg}

	url := m.GetShyftWSURL()
	expected := "wss://rpc.shyft.to?api_key=test-ws-key"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}
