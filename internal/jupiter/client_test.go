package jupiter

import (
	"context"
	"testing"
	"time"
)

func TestGetSwapTransaction_SimulationMode(t *testing.T) {
	// Setup client with simulation mode enabled
	client := NewClient("https://api.jup.ag/swap/v1", 50, 10*time.Second)
	client.SetSimulation(true, 1.0)

	ctx := context.Background()
	inputMint := "So11111111111111111111111111111111111111112"
	outputMint := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v" // USDC
	userPubkey := "DstF19y19y19y19y19y19y19y19y19y19y19y19y19y"
	amount := uint64(1000000)

	// Call GetSwapTransaction
	txStr, err := client.GetSwapTransaction(ctx, inputMint, outputMint, userPubkey, amount)
	if err != nil {
		t.Fatalf("GetSwapTransaction failed in simulation mode: %v", err)
	}

	expected := "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAA=="
	if txStr != expected {
		t.Errorf("Expected dummy transaction %q, got %q", expected, txStr)
	}
}
