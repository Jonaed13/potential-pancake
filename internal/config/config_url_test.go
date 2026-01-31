package config

import (
	"os"
	"testing"
)

func TestURLInjection(t *testing.T) {
	// Setup test config file
	tmpConfig := "test_config_url.yaml"
	content := []byte(`
rpc:
    shyft_url: https://rpc.shyft.to
    fallback_url: https://mainnet.helius-rpc.com
    shyft_api_key_env: TEST_SHYFT_KEY
    fallback_api_key_env: TEST_HELIUS_KEY
websocket:
    shyft_url: wss://rpc.shyft.to
`)
	if err := os.WriteFile(tmpConfig, content, 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpConfig)

	// Set env vars
	os.Setenv("TEST_SHYFT_KEY", "shyft-123")
	os.Setenv("TEST_HELIUS_KEY", "helius-456")
	defer os.Unsetenv("TEST_SHYFT_KEY")
	defer os.Unsetenv("TEST_HELIUS_KEY")

	m, err := NewManager(tmpConfig)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Test Shyft RPC
	got := m.GetShyftRPCURL()
	want := "https://rpc.shyft.to?api_key=shyft-123"
	if got != want {
		t.Errorf("GetShyftRPCURL() = %q, want %q", got, want)
	}

	// Test Fallback RPC
	got = m.GetFallbackRPCURL()
	want = "https://mainnet.helius-rpc.com?api-key=helius-456"
	if got != want {
		t.Errorf("GetFallbackRPCURL() = %q, want %q", got, want)
	}

	// Test Shyft WS
	got = m.GetShyftWSURL()
	want = "wss://rpc.shyft.to?api_key=shyft-123"
	if got != want {
		t.Errorf("GetShyftWSURL() = %q, want %q", got, want)
	}
}

func TestURLInjection_ExistingParams(t *testing.T) {
	// Setup test config file
	tmpConfig := "test_config_existing.yaml"
	content := []byte(`
rpc:
    shyft_url: https://rpc.shyft.to?foo=bar
    fallback_url: https://mainnet.helius-rpc.com
    shyft_api_key_env: TEST_SHYFT_KEY_2
    fallback_api_key_env: TEST_HELIUS_KEY_2
websocket:
    shyft_url: wss://rpc.shyft.to
`)
	if err := os.WriteFile(tmpConfig, content, 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpConfig)

	// Set env vars
	os.Setenv("TEST_SHYFT_KEY_2", "shyft-789")
	defer os.Unsetenv("TEST_SHYFT_KEY_2")

	m, err := NewManager(tmpConfig)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Test Shyft RPC with existing param
	got := m.GetShyftRPCURL()
	want := "https://rpc.shyft.to?foo=bar&api_key=shyft-789"
	if got != want {
		t.Errorf("GetShyftRPCURL() = %q, want %q", got, want)
	}
}

func TestURLInjection_NoEnvKey(t *testing.T) {
	// Setup test config file
	tmpConfig := "test_config_no_key.yaml"
	content := []byte(`
rpc:
    shyft_url: https://rpc.shyft.to
    fallback_url: https://mainnet.helius-rpc.com
    shyft_api_key_env: TEST_SHYFT_KEY_MISSING
    fallback_api_key_env: TEST_HELIUS_KEY_MISSING
websocket:
    shyft_url: wss://rpc.shyft.to
`)
	if err := os.WriteFile(tmpConfig, content, 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpConfig)

	// Ensure env vars are unset
	os.Unsetenv("TEST_SHYFT_KEY_MISSING")
	os.Unsetenv("TEST_HELIUS_KEY_MISSING")

	m, err := NewManager(tmpConfig)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Test Shyft RPC - should be unchanged
	got := m.GetShyftRPCURL()
	want := "https://rpc.shyft.to"
	if got != want {
		t.Errorf("GetShyftRPCURL() = %q, want %q", got, want)
	}
}
