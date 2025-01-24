package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"bytes"
	"testing"
	"strings"
)

// Test default wallet configuration creation
func TestNewWalletWithMnemonic_DefaultConfig(t *testing.T) {
	config := DefaultConfig()
	wallet, err := NewWalletWithMnemonic(config)
	
	if err != nil {
		t.Fatalf("Failed to create wallet with default config: %v", err)
	}

	// Validate wallet components
	if wallet.Mnemonic == "" {
		t.Error("Mnemonic should not be empty")
	}
	
	if wallet.PrivateKey == nil {
		t.Error("Private key should not be nil")
	}
	
	if wallet.PublicKey == nil {
		t.Error("Public key should not be nil")
	}
	
	if bytes.Equal(wallet.Address, []byte{}) {
		t.Error("Address should not be empty")
	}
}

// Test wallet creation with custom configuration
func TestNewWalletWithMnemonic_CustomConfig(t *testing.T) {
	testCases := []struct {
		name       string
		wordCount  int
		curve      elliptic.Curve
		passphrase string
	}{
		{
			name:       "12-word mnemonic with P256",
			wordCount:  12,
			curve:      elliptic.P256(),
			passphrase: "",
		},
		{
			name:       "24-word mnemonic with P256",
			wordCount:  24,
			curve:      elliptic.P256(),
			passphrase: "test passphrase",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &WalletConfig{
				WordCount:   tc.wordCount,
				Curve:       tc.curve,
				Passphrase:  tc.passphrase,
				UseChecksum: true,
			}

			wallet, err := NewWalletWithMnemonic(config)
			if err != nil {
				t.Fatalf("Failed to create wallet: %v", err)
			}

			// Validate mnemonic word count
			mnemonicWords := len(strings.Split(wallet.Mnemonic, " "))
			if mnemonicWords != tc.wordCount {
				t.Errorf("Expected %d word mnemonic, got %d words", tc.wordCount, mnemonicWords)
			}
		})
	}
}

// Test wallet recovery from mnemonic
func TestRecoverWalletFromMnemonic(t *testing.T) {
	// Create an initial wallet
	config := DefaultConfig()
	originalWallet, err := NewWalletWithMnemonic(config)
	if err != nil {
		t.Fatalf("Failed to create original wallet: %v", err)
	}

	// Recover wallet using the same mnemonic
	recoveredWallet, err := RecoverWalletFromMnemonic(originalWallet.Mnemonic, config)
	if err != nil {
		t.Fatalf("Failed to recover wallet: %v", err)
	}

	// Compare critical wallet components
	compareWallets(t, originalWallet, recoveredWallet)
}

// Test wallet recovery with different configurations
func TestRecoverWalletFromMnemonic_DifferentConfigs(t *testing.T) {
	// Create an initial wallet
	originalConfig := DefaultConfig()
	originalWallet, err := NewWalletWithMnemonic(originalConfig)
	if err != nil {
		t.Fatalf("Failed to create original wallet: %v", err)
	}

	testCases := []struct {
		name           string
		modifyConfig   func(*WalletConfig)
		expectDifferent bool
	}{
		{
			name: "Different Passphrase",
			modifyConfig: func(cfg *WalletConfig) {
				cfg.Passphrase = "different passphrase"
			},
			expectDifferent: true,
		},
		{
			name: "Different Curve",
			modifyConfig: func(cfg *WalletConfig) {
				cfg.Curve = elliptic.P384()
			},
			expectDifferent: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clone the original config and modify
			config := &WalletConfig{
				WordCount:   originalConfig.WordCount,
				Curve:       originalConfig.Curve,
				Passphrase:  originalConfig.Passphrase,
				UseChecksum: originalConfig.UseChecksum,
			}
			tc.modifyConfig(config)

			// Attempt to recover wallet
			recoveredWallet, err := RecoverWalletFromMnemonic(originalWallet.Mnemonic, config)
			if err != nil {
				t.Fatalf("Failed to recover wallet: %v", err)
			}

			// Check if wallets are different
			if tc.expectDifferent {
				if bytes.Equal(originalWallet.Address, recoveredWallet.Address) {
					t.Errorf("Expected different wallet addresses")
				}
			}
		})
	}
}

// Test invalid mnemonic scenarios
func TestRecoverWalletFromMnemonic_InvalidMnemonic(t *testing.T) {
	config := DefaultConfig()
	
	testCases := []struct {
		name     string
		mnemonic string
	}{
		{
			name:     "Empty Mnemonic",
			mnemonic: "",
		},
		{
			name:     "Invalid Mnemonic Words",
			mnemonic: "invalid mnemonic phrase that does not make sense",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := RecoverWalletFromMnemonic(tc.mnemonic, config)
			if err == nil {
				t.Errorf("Expected error for invalid mnemonic, got nil")
			}
		})
	}
}

// Utility function to compare wallets
func compareWallets(t *testing.T, w1, w2 *Wallet) {
	// Compare mnemonics
	if w1.Mnemonic != w2.Mnemonic {
		t.Errorf("Mnemonic mismatch: %s != %s", w1.Mnemonic, w2.Mnemonic)
	}

	// Compare addresses
	if !bytes.Equal(w1.Address, w2.Address) {
		t.Errorf("Address mismatch: %s != %s", w1.Address, w2.Address)
	}

	// Compare public keys
	if !publicKeysEqual(w1.PublicKey, w2.PublicKey) {
		t.Error("Public keys do not match")
	}
}

// Utility function to compare public keys
func publicKeysEqual(pk1, pk2 *ecdsa.PublicKey) bool {
	if pk1 == nil || pk2 == nil {
		return pk1 == pk2
	}

	return pk1.Curve == pk2.Curve &&
		pk1.X.Cmp(pk2.X) == 0 &&
		pk1.Y.Cmp(pk2.Y) == 0
}

// Benchmark wallet creation
func BenchmarkNewWalletWithMnemonic(b *testing.B) {
	config := DefaultConfig()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewWalletWithMnemonic(config)
		if err != nil {
			b.Fatalf("Failed to create wallet: %v", err)
		}
	}
}

// Benchmark wallet recovery
func BenchmarkRecoverWalletFromMnemonic(b *testing.B) {
	config := DefaultConfig()
	wallet, err := NewWalletWithMnemonic(config)
	if err != nil {
		b.Fatalf("Failed to create initial wallet: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := RecoverWalletFromMnemonic(wallet.Mnemonic, config)
		if err != nil {
			b.Fatalf("Failed to recover wallet: %v", err)
		}
	}
}
