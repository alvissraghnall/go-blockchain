package transaction

import (
	"testing"

	"go-blockchain/wallet"
)

func TestTransactionSigningAndVerification(t *testing.T) {
	// Generate a key pair for signing and verification
	privateKey, publicKey, err := wallet.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	// Define transaction inputs and outputs
	inputs := []TransactionInput{
		{PreviousTxID: "abcd1234", OutputIndex: 0},
	}
	outputs := []TransactionOutput{
		{Address: "recipient_address", Amount: 10.0},
	}

	// Create a new transaction
	tx := NewTransaction(inputs, outputs)

	// Ensure the transaction is initialized correctly
	if len(tx.Inputs) != 1 || len(tx.Outputs) != 1 {
		t.Errorf("Transaction initialization failed: expected 1 input and 1 output")
	}

	// Sign the transaction
	err = tx.SignTransaction(privateKey)
	if err != nil {
		t.Fatalf("SignTransaction failed: %v", err)
	}

	// Ensure the signature is added
	if len(tx.SignatureR) == 0 || len(tx.SignatureS) == 0 {
		t.Errorf("Transaction signature missing after signing")
	}

	// Verify the transaction
	if !tx.VerifyTransaction(publicKey) {
		t.Errorf("Transaction verification failed: signature is invalid")
	}

	// Modify the transaction to simulate tampering
	tx.Outputs[0].Amount = 20.0

	// Verify the tampered transaction
	if tx.VerifyTransaction(publicKey) {
		t.Errorf("Tampered transaction verification should have failed")
	}
}

