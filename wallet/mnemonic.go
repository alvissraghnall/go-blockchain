package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
//"	"crypto/rand"
	"math/big"
	"fmt"
	"time"
	"github.com/tyler-smith/go-bip39"
)

type WalletConfig struct {
	Curve       elliptic.Curve
	WordCount   int  // 12 or 24 words
	UseChecksum bool // Add optional address checksum
	Passphrase  string
	Alias 	    string
}

func DefaultConfig() *WalletConfig {
	alias, err := GenerateAlias()
	if err != nil {
		alias += time.Now().Format("20060102150405")
	}
	return &WalletConfig{
		Curve:       elliptic.P256(),
		WordCount:   12,
		UseChecksum: true,
		Passphrase:  "",
		Alias: alias,
	}
}

func GenerateMnemonic(wordCount int) (string, error) {
	var entropyBits int
	switch wordCount {
	case 12:
		entropyBits = 128
	case 24:
		entropyBits = 256
	default:
		return "", fmt.Errorf("invalid word count. must be 12 or 24")
	}

	entropy, err := bip39.NewEntropy(entropyBits)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

func PrivateKeyFromMnemonic(mnemonic string, config *WalletConfig) (*ecdsa.PrivateKey, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, config.Passphrase)

	// Derive private key using curve from config
	d := new(big.Int).SetBytes(seed[:32])
	privateKey := &ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: config.Curve,
			X:     nil,
			Y:     nil,
		},
	}

	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())

	return privateKey, nil
}
