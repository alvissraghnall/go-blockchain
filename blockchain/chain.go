package blockchain

import (
	"blockchain/types"
	"fmt"
 	"errors"
	"reflect"
)

// Blockchain represents the full blockchain.
type Blockchain struct {
	Blocks []*types.Block
}

type BlockWithHeight struct {
	types.Block
	height uint64
}

// NewBlockchain initializes a blockchain with a genesis block.
func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock(0, []types.Transaction{}, []byte{0x00})
	return &Blockchain{Blocks: []*types.Block{genesisBlock}}
}

// AddBlock adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(newBlock types.Block) {
	bc.Blocks = append(bc.Blocks, &newBlock)
}

func (bc *Blockchain) getHead () types.Block {
	return *bc.Blocks[len(bc.Blocks) - 1]
}

func (bc *Blockchain) GetLatestBlock () types.Block {
	return bc.getHead()
}

func (bc *Blockchain) GetHeight() uint64 {
	return uint64(len(bc.Blocks))
}

func (bc *Blockchain) GetBlock (hash []byte) (BlockWithHeight, error) {
	var i uint64
	for i = uint64(len(bc.Blocks) - 1); i >= 0; i-- {
		if reflect.DeepEqual(hash, bc.Blocks[i].Hash) {
			return BlockWithHeight { height: i, Block: *bc.Blocks[i] }, nil
		}
	}
	return BlockWithHeight {}, errors.New(fmt.Sprintf("No block with hash: %v found!", hash))
}

// IsValid checks the validity of the blockchain.
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// Validate current block hash
		if !reflect.DeepEqual(currentBlock.Hash, currentBlock.CalculateHash()) {
			return false
		}

		// Validate previous hash linkage
		if !reflect.DeepEqual(currentBlock.PrevHash, prevBlock.Hash) {
			return false
		}
	}

	return true
}
