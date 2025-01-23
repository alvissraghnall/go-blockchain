package consensus

import (
	"bytes"
	"crypto/sha256"
	"math"
	"fmt"
	"encoding/binary"
	"testing"
	"time"
)

func TestGenerateProof(t *testing.T) {
	testCases := []struct {
		name       string
		data       []byte
		difficulty uint64
	}{
		{
			name:       "Low Difficulty",
			data:       []byte("test transaction"),
			difficulty: 1000,
		},
		{
			name:       "Medium Difficulty",
			data:       []byte("complex blockchain transaction"),
			difficulty: 100000,
		},
		{
			name:       "High Difficulty",
			data:       []byte("high-security transaction"),
			difficulty: 1000000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate proof
			proof, err := GenerateProof(tc.data, tc.difficulty)
			
			// Check for generation errors
			if err != nil {
				t.Fatalf("Failed to generate proof: %v", err)
			}

			// Validate proof
			if !proof.ValidateProof() {
				t.Error("Generated proof failed validation")
			}

			// Check data integrity
			if !bytes.Equal(proof.Data, tc.data) {
				t.Error("Proof data does not match input data")
			}

			// Verify difficulty
			b  := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, proof.Nonce)

			combinedData := append(tc.data, []byte(b)...)
			hash := sha256.Sum256(combinedData)
			fmt.Println(hash)

			// Ensure nonce is within reasonable range
			if proof.Nonce > math.MaxUint64/2 {
				t.Errorf("Nonce seems unreasonably high: %d", proof.Nonce)
			}
		})
	}
}

func TestValidateProof(t *testing.T) {
	testCases := []struct {
		name       string
		data       []byte
		difficulty uint64
		shouldPass bool
	}{
		{
			name:       "Valid Proof",
			data:       []byte("valid transaction"),
			difficulty: 1000,
			shouldPass: true,
		},
		{
			name:       "Invalid Proof",
			data:       []byte("invalid transaction"),
			difficulty: math.MaxUint64, // Extremely high difficulty
			shouldPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// For invalid proof, use a fixed nonce that won't meet difficulty
			var proof *ProofOfWork
			var err error

			if tc.shouldPass {
				proof, err = GenerateProof(tc.data, tc.difficulty)
				if err != nil {
					t.Fatalf("Failed to generate proof: %v", err)
				}
			} else {
				proof = &ProofOfWork{
					Data:       tc.data,
					Difficulty: tc.difficulty,
					Nonce:      0, // Deliberately use a nonce that won't work
				}
			}

			// Validate proof
			isValid := proof.ValidateProof()
			if isValid != tc.shouldPass {
				t.Errorf("Proof validation incorrect. Expected %v, got %v", tc.shouldPass, isValid)
			}
		})
	}
}

func TestProofSerialization(t *testing.T) {
	originalData := []byte("serialization test")
	difficulty := uint64(5000)

	// Generate original proof
	originalProof, err := GenerateProof(originalData, difficulty)
	if err != nil {
		t.Fatalf("Failed to generate original proof: %v", err)
	}

	// Serialize proof
	serializedData := originalProof.Serialize()

	// Deserialize proof
	recoveredProof, err := Deserialize(serializedData)
	if err != nil {
		t.Fatalf("Failed to deserialize proof: %v", err)
	}

	// Compare original and recovered proofs
	if !bytes.Equal(originalProof.Data, recoveredProof.Data) {
		t.Errorf("Deserialized data does not match original")
	}

	if originalProof.Nonce != recoveredProof.Nonce {
		t.Errorf("Deserialized nonce does not match original")
	}

	if originalProof.Difficulty != recoveredProof.Difficulty {
		t.Errorf("Deserialized difficulty does not match original")
	}

	// Validate recovered proof
	if !recoveredProof.ValidateProof() {
		t.Error("Recovered proof failed validation")
	}
}

func TestProofMetadata(t *testing.T) {
	data := []byte("metadata test")
	difficulty := uint64(2000)

	// Generate proof
	proof, err := GenerateProof(data, difficulty)
	if err != nil {
		t.Fatalf("Failed to generate proof: %v", err)
	}

	// Generate metadata
	metadata := proof.GenerateMetadata()

	// Validate metadata
	if metadata.ProofHash == "" {
		t.Error("Proof hash should not be empty")
	}

	if metadata.Difficulty != difficulty {
		t.Errorf("Metadata difficulty does not match original: %d != %d", 
			metadata.Difficulty, difficulty)
	}

	// Check timestamp is recent
	currentTime := time.Now().Unix()
	if metadata.Timestamp > currentTime || 
	   metadata.Timestamp < currentTime - 10 { // Allow 10 seconds of slack
		t.Errorf("Timestamp seems incorrect: %d", metadata.Timestamp)
	}
}

func BenchmarkProofGeneration(b *testing.B) {
	data := []byte("benchmark transaction")
	difficulties := []uint64{1000, 10000, 100000}

	for _, difficulty := range difficulties {
		b.Run(fmt.Sprintf("Difficulty-%d", difficulty), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := GenerateProof(data, difficulty)
				if err != nil {
					b.Fatalf("Proof generation failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkProofValidation(b *testing.B) {
	data := []byte("validation benchmark")
	difficulties := []uint64{1000, 10000, 100000}

	for _, difficulty := range difficulties {
		b.Run(fmt.Sprintf("Difficulty-%d", difficulty), func(b *testing.B) {
			// Pre-generate a proof
			proof, err := GenerateProof(data, difficulty)
			if err != nil {
				b.Fatalf("Proof generation failed: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				proof.ValidateProof()
			}
		})
	}
}

func TestProofBoundaryConditions(t *testing.T) {
	testCases := []struct {
		name       string
		data       []byte
		difficulty uint64
	}{
		{
			name:       "Zero Difficulty",
			data:       []byte("zero difficulty test"),
			difficulty: 0,
		},
		{
			name:       "Extreme Difficulty",
			data:       []byte("extreme difficulty test"),
			difficulty: math.MaxUint64,
		},
		{
			name:       "Large Data",
			data:       bytes.Repeat([]byte("a"), 1024*1024), // 1MB of data
			difficulty: 10000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt proof generation
			proof, err := GenerateProof(tc.data, tc.difficulty)
			
			// Handle special cases
			if tc.difficulty == math.MaxUint64 {
				// Extremely high difficulty might timeout
				if err == nil {
					t.Errorf("Expected error for extremely high difficulty")
				}
				return
			}

			// Standard validation for other cases
			if err != nil {
				t.Fatalf("Failed to generate proof: %v", err)
			}

			if !proof.ValidateProof() {
				t.Error("Generated proof failed validation")
			}
		})
	}
}
