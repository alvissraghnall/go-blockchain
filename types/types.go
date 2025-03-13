package types

import (
  "encoding/json"
  "fmt"
  "reflect"
  "encoding/gob"
  "bytes"
  "crypto/sha256"
  "encoding/binary"
//  "encoding/hex"
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
	Sequence       uint32  		`json:"sequence"`
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
    PrevHash     []byte                  `json:"prev_hash"`
    Hash         []byte                  `json:"hash"`
    Nonce        int                     `json:"nonce"`
    Miner        string              `json:"miner"`
    BlockSize    uint64              `json:"blocksize"`
    Difficulty   uint32              `json:"difficulty"`
}

// Transaction represents a blockchain transaction.
type Transaction struct {
  ID       int64              `json:"id"`
  Version  int32              `json:"version"`
  Locktime uint32             `json:"locktime"`
  Inputs   []Input  		  `json:"inputs"`
  Outputs  []Output  		  `json:"outputs"`
  hash     []byte   		  // cached hash
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

// Hash returns the hash of the transaction
func (tx *Transaction) Hash() []byte {
    // Return cached hash if it exists
    if tx.hash != nil {
        return tx.hash
    }

    var buf bytes.Buffer
    encoder := gob.NewEncoder(&buf)

    // Encode transaction fields
    binary.Write(&buf, binary.LittleEndian, tx.ID)
    binary.Write(&buf, binary.LittleEndian, tx.Version)
    binary.Write(&buf, binary.LittleEndian, tx.Locktime)

    // Encode inputs
    for _, input := range tx.Inputs {
        encoder.Encode(input)
    }

    // Encode outputs
    for _, output := range tx.Outputs {
        encoder.Encode(output)
    }

    // Calculate hash
    hash := sha256.Sum256(buf.Bytes())
    // Double SHA256 
    hash = sha256.Sum256(hash[:])
    
    // Cache the hash
    tx.hash = hash[:]
    
    return tx.hash
}

// InvalidateHash clears the cached hash
// Call this when modifying the transaction
func (tx *Transaction) InvalidateHash() {
    tx.hash = nil
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

func (b *Block) GetPrevHash() []byte {
    return b.PrevHash
}

func (b *Block) GetHash() []byte {
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

func (b *Block) GetDifficulty() uint32 {
    return b.Difficulty
}

// CalculateHash generates a hash for the block.
func (b *Block) CalculateHash() []byte {
    data := fmt.Sprintf("%d%d%s%s%d", b.Index, b.Timestamp, b.PrevHash, b.serializeTransactions(), b.Nonce)
    hash := sha256.Sum256([]byte(data))
	return hash[:]
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
    return reflect.DeepEqual(b.Hash, b.CalculateHash())
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
    Left  *MerkleNode
    Right *MerkleNode
    Data  []byte
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
    node := MerkleNode{}

    if left == nil && right == nil {
        // Leaf node - hash the transaction data
        hash := sha256.Sum256(data)
        node.Data = hash[:]
    } else {
        // Internal node - hash the concatenation of left and right children
        prevHashes := append(left.Data, right.Data...)
        hash := sha256.Sum256(prevHashes)
        node.Data = hash[:]
    }

    node.Left = left
    node.Right = right

    return &node
}

// NewMerkleTree creates a new Merkle tree from a list of transactions
func NewMerkleTree(txHashes [][]byte) *MerkleNode {
    var nodes []MerkleNode

    // Create leaf nodes
    for _, hash := range txHashes {
        node := NewMerkleNode(nil, nil, hash)
        nodes = append(nodes, *node)
    }

    // If we have an odd number of transactions, duplicate the last one
    if len(nodes) % 2 != 0 {
        nodes = append(nodes, nodes[len(nodes)-1])
    }

    // Build tree by pairing nodes
    for len(nodes) > 1 {
        var level []MerkleNode

        // Process nodes two at a time to create parent nodes
        for i := 0; i < len(nodes); i += 2 {
            node := NewMerkleNode(&nodes[i], &nodes[i+1], nil)
            level = append(level, *node)
        }

        nodes = level
    }

    return &nodes[0]
}

func (b *Block) SerializeHeader() []byte {
    var buf bytes.Buffer
    encoder := gob.NewEncoder(&buf)

    encoder.Encode(b.PrevHash)
    encoder.Encode(b.Timestamp)
    encoder.Encode(b.Difficulty)

    // Calculate Merkle root from transactions
    var txHashes [][]byte
    for _, tx := range b.Transactions {
        txHashes = append(txHashes, tx.Hash())
    }
    
    merkleRoot := NewMerkleTree(txHashes).Data
    encoder.Encode(merkleRoot)

    return buf.Bytes()
}
