package db

import (
    "blockchain/types"
//    "database/sql"
    "os"
    "testing"
    "time"
)

// TestInitDatabase tests the InitDatabase function
func TestInitDatabase(t *testing.T) {
    dbPath := "test.db"
    defer os.Remove(dbPath) // Clean up after the test

    db, err := InitDatabase(dbPath)
    if err != nil {
        t.Fatalf("InitDatabase failed: %v", err)
    }
    defer db.db.Close()

    // Verify that tables were created
    tables := []string{"blocks", "transactions", "transaction_inputs", "transaction_outputs"}
    for _, table := range tables {
        var exists bool
        err = db.db.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ?)", table).Scan(&exists)
        if err != nil || !exists {
            t.Errorf("Table %s was not created", table)
        }
    }
}

// TestAddBlock tests the AddBlock function
func TestAddBlock(t *testing.T) {
    dbPath := "test.db"
    defer os.Remove(dbPath) // Clean up after the test

    db, err := InitDatabase(dbPath)
    if err != nil {
        t.Fatalf("InitDatabase failed: %v", err)
    }
    defer db.db.Close()

    // Create a sample block
    block := &types.Block{
        Index:        1,
        Timestamp:    time.Now().Unix(),
        Transactions: []types.Transaction{},
        PrevHash:     "prevhash",
        Hash:         "blockhash",
        Nonce:        123,
        Miner:        "miner",
        BlockSize:    100,
        Difficulty:   1.0,
    }

    // Add the block to the database
    err = db.AddBlock(block)
    if err != nil {
        t.Fatalf("AddBlock failed: %v", err)
    }

    // Verify that the block was added
    var count int
    err = db.db.QueryRow("SELECT COUNT(*) FROM blocks WHERE hash = ?", block.Hash).Scan(&count)
    if err != nil || count != 1 {
        t.Errorf("Block was not added to the database")
    }
}

// TestAddBlockWithTransactions tests adding a block with transactions
func TestAddBlockWithTransactions(t *testing.T) {
    dbPath := "test.db"
    defer os.Remove(dbPath) // Clean up after the test

    db, err := InitDatabase(dbPath)
    if err != nil {
        t.Fatalf("InitDatabase failed: %v", err)
    }
    defer db.db.Close()

    // Create a sample transaction
    tx := types.Transaction{
        ID:       1,
        Version:  1,
        Locktime: 0,
        Inputs: []types.Input{
            {
                ID:             1,
                PreviousTxHash: []byte("prevhash"),
                OutputIndex:    0,
                ScriptSig:      []byte("scriptsig"),
                Sequence:       0xFFFFFFFF,
            },
        },
        Outputs: []types.Output{
            {
                ID:           1,
                Amount:       10.0,
                ScriptPubKey: []byte("scriptpubkey"),
                ScriptType:   "P2PKH",
                Address:      []byte("address"),
            },
        },
    }

    // Create a sample block with the transaction
    block := &types.Block{
        Index:        1,
        Timestamp:    time.Now().Unix(),
        Transactions: []types.Transaction{tx},
        PrevHash:     "prevhash",
        Hash:         "blockhash",
        Nonce:        123,
        Miner:        "miner",
        BlockSize:    100,
        Difficulty:   1.0,
    }

    // Add the block to the database
    err = db.AddBlock(block)
    if err != nil {
        t.Fatalf("AddBlock failed: %v", err)
    }

    // Verify that the block and transaction were added
    var blockCount, txCount int
    err = db.db.QueryRow("SELECT COUNT(*) FROM blocks WHERE hash = ?", block.Hash).Scan(&blockCount)
    if err != nil || blockCount != 1 {
        t.Errorf("Block was not added to the database")
    }

    err = db.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE id = ?", tx.ID).Scan(&txCount)
    if err != nil || txCount != 1 {
        t.Errorf("Transaction was not added to the database")
    }
}

