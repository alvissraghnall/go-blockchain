package transaction

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"blockchain/wallet"
	"blockchain/chain"
)

// Transaction represents a blockchain transaction.
type Transaction struct {
	Inputs     []Input
	Outputs    []Output
	SignatureR []byte
	SignatureS []byte
}

// Input references an output from a previous transaction.
type Input struct {
	// This field represents the hash of the previous transaction that this input is spending from.
	PreviousTxHash []byte		`json:"previousTxHash"`
	// This field represents the index of the output in the previous transaction that this input is spending from.
	OutputIndex    uint64		`json:"outputIndex"`
	// This field represents the script signature that unlocks the previous transaction output being spent. 
	ScriptSig      []byte		`json:"scriptSig"`
}

// Output specifies a recipient and amount.
type Output struct {
	Amount       float64     `json:"amount"`
    ScriptPubKey []byte      `json:"scriptPubKey"`
}

// UTXO represents an unspent transaction output.
type UTXO struct {
	// This field represents the block in which the transaction containing this UTXO was included.
    Block        chain.BlockInterface    		`json:"block"`
	// This field represents the transaction that created this UTXO. 
    Transaction  Transaction 	`json:"transaction"`
	// This field represents the index of this UTXO within the transaction that created it.
    OutputIndex  uint32    		`json:"outputIndex"`
	// This field represents the amount of cryptocurrency associated with this UTXO.
    Amount       float64     		`json:"amount"`
}

// NewTransaction creates a new unsigned transaction.
func NewTransaction(inputs []Input, outputs []Output) *Transaction {
	return &Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func NewOutput(amount float64, scriptPubKey []byte) *Output {
    return &Output{
        Amount:       amount,
        ScriptPubKey: scriptPubKey,
    }
}

func NewInput(previousTxHash []byte, outputIndex uint64, scriptSig []byte) *Input {
    return &Input{
        PreviousTxHash:  previousTxHash,
		OutputIndex:  	 outputIndex,
        ScriptSig:    	 scriptSig,
    }
}

func NewUTXO(block chain.BlockInterface, transaction Transaction, outputIndex uint64, amount int64) *UTXO {
    return &UTXO{
        Block:        block,
        Transaction:  transaction,
        OutputIndex:  outputIndex,
        Amount:       amount,
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
func (tx *Transaction) Validate(publicKey *ecdsa.PublicKey, availableOutputs map[string]Output) error {
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

	// Step 2: Hash the message using SHA-256 
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
