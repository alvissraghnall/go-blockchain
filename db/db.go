package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"blockchain/blockchain"
)

// BlockchainDB manages the SQLite database for blockchain data
type BlockchainDB struct {
	db *sql.DB
}

// InitDatabase creates a new database connection and sets up tables
func InitDatabase(dbPath string) (*BlockchainDB, error) {
	// Open the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create blocks table
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS blocks (
        index INTEGER PRIMARY KEY,
        timestamp INTEGER,
        transactions BLOB,
        prev_hash TEXT,
        hash TEXT,
        nonce INTEGER,
        miner TEXT,
        blocksize INTEGER,
        difficulty REAL
    );
	`)


	if err != nil {
		return nil, fmt.Errorf("failed to create blocks table: %v", err)
	}

	return &BlockchainDB{db: db}, nil
}
