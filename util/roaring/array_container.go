package roaring

import (
	"errors"
	"fmt"
	"slices"
)

const MAX_KEYS int = 4096 * 2

type ArrayContainer struct {
	data []uint32
}

// CLASS METHODS
// ----------------

func NewArrayContainerWithCapacity(capacity int) ArrayContainer {
	return ArrayContainer{
		data: make([]uint32, 0, capacity),
	}
}

func NewArrayContainerWithLength(length int) ArrayContainer {
  return ArrayContainer {
    data: make([]uint32, length),
  }
}

// INSTANCE METHODS
// ----------------

func (arr *ArrayContainer) IntoBitset() (BitsetContainer, error) {
  bitset := NewBitsetContainer()

  for _, n := range arr.data {
    error := bitset.InsertOne(n)
    if error != nil {
      return bitset, error
    }
  }

  return bitset, nil
}

func (arr *ArrayContainer) Cardinality() int {
	return len(arr.data)
}

func (arr *ArrayContainer) Has(n uint32) bool {
  _, isPresent := slices.BinarySearch(arr.data, n)

  return isPresent
}

func (arr *ArrayContainer) InsertOne(n uint32) error {
	if len(arr.data) == MAX_KEYS {
		fmtString := "Attempt to insert [%d] into ArrayContainer which already has size 4096"
		errorMsg := fmt.Sprintf(fmtString, n)

		return errors.New(errorMsg)
	}

	neededCap := arr.expandHowMuch(1)

	insertionPoint, present := slices.BinarySearch(arr.data, n)
	if present {
		fmtString := "Attempt to insert [%d] into ArrayContainer which already contains it"
		errorMsg := fmt.Sprintf(fmtString, n)

		return errors.New(errorMsg)
	}

	if neededCap == 0 {
    arr.data = slices.Insert(arr.data, insertionPoint, n)

		return nil
	} else {
		newData := make([]uint32, len(arr.data), neededCap)
    copy(newData, arr.data)

    arr.data = slices.Insert(newData, insertionPoint, n)

		return nil
	}
}

func (arr *ArrayContainer) expandHowMuch(numNewElements int) int {
	newSize := len(arr.data) + numNewElements
	capacity := cap(arr.data)
	mustExpand := newSize > capacity

	if !mustExpand {
		return 0
	}

  // TODO: Compute what the right values for these cutoffs should be
	var newCapacity int
	if capacity < 64 {
		newCapacity = capacity << 1
	} else if capacity < 1067 {
		newCapacity = capacity + (capacity >> 1)
	} else if capacity <= 3840 {
		newCapacity = capacity + (capacity >> 2)
	} else {
		newCapacity = MAX_KEYS
	}

	return newCapacity
}
 
func (left *ArrayContainer) IntersectArray(right *ArrayContainer) ArrayContainer {
  var intersectMaxSize int 

  if left.Cardinality() < right.Cardinality() {
    intersectMaxSize = left.Cardinality()
  } else {
    intersectMaxSize = right.Cardinality()
  }

  intersect := NewArrayContainerWithLength(intersectMaxSize)
  diff := left.Cardinality() - right.Cardinality()

  if diff < 0 {
    diff = -diff
  }

  l := 0
  r := 0
  if diff > 8 {

    for l < left.Cardinality() && r < right.Cardinality() {
      if left.data[l] == right.data[r] {
        intersect.InsertOne(left.data[l])
      } else if left.data[l] < right.data[r] {
        l++
      } else {
        r++
      }
    }
  } else {
    // Apparently this is "galloping"
    var smaller, larger *ArrayContainer

    if left.Cardinality() < right.Cardinality() {
      smaller = left
      larger = right
    } else {
      smaller = right
      larger = left
    }

    si := smaller.Cardinality()
    li := larger.Cardinality()
    
    for si < smaller.Cardinality() && li < larger.Cardinality() {
      nextAvail := smaller.data[si]

      closest, present := slices.BinarySearch(larger.data, nextAvail)


    }
  }

  return intersect
}
