package wallet

import (
	"crypto/ecdsa"
	"fmt"
)

type Wallet struct {
	Mnemonic   string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    []byte
	Alias	   string
}

func NewWalletWithMnemonic(config *WalletConfig) (*Wallet, error) {
	mnemonic, err := GenerateMnemonic(config.WordCount)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	privateKey, publicKey, err := GenerateKeyPair(mnemonic, config)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	address := AddressFromPublicKey(publicKey, config.UseChecksum)

	return &Wallet{
		Mnemonic:   mnemonic,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		Alias:      config.Alias,
	}, nil
}

func RecoverWalletFromMnemonic(mnemonic string, config *WalletConfig) (*Wallet, error) {
	privateKey, err := PrivateKeyFromMnemonic(mnemonic, config)
	if err != nil {
		return nil, fmt.Errorf("failed to recover wallet: %w", err)
	}

	publicKey := &privateKey.PublicKey
	address := AddressFromPublicKey(publicKey, config.UseChecksum)

	return &Wallet{
		Mnemonic:   mnemonic,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		Alias: 	    config.Alias,
	}, nil
}
