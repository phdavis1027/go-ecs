package main

import (
	"fmt"
	"github.com/phdavis1027/goecs/generational"
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
