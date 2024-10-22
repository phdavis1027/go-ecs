package generational

import (
	"errors"
	"fmt"

	"github.com/phdavis1027/goecs/entity"
	"github.com/phdavis1027/goecs/util"
)

type GenIndex struct {
  Index        int `json:"index"`
  Generation   int `json:"generation"`
}

type GenArrayEntry[T any] struct {
  Val          T   `json:"val"`
  Generation   int `json:"generation"`
}

type GenArray[T any] []GenArrayEntry[util.Optional[T]] 

func CreateGenArrayOfSize[T any](n int) GenArray[T] {
  return make(GenArray[T], n)
}

func (genArray GenArray[T]) Get(genIndex GenIndex) (*T, error) {
  n := len(genArray) 

  if (genIndex.Index >= n) {
    fmtString := "Attempt to retrieve item at out-of-bounds index [%d] from genArray of size [%d]"
    errorMsg  := fmt.Sprintf(fmtString, genIndex.Index, n)

    return nil, errors.New(errorMsg) 
  }

  indexedGen := genArray[genIndex.Index].Generation

  if (genIndex.Generation != indexedGen) {
    fmtString := "Attempt to access item from genArray at index [%d] with generation [%d], found generation [%d]"
    errorMsg := fmt.Sprintf(fmtString, genIndex.Index, genIndex.Generation, indexedGen)

    return nil, errors.New(errorMsg)
    // Good case
  } 

  if (!genArray[genIndex.Index].Val.IsSome) {
    fmtString := "Attempt to access item from genArray at index [%d] with generation [%d], found generation [%d]"
    errorMsg := fmt.Sprintf(fmtString, genIndex.Index, genIndex.Generation, indexedGen)

    return nil, errors.New(errorMsg) 
  }

  return &genArray[genIndex.Index].Val.Inner, nil
}

func (genArray GenArray[T]) Set(genIndex entity.Entity, newVal T) error {
  n := len(genArray)

  if (genIndex.Index() >= n) {
    fmtString := "Attempt to set index [%d] when genArray has size [%d]"
    errorMsg := fmt.Sprintf(fmtString, genIndex.Index, n)

    return errors.New(errorMsg) 
  }

  indexedGen := genArray[genIndex.Index()].Generation

  if (genIndex.Generation() < indexedGen) {
    fmtString := "Attempt to overwrite genArray entry with generation [%d], attempted generation was only [%d]"
    errorMsg := fmt.Sprintf(fmtString, indexedGen, genIndex.Generation())

    return errors.New(errorMsg)
  }

  genArray[genIndex.Index()].Val.Inner = newVal
  genArray[genIndex.Index()].Val.IsSome = true

  genArray[genIndex.Index()].Generation = genIndex.Generation()

  return nil
}

// TODO: I think the delete function belongs on the actual ECS
func (genArray GenArray[T]) Delete(genIndex GenIndex) error {
  n := len(genArray)

  if (genIndex.Index >= n) {
    fmtString := "Attempt to delete out-of-range index [%d] from genArray of size [%d]" 
    errorMsg := fmt.Sprintf(fmtString, genIndex.Index, n)

    return errors.New(errorMsg)
  }

  return nil
}

/* UTILITY FUNCTIONS */
/*********************/

func (genArray GenArray[T]) String() string {
  n := len(genArray)

  if (n == 0) {
    return "[]"
  }

  var str string = "["

  for i := 0; i < n; i++ {
    if (genArray[i].Val.IsSome) {
      str += fmt.Sprintf("%v", genArray[i].Val.Inner)
    } else {
      str += "nil"
    }

    if (i != n - 1) {
      str += ", "
    }
  }

  str += "]"

  return str
}
