package utils

import (
  "fmt"
  "crypto/sha256"
)

type MerkleTree interface {
  
}

func MerkleRoot () {
  
}

func HashLeaves(leaves []*Transaction) []string {
  hashBook := make([]string, len(leaves))
  for i,v := range leaves {
    hash := sha256.New()
    hash.Write([]byte(fmt.Sprintf("%v", v)))
    
    //return fmt.Sprintf("%x", hash.Sum(nil))
    hashBook[i] = hash.Sum(nil)
  }
  return hashBook
}

func Hash(txt string) []byte {
  hash := sha256.New()
  hash.Write([]byte(txt))
  return hash.Sum(nil)
}

func main() {
  fmt.Println("x")
}