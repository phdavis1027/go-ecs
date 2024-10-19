package entity

import (
  "github.com/phdavis1027/goecs/util/roaring"
)


type System struct {
  queries         []EntityType // 256 kinds of allowed entity types
  entities        []roaring.RoaringBitset
  OnEntityCreated func(*ECS, Entity, EntityType) (*any, error)
  // NOTE: Here is a list of functions that it is safe to call from CustumOnTick
  // - ecs.DestroyEntity 
  // - ecs.CreateEntity
  CustumOnTick    func(*ECS, []EntityType, []roaring.RoaringBitset)
}

func (s *System) AddEntityTypeToQuery(et EntityType) {
  s.queries  = append(s.queries, et)
  s.entities = append(s.entities, roaring.NewRoaringBitset()) 
}

func (s *System) MatchesQuery(e EntityType) bool {
  // NOTE: Brute-force, but I don't expect there 
  // to more than a couple of queries per system
  for _, et := range s.queries {
    if et == e {
      return true
    }
  }
  return false
}

func (s *System) AddEntity(e Entity, et EntityType) {
  for i, query := range s.queries {
    if query == et {
      s.entities[i].InsertOne(uint64(e))
    }
  }
}

func (s *System) OnTick(ecs *ECS) {
  for i, et := range s.queries {
    if (ecs.dirtyMap[et]) {
      s.entities[i] = 
    }
  }

  if s.OnTick != nil {
    s.CustumOnTick(ecs, s.queries, s.entities)
  }
}
