package consensus

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
//	"fmt"
	"math/big"
	"sync"
	"time"
)

const (
	TargetBits         = 24    // Initial difficulty (24 = ~16.7 million hashes)
	MaxNonce           = ^uint64(0) // Maximum nonce value
	DifficultyInterval = 2016  // Blocks between difficulty adjustments
	TargetDuration     = 14 * time.Minute // Expected time per block
)

type ProofOfWork struct {
	Data       []byte
	Timestamp  int64
	Nonce      uint64
	Difficulty uint32 // Stored in compact "bits" format
}

// NewProof creates a new ProofOfWork instance
func NewProof(data []byte) *ProofOfWork {
	return &ProofOfWork{
		Data:       data,
		Timestamp:  time.Now().Unix(),
		Difficulty: CalculateDifficultyBits(TargetBits),
	}
}

// CalculateTarget converts compact difficulty bits to target
func CalculateTarget(bits uint32) *big.Int {
	exponent := bits >> 24
	mantissa := bits & 0x007fffff

	target := big.NewInt(int64(mantissa))
	target.Lsh(target, 8*(uint(exponent)-3))
	return target
}

// CalculateDifficultyBits converts difficulty to compact "bits" format
func CalculateDifficultyBits(difficulty uint32) uint32 {
	target := big.NewInt(1)
	target.Lsh(target, 256-uint(difficulty))
	
	// Convert to Bitcoin-like compact format
	bits := target.Bytes()
	exponent := uint32(len(bits))
	mantissa := uint32(0)
	if exponent > 3 {
		mantissa = uint32(binary.BigEndian.Uint32(bits[:4]))
	} else {
		mantissa = uint32(binary.BigEndian.Uint32(append(make([]byte, 4-exponent), bits...)))
	}
	
	return (exponent << 24) | mantissa
}

// Validate validates the proof against current difficulty
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	hash := pow.calculateHash()
	hashInt.SetBytes(hash[:])

	target := CalculateTarget(pow.Difficulty)
	return hashInt.Cmp(target) < 0
}

// calculateHash performs double SHA-256
func (pow *ProofOfWork) calculateHash() []byte {
	data := pow.prepareData()
	firstHash := sha256.Sum256(data)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:]
}

// prepareData serializes all fields for hashing
func (pow *ProofOfWork) prepareData() []byte {
	data := []byte{}
	
	// Convert timestamp to bytes
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts, uint64(pow.Timestamp))
	
	// Convert nonce to bytes
	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, pow.Nonce)
	
	// Combine all components
	data = append(data, pow.Data...)
	data = append(data, ts...)
	data = append(data, nonce...)
	data = append(data, byte(pow.Difficulty>>24), byte(pow.Difficulty>>16), 
		byte(pow.Difficulty>>8), byte(pow.Difficulty))
	
	return data
}

// Mine performs the proof-of-work computation
func (pow *ProofOfWork) Mine(timeout time.Duration) error {
	var wg sync.WaitGroup
	var found bool
	var mutex sync.Mutex
	
	target := CalculateTarget(pow.Difficulty)
	startTime := time.Now()
	
	workers := 4 // Number of parallel workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(offset uint64) {
			defer wg.Done()
			
			localNonce := offset
			for !found {
				// Check timeout
				if time.Since(startTime) > timeout {
					return
				}
				
				// Calculate hash
				hash := pow.calculateHashWithNonce(localNonce)
				
				// Check against target
				var hashInt big.Int
				hashInt.SetBytes(hash)
				if hashInt.Cmp(target) < 0 {
					mutex.Lock()
					if !found {
						found = true
						pow.Nonce = localNonce
					}
					mutex.Unlock()
					return
				}
				
				localNonce += uint64(workers)
			}
		}(uint64(i))
	}
	
	wg.Wait()
	if !found {
		return errors.New("proof-of-work not found within timeout")
	}
	return nil
}

func (pow *ProofOfWork) calculateHashWithNonce(nonce uint64) []byte {
	// Save current nonce and restore after calculation
	defer func(original uint64) { pow.Nonce = original }(pow.Nonce)
	pow.Nonce = nonce
	return pow.calculateHash()
}

// Serialize converts proof to byte format
func (pow *ProofOfWork) Serialize() []byte {
	data := []byte{}
	
	// Add timestamp
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts, uint64(pow.Timestamp))
	data = append(data, ts...)
	
	// Add nonce
	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, pow.Nonce)
	data = append(data, nonce...)
	
	// Add difficulty bits
	bits := make([]byte, 4)
	binary.BigEndian.PutUint32(bits, pow.Difficulty)
	data = append(data, bits...)
	
	// Add original data
	data = append(data, pow.Data...)
	
	return data
}

// Deserialize recreates ProofOfWork from byte data
func Deserialize(data []byte) (*ProofOfWork, error) {
	if len(data) < 20 {
		return nil, errors.New("invalid serialized data")
	}
	
	pow := &ProofOfWork{}
	
	// Read timestamp
	pow.Timestamp = int64(binary.BigEndian.Uint64(data[:8]))
	
	// Read nonce
	pow.Nonce = binary.BigEndian.Uint64(data[8:16])
	
	// Read difficulty bits
	pow.Difficulty = binary.BigEndian.Uint32(data[16:20])
	
	// Read original data
	pow.Data = data[20:]
	
	return pow, nil
}

// ProofMetadata contains mining information
type ProofMetadata struct {
	Hash       string
	Difficulty uint32
	Timestamp  int64
	Nonce      uint64
	Duration   time.Duration
}

// Metadata generates proof metadata
func (pow *ProofOfWork) Metadata(start time.Time) ProofMetadata {
	hash := pow.calculateHash()
	return ProofMetadata{
		Hash:       hex.EncodeToString(hash),
		Difficulty: pow.Difficulty,
		Timestamp:  pow.Timestamp,
		Nonce:      pow.Nonce,
		Duration:   time.Since(start),
	}
}

// AdjustDifficulty adjusts difficulty based on time taken
func (pow *ProofOfWork) AdjustDifficulty(actualDuration time.Duration) {
	ratio := float64(actualDuration) / float64(TargetDuration)
	if ratio < 0.75 {
		pow.Difficulty += 1
	} else if ratio > 1.25 {
		pow.Difficulty -= 1
	}
}
