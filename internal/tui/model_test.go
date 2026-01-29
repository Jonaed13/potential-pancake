package tui

import (
	"testing"
)

func TestNewConfigModal(t *testing.T) {
	modal := NewConfigModal(nil)

	if len(modal.Fields) != len(modal.Descriptions) {
		t.Errorf("Fields length (%d) does not match Descriptions length (%d)", len(modal.Fields), len(modal.Descriptions))
	}

	expectedFields := []string{"MinEntry", "TakeProfit", "MaxAlloc", "MaxPos", "PrioFee", "AutoTrade"}
	for i, f := range expectedFields {
		if modal.Fields[i] != f {
			t.Errorf("Expected field %s at index %d, got %s", f, i, modal.Fields[i])
		}
	}

	// Verify descriptions are populated
	for i, d := range modal.Descriptions {
		if d == "" {
			t.Errorf("Description at index %d is empty", i)
		}
	}
}
