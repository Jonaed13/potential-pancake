package blockchain

import (
	"crypto/ed25519"
	"testing"
	"github.com/mr-tron/base58"
)

func TestSignSerializedTransaction_SimulationDummy(t *testing.T) {
	// Create a dummy wallet
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Create wallet using public constructor or by mocking if possible.
	// Since fields are private, we use NewWallet with base58 string.
	privKeyBase58 := base58.Encode(privKey)
	wallet, err := NewWallet(privKeyBase58)
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Ensure public key matches
	if base58.Encode(pubKey) != wallet.Address() {
		t.Errorf("Wallet address mismatch. Got %s, want %s", wallet.Address(), base58.Encode(pubKey))
	}

	// Create TransactionBuilder (dependencies can be nil for this test as we only use SignSerializedTransaction)
	tb := NewTransactionBuilder(wallet, nil, 0)

	// The dummy string returned by simulation mode
	dummyTx := "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAA=="

	// Try to sign it
	signedTx, err := tb.SignSerializedTransaction(dummyTx)
	if err != nil {
		t.Fatalf("SignSerializedTransaction failed with dummy tx: %v", err)
	}

	if signedTx == "" {
		t.Error("SignSerializedTransaction returned empty string")
	}

	// Ensure it didn't crash and returned something
	t.Logf("Signed dummy tx: %s", signedTx)
}
