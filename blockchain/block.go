package blockchain

import (
	"time"
	"blockchain/types"
    "blockchain/consensus"
)

// NewBlock creates a new block
func NewBlock(index int, transactions []types.Transaction, prevHash []byte) *types.Block {
	block := &types.Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Nonce:        0, // Default nonce before mining
		Miner:        "",
	    Difficulty:   consensus.CalculateDifficultyBits(consensus.TargetBits),
	}

	block.Hash = block.CalculateHash()
	/* CALCULATE BLOCK SIZE && DIFFICULTY */

	return block
}
