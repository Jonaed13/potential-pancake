package blockchain

import (
	"fmt"
	"strconv"
	"testing"
)

// Ensure ParseUint behaves as expected (returns 0 on error)
func TestParseAmount_Behavior(t *testing.T) {
	cases := []struct {
		input    string
		expected uint64
	}{
		{"12345", 12345},
		{"0", 0},
		{"invalid", 0},
		{"", 0},
		{"-1", 0}, // ParseUint fails on negative sign, returns 0
	}

	for _, tc := range cases {
		// Old behavior simulation
		var sscanfAmount uint64
		fmt.Sscanf(tc.input, "%d", &sscanfAmount)

		// New behavior
		strconvAmount, _ := strconv.ParseUint(tc.input, 10, 64)

		if sscanfAmount != tc.expected {
			t.Errorf("Sscanf assumption failed for '%s': got %d, want %d", tc.input, sscanfAmount, tc.expected)
		}
		if strconvAmount != tc.expected {
			t.Errorf("Strconv behavior mismatch for '%s': got %d, want %d", tc.input, strconvAmount, tc.expected)
		}
	}
}

func BenchmarkParseAmount_Sscanf(b *testing.B) {
	input := "1234567890"
	var amount uint64
	for i := 0; i < b.N; i++ {
		fmt.Sscanf(input, "%d", &amount)
	}
}

func BenchmarkParseAmount_Strconv(b *testing.B) {
	input := "1234567890"
	for i := 0; i < b.N; i++ {
		_, _ = strconv.ParseUint(input, 10, 64)
	}
}
