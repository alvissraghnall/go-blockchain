package node

import (
	"go-blockchain/blockchain"
	"go-blockchain/transaction"
	"fmt"
)

// Node represents a blockchain node.
type Node struct {
	ID          string
	Blockchain  *blockchain.Blockchain
	TransactionPool []transaction.Transaction
	Peers       []string
}

// NewNode initializes a new blockchain node.
func NewNode(id string) *Node {
	return &Node{
		ID:          id,
		Blockchain:  &blockchain.Blockchain{},
		TransactionPool: make([]transaction.Transaction, 0),
		Peers:       make([]string, 0),
	}
}

// AddTransaction adds a transaction to the node's transaction pool if valid.
func (n *Node) AddTransaction(tx transaction.Transaction) error {
	// Validate the transaction
	if !tx.IsValid() {
		return fmt.Errorf("invalid transaction")
	}

	// Check for duplicate transactions
	for _, poolTx := range n.TransactionPool {
		if poolTx.TransactionID() == tx.TransactionID() {
			return fmt.Errorf("transaction already exists in pool")
		}
	}

	// Add to the pool
	n.TransactionPool = append(n.TransactionPool, tx)
	return nil
}

// AddPeer connects the node to a new peer.
func (n *Node) AddPeer(peer string) {
	n.Peers = append(n.Peers, peer)
}


// AddBlock adds a block to the blockchain if valid.
func (n *Node) AddBlock(block blockchain.Block) error {
	// Validate the block
	if !block.IsValid() {
		return fmt.Errorf("invalid block")
	}

	// Ensure the block links correctly
	lastBlock := n.Blockchain.Blocks[len(n.Blockchain.Blocks)-1]
	if block.PrevHash != lastBlock.Hash {
		return fmt.Errorf("block does not link to the chain")
	}

	// Add to the blockchain
	n.Blockchain.Blocks = append(n.Blockchain.Blocks, block)
	return nil
}

// MineBlock mines a new block using transactions from the pool.
func (n *Node) MineBlock() (*blockchain.Block, error) {
	if len(n.TransactionPool) == 0 {
		return nil, fmt.Errorf("no transactions to mine")
	}

	// Use transactions in the pool to create a new block
	lastBlock := n.Blockchain.Blocks[len(n.Blockchain.Blocks)-1]
	newBlock := blockchain.NewBlock(lastBlock.Index + 1, n.TransactionPool, lastBlock.Hash)

	// Mine the block
	newBlock.MineBlock(75)

	// Add the block to the blockchain
	err := n.AddBlock(*newBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to add mined block: %w", err)
	}

	// Clear the transaction pool
	n.TransactionPool = []transaction.Transaction{}

	// Broadcast the block to peers
	n.BroadcastMessage(fmt.Sprintf("New block mined: %s", newBlock.Hash))

	return newBlock, nil
}
