package wallet

import (
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	message := "Hello, Blockchain!"
	r, s, err := SignMessage(privateKey, message)
	if err != nil {
		t.Fatalf("SignMessage failed: %v", err)
	}

	valid := VerifySignature(publicKey, message, r, s)
	if !valid {
		t.Errorf("Signature verification failed for valid signature")
	}
}
