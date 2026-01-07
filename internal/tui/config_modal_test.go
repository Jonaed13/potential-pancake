package tui

import (
	"testing"
)

func TestConfigModal_GetDescription(t *testing.T) {
	// Create a nil config manager is fine for GetDescription as it doesn't use it
	cm := NewConfigModal(nil)

	tests := []struct {
		selected int
		want     string
	}{
		{0, "Minimum signal strength to enter a trade"},
		{1, "Multiplier target to sell (e.g. 2.0x)"},
		{2, "Max % of wallet to use per trade"},
		{3, "Maximum concurrent open positions"},
		{4, "SOL fee paid to miners for speed"},
		{5, "Master switch for automated trading"},
		{99, "Adjust settings"},
	}

	for _, tt := range tests {
		cm.Selected = tt.selected
		got := cm.GetDescription()
		if got != tt.want {
			t.Errorf("GetDescription(%d) = %q, want %q", tt.selected, got, tt.want)
		}
	}
}
