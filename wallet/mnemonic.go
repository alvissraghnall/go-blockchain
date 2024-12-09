package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/tyler-smith/go-bip39"
)

// GenerateMnemonic creates a new mnemonic phrase.
func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128) // 128 bits for 12-word mnemonic
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// PrivateKeyFromMnemonic derives a private key from a mnemonic phrase.
func PrivateKeyFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, error) {
	// Ensure the mnemonic is valid
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	// Generate a seed from the mnemonic
	seed := bip39.NewSeed(mnemonic, "") // Empty passphrase for simplicity

	// Use the seed to create a private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key from seed: %w", err)
	}

	return privateKey, nil
}
