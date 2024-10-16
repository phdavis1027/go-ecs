package entity

import (
  "errors"
	"fmt"

  // Local libs
	"github.com/phdavis1027/goecs/entity/generational"
  "github.com/phdavis1027/goecs/util/roaring"
)

type Entity int64 

func (entity Entity) Index() int {
  return int(entity) & 0x00000000FFFFFFFF
}

func (entity Entity) Generation() int {
  return int(entity) >> 32
}

type ECS struct {
  capacity        int
    
  genAlloc        generational.GenAllocator 

  // Add entities here
  healthComponent generational.GenArray[int]

  entities        roaring.RoaringBitset
}

// Allocs and owns the memory
func CreateEcsOfCapacity(capacity int) *ECS {
  ecs := new(ECS)

  ecs.genAlloc = generational.CreateGenAllocatorOfSize(capacity)

  ecs.healthComponent = generational.CreateGenArrayOfSize[int](capacity)

  ecs.entities = roaring.NewRoaringBitset() 

  ecs.capacity = capacity

  return ecs
}

func (ecs *ECS) createEntity() (Entity, error) {
  genIndex := ecs.genAlloc.Allocate()

  // A poor man's typecast

  entity := Entity(genIndex.Index | (genIndex.Generation << 32))

  ecs.entities.InsertOne(entity)

  return entity, nil
}

// NOTE: This should never be called directly
// by game-logic code. It should only be called
// by the game engine for cleaning up specific types of entities.
// Otherwise, we risk leaving orphan components.

func (ecs *ECS) AddHealthComponent(entity Entity, initialValue int) error {
  if (entity.Index >= ecs.capacity) {
    fmtString := "Entity index %d is out of bounds for ECS of capacity %d"
    errorMsg := fmt.Sprintf(fmtString, entity.Index, ecs.capacity)

    return errors.New(errorMsg)
  }

  if (ecs.healthComponent[entity.Index].Val.IsSome) {
    fmtString := "Attempt to add a healthComponent to an entity `a`, but that healthComponent" 
    fmtString += "is already claimed by another entity. Found healthComponent: [%s]" 
    errorMsg  := fmt.Sprintf(fmtString, entity, ecs.healthComponent[entity.Index])

    return errors.New(errorMsg)
  }


  ecs.healthComponent[entity.Index].Val.IsSome = true
  ecs.healthComponent[entity.Index].Val.Inner  = initialValue
  ecs.healthComponent[entity.Index].Generation = entity.Generation

  return nil
}
