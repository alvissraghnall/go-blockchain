package blockchain

import (
	"fmt"
	"testing"
	"blockchain/transaction"
)


func TestBlockchain(t *testing.T) {
    // Initialize blockchain
    chain := NewBlockchain()

    // Add blocks
    block1 := NewBlock(1, []transaction.Transaction{
        {
            Inputs: []transaction.TransactionInput{
                {"tx1", 0},
            },
            Outputs: []transaction.TransactionOutput{
                {"address1", 10.0},
            },
        },
    }, chain.GetLatestBlock().Hash)

    chain.AddBlock(*block1)

    block2 := NewBlock(2, []transaction.Transaction{
        {
            Inputs: []transaction.TransactionInput{
                {"tx2", 1},
            },
            Outputs: []transaction.TransactionOutput{
                {"address2", 20.0},
            },
        },
    }, block1.Hash)

    chain.AddBlock(*block2)

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
