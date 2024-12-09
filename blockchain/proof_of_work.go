package blockchain

import (
	"strings"
)

// MineBlock performs proof-of-work to mine a block.
func (b *Block) MineBlock(difficulty int) {
	prefix := strings.Repeat("0", difficulty)

	for {
		b.Hash = b.CalculateHash()
		if strings.HasPrefix(b.Hash, prefix) {
			break
		}
		b.Nonce++
	}
}
