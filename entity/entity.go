package entity

import (
	// Local libs
	"fmt"
	"sync"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/dominikbraun/graph"
)

type Entity int64

type SystemHandle int

type EntityType uint8
const (
	TILE EntityType = iota
	LEVEL 
)

func (entity Entity) Index() int {
	return int(entity) & 0x00000000FFFFFFFF
}

func (entity Entity) Generation() int {
	return int(entity) >> 32
}

type ECS struct {
	capacity 	      int

	genAlloc 	      GenAllocator

	// Add entities here
	layerComponent          GenArray[int]
	tileMapComponent        GenArray[TileMapComponent]
	positionComponent       GenArray[PositionComponent]
	renderableQuadComponent GenArray[RenderableQuadComponent]


	entities          [256]roaring64.Bitmap
	dirtyMap          [256]bool

	Systems           graph.Graph[string, *System]
	numSystems        int
}

func (ecs *ECS) SendEvent(systemName string, event Event, priority EventPriority) error {
	sys, err := ecs.Systems.Vertex(systemName)
	if err != nil {
		return err
	}

	sys.Mailboxes.SendToPriority(priority, event)

	return nil
}

func (ecs *ECS) AttachRenderableQuadComponent(entity Entity, entityType EntityType) error {
	if !ecs.isValidEntryOfType(entityType, entity) {
		return fmt.Errorf("invalid entity")
	}

	err := ecs.renderableQuadComponent.Set(entity, RenderableQuadComponent{})
	if err != nil {
		return err
	}

	return nil
}

func (ecs *ECS) GetRenderableQuadComponent(entity Entity) (*RenderableQuadComponent, error) {
	_, ok := ecs.isValidEntry(entity)
	if !ok {
		return nil, fmt.Errorf("invalid entity")
	}

	genIndex := GenIndex {
		Index: entity.Index(),
		Generation: entity.Generation(),
	}

	renderableQuadComponent, err := ecs.renderableQuadComponent.Get(genIndex)
	if err != nil {
		return nil, err
	}

	return renderableQuadComponent, nil
}

func (ecs *ECS) GetPositionComponent(entity Entity) (*PositionComponent, error) {
	_, ok := ecs.isValidEntry(entity)
	if !ok {
		return nil, fmt.Errorf("invalid entity")
	}

	genIndex := GenIndex {
		Index: entity.Index(),
		Generation: entity.Generation(),
	}

	positionComponent, err := ecs.positionComponent.Get(genIndex)
	if err != nil {
		return nil, err
	}

	return positionComponent, nil
}

func (ecs *ECS) AttachPositionComponent(entity Entity, entityType EntityType) error {
	if !ecs.isValidEntryOfType(entityType, entity) {
		return fmt.Errorf("invalid entity")
	}

	err := ecs.positionComponent.Set(entity, PositionComponent{})
	if err != nil {
		return err
	}

	return nil
}

func (ecs *ECS) GetTileMapComponent(entity Entity) (*TileMapComponent, error) {
	_, ok := ecs.isValidEntry(entity)
	if !ok {
		return nil, fmt.Errorf("invalid entity")
	}

	genIndex := GenIndex {
		Index: entity.Index(),
		Generation: entity.Generation(),
	}

	tileMapComponent, err := ecs.tileMapComponent.Get(genIndex)
	if err != nil {
		return nil, err
	}

	return tileMapComponent, nil
}

func (ecs *ECS) AttachTileMapComponent(entity Entity, entityType EntityType) error {
	if !ecs.isValidEntryOfType(entityType, entity) {
		return fmt.Errorf("invalid entity")
	}

	err := ecs.tileMapComponent.Set(entity, TileMapComponent{})
	if err != nil {
		return err
	}

	return nil
}

func (ecs *ECS) AttachLayerComponent(entity Entity, entityType EntityType) error {
	if !ecs.isValidEntryOfType(entityType, entity) {
		return fmt.Errorf("invalid entity")
	}

	err := ecs.layerComponent.Set(entity, -1)
	if err != nil {
		return err
	}

	return nil
}

func (ecs *ECS) RegisterSystem(name string, cb SystemFunc) {
	system := new(System)

	system.name = name
	system.CustumOnTick = cb
	system.Mailboxes = NewMailboxesWithCapacity(100)

	ecs.Systems.AddVertex(system)
	ecs.numSystems++
}

