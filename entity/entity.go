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

  entityTypes     []EntityType
  entities        []roaring.RoaringBitset
  Systems         []System
  dirtyMap        [256]bool
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

func (ecs *ECS) SetDirty(typ EntityType, dirty bool) {
  ecs.dirtyMap[typ] = dirty 
}

func (ecs *ECS) createEntity() (Entity, error) {
  genIndex := ecs.genAlloc.Allocate()

  entity := Entity(genIndex.Index | (genIndex.Generation << 32))
  ecs.entities.InsertOne(uint64(entity))

  return entity, nil
}

// WARNING: Doesn't work, just written out to convince myself that 
// the `dirty` logic is sane
func (ecs *ECS) DestroyEntity(entity Entity, entityType EntityType) {
  if !ecs.isValidEntry(entity) {
    return
  }

  ecs.entities.DeleteOne(uint64(entity))
  ecs.genAlloc.Deallocate(entity.Index())
  ecs.dirtyMap[entityType] = true 
}

// NOTE: This should never be called directly
// by game-logic code. It should only be called
// by the game engine for cleaning up specific types of entities.
// Otherwise, we risk leaving orphan components.

func (ecs *ECS) CreateEntityOfType(entityType EntityType) (Entity, error) {
  entity, err := ecs.createEntity()
  if err != nil {
    return entity, err
  }

  for _, system := range ecs.Systems {
    if system.MatchesQuery(entityType) {
      system.OnEntityCreated(ecs, entity, entityType)
      system.AddEntity(entity, entityType)
    }
  }

  return entity, nil
}
