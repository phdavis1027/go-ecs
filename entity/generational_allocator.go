package entity

import (
	"errors"
	"fmt"
)

type GenAllocatorEntry struct {
	// The generation of the entry.
	Generation int  `json:"generation"`
	IsLive     bool `json:"isLive"`
}

type GenAllocator struct {
	FreeList []int               `json:"freeList"`
	Entries  []GenAllocatorEntry `json:"entries"`
}

func CreateGenAllocatorOfSize(size int) GenAllocator {
	alloc := GenAllocator{
		FreeList: make([]int, size, size),
		Entries:  make([]GenAllocatorEntry, size, size),
	}

	for i := range size {
		alloc.FreeList[i] = i
	}

	return alloc
}

func (alloc *GenAllocator) Allocate() GenIndex {
	top := alloc.FreeList[len(alloc.FreeList)-1]

	alloc.FreeList = alloc.FreeList[:len(alloc.FreeList)-1]

	generation := alloc.Entries[top].Generation

	alloc.Entries[top].IsLive = true

	return GenIndex{
		Index:      top,
		Generation: generation,
	}
}

func (alloc *GenAllocator) Deallocate(genIndex GenIndex) error {
	n := len(alloc.FreeList)

	if genIndex.Index >= n {
		fmtString := "Attempt to deallocate out-of-bounds genIndex [%d] from array of length [%d]"
		errorMsg := fmt.Sprintf(fmtString, genIndex.Index, n)

		return errors.New(errorMsg)
	}

	if !alloc.Entries[genIndex.Index].IsLive {
		fmtString := "Attempt to deallocate freed genIndex [%d]"
		errorMsg := fmt.Sprint(fmtString, genIndex.Index)

		return errors.New(errorMsg)
	}

	alloc.Entries[genIndex.Index].IsLive = false
	alloc.Entries[genIndex.Index].Generation += 1

	alloc.FreeList = append(alloc.FreeList, genIndex.Index)

	return nil
}
