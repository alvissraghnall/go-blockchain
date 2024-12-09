package blockchain

import "go-blockchain/transaction"

// Blockchain represents the full blockchain.
type Blockchain struct {
	Blocks []*Block
}

// NewBlockchain initializes a blockchain with a genesis block.
func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock(0, []transaction.Transaction{}, "0")
	return &Blockchain{Blocks: []*Block{genesisBlock}}
}

// AddBlock adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(transactions []transaction.Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(len(bc.Blocks), transactions, prevBlock.Hash)

	// Mine the block before adding
	newBlock.MineBlock(4)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// IsValid checks the validity of the blockchain.
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// Validate current block hash
		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		// Validate previous hash linkage
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}

	return true
}
