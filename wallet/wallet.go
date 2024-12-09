package wallet

import (
	"crypto/ecdsa"
	"fmt"
)

// Wallet represents a user's wallet with mnemonic, public key, and address.
type Wallet struct {
	Mnemonic   string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
}

// NewWalletWithMnemonic creates a new wallet and generates a mnemonic.
func NewWalletWithMnemonic() (*Wallet, error) {
	// Generate mnemonic
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// Derive private key from mnemonic
	privateKey, err := PrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	publicKey := &privateKey.PublicKey
	address := AddressFromPublicKey(publicKey)

	return &Wallet{
		Mnemonic:   mnemonic,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// RecoverWalletFromMnemonic recreates a wallet using a mnemonic phrase.
func RecoverWalletFromMnemonic(mnemonic string) (*Wallet, error) {
	privateKey, err := PrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to recover wallet: %w", err)
	}

	publicKey := &privateKey.PublicKey
	address := AddressFromPublicKey(publicKey)

	return &Wallet{
		Mnemonic:   mnemonic,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}
