package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"blockchain/types"
	"strconv"
)

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
		INSERT INTO blocks (id, timestamp, transactions, prev_hash, hash, nonce, miner, blocksize, difficulty)
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
/**
// GetBlock retrieves a block and all its transactions by hash
func (bdb *BlockchainDB) GetBlock(hash string) (*types.Block, error) {
	var block types.Block
	var txJSON []byte

	err := bdb.db.QueryRow(`
		SELECT id, timestamp, transactions, prev_hash, hash, 
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
*/

func (bdb *BlockchainDB) GetBlock(hash string) (*types.Block, error) {
    var block types.Block

    // Construct query that joins blocks, transactions, inputs and outputs
    query := `
        SELECT
            b.id, b.timestamp, b.prev_hash, b.hash, b.nonce,
            b.miner, b.blocksize, b.difficulty,
            t.id as tx_id, t.version, t.locktime,
            ti.id as input_id, ti.previous_tx_hash, ti.output_index,
            ti.script_sig, ti.sequence,
            tx.id as output_id, tx.value, tx.script_pubkey,
            tx.script_type, tx.address
        FROM blocks b
        LEFT JOIN transactions t ON t.block_hash = b.hash
        LEFT JOIN transaction_inputs ti ON ti.transaction_id = t.id
        LEFT JOIN transaction_outputs tx ON tx.transaction_id = t.id
        WHERE b.hash = ?
        ORDER BY t.id, ti.id, tx.id`

    rows, err := bdb.db.Query(query, hash)
    if err != nil {
        return nil, fmt.Errorf("failed to query block data: %v", err)
    }
    defer rows.Close()

    // Maps to keep track of transactions we've already seen
    txMap := make(map[int64]*types.Transaction)
    var currentTxID int64
    var tx *types.Transaction

    // Scan the first row to get block data
    first := true
    for rows.Next() {
        var (
            txID, inputID, outputID                               sql.NullInt64
            version, locktime                                     sql.NullInt64
            prevTxHash, scriptType, address                       sql.NullString
            outputIndex, sequence                                 sql.NullString
            value                                                sql.NullFloat64
            scriptSig, scriptPubKey                               []byte 
        )

        err := rows.Scan(
            &block.Index, &block.Timestamp, &block.PrevHash, &block.Hash,
            &block.Nonce, &block.Miner, &block.BlockSize, &block.Difficulty,
            &txID, &version, &locktime,
            &inputID, &prevTxHash, &outputIndex, &scriptSig, &sequence,
            &outputID, &value, &scriptPubKey, &scriptType, &address,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %v", err)
        }

        // Skip block initialization after first row
        if first {
            first = false
        }

        // If this row contains transaction data
        if txID.Valid {
            // If this is a new transaction
            if txID.Int64 != currentTxID {
                currentTxID = txID.Int64
                tx = &types.Transaction{
                    ID:       txID.Int64,
                    Version:  int32(version.Int64),
                    Locktime: uint32(locktime.Int64),
                }
                txMap[currentTxID] = tx
                block.Transactions = append(block.Transactions, *tx)
            }

            // Add input if present
            if inputID.Valid {
                // Convert outputIndex and sequence from string to uint64 and uint32
                outputIndexUint, err := strconv.ParseUint(outputIndex.String, 10, 64)
                if err != nil {
                    return nil, fmt.Errorf("failed to parse outputIndex: %v", err)
                }
                sequenceUint, err := strconv.ParseUint(sequence.String, 10, 32)
                if err != nil {
                    return nil, fmt.Errorf("failed to parse sequence: %v", err)
                }

                input := types.Input{
                    ID:             inputID.Int64,
                    PreviousTxHash: []byte(prevTxHash.String),
                    OutputIndex:    outputIndexUint,
                    ScriptSig:      scriptSig, // Use []byte directly
                    Sequence:       uint32(sequenceUint),
                }
                tx.Inputs = append(tx.Inputs, input)
            }

            // Add output if present
            if outputID.Valid {
                output := types.Output{
                    ID:           uint64(outputID.Int64), // Convert int64 to uint64
                    Amount:       value.Float64,
                    ScriptPubKey: scriptPubKey, // Use []byte directly
                    ScriptType:   scriptType.String,
                    Address:      []byte(address.String),
                }
                tx.Outputs = append(tx.Outputs, output)
            }
        }
    }

    if first {
        return nil, fmt.Errorf("block not found")
    }

    return &block, nil
}
