package tui

import (
	"strings"
	"testing"
	"solana-pump-bot/internal/config"
)

// TestConfigModalDescription verifies that GetDescription returns the correct help text
// for each configuration option.
func TestConfigModalDescription(t *testing.T) {
	// Create a dummy config manager to satisfy the struct type, though we pass nil below.
	var cfg *config.Manager = nil

	cm := ConfigModal{
		Cfg: cfg,
		Fields: []string{"MinEntry", "TakeProfit", "MaxAlloc", "MaxPos", "PrioFee", "AutoTrade"},
		Selected: 0,
	}

	tests := []struct {
		index    int
		expected string
	}{
		{0, "Minimum signal confidence"},
		{1, "Multiplier to exit position"},
		{2, "Max percentage of wallet"},
		{3, "Maximum concurrent open positions"},
		{4, "Jito/Priority fee"},
		{5, "Toggle automated trading"},
		{99, "Adjust configuration values"}, // Default case
	}

	for _, tc := range tests {
		cm.Selected = tc.index
		desc := cm.GetDescription()

		if len(desc) < 5 {
			t.Errorf("Description for index %d is too short: %s", tc.index, desc)
		}

		if !strings.Contains(desc, tc.expected) {
			t.Errorf("For index %d, expected description containing '%s', got '%s'", tc.index, tc.expected, desc)
		}
	}
}
