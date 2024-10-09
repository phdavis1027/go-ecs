package entity

import (
	"errors"
	"fmt"

	"github.com/phdavis1027/goecs/entity/generational"
)

type Entity generational.GenIndex

type ECS struct {
  capacity        int
    
  genAlloc        generational.GenAllocator 

  // Add entities here
  healthComponent generational.GenArray[int]

  entities        []Entity
}

// Allocs and owns the memory
func CreateEcsOfCapacity(capacity int) *ECS {
  ecs := new(ECS)

  ecs.genAlloc = generational.CreateGenAllocatorOfSize(capacity)

  ecs.healthComponent = generational.CreateGenArrayOfSize[int](capacity)

  ecs.entities = make([]Entity, capacity)

  ecs.capacity = capacity

  return ecs
}

func (ecs *ECS) CreateEntity() (Entity, error) {
  genIndex := ecs.genAlloc.Allocate()

  // A poor man's typecast
  entity := Entity {
    Index: genIndex.Index, 
    Generation: genIndex.Generation,
  }

  return entity, nil
}

type UnderConstructionECS struct {
  inner     ECS
}

func (ecs *UnderConstructionECS) AddHealthComponent(entity Entity, initialValue int) error {
  if (entity.Index >= ecs.capacity) {
    fmtString := "Entity index %d is out of bounds for ECS of capacity %d"
    errorMsg := fmt.Sprintf(fmtString, entity.Index, ecs.capacity)

    return errors.New(errorMsg)
  }


  if (ecs.healthComponent[entity.Index].Val.IsSome) {
    fmtString := "Attempt to add a healthComponent to an entity `a`, but that healthComponent" 
    fmtString += "is already claimed by another entity. Found healthComponent: [%s]" 
    errorMsg  := fmt.Sprintf(fmtString, entity, ecs.inner.healthComponent[entity.Index])

    return errors.New(errorMsg)
  }


  ecs.healthComponent[entity.Index].Val.IsSome = true
  ecs.healthComponent[entity.Index].Val.Inner  = initialValue

  if (entity.Generation )
}
