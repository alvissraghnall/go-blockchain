package wallet

import (
    "strings"
    "github.com/tyler-smith/go-bip39"
    "fmt"
)

func GenerateAlias() (string, error) {
  entropyBits := 256

  entropy, err := bip39.NewEntropy(entropyBits)
  if err != nil {
    return "", fmt.Errorf("failed to generate entropy: %w", err)
  }

  mnem, err := bip39.NewMnemonic(entropy)
  if err != nil {
    return "", fmt.Errorf("failed to generate alias: %w", err)
  }

  aliasArr := strings.Fields(mnem)[:2]

  alias := strings.Join(aliasArr, "-")

  return alias, nil
}
