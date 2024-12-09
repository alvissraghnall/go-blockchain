package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
//"	"crypto/rand"
	"math/big"
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

	// Derive private key from seed (simplified example)
	d := new(big.Int).SetBytes(seed[:32]) // Use the first 32 bytes of the seed
	privateKey := &ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     nil,
			Y:     nil,
		},
	}

	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())

	return privateKey, nil
}

