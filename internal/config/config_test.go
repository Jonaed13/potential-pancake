package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDynamicURLGeneration(t *testing.T) {
	// Setup environment
	os.Setenv("SHYFT_API_KEY", "test-shyft-key")
	os.Setenv("HELIUS_API_KEY", "test-helius-key")
	defer os.Unsetenv("SHYFT_API_KEY")
	defer os.Unsetenv("HELIUS_API_KEY")

	// Create temp config file
	content := `
rpc:
    shyft_url: https://rpc.shyft.to
    fallback_url: https://mainnet.helius-rpc.com
    shyft_api_key_env: SHYFT_API_KEY
    helius_api_key_env: HELIUS_API_KEY
websocket:
    shyft_url: wss://rpc.shyft.to
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	// Load manager
	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Test Shyft RPC URL
	shyftURL := m.GetShyftRPCURL()
	expectedShyft := "https://rpc.shyft.to?api_key=test-shyft-key"
	if shyftURL != expectedShyft {
		t.Errorf("GetShyftRPCURL = %q, want %q", shyftURL, expectedShyft)
	}

	// Test Fallback RPC URL (Helius)
	fallbackURL := m.GetFallbackRPCURL()
	// Check contains both base and key param (order might vary if logic changes, but here it's simple append)
	if !strings.Contains(fallbackURL, "https://mainnet.helius-rpc.com") || !strings.Contains(fallbackURL, "api-key=test-helius-key") {
		t.Errorf("GetFallbackRPCURL = %q, want it to contain base url and api key", fallbackURL)
	}

	// Test WebSocket URL
	wsURL := m.GetShyftWSURL()
	expectedWS := "wss://rpc.shyft.to?api_key=test-shyft-key"
	if wsURL != expectedWS {
		t.Errorf("GetShyftWSURL = %q, want %q", wsURL, expectedWS)
	}
}

func TestDynamicURLGeneration_ExistingQueryParams(t *testing.T) {
	// Setup environment
	os.Setenv("SHYFT_API_KEY", "test-shyft-key")
	defer os.Unsetenv("SHYFT_API_KEY")

	// Create temp config file with existing query params
	content := `
rpc:
    shyft_url: https://rpc.shyft.to?param=value
    shyft_api_key_env: SHYFT_API_KEY
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	// Load manager
	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Test Shyft RPC URL
	shyftURL := m.GetShyftRPCURL()
	expectedShyft := "https://rpc.shyft.to?param=value&api_key=test-shyft-key"
	if shyftURL != expectedShyft {
		t.Errorf("GetShyftRPCURL = %q, want %q", shyftURL, expectedShyft)
	}
}
