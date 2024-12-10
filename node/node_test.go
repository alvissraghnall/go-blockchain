package node

import (
//	"go-blockchain/blockchain"
	"go-blockchain/transaction"
	"testing"
)

func TestNode(t *testing.T) {
	// Initialize a node
	node := NewNode("node-1")

	// Test adding a peer
	node.AddPeer("node-2")
	if len(node.Peers) != 1 {
		t.Errorf("expected 1 peer, got %d", len(node.Peers))
	}

	// Test adding a transaction
	tx := transaction.NewTransaction(
		[]transaction.TransactionInput{{PreviousTxID: "tx1", OutputIndex: 0}},
		[]transaction.TransactionOutput{{Address: "recipient", Amount: 10.0}},
	)
	node.AddTransaction(*tx)

	if len(node.TransactionPool) != 1 {
		t.Errorf("expected 1 transaction in pool, got %d", len(node.TransactionPool))
	}
}

func TestBroadcastAndReceiveMessage(t *testing.T) {
	node := NewNode("node-1")
	node.AddPeer("node-2")

	// Test broadcasting a message
	node.BroadcastMessage("Test Message")
}

func TestNodeMining(t *testing.T) {
	node := NewNode("node-1")

	// Create and add valid transactions
	tx1 := transaction.NewTransaction(
		[]transaction.TransactionInput{{PreviousTxID: "tx1", OutputIndex: 0}},
		[]transaction.TransactionOutput{{Address: "recipient1", Amount: 5.0}},
	)
	tx2 := transaction.NewTransaction(
		[]transaction.TransactionInput{{PreviousTxID: "tx2", OutputIndex: 1}},
		[]transaction.TransactionOutput{{Address: "recipient2", Amount: 2.0}},
	)

	_ = node.AddTransaction(*tx1)
	_ = node.AddTransaction(*tx2)

	// Mine a new block
	block, err := node.MineBlock()
	if err != nil {
		t.Fatalf("failed to mine block: %v", err)
	}

	if len(block.Transactions) != 2 {
		t.Errorf("expected 2 transactions in the mined block, got %d", len(block.Transactions))
	}

	// Verify the block was added to the chain
	if len(node.Blockchain.Blocks) != 2 { // Including genesis block
		t.Errorf("expected blockchain to have 2 blocks, got %d", len(node.Blockchain.Blocks))
	}
}
