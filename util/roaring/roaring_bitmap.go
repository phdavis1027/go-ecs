package roaring

import (
	"fmt"
	"slices"
)

const ARRAY_MASK int = 0x0FFF0000
const BITSET_MASK int = 0x0EEE0000

type RoaringBitset struct {
	arrayKeys []uint16
	arrayVals []ArrayContainer

	bitsetKeys []uint16
	bitsetVals []BitsetContainer
}

// CLASS METHODS
func NewRoaringBitset() RoaringBitset {
  return RoaringBitset{
    arrayKeys: make([]uint16, 0),
    arrayVals: make([]ArrayContainer, 0),

    bitsetKeys: make([]uint16, 0),
    bitsetVals: make([]BitsetContainer, 0),
  }
}

func (r *RoaringBitset) InsertOne(value int32) error {
	least := leastSignificantBits16(uint32(value))
	most := mostSignificantBits16(uint32(value))

	// Check arrays first
  index, isPresent := slices.BinarySearch(r.arrayKeys, most)
  if isPresent {
    //NOTE: Apparently indexing without taking an address 
    // does an implicit copy
		arr    := &r.arrayVals[index]

    if arr.Cardinality() >= 4096 {
      bitset, error := arr.IntoBitset() 
      if error != nil {
        return error
      }

      // Remove the old array
      r.arrayKeys = slices.Delete(r.arrayKeys, index, index + 1)
      r.arrayVals = slices.Delete(r.arrayVals, index, index + 1)

      bitset.InsertOne(least)

      insertionPoint, _ := slices.BinarySearch(r.bitsetKeys, most)

      r.bitsetKeys = slices.Insert(r.bitsetKeys, insertionPoint, most)
      r.bitsetVals = slices.Insert(r.bitsetVals, insertionPoint, bitset)

      return nil
    } else {
      return arr.InsertOne(least)
    }
  }

	// Next, check bitset containers
  index, isPresent = slices.BinarySearch(r.bitsetKeys, most)
  if isPresent {
		bitset := &r.bitsetVals[index]
    return bitset.InsertOne(least)
  }

	// No appropriate chunk found, gotta make a new container
  // We always start out a container as an array container
  insertionPoint, _ := slices.BinarySearch(r.arrayKeys, most)
  r.arrayKeys = slices.Insert(r.arrayKeys, insertionPoint, most)

  arrayContainer := NewArrayContainerWithCapacity(8)
  error := arrayContainer.InsertOne(least)
  if error != nil {
    return error
  }

  r.arrayVals = slices.Insert(r.arrayVals, insertionPoint, arrayContainer)

  return nil
}

func (r *RoaringBitset) Has(n int32) bool {
  least := leastSignificantBits16(uint32(n))
  most  := mostSignificantBits16(uint32(n))

  i, isPresent := slices.BinarySearch(r.arrayKeys, most)
  if isPresent {
    return r.arrayVals[i].Has(least)
  } 

  i, isPresent = slices.BinarySearch(r.bitsetKeys, most)
  if !isPresent {
    return false
  }

  return r.bitsetVals[i].Has(least)
}

// UTILITY FUNCTIONS
// -----------------

func mostSignificantBits16(n uint32) uint16 {
	// SAFETY: we mask the 16 most significant bits
	return uint16((n & 0xFFFF0000) >> 16)
}

func leastSignificantBits16(n uint32) uint16 {
	// SAFETY: we mask the 16 least significant bits
	return uint16(n & 0x0000FFFF)
}

func (r *RoaringBitset) String() string {
  return fmt.Sprintf("RoaringBitset with %d array containers and %d bitset containers", len(r.arrayKeys), len(r.bitsetKeys))
}
