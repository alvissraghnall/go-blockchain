package cli

import (
    "fmt"
    "os"
    "encoding/json"
    "go-blockchain/wallet"
    "github.com/spf13/cobra"
)

var wallets []*wallet.Wallet

var createWalletCmd = &cobra.Command{
    Use:   "createwallet",
    Short: "Create a new wallet",
    Run: func(cmd *cobra.Command, args []string) {
        wallet := createWallet()
	loadWallets()
        wallets = append(wallets, wallet)
        saveWallets()
        fmt.Printf("Wallet created: %s\n", wallet.PublicKey)
    },
}

var listWalletsCmd = &cobra.Command{
    Use:   "listwallets",
    Short: "List all wallets",
    Run: func(cmd *cobra.Command, args []string) {
        /**for _, wallet := range wallets {
            fmt.Printf("ID: %s, Address: %s\n", wallet.ID, wallet.Address)
        }*/
    },
}

var setDefaultWalletCmd = &cobra.Command{
    Use:   "setdefaultwallet <walletID>",
    Short: "Set the default wallet",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        // walletID := args[0]
        /** for _, wallet := range wallets {
            if wallet.alias == walletID {
                defaultWalletID = walletID
                fmt.Printf("Default wallet set to: %s\n", walletID)
                return
            }
        fmt.Println("Wallet not found")
	*/
    },
}

var getBalanceCmd = &cobra.Command{
    Use:   "getbalance <walletID>",
    Short: "Get the balance of a wallet",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        
    },
}

func init() {
    rootCmd.AddCommand(createWalletCmd)
    rootCmd.AddCommand(listWalletsCmd)
    rootCmd.AddCommand(setDefaultWalletCmd)
    rootCmd.AddCommand(getBalanceCmd)
}

func createWallet() *wallet.Wallet {
    config := wallet.DefaultConfig()

    wallet, err := wallet.NewWalletWithMnemonic(config)

    if err != nil {
        fmt.Println("Wallet generation failed!")
	os.Exit(1)
    }

    return wallet
}

func saveWallets() {
    file, _ := json.MarshalIndent(wallets, "", " ")
    _ = os.WriteFile("wallets.dat", file, 0644)
}

func loadWallets() {
    file, _ := os.ReadFile("wallets.dat")
    _ = json.Unmarshal(file, &wallets)
}
