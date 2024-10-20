package entity

import (
	// Local libs
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/dominikbraun/graph"
	"github.com/phdavis1027/goecs/entity/generational"
)

type Entity int64

type EntityType uint8
type SystemHandle int

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
	capacity int

	genAlloc generational.GenAllocator

	// Add entities here
	healthComponent generational.GenArray[int]

	entities [256]roaring64.Bitmap
	dirtyMap [256]bool

	Systems graph.Graph[string, *System]
}

func (ecs *ECS) RegisterSystem(name string, cb func(*ECS, []EntityType, []roaring64.Bitmap), queries, queriesMut []EntityType) {
	system := new(System)

	system.name = name
	system.CustumOnTick = cb
	system.queries = queries
	system.queriesMut = queriesMut

	ecs.Systems.AddVertex(system)
}

func (ecs *ECS) CompileSchedule() error {
	// TODO: I need to implement this
	return nil
}

func (ecs *ECS) RegisterQueries(system string, queries []EntityType) error {
	vertex, err := ecs.Systems.Vertex(system)
	if err != nil {
		return err
	}

	vertex.queries = append(vertex.queries, queries...)

	return nil
}

func (ecs *ECS) RegisterMutQueries(system string, queries []EntityType) error {
	vertex, err := ecs.Systems.Vertex(system)
	if err != nil {
		return err
	}

	vertex.queriesMut = append(vertex.queriesMut, queries...)
	return nil
}

func (ecs *ECS) ScheduleDependency(before, after string) error {
	return ecs.Systems.AddEdge(before, after)
}

func (ecs *ECS) isValidEntryOfType(typ EntityType, entity Entity) bool {
	return ecs.entities[typ].Contains(uint64(entity))
}

func (ecs *ECS) isValidEntry(entity Entity) (EntityType, bool) {
	for typ := 0; typ < len(ecs.entities); typ++ {
		if ecs.isValidEntryOfType(EntityType(typ), entity) {
			return EntityType(typ), true
		}
	}
	return 0, false
}

// Allocs and owns the memory
func CreateEcsOfCapacity(capacity int) *ECS {
	ecs := new(ECS)

	ecs.genAlloc = generational.CreateGenAllocatorOfSize(capacity)

	ecs.healthComponent = generational.CreateGenArrayOfSize[int](capacity)

	ecs.capacity = capacity

	ecs.Systems = graph.New(SystemHash, graph.Directed(), graph.Acyclic())

	return ecs
}

func (ecs *ECS) SetDirty(typ EntityType, dirty bool) {
	ecs.dirtyMap[typ] = dirty
}

// WARNING: Doesn't work, just written out to convince myself that
// the `dirty` logic is sane
// func (ecs *ECS) DestroyEntity(entity Entity, entityType EntityType) {
//   if !ecs.isValidEntryOfType(entityType, entity) {
//     return
//   }
//
//   ecs.entities.DeleteOne(uint64(entity))
//   ecs.genAlloc.Deallocate(entity.Index())
//   ecs.dirtyMap[entityType] = true
// }

// NOTE: This should never be called directly
// by game-logic code. It should only be called
// by the game engine for cleaning up specific types of entities.
// Otherwise, we risk leaving orphan components.

func (ecs *ECS) CreateEntityOfType(entityType EntityType) (Entity, error) {
	genIndex := ecs.genAlloc.Allocate()

	entity := Entity(genIndex.Index | (genIndex.Generation << 32))
	ecs.entities[entityType].Add(uint64(entity))

	adjMap, err := ecs.Systems.AdjacencyMap()
	if err != nil {
		return entity, err
	}

	for hash := range adjMap {
		vertex, err := ecs.Systems.Vertex(hash)
		if err != nil {
			return entity, err
		}

		vertex.AddIfMatches(entity, entityType)
	}

	return entity, nil
}
