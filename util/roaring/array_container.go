package roaring

import (
	"errors"
	"fmt"
	"slices"
)

const MAX_KEYS int = 4096

type ArrayContainer struct {
	data []uint16
}

// CLASS METHODS
// ----------------

func NewArrayContainerWithCapacity(capacity int) ArrayContainer {
	return ArrayContainer{
		data: make([]uint16, 0, capacity),
	}
}

func NewArrayContainerWithLength(length int) ArrayContainer {
  return ArrayContainer {
    data: make([]uint16, length),
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

func (arr *ArrayContainer) Has(n uint16) bool {
  _, isPresent := slices.BinarySearch(arr.data, n)

  return isPresent
}

func (arr *ArrayContainer) InsertOne(n uint16) error {
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
		newData := make([]uint16, len(arr.data), neededCap)
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
