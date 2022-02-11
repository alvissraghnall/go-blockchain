package test

import (
  "../chain/utils/merkletree"
  "../chain/utils/transaction.go"
)

func main() {
  in1 := transaction.Input{
    "00000000",
    6,
    "0xf177b728"
  }
  out1 := transaction.Output{
    1000,
    "0xc25d2893e72cc3772ab2872d1616f7b8a"
  }
  txn1 := transaction.Transaction{
    in1, out1
  }
  
  fmt.Println(merkletree.HashLeaves(txn1))
}