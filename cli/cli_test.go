package cli

import (
    "testing"
    "go-blockchain/blockchain"
    "go-blockchain/transaction"
    "go-blockchain/wallet"
)

func TestGenerateWallet(t *testing.T) {
    // Generate a wallet and check if the mnemonic and keys are produced
    privateKey, publicKey, err := wallet.GenerateKeyPair()
    if err != nil {
        t.Fatalf("GenerateKeyPair failed: %v", err)
    }

    mnemonic := wallet.GenerateMnemonic()
    if len(mnemonic) == 0 {
        t.Fatalf("Mnemonic generation failed")
    }

    if privateKey == nil || publicKey == nil {
        t.Fatalf("Keys should not be nil")
    }

    t.Logf("Mnemonic: %s", mnemonic)
    t.Logf("Private Key: %s", wallet.PrivateKeyToHex(privateKey))
    t.Logf("Public Key: %s", wallet.PublicKeyToHex(publicKey))
}

func TestMineBlock(t *testing.T) {
    blockchain := blockchain.NewBlockchain()
    transactionPool := transaction.NewTransactionPool()

    miner := blockchain.NewMiner(blockchain, transactionPool, 4)

    block, err := miner.Mine()
    if err != nil {
        t.Fatalf("Mining failed: %v", err)
    }

    if block == nil {
        t.Fatal("Mined block should not be nil")
    }

    t.Logf("Mined Block: %+v", block)
}
