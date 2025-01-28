package types

import (
  "crypto/sha256"
  "encoding/json"
  "fmt"
  "encoding/hex"
)

// Input references an output from a previous transaction.
type Input struct {
    ID             int64   `json:"id"`
    // This field represents the hash of the previous transaction that this inp>
    PreviousTxHash []byte       `json:"previousTxHash"`
    // This field represents the index of the output in the previous transactio>
    OutputIndex    uint64       `json:"outputIndex"`
    // This field represents the script signature that unlocks the previous tra>
    ScriptSig      []byte       `json:"scriptSig"`
	Sequence       uint32  `json:"sequence"`
}

// Output specifies a recipient and amount.
type Output struct {
    ID           uint64   `json:"id"`
    Amount       float64     `json:"amount"`
    ScriptPubKey []byte      `json:"scriptPubKey"`
    ScriptType   string  `json:"script_type"`
    Address      []byte  `json:"address,omitempty"`
}

// UTXO represents an unspent transaction output.
type UTXO struct {
    // This field represents the block in which the transaction containing this>
    Block        Block           `json:"block"`
    // This field represents the transaction that created this UTXO.
    Transaction  Transaction    `json:"transaction"`
    // This field represents the index of this UTXO within the transaction that>
    OutputIndex  uint64         `json:"outputIndex"`
    // This field represents the amount of cryptocurrency associated with this >
    Amount       float64            `json:"amount"`
}

// Block represents a single block in the blockchain.
type Block struct {
    Index        int                 `json:"index"`
    Timestamp    int64               `json:"timestamp"`
    Transactions []Transaction   `json:"transactions"`
    PrevHash     string                  `json:"prev_hash"`
    Hash         string                  `json:"hash"`
    Nonce        int                     `json:"nonce"`
    Miner        string              `json:"miner"`
    BlockSize    uint64              `json:"blocksize"`
    Difficulty   float64             `json:"difficulty"`
}

// Transaction represents a blockchain transaction.
type Transaction struct {
  ID       int64              `json:"id"`
  Version  int32              `json:"version"`
  Locktime uint32             `json:"locktime"`
  Inputs     []Input
  Outputs    []Output
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
        OutputIndex:     outputIndex,
        ScriptSig:       scriptSig,
    }
}

func NewUTXO(block Block, transaction Transaction, outputIndex uint64, amount float64) *UTXO {
    return &UTXO{
        Block:        block,
        Transaction:  transaction,
        OutputIndex:  outputIndex,
        Amount:       amount,
    }
}

type BlockInterface interface {
    GetIndex() int
    GetTimestamp() int64
    GetTransactions() []Transaction
    GetPrevHash() string
    GetHash() string
    GetNonce() int
    GetMiner() string
    GetBlockSize() uint64
    GetDifficulty() float64
}

func (b *Block) GetIndex() int {
    return b.Index
}

func (b *Block) GetTimestamp() int64 {
    return b.Timestamp
}

func (b *Block) GetTransactions() []Transaction {
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
