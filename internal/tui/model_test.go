package tui

import (
	"solana-pump-bot/internal/config"
	"testing"
)

func TestNewConfigModal(t *testing.T) {
	// Create a dummy config manager
	// Since NewManager requires a file, we can just pass nil to NewConfigModal
	// because NewConfigModal doesn't use the config manager immediately for initialization
	// of the struct fields we care about (Fields, Descriptions).
	// But it stores it in Cfg.

	cfg := &config.Manager{}

	cm := NewConfigModal(cfg)

	// Check Fields
	if len(cm.Fields) != 6 {
		t.Errorf("Expected 6 fields, got %d", len(cm.Fields))
	}

	// Check Descriptions
	if len(cm.Descriptions) != 6 {
		t.Errorf("Expected 6 descriptions, got %d", len(cm.Descriptions))
	}

	if len(cm.Descriptions) != len(cm.Fields) {
		t.Errorf("Mismatch between fields count (%d) and descriptions count (%d)", len(cm.Fields), len(cm.Descriptions))
	}

	expectedDesc := "Buy when token price increases by this % (e.g., 50.0)"
	if cm.Descriptions[0] != expectedDesc {
		t.Errorf("Expected description[0] to be %q, got %q", expectedDesc, cm.Descriptions[0])
	}
}
