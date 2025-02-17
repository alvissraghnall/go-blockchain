package wallet

import (
	"crypto/ecdsa"
//	"crypto/elliptic"
//	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

// GenerateKeyPair generates a secp256k1 private/public key pair.
func GenerateKeyPair(mnemonic string, config *WalletConfig) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := PrivateKeyFromMnemonic(mnemonic, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

// PrivateKeyToHex converts a private key to a hexadecimal string.
func PrivateKeyToHex(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(crypto.FromECDSA(privateKey))
}

// PublicKeyToHex converts a public key to a hexadecimal string.
func PublicKeyToHex(publicKey *ecdsa.PublicKey) string {
	return hex.EncodeToString(crypto.FromECDSAPub(publicKey))
}
