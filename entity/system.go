package entity

import (
  "github.com/phdavis1027/goecs/util/roaring"
)


type System struct {
  queries         []EntityType // 256 kinds of allowed entity types
  entities        []roaring.RoaringBitset
  OnEntityCreated func(Entity, EntityType) (*any, error)
  CustumOnTick    func([]EntityType, []roaring.RoaringBitset)
}

func (s *System) AddEntityTypeToQuery(et EntityType) {
  s.queries  = append(s.queries, et)
  s.entities = append(s.entities, roaring.NewRoaringBitset()) 
}

func (s *System) MatchesQuery(e EntityType) bool {
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
  if (ecs.dirty) {
    // We only ever have to delete entities from the system
    // to keep ourselves from getting out of sync with the ECS
    // since when we add entities we always add them to the matching systems
    for i, r := range s.entities {
      s.entities[i] = r.IntersectWith(&ecs.entities)
    }
  }

  if s.OnTick != nil {
    s.CustumOnTick(s.queries, s.entities)
  }
}
