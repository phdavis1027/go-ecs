package main

import (
  "github.com/phdavis1027/goecs/entity"
) 

func main() {
  ecs := entity.CreateEcsOfCapacity(10)

  humans := make([]entity.Entity, 8)
  orcs   := make([]entity.Entity, 2)
    
  // Create an entity  
  for i := 0; i < 8; i++ {
    human, err := ecs.CreateEntityOfType(entity.Human)
    if err != nil {
      panic(err)
    }
    humans[i] = human
  }
  
  for i := 0; i < 2; i++ {
    orc, err := ecs.CreateEntityOfType(entity.Orc)
    if err != nil {
      panic(err)
    }
    orcs[i] = orc
  }  

}
