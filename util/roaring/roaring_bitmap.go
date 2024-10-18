package roaring

import (
	"fmt"
	"slices"
)

type RoaringBitset struct {
	arrayKeys []uint32
	arrayVals []ArrayContainer

	bitsetKeys []uint32
	bitsetVals []BitsetContainer
}

// CLASS METHODS
func NewRoaringBitset() RoaringBitset {
  return RoaringBitset{
    arrayKeys: make([]uint32, 0),
    arrayVals: make([]ArrayContainer, 0),

    bitsetKeys: make([]uint32, 0),
    bitsetVals: make([]BitsetContainer, 0),
  }
}

func (r *RoaringBitset) InsertOne(value uint64) error {
	least := leastSignificantBits32(value)
	most := mostSignificantBits32(value)

	// Check arrays first
  index, isPresent := slices.BinarySearch(r.arrayKeys, most)
  if isPresent {
    //NOTE: Apparently indexing without taking an address 
    // does an implicit copy
		arr    := &r.arrayVals[index]

    // MAX_KEYS is a constant defined in array_container.go
    if arr.Cardinality() >= MAX_KEYS {
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

func (r *RoaringBitset) Has(n uint64) bool {
  least := leastSignificantBits32(n)
  most  := mostSignificantBits32(n)

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

// In-place modify r to be the intersection into `left`
func (left *RoaringBitset) IntersectWith(right *RoaringBitset) RoaringBitset {
  i := 0
  j := 0

  intersect := NewRoaringBitset()

  for {
    if i >= len(left.arrayKeys) || j >= len(right.arrayKeys) {
      break
    }

    // Compare left arrays with right arrays
    if left.arrayKeys[i] == right.arrayKeys[j] {
      leftArray    := left.arrayVals[i]
      rightArray   := right.arrayVals[i]

      arrayIntserct:=  leftArray.IntersectArray(&rightArray)

      insertionPoint, _ := slices.BinarySearch(intersect.arrayKeys, left.arrayKeys[i])

      intersect.arrayKeys = slices.Insert(intersect.arrayKeys, insertionPoint, left.arrayKeys[i])
      intersect.arrayVals = slices.Insert(intersect.arrayVals, insertionPoint, arrayIntserct)

      i++
      j++
    } else if left.arrayKeys[i] < right.arrayKeys[j] {
      i++
    } else {
      j++
    }
  }

  i = 0
  j = 0
  // Compare left arrays with right bitsets
  for {
    if i >= len(left.arrayKeys) || j >= len(right.bitsetKeys) {
      break
    }

    if left.arrayKeys[i] == right.bitsetKeys[j] {
      arr := left.arrayVals[i]
      bitset := right.bitsetVals[j]


      arrayIntserct := arr.IntersectBitset(&bitset)

      insertionPoint, _ := slices.BinarySearch(intersect.arrayKeys, left.arrayKeys[i])

      intersect.arrayKeys = slices.Insert(intersect.arrayKeys, insertionPoint, left.arrayKeys[i])
      intersect.arrayVals = slices.Insert(intersect.arrayVals, insertionPoint, arrayIntserct)

      i++
      j++
    } else if left.arrayKeys[i] < right.bitsetKeys[j] {
      i++
    } else {
      j++
    }
  }

  i = 0
  j = 0

  return intersect
}

// UTILITY FUNCTIONS
// -----------------

func mostSignificantBits32(n uint64) uint32 {
	// SAFETY: we mask the 16 most significant bits
  return uint32(n >> 32)
}

func leastSignificantBits32(n uint64) uint32 {
	// SAFETY: we mask the 16 least significant bits
	return uint32(n & 0x00000000FFFFFFFF)
}

func (r *RoaringBitset) String() string {
  return fmt.Sprintf("RoaringBitset with %d array containers and %d bitset containers", len(r.arrayKeys), len(r.bitsetKeys))
}
