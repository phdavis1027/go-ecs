package roaring

import (
	"errors"
	"fmt"
	"slices"
)

const ARRAY_MASK int = 0x0FFF0000
const BITSET_MASK int = 0x0EEE0000

type RoaringBitset struct {
	array_keys []uint16
	array_vals []ArrayContainer

	bitset_keys []uint16
	bitset_vals []BitsetContainer
}

func (r *RoaringBitset) InsertOne(value int32) error {
	least := leastSignificantBits16(uint32(value))
	most := mostSignificantBits16(uint32(value))

	// Find appropriate container, if exists
	container_index := -1

	// Check arrays first
	for index, key := range r.array_keys {
		if least == key {
			container_index = index | ARRAY_MASK
			break
		}
	}

	// Next, check bitset containers
	if container_index == -1 {
		for index, key := range r.bitset_keys {
			if least == key {
				container_index = index | BITSET_MASK
			}

			break
		}
	}

	// No appropriate chunk found, gotta make a new container
	if container_index == -1 {
		// We always start out a container as an array container

		insertion_point, _ := slices.BinarySearch(r.array_keys, least)
		r.array_keys = slices.Insert(r.array_keys, insertion_point, least)

		arrayContainer := NewArrayContainerWithCapacity(8)
		error := arrayContainer.InsertOne(most)
		if error != nil {
			return error
		}

		r.array_vals = slices.Insert(r.array_vals, insertion_point, arrayContainer)

		return nil
	}

	whichContainer := container_index & 0xFFFF0000
	realIndex := container_index & 0x0000FFFF

	if whichContainer == ARRAY_MASK {
		arr    := r.array_vals[realIndex]

	} else if whichContainer == BITSET_MASK {
		bitset := r.bitset_vals[realIndex]
	} else {
		formatString := "Found unexpected container index in RoaringBitset: [%d]"
		errorMsg := fmt.Sprint(formatString, container_index)

		return errors.New(errorMsg)
	}

	return nil
}

// UTILITY FUNCTIONS

func mostSignificantBits16(n uint32) uint16 {
	// SAFETY: we mask the 16 most significant bits
	return uint16((n & 0xFFFF0000) >> 16)
}

func leastSignificantBits16(n uint32) uint16 {
	// SAFETY: we mask the 16 least significant bits
	return uint16(n & 0x0000FFFF)
}
