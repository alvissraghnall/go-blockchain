package wallet

import (
	"crypto/ecdsa"
	"golang.org/x/crypto/ripemd160"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"github.com/ethereum/go-ethereum/crypto"
)

// AddressFromPublicKey generates a Bitcoin-like wallet address from a public key.
// Uses SHA256 followed by RIPEMD-160 hashing (Bitcoin style).
func AddressFromPublicKey(publicKey *ecdsa.PublicKey, useChecksum bool) []byte {
    publicKeyBytes := crypto.FromECDSAPub(publicKey)[1:] // Remove 0x04 prefix for uncompressed keys

    sha256Hash := sha256.Sum256(publicKeyBytes)

    secondHash := ripemd160.New()
    secondHash.Write(sha256Hash[:])

    ripemd160Hash := secondHash.Sum(nil)

    address := ripemd160Hash
    /*  if useChecksum {
        return calculateChecksumAddress(address)
    }*/
    //return base58.Encode(address)
    return address
}

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
