package utils

import (
  "testing"
  "regexp"
  "fmt"
)


func TestTxn(t *testing.T) {
  in1 := Input{
    "00000000",
    6,
    "0xf177b728",
  }
  in2 := Input{
    "0x161f16ac2727d6b7",
    18,
    "0xac23aa56",
  }
  
  out1 := &Output{
    1000,
    "0xc25d2893e72cc3772ab2872d1616f7b8a",
  }
  
  inArr := [2]*Input{ *in1, *in2}
  
  txn1 := []*Transaction{
    inArr, []*Output { out1 },
  }
  
  fmt.Println(HashLeaves(txn1))
}

func TestHash(t *testing.T) {
  testData := "A"
  want := regexp.MustCompile(`\b`+testData+`\b`)
  hash, err := Hash("A")
  
  if !want.MatchString(hash) || err != nil {
    t.Fatalf(`Hash("A") = %q, %v, want match for %#q, nil`, hash, err, want)
  }
  
}

//fmt.Printf("%x\n", bs)