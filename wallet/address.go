package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

// AddressFromPublicKey generates a Bitcoin-like wallet address from a public key.
// Uses SHA256 followed by RIPEMD-160 hashing (Bitcoin style).
func AddressFromPublicKey(publicKey *ecdsa.PublicKey) string {
	publicKeyBytes := crypto.FromECDSAPub(publicKey)[1:] // Remove 0x04 prefix for uncompressed key
	hash := crypto.Keccak256(publicKeyBytes)            // Keccak-256 hash of public key
	address := hash[len(hash)-20:]                      // Use last 20 bytes (Bitcoin uses RIPEMD-160 after SHA256)
	return hex.EncodeToString(address)
}
