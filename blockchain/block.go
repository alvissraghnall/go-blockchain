package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	"blockchain/transaction"
)

type BlockInterface interface {
    GetIndex() int
    GetTimestamp() int64
    GetTransactions() []transaction.Transaction
    GetPrevHash() string
    GetHash() string
    GetNonce() int
    GetMiner() string
    GetBlockSize() uint64
    GetDifficulty() float64
}

// Block represents a single block in the blockchain.
type Block struct {
	Index        int           		 `json:"index"`
	Timestamp    int64         		 `json:"timestamp"`
	Transactions []transaction.Transaction   `json:"transactions"`
	PrevHash     string         		 `json:"prev_hash"`
	Hash         string         		 `json:"hash"`
	Nonce        int            	 	 `json:"nonce"`
	Miner		 string				 `json:"miner"`
	BlockSize	 uint64				 `json:"blocksize"`
	Difficulty	 float64			 `json:"difficulty"`
}

func (b *Block) GetIndex() int {
    return b.Index
}

func (b *Block) GetTimestamp() int64 {
    return b.Timestamp
}

func (b *Block) GetTransactions() []transaction.TransactionInterface {
    return b.Transactions
}

func (b *Block) GetPrevHash() string {
    return b.PrevHash
}

func (b *Block) GetHash() string {
    return b.Hash
}

func (b *Block) GetNonce() int {
    return b.Nonce
}

func (b *Block) GetMiner() string {
    return b.Miner
}

func (b *Block) GetBlockSize() uint64 {
    return b.BlockSize
}

func (b *Block) GetDifficulty() float64 {
    return b.Difficulty
}

// NewBlock creates a new block.
func NewBlock(index int, transactions []transaction.Transaction, prevHash string) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Nonce:        0, // Default nonce before mining
		Miner:        "",
	}

	block.Hash = block.CalculateHash()
	/* CALCULATE BLOCK SIZE && DIFFICULTY */

	return block
}

// CalculateHash generates a hash for the block.
func (b *Block) CalculateHash() string {
	data := fmt.Sprintf("%d%d%s%s%d", b.Index, b.Timestamp, b.PrevHash, b.serializeTransactions(), b.Nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// serializeTransactions converts transactions to a string for hashing.
func (b *Block) serializeTransactions() string {
	bytes, err := json.Marshal(b.Transactions)
	if err != nil {
		return ""
	}
	return string(bytes)
}


// IsValid checks if the block is valid
func (b *Block) IsValid() bool {
	// Simple validation: check if the block's hash matches the calculated hash.
	return b.Hash == b.CalculateHash()
}
