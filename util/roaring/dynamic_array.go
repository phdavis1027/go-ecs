package roaring
import (
	"errors"
	"fmt"
)

const MAX_KEYS int = 4096

type DynamicArray[T any] struct {
  data      []T
}

func NewDynamicArrayWithCapacity[T any] (capacity int) DynamicArray[T] {
  return DynamicArray[T] {
    data:  make([]T, 0, capacity),
  };
}

func (dua *DynamicArray[T]) Get(i int) (T, error) {
  if (i < 0 || i >= len(dua.data)) {
    fmtString := "Attempt to get out-of-range index [%d] for DynamicArray with length %d"
    errorMsg  := fmt.Sprintf(fmtString, i, len(dua.data))
    
    return *new(T), errors.New(errorMsg)
  }

  return dua.data[i], nil
}


func (dua *DynamicArray[T]) Insert(n T) error {
  if (len(dua.data) >= MAX_KEYS) {
    fmtString := "Attempt to insert key [%d] when array already has max size (4096)"
    errorMSg  := fmt.Sprintf(fmtString, n)

    return errors.New(errorMSg)
  } 

  if (len(dua.data) == cap(dua.data)) {
    dua.expandAndInsert(n)
  }

  dua.data = append(dua.data, n)

  return nil
}

// Assumes new length won't exceed 4096,
// maybe add more proof of that if we ever introduce
// a callstie besides the one in DynamicArray.Insert 
func (dua *DynamicArray[T]) expandAndInsert(n T) {
  newCap := 0

  oldCap := cap(dua.data)

  if (cap(dua.data) < 64) {
    newCap = oldCap * 2  
  } else if (oldCap < 1067) {
    newCap = oldCap + ( oldCap >> 1 ) // newCap = oldCap * 1.5
  } else if (oldCap > 3840) {
    newCap = MAX_KEYS 
  } else {
    newCap = oldCap * ( oldCap >> 2 ) // newCap = oldCap + 1.25
  }

  newCap = min(newCap, MAX_KEYS)

  newData := make([]T, len(dua.data) + 1, newCap)

  copy(newData, dua.data)

  newData[len(dua.data)] = n
}
