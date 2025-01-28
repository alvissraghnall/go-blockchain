package db

import (
	"database/sql"
	"fmt"
//	"log"
//	"time"

	_ "github.com/mattn/go-sqlite3"
)

// BlockchainDB manages the SQLite database for blockchain data
type BlockchainDB struct {
	db *sql.DB
}

// InitDatabase creates a new database connection and sets up tables
func InitDatabase(dbPath string) (*BlockchainDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Create blocks table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS blocks (
			index INTEGER PRIMARY KEY,
			timestamp INTEGER NOT NULL,
			transactions BLOB NOT NULL,
			prev_hash TEXT NOT NULL,
			hash TEXT NOT NULL UNIQUE,
			nonce INTEGER NOT NULL,
			miner TEXT NOT NULL,
			blocksize INTEGER NOT NULL,
			difficulty REAL NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create blocks table: %v", err)
	}

	// Create transactions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY,
			block_hash TEXT NOT NULL,
			block_index INTEGER NOT NULL,
			version INTEGER NOT NULL,
			locktime INTEGER NOT NULL,
			FOREIGN KEY (block_hash) REFERENCES blocks(hash),
			FOREIGN KEY (block_index) REFERENCES blocks(index)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions table: %v", err)
	}

	// Create transaction inputs table with script fields
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transaction_inputs (
			id INTEGER PRIMARY KEY,
			transaction_id INTEGER NOT NULL,
			previous_tx_hash TEXT NOT NULL,
			output_index INTEGER NOT NULL,
			script_sig BLOB NOT NULL,      -- The unlocking script
			sequence INTEGER NOT NULL,      -- Input sequence number
			FOREIGN KEY (transaction_id) REFERENCES transactions(id)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction_inputs table: %v", err)
	}

	// Create transaction outputs table with script fields
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transaction_outputs (
			id INTEGER PRIMARY KEY,
			transaction_id INTEGER NOT NULL,
			value REAL NOT NULL,            -- Amount in cryptocurrency
			script_pubkey BLOB NOT NULL,    -- The locking script
			script_type TEXT NOT NULL,      -- P2PKH, P2SH, etc.
			address TEXT,                   -- Optional derived address
			FOREIGN KEY (transaction_id) REFERENCES transactions(id)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction_outputs table: %v", err)
	}

	return &BlockchainDB{db: db}, nil
}
