package cli
/**
import (
    "fmt"
    "go-blockchain/blockchain"
    "go-blockchain/transaction"
    "go-blockchain/wallet"
//    "os"
    "flag"
    "strconv"
)

func PrintBlockchain(blockchain *blockchain.Blockchain) {
    for _, block := range blockchain.Blocks {
        fmt.Printf("Block Index: %d\n", block.Index)
        fmt.Printf("Block Hash: %s\n", block.Hash)
        fmt.Printf("Previous Block Hash: %s\n", block.PrevHash)
        fmt.Printf("Transactions:\n")
        for _, tx := range block.Transactions {
            fmt.Printf("  TxID: %s\n", tx.TransactionID())
            for _, input := range tx.Inputs {
                fmt.Printf("    Input: %s:%d\n", input.PreviousTxID, input.OutputIndex)
            }
            for _, output := range tx.Outputs {
                fmt.Printf("    Output: %s -> %.2f\n", output.Address, output.Amount)
            }
        }
        fmt.Println()
    }
}

func GenerateWallet() {
    /**
privateKey, publicKey, err := wallet.GenerateKeyPair()
    if err != nil {
        fmt.Printf("Error generating key pair: %v\n", err)
        return
    } 

    new_wallet, err := wallet.NewWalletWithMnemonic()
    
    if err != nil {
	fmt.Printf("Error generating wallet mnemonic: %v\n", err)
	return
    }

    fmt.Println("Mnemonic:", new_wallet.Mnemonic)
    fmt.Println("Private Key:", wallet.PrivateKeyToHex(new_wallet.PrivateKey))
    fmt.Println("Public Key:", wallet.PublicKeyToHex(new_wallet.PublicKey))
}

func MineBlock(blockchain *blockchain.Blockchain, transactionPool *transaction.TransactionPool, difficulty int) {
    miner := blockchain.NewMiner(blockchain, transactionPool, difficulty)
    newBlock, err := miner.Mine()
    if err != nil {
        fmt.Printf("Error mining block: %v\n", err)
        return
    }
    blockchain.Blocks = append(blockchain.Blocks, newBlock)
    fmt.Printf("New block mined: %v\n", newBlock)
}

func SendTransaction(transactionPool *transaction.TransactionPool, fromPrivateKey *wallet.PrivateKey, toAddress string, amount float64) {
    // Create the transaction
    tx := transaction.NewTransaction([]transaction.TransactionInput{}, []transaction.TransactionOutput{
        {Address: toAddress, Amount: amount},
    })

    // Sign the transaction
    err := tx.SignTransaction(fromPrivateKey)
    if err != nil {
        fmt.Printf("Error signing transaction: %v\n", err)
        return
    }

    // Add the transaction to the pool
    transactionPool.AddTransaction(*tx)
    fmt.Printf("Transaction from %s to %s for %.2f added to the pool\n", wallet.PrivateKeyToHex(fromPrivateKey), toAddress, amount)
}

func StartCLI() {
    var viewChain bool
    var mine bool
    var difficulty int
    var generateWallet bool
    var sendTx bool
    var fromPrivateKey string
    var toAddress string
    var amount float64

    // Define the commands
    flag.BoolVar(&viewChain, "view-chain", false, "View the blockchain")
    flag.BoolVar(&mine, "mine", false, "Mine a block")
    flag.IntVar(&difficulty, "difficulty", 4, "Set mining difficulty (default: 4)")
    flag.BoolVar(&generateWallet, "generate-wallet", false, "Generate a new wallet")
    flag.BoolVar(&sendTx, "send-tx", false, "Send a transaction")
    flag.StringVar(&fromPrivateKey, "from", "", "Sender's private key for the transaction")
    flag.StringVar(&toAddress, "to", "", "Recipient address for the transaction")
    flag.Float64Var(&amount, "amount", 0, "Amount for the transaction")
    flag.Parse()

    // Create a simple blockchain and transaction pool
    blockchain := blockchain.NewBlockchain()
    transactionPool := transaction.NewTransactionPool()

    // Handle commands
    if viewChain {
        PrintBlockchain(blockchain)
    } else if generateWallet {
        GenerateWallet()
    } else if mine {
        MineBlock(blockchain, transactionPool, difficulty)
    } else if sendTx {
        if fromPrivateKey == "" || toAddress == "" || amount <= 0 {
            fmt.Println("Error: Invalid arguments for sending transaction")
            return
        }

        privateKey, err := wallet.HexToPrivateKey(fromPrivateKey)
        if err != nil {
            fmt.Printf("Error converting private key: %v\n", err)
            return
        }

        SendTransaction(transactionPool, privateKey, toAddress, amount)
    } else {
        fmt.Println("Error: No valid command provided. Use --help for options.")
    }
}

func main() {
    StartCLI()
}

*/
