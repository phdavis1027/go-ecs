package main

import (
  "github.com/phdavis1027/goecs/generational"
  "fmt"
)

func main() {
  genAlloc := generational.CreateGenAllocatorOfSize(10)
  healthComponent := generational.CreateGenArrayOfSize[int](10)


  healthComponent.Set(genIndex, 100)


  value, err := healthComponent.Get(genIndex)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(*value)
  }
} 
