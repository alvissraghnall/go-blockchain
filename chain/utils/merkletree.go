package "utils"

import (
  "fmt",
  "crypto/sha256"
)

type MerkleTree interface {
  
}

func MerkleRoot () {
  if()
}

func hashLeaves(leaves []Transaction) []string {
  hashBook := [leaves.len]
  for _,v := range leaves {
    hash := sha256.New()
    hash.Write([]byte(fmt.Sprintf("%v", leaves)))
    
    return fmt.Sprintf("%x", h.Sum(nil))
    
  }
}