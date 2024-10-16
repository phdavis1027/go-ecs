package roaring

import (
	"errors"
	"fmt"
	"math/bits"
)

type BitsetContainer struct {
	data        [2048]uint64
	cardinality int
}

func NewBitsetContainer() BitsetContainer {
  return BitsetContainer{}
}

func (bc *BitsetContainer) InsertOne(n uint32) error {
  // Neat! Switching to 64 bits doesn't affect this calculation
  // since it's entirely dependent on "row-length"
  index := n >> 6
	bucket := bc.data[index]
	n %= 64

  // I think this should also work unmodified?
	mask := uint64(1 << n)
	found := bucket & uint64(mask)

	if found>>n == 1 {
		formatString := "Attempt to insert value [%d] to a BitsetContainer already containing it"
		errorMsg := fmt.Sprintf(formatString, n)

		return errors.New(errorMsg)
	}

	bc.cardinality++

	res := mask | bucket
	bc.data[index] = res

	return nil
}

func (bc *BitsetContainer) Has(n uint32) bool {
	// index is floor(n / 64)
	bucket := bc.data[n>>6]
	n %= 64

	// Get the nth bit
	mask := 1 << n
	res := (bucket & uint64(mask)) >> n

	return res == 1
}

func (bc *BitsetContainer) Cardinality() int {
	return bc.cardinality
}

// TODO: implement
// func (bc *BitsetContainer) AndArrayToArray() ArrayContainer {
// 
// }
// TODO: implement
// func (bc *BitsetContainer) AndArrayToBitmap() BitsetContainer {
// 
// }

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
      arr.data[numAdded] = uint32(offset << 6) + uint32(bits.TrailingZeros32(uint32(temp - 1)))
      numAdded++
      bucket &= bucket - 1
    }
    
    offset++
  }

  return arr
}

