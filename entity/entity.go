package entity

import (
  // Local libs
	"github.com/phdavis1027/goecs/entity/generational"
  "github.com/phdavis1027/goecs/util/roaring"
)

type Entity int64 

type EntityType uint8 
const (
  Human EntityType = iota
  Orc
)

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
  Systems         []System
  dirty           bool
}

func (ecs *ECS) isValidEntry(entity Entity) bool {
  return ecs.entities.Has(uint64(entity))
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

  entity := Entity(genIndex.Index | (genIndex.Generation << 32))
  ecs.entities.InsertOne(uint64(entity))

  return entity, nil
}

// NOTE: This should never be called directly
// by game-logic code. It should only be called
// by the game engine for cleaning up specific types of entities.
// Otherwise, we risk leaving orphan components.

func (ecs *ECS) CreateEntityOfType(entityType EntityType) (Entity, error) {
  entity, err := ecs.createEntity()

  if (err != nil) {
    return entity, err
  }

  for _, system := range ecs.Systems {
    if system.MatchesQuery(entityType) {
      system.OnEntityCreated(entity, entityType)
      system.AddEntity(entity)
    }
  }

  return entity, nil
}
