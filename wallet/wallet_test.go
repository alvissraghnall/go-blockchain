package wallet

import (
//	"crypto/ecdsa"
	"testing"
	"fmt"
)

// TestGenerateKeyPair checks if the key pair generation works.
func TestGenerateKeyPair(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	if privateKey == nil || publicKey == nil {
		t.Fatalf("Generated keys should not be nil")
	}

	fmt.Printf("Private Key: %s\n", PrivateKeyToHex(privateKey))
	fmt.Printf("Public Key: %s\n", PublicKeyToHex(publicKey))

}

// TestPrivateKeyToHex ensures private key serialization is correct.
func TestPrivateKeyToHex(t *testing.T) {
	privateKey, _, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	hexKey := PrivateKeyToHex(privateKey)
	if len(hexKey) == 0 {
		t.Errorf("PrivateKeyToHex returned empty string")
	}
}

// TestPublicKeyToHex ensures public key serialization is correct.
func TestPublicKeyToHex(t *testing.T) {
	_, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	hexKey := PublicKeyToHex(publicKey)
	if len(hexKey) == 0 {
		t.Errorf("PublicKeyToHex returned empty string")
	}
}

// TestAddressFromPublicKey checks if the address generation logic is valid.
func TestAddressFromPublicKey(t *testing.T) {
	_, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	address := AddressFromPublicKey(publicKey)
	if len(address) == 0 {
		t.Errorf("AddressFromPublicKey returned empty address")
	}
}

// TestGenerateMnemonic verifies that mnemonics are generated correctly.
func TestGenerateMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("GenerateMnemonic failed: %v", err)
	}

	if len(mnemonic) == 0 {
		t.Errorf("Generated mnemonic is empty")
	}
}

// TestPrivateKeyFromMnemonic checks if private key derivation from mnemonic works.
func TestPrivateKeyFromMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("GenerateMnemonic failed: %v", err)
	}

	privateKey, err := PrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		t.Fatalf("PrivateKeyFromMnemonic failed: %v", err)
	}

	if privateKey == nil {
		t.Errorf("PrivateKeyFromMnemonic returned nil private key")
	}
}

// TestNewWalletWithMnemonic checks if wallet creation with mnemonic works.
func TestNewWalletWithMnemonic(t *testing.T) {
	wallet, err := NewWalletWithMnemonic()
	if err != nil {
		t.Fatalf("NewWalletWithMnemonic failed: %v", err)
	}

	if wallet.PrivateKey == nil || wallet.PublicKey == nil {
		t.Fatalf("Wallet keys should not be nil")
	}

	if len(wallet.Mnemonic) == 0 {
		t.Errorf("Wallet mnemonic is empty")
	}

	if len(wallet.Address) == 0 {
		t.Errorf("Wallet address is empty")
	}
}

// TestRecoverWalletFromMnemonic checks if wallet recovery works.
func TestRecoverWalletFromMnemonic(t *testing.T) {
	// Create a new wallet to get a valid mnemonic
	wallet, err := NewWalletWithMnemonic()
	if err != nil {
		t.Fatalf("NewWalletWithMnemonic failed: %v", err)
	}

	// Recover the wallet using the mnemonic
	recoveredWallet, err := RecoverWalletFromMnemonic(wallet.Mnemonic)
	if err != nil {
		t.Fatalf("RecoverWalletFromMnemonic failed: %v", err)
	}

	if wallet.Address != recoveredWallet.Address {
		t.Errorf("Recovered wallet address does not match original address")
	}
}
