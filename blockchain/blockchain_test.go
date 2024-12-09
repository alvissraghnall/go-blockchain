package blockchain

import (
	"fmt"
	"testing"
	"go-blockchain/transaction"
)

func TestBlockchain(t *testing.T) {
	// Initialize blockchain
	chain := NewBlockchain()

	// Add blocks
	chain.AddBlock([]transaction.Transaction{
		{Inputs: []transaction.TransactionInput{{"tx1", 0}}, Outputs: []transaction.TransactionOutput{{"address1", 10.0}}},
	})
	chain.AddBlock([]transaction.Transaction{
		{Inputs: []transaction.TransactionInput{{"tx2", 1}}, Outputs: []transaction.TransactionOutput{{"address2", 20.0}}},
	})

	// Print the blockchain
	for _, block := range chain.Blocks {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Nonce: %d\n\n", block.Nonce)
	}

	// Validate the blockchain
	if !chain.IsValid() {
		t.Errorf("Blockchain is invalid")
	}
}
