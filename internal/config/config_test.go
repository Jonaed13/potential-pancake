package config

import (
	"os"
	"testing"
)

func TestGetHeliusAPIKey(t *testing.T) {
	// Create a temporary config file
	configContent := `
rpc:
  helius_api_key_env: TEST_HELIUS_KEY
`
	tmpfile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Set environment variable
	os.Setenv("TEST_HELIUS_KEY", "secret-helius-key")
	defer os.Unsetenv("TEST_HELIUS_KEY")

	// Load config
	mgr, err := NewManager(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Verify GetHeliusAPIKey
	if key := mgr.GetHeliusAPIKey(); key != "secret-helius-key" {
		t.Errorf("expected 'secret-helius-key', got '%s'", key)
	}
}

func TestGetHeliusAPIKeyDefault(t *testing.T) {
	// Create a temporary config file without helius_api_key_env
	configContent := `
rpc:
  shyft_url: "http://example.com"
`
	tmpfile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Set default environment variable
	os.Setenv("HELIUS_API_KEY", "default-helius-key")
	defer os.Unsetenv("HELIUS_API_KEY")

	// Load config
	mgr, err := NewManager(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Verify GetHeliusAPIKey returns value from default env var
	if key := mgr.GetHeliusAPIKey(); key != "default-helius-key" {
		t.Errorf("expected 'default-helius-key', got '%s'", key)
	}
}
