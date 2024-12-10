package node

import (
	"fmt"
	"go-blockchain/transaction"
)

// BroadcastMessage sends a message to all peers.
func (n *Node) BroadcastMessage(message string) {
	for _, peer := range n.Peers {
		fmt.Printf("Sending message to peer %s: %s\n", peer, message)
	}
}

// ReceiveMessage simulates receiving a message from a peer.
func (n *Node) ReceiveMessage(peer string, message string) {
	fmt.Printf("Received message from peer %s: %s\n", peer, message)
}

// BroadcastTransaction broadcasts a transaction to all peers.
func (n *Node) BroadcastTransaction(tx transaction.Transaction) {
	message := fmt.Sprintf("Transaction: %s", tx.TransactionID())
	n.BroadcastMessage(message)
}
