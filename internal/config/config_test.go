package config

import (
	"os"
	"testing"
)

func TestGetShyftRPCURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		envKey   string
		envValue string
		want     string
	}{
		{
			name:     "No key in URL, Env set",
			url:      "https://rpc.shyft.to",
			envKey:   "SHYFT_API_KEY",
			envValue: "secret123",
			want:     "https://rpc.shyft.to?api_key=secret123",
		},
		{
			name:     "Key already in URL",
			url:      "https://rpc.shyft.to?api_key=embedded",
			envKey:   "SHYFT_API_KEY",
			envValue: "secret123",
			want:     "https://rpc.shyft.to?api_key=embedded",
		},
		{
			name:     "No key in URL, Env empty",
			url:      "https://rpc.shyft.to",
			envKey:   "SHYFT_API_KEY",
			envValue: "",
			want:     "https://rpc.shyft.to",
		},
		{
			name:     "URL with existing params, Env set",
			url:      "https://rpc.shyft.to?foo=bar",
			envKey:   "SHYFT_API_KEY",
			envValue: "secret123",
			want:     "https://rpc.shyft.to?foo=bar&api_key=secret123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock Env
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			cfg := &Config{
				RPC: RPCConfig{
					ShyftURL:       tt.url,
					ShyftAPIKeyEnv: tt.envKey,
				},
			}
			m := &Manager{config: cfg}

			got := m.GetShyftRPCURL()
			if got != tt.want {
				t.Errorf("GetShyftRPCURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFallbackRPCURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		envKey   string
		envValue string
		want     string
	}{
		{
			name:     "No key in URL, Env set",
			url:      "https://helius.xyz",
			envKey:   "HELIUS_API_KEY",
			envValue: "hkey",
			want:     "https://helius.xyz?api-key=hkey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			cfg := &Config{
				RPC: RPCConfig{
					FallbackURL:     tt.url,
					HeliusAPIKeyEnv: tt.envKey,
				},
			}
			m := &Manager{config: cfg}

			got := m.GetFallbackRPCURL()
			if got != tt.want {
				t.Errorf("GetFallbackRPCURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
