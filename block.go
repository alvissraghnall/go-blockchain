package "block"

import "utils"

type Block struct {
  header *Header
  body []*txn.Transaction
}

type Header struct {
  prevHash string
  merkle_root string
  time Date
  difficulty float32
  nonce uint64
}