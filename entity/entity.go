package entity

import (
	// Local libs
	"fmt"
	"os"
	"sync"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/phdavis1027/goecs/entity/generational"
)

type Entity int64

type EntityType uint8
type SystemHandle int

const (
	Zero EntityType = iota
  One	
  Two 
  Three 
  Tile
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
  rectComponent   generational.GenArray[Rect]        
  layerComponent  generational.GenArray[int]

	entities [256]roaring64.Bitmap
	dirtyMap [256]bool

	Systems graph.Graph[string, *System]
  numSystems int
}

func (ecs *ECS) AttachLayerComponent(entity Entity, entityType EntityType, layer int) error {
  if !ecs.isValidEntryOfType(entityType, entity) {
    return fmt.Errorf("invalid entity")
  }

  err := ecs.layerComponent.Set(entity, layer)
  if err != nil {
    return err
  }

  return nil
}

func (ecs *ECS) AttachRectComponent(entity Entity, entityType EntityType, rect Rect ) error {
  if !ecs.isValidEntryOfType(entityType, entity) {
    return fmt.Errorf("invalid entity")
  }

  err := ecs.rectComponent.Set(entity, rect)
  if err != nil {
    return err
  }

  return nil
}

func (ecs *ECS) RegisterSystem(name string, cb func(*ECS, []EntityType, []roaring64.Bitmap, []EntityType, []roaring64.Bitmap)) {
	system := new(System)

	system.name = name
	system.CustumOnTick = cb

	ecs.Systems.AddVertex(system)
  ecs.numSystems++
}

type DSetEntry struct {
  name     string
  mut      bool
}

func (ecs *ECS) RunSchedule() error {
  quitter := make(chan struct{})

  for {
    select {
      case <- quitter:
        break
      default: {
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

              systems.RemoveVertex( k )
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

func (ecs *ECS) CompileSchedule(debug bool) error {
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
        entry := DSetEntry{ name: qHash, mut: false }
        dSets[query] = append(dSets[query], entry)
      } else if len(*dSet) == 1 && (*dSet)[0].mut {
        rHash := (*dSet)[0].name

        err := ecs.Systems.AddEdge(qHash, rHash)
        if err != nil {
          return err
        }

        dSets[query][0] = DSetEntry{ name: qHash, mut: false }
      } else {
        maybeLk := dSets[query][0]

        if (maybeLk.mut) {
          lkHash := maybeLk.name

          err := ecs.Systems.AddEdge(qHash, lkHash)
          if err != nil {
            return err
          }

        }

        entry := DSetEntry{ name: qHash, mut: false }
        dSets[query] = append(dSets[query], entry)
      }
    } // for query

    for _, mutQuery := range q.queriesMut {
      dSet := &(dSets[mutQuery])

      if len(*dSet) == 0  {
        entry := DSetEntry{ name: qHash, mut: true }
        dSets[mutQuery] = append(dSets[mutQuery], entry)
      } else if len(*dSet) == 1 && (*dSet)[0].mut {
        rHash := (*dSet)[0].name

        err := ecs.Systems.AddEdge(qHash, rHash)
        if err != nil {
          return err
        }

        dSets[mutQuery][0] = DSetEntry{ name: qHash, mut: true }
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

        newEntry := DSetEntry{ name: qHash, mut: true }
        dSets[mutQuery] = nil 
        dSets[mutQuery] = append(dSets[mutQuery], newEntry)
      }
    } // for mutQuery
  } // for q

  if debug {
    file, _ := os.Create("debug.txt")
    defer file.Close()
    draw.DOT(ecs.Systems, file)
  }

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

func (ecs *ECS) ScheduleManualDependency(before, after string) error {
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
