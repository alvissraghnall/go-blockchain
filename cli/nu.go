package cli

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"go-blockchain/wallet"
	"os"
	"strings"
)

func StartCLI() {
	// Greet the user and present options
	fmt.Println("Welcome to the Go Blockchain CLI!")
	fmt.Println("Choose an option:")
	fmt.Println("1. Enter existing mnemonic")
	fmt.Println("2. Generate new wallet")

	// Read the user's choice
	reader := bufio.NewReader(os.Stdin)
	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	var privateKey *ecdsa.PrivateKey
	var publicKey *ecdsa.PublicKey
	var mnemonic string
//	var err error

	switch option {
	case "1":
		// User wants to input an existing mnemonic
		fmt.Println("Enter your 12-word mnemonic phrase:")

		mnemonic, _ = reader.ReadString('\n')
		mnemonic = strings.TrimSpace(mnemonic)

		// Validate mnemonic and recover the wallet
/*		privateKey, publicKey, err = wallet.RecoverWalletFromMnemonic(mnemonic)
		if err != nil {
			fmt.Printf("Error recovering wallet: %v\n", err)
			return
		}*/
		fmt.Println("Wallet recovered successfully!")
	case "2":
		// User wants to generate a new wallet
		new_wallet, err := wallet.NewWalletWithMnemonic()
		if err != nil {
			fmt.Printf("Error generating wallet: %v\n", err)
			return
		}
		privateKey = new_wallet.PrivateKey
		publicKey = new_wallet.PublicKey
		fmt.Printf("New wallet generated successfully!\n")
		fmt.Printf("Your 12-word mnemonic: %s\n", new_wallet.Mnemonic)
	default:
		fmt.Println("Invalid option, please select 1 or 2.")
		return
	}

	// Display wallet information
	fmt.Printf("Private Key: %s\n", wallet.PrivateKeyToHex(privateKey))
	fmt.Printf("Public Key: %s\n", wallet.PublicKeyToHex(publicKey))

	
}
