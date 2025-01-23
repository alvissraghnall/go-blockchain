package consensus

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"time"
	"math/big"
)

// ProofOfWork represents the non-interactive proof of work structure
type ProofOfWork struct {
	Data      []byte // Data to be proven
	Difficulty uint64 // Mining difficulty
	Nonce     uint64 // Nonce value used in proof generation
}

func GenerateProof(data []byte, difficulty uint64) (*ProofOfWork, error) {
	pow := &ProofOfWork{
		Data:       data,
		Difficulty: difficulty,
		Nonce:      0,
	}

	// Find a valid nonce that meets the difficulty requirement
	for pow.Nonce < math.MaxUint64 {
		if pow.ValidateProof() {
			return pow, nil
		}
		pow.Nonce++
	}

	return nil, fmt.Errorf("could not generate proof within allowed attempts")
}

// ValidateProof checks if the current proof meets the difficulty requirement
func (pow *ProofOfWork) ValidateProof() bool {
	// Combine data and nonce
	combinedData := append(pow.Data, []byte(fmt.Sprintf("%d", pow.Nonce))...)
	
	// Calculate hash
	hash := sha256.Sum256(combinedData)
	
	// Convert hash to big integer
	hashInt := new(big.Int).SetBytes(hash[:])
	
	// Calculate target (inverse of difficulty)
	var target *big.Int

	if pow.Difficulty == 0 {
	    target = big.NewInt(1)
	} else {
	    target = new(big.Int).Div(
	        new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
	        big.NewInt(int64(pow.Difficulty)),
	    )
	}

	// Check if hash meets difficulty requirement
	return hashInt.Cmp(target) < 0
}

// Serialize converts the proof to a byte slice for storage/transmission
func (pow *ProofOfWork) Serialize() []byte {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, pow.Nonce)
	
	serialized := append(pow.Data, nonceBytes...)
	serialized = append(serialized, byte(pow.Difficulty))
	
	return serialized
}

// Deserialize recreates a ProofOfWork from a byte slice
func Deserialize(serializedData []byte) (*ProofOfWork, error) {
	if len(serializedData) < 9 {
		return nil, fmt.Errorf("invalid serialized data")
	}

	// Extract difficulty (last byte)
	difficulty := uint64(serializedData[len(serializedData)-1])
	
	// Extract nonce (last 8 bytes before difficulty)
	nonce := binary.BigEndian.Uint64(serializedData[len(serializedData)-9 : len(serializedData)-1])
	
	// Extract original data
	data := serializedData[:len(serializedData)-9]

	pow := &ProofOfWork{
		Data:       data,
		Difficulty: difficulty,
		Nonce:      nonce,
	}

	return pow, nil
}

// ProofMetadata provides additional proof information
type ProofMetadata struct {
	ProofHash string
	Difficulty uint64
	Timestamp int64
	Nonce uint64
}

// GenerateProofMetadata creates metadata for the proof
func (pow *ProofOfWork) GenerateMetadata() ProofMetadata {
	combinedData := append(pow.Data, []byte(fmt.Sprintf("%d", pow.Nonce))...)
	hash := sha256.Sum256(combinedData)

	return ProofMetadata{
		ProofHash:  hex.EncodeToString(hash[:]),
		Difficulty: pow.Difficulty,
		Timestamp:  time.Now().Unix(),
		Nonce: 		pow.Nonce,
	}
}
