package roaring

import (
	"errors"
	"fmt"
)

type BitsetContainer struct {
  data          [1024]uint64
  cardinality   int
}

func (bc *BitsetContainer) InsertOne(n uint16) error {
  bucket := bc.data[n >> 10] 
  n %= 64

  mask := uint64(1 << n)
  found := bucket & uint64(mask)

  if found >> n == 1 {
    formatString := "Attempt to insert value [%d] to a BitsetContainer already containing it"
    errorMsg := fmt.Sprintf(formatString, n)

    return errors.New(errorMsg)
  }

  bc.cardinality++
  
  res := mask | bucket
  bc.data[n >> 10] = res
  
  return nil
}

func (bc *BitsetContainer) Has(n uint16) bool {
  // index is floor(n / 1024)
  bucket := bc.data[n >> 10]
  n %= 64

  // Get the nth bit
  mask := 1 << n
  res := bucket & uint64(mask)

  return res == 1
}

func (bc *BitsetContainer) Cardinality() int {
  return bc.cardinality
}

