package blockchain

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkParseUint_Sscanf(b *testing.B) {
	s := "1234567890123456789"
	var v uint64
	for i := 0; i < b.N; i++ {
		fmt.Sscanf(s, "%d", &v)
	}
}

func BenchmarkParseUint_Strconv(b *testing.B) {
	s := "1234567890123456789"
	for i := 0; i < b.N; i++ {
		strconv.ParseUint(s, 10, 64)
	}
}
