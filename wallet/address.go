package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strings"
	"github.com/ethereum/go-ethereum/crypto"
)

// AddressFromPublicKey generates a Bitcoin-like wallet address from a public key.
// Uses SHA256 followed by RIPEMD-160 hashing (Bitcoin style).
func AddressFromPublicKey(publicKey *ecdsa.PublicKey, useChecksum bool) string {
	publicKeyBytes := crypto.FromECDSAPub(publicKey)[1:] // Remove 0x04 prefix for uncompressed key
	hash := crypto.Keccak256(publicKeyBytes)            // Keccak-256 hash of public key
	address := hash[len(hash)-20:]                      // Use last 20 bytes (Bitcoin uses RIPEMD-160 after SHA256)
	
	if useChecksum {
		return calculateChecksumAddress(address)
	}

	return hex.EncodeToString(address)
}

// calculateChecksumAddress adds a simple checksum mechanism
func calculateChecksumAddress(address []byte) string {
	hexAddress := hex.EncodeToString(address)
	checksummed := ""
	
	for i, char := range hexAddress {
		if char >= '0' && char <= '9' {
			checksummed += string(char)
		} else {
			// Capitalize if corresponding bit in hash is 1
			if hashBit := (address[i/2] >> (4 * (1 - uint(i)%2))) & 1; hashBit == 1 {
				checksummed += strings.ToUpper(string(char))
			} else {
				checksummed += string(char)
			}
		}
	}
	
	return checksummed
}
