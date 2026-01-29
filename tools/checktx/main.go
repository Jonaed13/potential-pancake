package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"solana-pump-bot/internal/blockchain"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./tools/checktx <TX_SIGNATURE>")
		fmt.Println("Example: go run ./tools/checktx 2gHc4gtPHJgVJhccGytQqivvETZoyfiAu12UTE3vN4v6WPz3mGmPGmwxS7NwbXcv28NAQP6Re8rdi2XS9tU6rMRs")
		os.Exit(1)
	}

	txSig := os.Args[1]

	fmt.Println("📊 TX STATUS CHECKER")
	fmt.Println("===================")
	fmt.Printf("TX: %s\n\n", txSig)

	// RPC (Load from environment since this is a tool)
	shyftKey := os.Getenv("SHYFT_API_KEY")
	heliusKey := os.Getenv("HELIUS_API_KEY")

	// Fallback to hardcoded public nodes if keys missing (for manual testing)
	shyftURL := "https://rpc.shyft.to?api_key=" + shyftKey
	heliusURL := "https://mainnet.helius-rpc.com/?api-key=" + heliusKey

	rpc := blockchain.NewRPCClient(shyftURL, heliusURL, "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := rpc.CheckTransaction(ctx, txSig)
	if err != nil {
		fmt.Printf("❌ RPC Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result.String())
	fmt.Println("")

	if result.Status == "FAILED" {
		fmt.Println("📋 ERROR DETAILS:")
		fmt.Printf("%+v\n", result.ErrorDetails)
	}

	if result.Status == "SUCCESS" {
		fmt.Printf("Slot: %d\n", result.Slot)
		fmt.Printf("Confirmations: %d\n", result.Confirmations)
		fmt.Printf("Status: %s\n", result.ConfirmationStatus)
	}
}
