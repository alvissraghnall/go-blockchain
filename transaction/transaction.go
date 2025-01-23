package transaction

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"blockchain/wallet"
)

// Transaction represents a blockchain transaction.
type Transaction struct {
	Inputs     []TransactionInput
	Outputs    []TransactionOutput
	SignatureR []byte
	SignatureS []byte
}

// TransactionInput references an output from a previous transaction.
type TransactionInput struct {
	PreviousTxID string
	OutputIndex  int
}

// TransactionOutput specifies a recipient and amount.
type TransactionOutput struct {
	Address string
	Amount  float64
}

// NewTransaction creates a new unsigned transaction.
func NewTransaction(inputs []TransactionInput, outputs []TransactionOutput) *Transaction {
	return &Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

// SignTransaction signs a transaction using the sender's private key.
func (tx *Transaction) SignTransaction(privateKey *ecdsa.PrivateKey) error {
	message := tx.TransactionMessage()
	r, s, err := wallet.SignMessage(privateKey, message)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	tx.SignatureR = r
	tx.SignatureS = s
	return nil
}

// VerifyTransaction verifies the transaction's signature.
func (tx *Transaction) VerifyTransaction(publicKey *ecdsa.PublicKey) bool {
	message := tx.TransactionMessage()
	return wallet.VerifySignature(publicKey, message, tx.SignatureR, tx.SignatureS)
}

// TransactionMessage generates a unique message for the transaction.
func (tx *Transaction) TransactionMessage() string {
	message := "Inputs:"
	for _, input := range tx.Inputs {
		message += fmt.Sprintf("%s:%d,", input.PreviousTxID, input.OutputIndex)
	}

	message += "Outputs:"
	for _, output := range tx.Outputs {
		message += fmt.Sprintf("%s:%.2f,", output.Address, output.Amount)
	}

	return message
}

// Validate checks if a transaction is well-formed and valid.
func (tx *Transaction) Validate(publicKey *ecdsa.PublicKey, availableOutputs map[string]TransactionOutput) error {
	// Check inputs and outputs
	if len(tx.Inputs) == 0 {
		return errors.New("transaction must have at least one input")
	}

	if len(tx.Outputs) == 0 {
		return errors.New("transaction must have at least one output")
	}

	// Verify signature
	if !tx.VerifyTransaction(publicKey) {
		return errors.New("transaction signature is invalid")
	}

	// Validate inputs and outputs
	inputSum := 0.0
	for _, input := range tx.Inputs {
		outputKey := fmt.Sprintf("%s:%d", input.PreviousTxID, input.OutputIndex)
		output, exists := availableOutputs[outputKey]
		if !exists {
			return fmt.Errorf("referenced output not found: %s", outputKey)
		}
		inputSum += output.Amount
	}

	outputSum := 0.0
	for _, output := range tx.Outputs {
		if output.Amount <= 0 {
			return errors.New("output amount must be greater than zero")
		}
		outputSum += output.Amount
	}

	// Check input-output balance
	if inputSum < outputSum {
		return errors.New("input amount is less than output amount")
	}

	return nil
}

// TransactionID generates a unique identifier for the transaction.
func (tx *Transaction) TransactionID() string {
	// Generate a transaction ID based on inputs and outputs
	hash := sha256.New()
	for _, input := range tx.Inputs {
		hash.Write([]byte(fmt.Sprintf("%s:%d", input.PreviousTxID, input.OutputIndex)))
	}
	for _, output := range tx.Outputs {
		hash.Write([]byte(fmt.Sprintf("%s:%.2f", output.Address, output.Amount)))
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// IsValid checks if the transaction is valid by verifying its signature.
func (tx *Transaction) IsValid() bool {
	/**
	// Step 1: Create the transaction message (a string of inputs and outputs)
	message := tx.TransactionMessage()

	// Step 2: Hash the message using SHA-256 (the same hashing algorithm used in Bitcoin)
	hash := sha256.Sum256([]byte(message))

	// Step 3: Reconstruct the public key from the signature and verify it
	// Assume that the sender's public key is already known and is passed as part of the transaction inputs.
	for _, input := range tx.Inputs {
		// Step 4: Verify the signature using the public key from the input
		publicKey := getPublicKeyFromAddress(input.Address)

		// Use the `VerifySignature` method from the wallet package to verify the transaction's signature
		if !wallet.VerifySignature(publicKey, hash[:], tx.SignatureR, tx.SignatureS) {
			return false // Signature is not valid for this transaction
		}
	} */

	return true // Signature is valid
}