type DSetEntry struct {
	name string
	mut  bool
}

func (ecs *ECS) RunSchedule() error {
	quitter := make(chan struct{})

	// Main loop, one game tick
	for {

		select {
		case <-quitter:
			break
		default:
			{
				n := ecs.numSystems


				systems, err := ecs.Systems.Clone()
				if err != nil {
					return err
				}

				for n > 0 {
					predMap, err := systems.PredecessorMap()
					if err != nil {
						return err
					}

					adjMap, err := systems.AdjacencyMap()
					if err != nil {
						return err
					}

					var wg sync.WaitGroup

					for k := range predMap {
						if len(predMap[k]) == 0 {

							// Rip out 0 in-degree system
							for _, v := range adjMap[k] {
								systems.RemoveEdge(v.Source, v.Target)
							}

							system, err := systems.Vertex(k)
							if err != nil {
								return err
							}

							systems.RemoveVertex(k)
							n--

							wg.Add(1)

							go func(hash string, system *System) {
								defer wg.Done()

								system.OnTick(ecs)
							}(k, system)
						}
					}
					wg.Wait()
				} // for n > 0

			} // default
		} // select

	} // for (infinite loop)
}

func (ecs *ECS) CompileSchedule() error {
	// Algorithm 1 from Yao et al
	// iterate over all vertices
	adjMap, err := ecs.Systems.AdjacencyMap()
	if err != nil {
		return err
	}

	dSets := [256][]DSetEntry{}

	for qHash := range adjMap {
		// fmt.Printf( "Processing %s\n", qHash)

		q, err := ecs.Systems.Vertex(qHash)
		if err != nil {
			return err
		}

		for _, query := range q.queries {
			// fmt.Printf("%v\n", dSets)

			dSet := &(dSets[query])

			if len(*dSet) == 0 {
				entry := DSetEntry{name: qHash, mut: false}
				dSets[query] = append(dSets[query], entry)
			} else if len(*dSet) == 1 && (*dSet)[0].mut {
				rHash := (*dSet)[0].name

				err := ecs.Systems.AddEdge(qHash, rHash)
				if err != nil {
					return err
				}

				dSets[query][0] = DSetEntry{name: qHash, mut: false}
			} else {
				maybeLk := dSets[query][0]

				if maybeLk.mut {
					lkHash := maybeLk.name

					err := ecs.Systems.AddEdge(qHash, lkHash)
					if err != nil {
						return err
					}

				}

				entry := DSetEntry{name: qHash, mut: false}
				dSets[query] = append(dSets[query], entry)
			}
		} // for query

		for _, mutQuery := range q.queriesMut {
			dSet := &(dSets[mutQuery])

			if len(*dSet) == 0 {
				entry := DSetEntry{name: qHash, mut: true}
				dSets[mutQuery] = append(dSets[mutQuery], entry)
			} else if len(*dSet) == 1 && (*dSet)[0].mut {
				rHash := (*dSet)[0].name

				err := ecs.Systems.AddEdge(qHash, rHash)
				if err != nil {
					return err
				}

				dSets[mutQuery][0] = DSetEntry{name: qHash, mut: true}
			} else {
				for _, entry := range dSets[mutQuery] {
					if entry.name == qHash {
						continue
					}

					rHash := entry.name

					err := ecs.Systems.AddEdge(qHash, rHash)
					if err != nil {
						return err
					}
				}

				newEntry := DSetEntry{name: qHash, mut: true}
				dSets[mutQuery] = nil
				dSets[mutQuery] = append(dSets[mutQuery], newEntry)
			}
		} // for mutQuery
	} // for q

	return nil
}

func (ecs *ECS) RegisterQueries(system string, queries ...EntityType) error {
	vertex, err := ecs.Systems.Vertex(system)
	if err != nil {
		return err
	}

	vertex.queries = append(vertex.queries, queries...)

	return nil
}

func (ecs *ECS) RegisterMutQueries(system string, queries ...EntityType) error {
	vertex, err := ecs.Systems.Vertex(system)
	if err != nil {
		return err
	}

	vertex.queriesMut = append(vertex.queriesMut, queries...)
	return nil
}

func (ecs *ECS) ScheduleBefore(before, after string) error {
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

	ecs.genAlloc = CreateGenAllocatorOfSize(capacity)

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
