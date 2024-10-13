package roaring

import (
	"bits"
	"errors"
	"fmt"
	"math/bits"
)

type BitsetContainer struct {
	data        [1024]uint64
	cardinality int
}

func (bc *BitsetContainer) InsertOne(n uint16) error {
	bucket := bc.data[n>>10]
	n %= 64

	mask := uint64(1 << n)
	found := bucket & uint64(mask)

	if found>>n == 1 {
		formatString := "Attempt to insert value [%d] to a BitsetContainer already containing it"
		errorMsg := fmt.Sprintf(formatString, n)

		return errors.New(errorMsg)
	}

	bc.cardinality++

	res := mask | bucket
	bc.data[n>>10] = res

	return nil
}

func (bc *BitsetContainer) Has(n uint16) bool {
	// index is floor(n / 1024)
	bucket := bc.data[n>>10]
	n %= 64

	// Get the nth bit
	mask := 1 << n
	res := bucket & uint64(mask)

	return res == 1
}

func (bc *BitsetContainer) Cardinality() int {
	return bc.cardinality
}

func (bc *BitsetContainer) AndArrayToArray() ArrayContainer {

}

func (bc *BitsetContainer) AndArrayToBitmap() BitsetContainer {

}

// NOTE: This is only used for the case where adding or removing
// an element from a bitset container would cause it to become
// an array container. For computing set operations, we can do 
// better by computing the result on the fly.
// This is Algorithm 2 in "Better bitmap performance with Roaring bitmaps"
// by Chambi et al.
func (bc *BitsetContainer) IntoArrayContainer() ArrayContainer {
	// Start out with higher capacity, since we know there
	// are likely a lot of elements
	arr := NewArrayContainerWithLength(bc.cardinality)

  offset := 0
  numAdded := 0
  for bucket := range bc.data {
    for bucket != 0 {
      temp := bucket & (^bucket + 1) 
      arr.data[numAdded] = (offset << 6) + bits.TrailingZeros16(temp - 1)
      numAdded++
      bucket &= bucket - 1
    }
    
    offset++
  }

  return arr
}

