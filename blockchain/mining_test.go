package blockchain

import (
	"testing"
	"go-blockchain/transaction"
)

func TestMining(t *testing.T) {
	// Initialize blockchain and transaction pool
	chain := NewBlockchain()
	txPool := transaction.NewTransactionPool()

	// Add some transactions to the pool
	txPool.AddTransaction(transaction.Transaction{
		Inputs:  []transaction.TransactionInput{{"tx1", 0}},
		Outputs: []transaction.TransactionOutput{{"address1", 5.0}},
	})
	txPool.AddTransaction(transaction.Transaction{
		Inputs:  []transaction.TransactionInput{{"tx2", 1}},
		Outputs: []transaction.TransactionOutput{{"address2", 3.0}},
	})

	// Initialize miner
	miner := NewMiner(chain, txPool, 4)

	// Mine a block
	block, err := miner.Mine()
	if err != nil {
		t.Fatalf("Mining failed: %v", err)
	}

	// Verify block content
	if len(block.Transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(block.Transactions))
	}

	// Verify blockchain length
	if len(chain.Blocks) != 2 {
		t.Errorf("Expected blockchain length 2, got %d", len(chain.Blocks))
	}
}
