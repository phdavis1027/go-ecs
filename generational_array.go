package goecs

import (
	"errors"
	"fmt"
)

type GenIndex struct {
  Index        int
  Generation   int
}

type GenArrayEntry[T any] struct {
  Val          T
  Generation   int
}

type GenArray[T any] []GenArrayEntry[T] 

func (genArray GenArray[T]) Get(genIndex GenIndex) (*T, error) {
  n := len(genArray) 

  if (genIndex.Index >= n) {
    fmtString := "Attempt to retrieve item at out-of-bounds index [%d] from genArray of size [%d]"
    errorMsg  := fmt.Sprintf(fmtString, genIndex.Index, n)

    return nil, errors.New(errorMsg) 
  }

  indexedGen := genArray[genIndex.Index].Generation

  if (genIndex.Generation == indexedGen) {
    // Good case
    return &genArray[genIndex.Index].Val, nil
  } else {
    fmtString := "Attempt to access item from genArray at index [%d] with generation [%d], found generation [%d]"
    errorMsg := fmt.Sprintf(fmtString, genIndex.Index, genIndex.Generation, indexedGen)

    return nil, errors.New(errorMsg)
  }
}

func (genArray GenArray[T]) Set(genIndex GenIndex, newVal T) error {
  n := len(genArray)

  if (genIndex.Index >= n) {
    fmtString := "Attempt to set index [%d] when genArray has size [%d]"
    errorMsg := fmt.Sprintf(fmtString, genIndex.Index, n)

    return errors.New(errorMsg) 
  }

  indexedGen := genArray[genIndex.Index].Generation

  if (genIndex.Generation <= indexedGen) {
    fmtString := "Attempt to overwrite genArray entry with generation [%d], attempted generation was only [%d]"
    errorMsg := fmt.Sprintf(fmtString, indexedGen, genIndex.Generation)

    return errors.New(errorMsg)
  }

  genArray[genIndex.Index].Val = newVal
  genArray[genIndex.Index].Generation = genIndex.Generation

  return nil
}
