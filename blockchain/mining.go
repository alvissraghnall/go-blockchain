package blockchain

import (
	"fmt"
	"blockchain/transaction"
)

// Miner represents a miner who mines blocks.
type Miner struct {
	Blockchain      *Blockchain
	TransactionPool *transaction.TransactionPool
	Difficulty      int
}

// NewMiner initializes a new miner.
func NewMiner(blockchain *Blockchain, transactionPool *transaction.TransactionPool, difficulty int) *Miner {
	return &Miner{
		Blockchain:      blockchain,
		TransactionPool: transactionPool,
		Difficulty:      difficulty,
	}
}

// Mine retrieves transactions, mines a block, and appends it to the chain.
func (m *Miner) Mine() (*Block, error) {
	transactions := m.TransactionPool.GetTransactions()

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transactions to mine")
	}

	// Get the previous block
	prevBlock := m.Blockchain.Blocks[len(m.Blockchain.Blocks)-1]

	// Create a new block
	newBlock := NewBlock(len(m.Blockchain.Blocks), transactions, prevBlock.Hash)

	// Perform mining (Proof-of-Work)
	newBlock.MineBlock(m.Difficulty)

	// Add the mined block to the blockchain
	m.Blockchain.Blocks = append(m.Blockchain.Blocks, newBlock)

	return newBlock, nil
}
