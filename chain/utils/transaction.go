package utils


type Transaction struct {
  inputs []*Input
  outputs []*Output
}

type Input struct {
  prevout_hash string
  prevout_index uint32
  scriptSig string
}

type Output struct {
  value uint64
  scriptPubKey string
}