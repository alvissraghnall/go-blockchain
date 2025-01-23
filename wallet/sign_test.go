package wallet

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "encoding/hex"
    "testing"
)

func TestSignMessage(t *testing.T) {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        t.Fatalf("Failed to generate private key: %v", err)
    }

    message := "Test message"
    rBytes, sBytes, err := SignMessage(privateKey, message)
    if err != nil {
        t.Fatalf("SignMessage returned an error: %v", err)
    }

    if len(rBytes) == 0 || len(sBytes) == 0 {
        t.Fatal("SignMessage returned empty signature components")
    }

    t.Logf("Signature generated successfully: r=%s, s=%s",
        hex.EncodeToString(rBytes), hex.EncodeToString(sBytes))
}

func TestVerifySignature(t *testing.T) {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        t.Fatalf("Failed to generate private key: %v", err)
    }

    publicKey := &privateKey.PublicKey
    message := "Test message"

    rBytes, sBytes, err := SignMessage(privateKey, message)
    if err != nil {
        t.Fatalf("SignMessage returned an error: %v", err)
    }

    isValid := VerifySignature(publicKey, message, rBytes, sBytes)
    if !isValid {
        t.Fatal("VerifySignature failed for a valid signature")
    }

    t.Log("VerifySignature successfully validated the signature")
}

func TestVerifySignatureWithInvalidSignature(t *testing.T) {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        t.Fatalf("Failed to generate private key: %v", err)
    }

    publicKey := &privateKey.PublicKey
    message := "Test message"

    // Sign the message
    _, _, err = SignMessage(privateKey, message)
    if err != nil {
        t.Fatalf("SignMessage returned an error: %v", err)
    }

    // Provide invalid signature values
    rBytes := []byte{0}
    sBytes := []byte{0}

    isValid := VerifySignature(publicKey, message, rBytes, sBytes)
    if isValid {
        t.Fatal("VerifySignature validated an invalid signature")
    }

    t.Log("VerifySignature correctly rejected an invalid signature")
}

func TestVerifySignatureWithTamperedMessage(t *testing.T) {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        t.Fatalf("Failed to generate private key: %v", err)
    }

    publicKey := &privateKey.PublicKey
    message := "Test message"
    tamperedMessage := "Tampered message"

    rBytes, sBytes, err := SignMessage(privateKey, message)
    if err != nil {
        t.Fatalf("SignMessage returned an error: %v", err)
    }

    isValid := VerifySignature(publicKey, tamperedMessage, rBytes, sBytes)
    if isValid {
        t.Fatal("VerifySignature validated a signature for a tampered message")
    }

    t.Log("VerifySignature correctly rejected a tampered message")
}

func TestSignAndVerifyMessage(t *testing.T) {
	// Generate a private key for all test cases
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Test cases
	tests := []struct {
		name       string
		message    string
		tampered   bool
		expectValid bool
	}{
		{
			name:       "Valid message, valid signature",
			message:    "Hello, world!",
			tampered:   false,
			expectValid: true,
		},
		{
			name:       "Empty message",
			message:    "",
			tampered:   false,
			expectValid: true,
		},
		{
			name:       "Long message",
			message:    string(make([]byte, 1000)), // 1KB of data
			tampered:   false,
			expectValid: true,
		},
		{
			name:       "Tampered message",
			message:    "Tamper test",
			tampered:   true,
			expectValid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Sign the message
			rBytes, sBytes, err := SignMessage(privateKey, tc.message)
			if err != nil {
				t.Fatalf("SignMessage failed: %v", err)
			}

			// Optionally tamper the message
			messageToVerify := tc.message
			if tc.tampered {
				messageToVerify += "tampered"
			}

			// Verify the signature
			isValid := VerifySignature(publicKey, messageToVerify, rBytes, sBytes)
			if isValid != tc.expectValid {
				t.Errorf("Expected validity: %v, got: %v", tc.expectValid, isValid)
			}
		})
	}
}

func TestInvalidSignatures(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Test cases with invalid signatures
	tests := []struct {
		name    string
		rBytes  []byte
		sBytes  []byte
		message string
	}{
		{
			name:    "Empty signature components",
			rBytes:  []byte{},
			sBytes:  []byte{},
			message: "Test message",
		},
		{
			name:    "Incorrect signature length",
			rBytes:  []byte{0x01, 0x02},
			sBytes:  []byte{0x03},
			message: "Test message",
		},
		{
			name:    "Random invalid signature",
			rBytes:  []byte{0x01, 0x02, 0x03, 0x04},
			sBytes:  []byte{0x05, 0x06, 0x07, 0x08},
			message: "Test message",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isValid := VerifySignature(publicKey, tc.message, tc.rBytes, tc.sBytes)
			if isValid {
				t.Error("VerifySignature should have returned false for invalid signature")
			}
		})
	}
}

// Benchmark for SignMessage
func BenchmarkSignMessage(b *testing.B) {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        b.Fatalf("Failed to generate private key: %v", err)
    }

    message := "Benchmark message"
    for i := 0; i < b.N; i++ {
        _, _, err := SignMessage(privateKey, message)
        if err != nil {
            b.Fatalf("SignMessage returned an error: %v", err)
        }
    }
}

// Benchmark for VerifySignature
func BenchmarkVerifySignature(b *testing.B) {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        b.Fatalf("Failed to generate private key: %v", err)
    }

    publicKey := &privateKey.PublicKey
    message := "Benchmark message"

    rBytes, sBytes, err := SignMessage(privateKey, message)
    if err != nil {
        b.Fatalf("SignMessage returned an error: %v", err)
    }

	b.ResetTimer()
    for i := 0; i < b.N; i++ {
        isValid := VerifySignature(publicKey, message, rBytes, sBytes)
        if !isValid {
            b.Fatal("VerifySignature failed during benchmark")
        }
    }
}
