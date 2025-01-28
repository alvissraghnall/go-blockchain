package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"blockchain/types"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	ID       int64              `json:"id"`
	Version  int32              `json:"version"`
	Locktime uint32             `json:"locktime"`
	Inputs   []TransactionInput  `json:"inputs"`
	Outputs  []TransactionOutput `json:"outputs"`
}

// TransactionInput represents a transaction input with script
type TransactionInput struct {
	ID             int64   `json:"id"`
	PreviousTxHash string  `json:"previous_tx_hash"`
	OutputIndex    int64   `json:"output_index"`
	ScriptSig      []byte  `json:"script_sig"`
	Sequence       uint32  `json:"sequence"`
}

// TransactionOutput represents a transaction output with script
type TransactionOutput struct {
	ID           int64   `json:"id"`
	Value        float64 `json:"value"`
	ScriptPubKey []byte  `json:"script_pubkey"`
	ScriptType   string  `json:"script_type"`
	Address      string  `json:"address,omitempty"`
}

// AddBlock stores a new block in the database
func (bdb *BlockchainDB) AddBlock(block *types.Block) error {
	tx, err := bdb.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Convert transactions to JSON for storage
	txJSON, err := json.Marshal(block.Transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %v", err)
	}

	// Insert block
	_, err = tx.Exec(`
		INSERT INTO blocks (index, timestamp, transactions, prev_hash, hash, nonce, miner, blocksize, difficulty)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		block.Index, block.Timestamp, txJSON, block.PrevHash, block.Hash,
		block.Nonce, block.Miner, block.BlockSize, block.Difficulty,
	)
	if err != nil {
		return fmt.Errorf("failed to insert block: %v", err)
	}

	// Store each transaction
	for _, txn := range block.Transactions {
		// Insert transaction
		result, err := tx.Exec(`
			INSERT INTO transactions (block_hash, block_index, version, locktime)
			VALUES (?, ?, ?, ?)
		`, block.Hash, block.Index, txn.Version, txn.Locktime)
		if err != nil {
			return fmt.Errorf("failed to insert transaction: %v", err)
		}

		txID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get transaction ID: %v", err)
		}

		// Insert inputs
		for _, input := range txn.Inputs {
			_, err = tx.Exec(`
				INSERT INTO transaction_inputs (
					transaction_id, previous_tx_hash, output_index, 
					script_sig, sequence
				)
				VALUES (?, ?, ?, ?, ?)
			`, txID, input.PreviousTxHash, input.OutputIndex,
				input.ScriptSig, input.Sequence)
			if err != nil {
				return fmt.Errorf("failed to insert transaction input: %v", err)
			}
		}

		// Insert outputs
		for _, output := range txn.Outputs {
			_, err = tx.Exec(`
				INSERT INTO transaction_outputs (
					transaction_id, value, script_pubkey, 
					script_type, address
				)
				VALUES (?, ?, ?, ?, ?)
			`, txID, output.Amount, output.ScriptPubKey,
				output.ScriptType, output.Address)
			if err != nil {
				return fmt.Errorf("failed to insert transaction output: %v", err)
			}
		}
	}

	return tx.Commit()
}

// GetTransaction retrieves a complete transaction by its ID
func (bdb *BlockchainDB) GetTransaction(txID int64) (*types.Transaction, error) {
	tx := &types.Transaction{ID: txID}

	// Get transaction details
	err := bdb.db.QueryRow(`
		SELECT version, locktime FROM transactions WHERE id = ?
	`, txID).Scan(&tx.Version, &tx.Locktime)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	// Get inputs
	rows, err := bdb.db.Query(`
		SELECT id, previous_tx_hash, output_index, script_sig, sequence
		FROM transaction_inputs WHERE transaction_id = ?
	`, txID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inputs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var input types.Input
		err := rows.Scan(
			&input.ID, &input.PreviousTxHash, &input.OutputIndex,
			&input.ScriptSig, &input.Sequence,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan input: %v", err)
		}
		tx.Inputs = append(tx.Inputs, input)
	}

	// Get outputs
	rows, err = bdb.db.Query(`
		SELECT id, value, script_pubkey, script_type, address
		FROM transaction_outputs WHERE transaction_id = ?
	`, txID)
	if err != nil {
		return nil, fmt.Errorf("failed to query outputs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var output types.Output
		err := rows.Scan(
			&output.ID, &output.Amount, &output.ScriptPubKey,
			&output.ScriptType, &output.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan output: %v", err)
		}
		tx.Outputs = append(tx.Outputs, output)
	}

	return tx, nil
}

// GetBlock retrieves a block and all its transactions by hash
func (bdb *BlockchainDB) GetBlock(hash string) (*types.Block, error) {
	var block types.Block
	var txJSON []byte

	err := bdb.db.QueryRow(`
		SELECT index, timestamp, transactions, prev_hash, hash, 
		       nonce, miner, blocksize, difficulty
		FROM blocks WHERE hash = ?
	`, hash).Scan(
		&block.Index, &block.Timestamp, &txJSON, &block.PrevHash, &block.Hash,
		&block.Nonce, &block.Miner, &block.BlockSize, &block.Difficulty,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("block not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query block: %v", err)
	}

	// Get all transactions for this block
	rows, err := bdb.db.Query(`
		SELECT id FROM transactions WHERE block_hash = ?
	`, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var txID int64
		if err := rows.Scan(&txID); err != nil {
			return nil, fmt.Errorf("failed to scan transaction ID: %v", err)
		}
		
		tx, err := bdb.GetTransaction(txID)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction %d: %v", txID, err)
		}
		block.Transactions = append(block.Transactions, *tx)
	}

	return &block, nil
}
