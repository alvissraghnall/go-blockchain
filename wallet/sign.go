package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

//"	"github.com/ethereum/go-ethereum/crypto"
)

// SignMessage signs a given message with the private key.
func SignMessage(privateKey *ecdsa.PrivateKey, message string) ([]byte, []byte, error) {
	hash := sha256.Sum256([]byte(message))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign message: %w", err)
	}

	return r.Bytes(), s.Bytes(), nil
}

// VerifySignature verifies a message signature using the public key.
func VerifySignature(publicKey *ecdsa.PublicKey, message string, rBytes, sBytes []byte) bool {
	hash := sha256.Sum256([]byte(message))

	var r, s big.Int
	r.SetBytes(rBytes)
	s.SetBytes(sBytes)

	return ecdsa.Verify(publicKey, hash[:], &r, &s)
}