// TestGetTransaction tests the GetTransaction function
func TestGetTransaction(t *testing.T) {
    dbPath := "test.db"
    defer os.Remove(dbPath) // Clean up after the test

    db, err := InitDatabase(dbPath)
    if err != nil {
        t.Fatalf("InitDatabase failed: %v", err)
    }
    defer db.db.Close()

    // Create a sample transaction
    tx := types.Transaction{
        ID:       1,
        Version:  1,
        Locktime: 0,
        Inputs: []types.Input{
            {
                ID:             1,
                PreviousTxHash: []byte("prevhash"),
                OutputIndex:    0,
                ScriptSig:      []byte("scriptsig"),
                Sequence:       0xFFFFFFFF,
            },
        },
        Outputs: []types.Output{
            {
                ID:           1,
                Amount:       10.0,
                ScriptPubKey: []byte("scriptpubkey"),
                ScriptType:   "P2PKH",
                Address:      []byte("address"),
            },
        },
    }

    // Add a block with the transaction
    block := &types.Block{
        Index:        1,
        Timestamp:    time.Now().Unix(),
        Transactions: []types.Transaction{tx},
        PrevHash:     "prevhash",
        Hash:         "blockhash",
        Nonce:        123,
        Miner:        "miner",
        BlockSize:    100,
        Difficulty:   1.0,
    }
    err = db.AddBlock(block)
    if err != nil {
        t.Fatalf("AddBlock failed: %v", err)
    }

    // Retrieve the transaction
    retrievedTx, err := db.GetTransaction(tx.ID)
    if err != nil {
        t.Fatalf("GetTransaction failed: %v", err)
    }

    // Verify that the retrieved transaction matches the original
    if retrievedTx.ID != tx.ID || retrievedTx.Version != tx.Version || retrievedTx.Locktime != tx.Locktime {
        t.Errorf("Retrieved transaction does not match the original")
    }
}

// TestGetBlock tests the GetBlock function
func TestGetBlock(t *testing.T) {
    dbPath := "test.db"
    defer os.Remove(dbPath) // Clean up after the test

    db, err := InitDatabase(dbPath)
    if err != nil {
        t.Fatalf("InitDatabase failed: %v", err)
    }
    defer db.db.Close()

    // Create a sample block with a transaction
    tx := types.Transaction{
        ID:       1,
        Version:  1,
        Locktime: 0,
        Inputs: []types.Input{
            {
                ID:             1,
                PreviousTxHash: []byte("prevhash"),
                OutputIndex:    0,
                ScriptSig:      []byte("scriptsig"),
                Sequence:       0xFFFFFFFF,
            },
        },
        Outputs: []types.Output{
            {
                ID:           1,
                Amount:       10.0,
                ScriptPubKey: []byte("scriptpubkey"),
                ScriptType:   "P2PKH",
                Address:      []byte("address"),
            },
        },
    }
    block := &types.Block{
        Index:        1,
        Timestamp:    time.Now().Unix(),
        Transactions: []types.Transaction{tx},
        PrevHash:     "prevhash",
        Hash:         "blockhash",
        Nonce:        123,
        Miner:        "miner",
        BlockSize:    100,
        Difficulty:   1.0,
    }

    // Add the block to the database
    err = db.AddBlock(block)
    if err != nil {
        t.Fatalf("AddBlock failed: %v", err)
    }

    // Retrieve the block
    retrievedBlock, err := db.GetBlock(block.Hash)
    if err != nil {
        t.Fatalf("GetBlock failed: %v", err)
    }

    // Verify that the retrieved block matches the original
    if retrievedBlock.Index != block.Index || retrievedBlock.Hash != block.Hash {
        t.Errorf("Retrieved block does not match the original")
    }

    // Verify that the transaction was retrieved correctly
    if len(retrievedBlock.Transactions) != 1 || retrievedBlock.Transactions[0].ID != tx.ID {
        t.Errorf("Retrieved transaction does not match the original")
    }
}
