package cli

import (
    "fmt"
    "os"
    "encoding/json"
    "go-blockchain/wallet"
    "github.com/spf13/cobra"
)

type CliWallet struct {
  DefaultWallet bool
  *wallet.Wallet
}

var wallets []CliWallet
var defaultWalletID string

var createWalletCmd = &cobra.Command{
    Use:   "createwallet",
    Short: "Create a new wallet",
    Run: func(cmd *cobra.Command, args []string) {
        wallet := createWallet()
        loadWallets()
        var defaultWallet bool
        if len(wallets) == 0 {
          defaultWallet = true
        } else {
          defaultWallet = false
        }
        wallets = append(wallets, CliWallet { defaultWallet, wallet })
        saveWallets()
        fmt.Printf("Wallet created: %s\n", wallet.PublicKey)
    },
}

var listWalletsCmd = &cobra.Command{
    Use:   "listwallets",
    Short: "List all wallets",
    Run: func(cmd *cobra.Command, args []string) {
      loadWallets()
      for _, wallet := range wallets {
        fmt.Printf("Alias: %s, Address: %x\n", wallet.Alias, wallet.Address)
      }
      printTable(wallets, []string{"DefaultWallet", "Alias", "Address"})
    },
}

// idx := slices.IndexFunc(wallets, func(c CliWallet) bool { return fmt.Sprintf("%x", c.Address) == defaultWalletID || c.Alias == defaultWalletID })

var setDefaultWalletCmd = &cobra.Command{
  Use:   "setdefaultwallet <walletID>",
  Short: "Set the default wallet",
  Args:  cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    addressOrAlias := args[0]
    loadWallets()
    
    for i := range wallets {
      wallets[i].DefaultWallet = false
    }
    
    found := false
    for i, wallet := range wallets {
      if wallet.Alias == addressOrAlias || fmt.Sprintf("%x", wallet.Address) == addressOrAlias {
        wallets[i].DefaultWallet = true
        defaultWalletID = addressOrAlias
        found = true
        break
      }
    }
    
    if !found {
      fmt.Println("Wallet not found")
      os.Exit(1)
    }
    
    saveWallets()
    fmt.Printf("Default wallet set to: %s\n", addressOrAlias)
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
