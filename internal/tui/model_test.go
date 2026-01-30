package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestKeyBindings(t *testing.T) {
	// Reset global theme state
	CurrentThemeIndex = 0

	// Initialize model with nil config (NewModel handles nil config for structure init)
	m := NewModel(nil)
	m.UIMode = 1 // Force Classic Mode where '3' is Metrics, not Theme

	// 1. Test 't' -> Trades
	// This ensures 't' is correctly routed to Trades and doesn't conflict or do nothing
	msgTrades := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}, Alt: false}

	// We need to simulate the tea.KeyMsg behavior where String() returns "t"
	// key.Matches relies on msg.String() matching one of the keys in WithKeys()

	updatedM, _ := m.Update(msgTrades)
	newM, ok := updatedM.(Model)
	if !ok {
		t.Fatalf("Update did not return a Model")
	}

	if newM.CurrentScreen != ScreenTrades {
		t.Errorf("Expected 't' to switch to ScreenTrades, got %v", newM.CurrentScreen)
	}

	// Verify Theme did NOT cycle (should still be 0)
	if CurrentThemeIndex != 0 {
		t.Errorf("Expected 't' NOT to cycle theme, but index changed to %d", CurrentThemeIndex)
	}

	// 2. Test 'T' -> Theme
	// This is the new binding we want to introduce.
	// Reset screen to Dashboard to ensure global keys are handled
	newM.CurrentScreen = ScreenDashboard

	msgTheme := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T'}, Alt: false}

	updatedM2, _ := newM.Update(msgTheme)
	_ = updatedM2.(Model)

	// verification: Theme should cycle to next one (index 1)
	if CurrentThemeIndex != 1 {
		t.Errorf("Expected 'T' to cycle theme to 1, got %d", CurrentThemeIndex)
	}
}

// Helper to simulate key.Matches behavior if needed, but tea.KeyMsg matches standard bubbletea behavior
