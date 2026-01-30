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

	// RPC
	shyftKey := os.Getenv("SHYFT_API_KEY")
	heliusKey := os.Getenv("HELIUS_API_KEY")

	rpcURL := "https://rpc.shyft.to?api_key=" + shyftKey
	fallbackURL := "https://mainnet.helius-rpc.com/?api-key=" + heliusKey

	if shyftKey == "" {
		fmt.Println("Warning: SHYFT_API_KEY not set")
		rpcURL = "https://rpc.shyft.to"
	}
	if heliusKey == "" {
		fmt.Println("Warning: HELIUS_API_KEY not set")
		fallbackURL = "https://mainnet.helius-rpc.com"
	}

	rpc := blockchain.NewRPCClient(rpcURL, fallbackURL, "")

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
