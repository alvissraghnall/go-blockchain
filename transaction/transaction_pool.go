package transaction

import "sync"

// TransactionPool stores unconfirmed transactions.
type TransactionPool struct {
	mu           sync.Mutex
	transactions []Transaction
}

// NewTransactionPool initializes a new transaction pool.
func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		transactions: []Transaction{},
	}
}

// AddTransaction adds a transaction to the pool.
func (tp *TransactionPool) AddTransaction(tx Transaction) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.transactions = append(tp.transactions, tx)
}

// GetTransactions retrieves all pending transactions and clears the pool.
func (tp *TransactionPool) GetTransactions() []Transaction {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	pending := tp.transactions
	tp.transactions = []Transaction{} // Clear the pool
	return pending
}

// Count returns the number of pending transactions in the pool.
func (tp *TransactionPool) Count() int {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	return len(tp.transactions)
}
