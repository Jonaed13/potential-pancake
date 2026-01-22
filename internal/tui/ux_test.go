package tui

import (
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"solana-pump-bot/internal/config"
)

func TestHelpKeybinding(t *testing.T) {
	// Create temp config file
	f, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	// Write dummy config content
	content := []byte(`
trading:
  min_entry_percent: 50
  take_profit_multiple: 2.0
  max_alloc_percent: 10
  max_open_positions: 5
  auto_trading_enabled: false
fees:
  static_priority_fee_sol: 0.001
`)
	if _, err := f.Write(content); err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Init Config Manager
	cfg, err := config.NewManager(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Init Model
	m := NewModel(cfg)
	m.Width = 80
	m.Height = 24

    // Set initial screen to Logs to test restoration
    m.CurrentScreen = ScreenLogs

	// Test 1: Press '?' should switch to ScreenHelp
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}

	// Update
	updatedModel, _ := m.Update(msg)

	finalModel, ok := updatedModel.(Model)
	if !ok {
		t.Fatal("Model type assertion failed")
	}

	if finalModel.CurrentScreen != ScreenHelp {
		t.Errorf("Expected ScreenHelp, got %v", finalModel.CurrentScreen)
	}

    if finalModel.PreviousScreen != ScreenLogs {
        t.Errorf("Expected PreviousScreen to be ScreenLogs, got %v", finalModel.PreviousScreen)
    }

	// Test 2: Press 'Esc' should switch back to ScreenLogs (Restoration)
	msgEsc := tea.KeyMsg{Type: tea.KeyEscape}
	updatedModel2, _ := finalModel.Update(msgEsc)
	finalModel2, ok := updatedModel2.(Model)
	if !ok {
		t.Fatal("Model type assertion failed")
	}

	if finalModel2.CurrentScreen != ScreenLogs {
		t.Errorf("Expected ScreenLogs after Esc, got %v", finalModel2.CurrentScreen)
	}
}
