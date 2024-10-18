package entity

import (
  "github.com/phdavis1027/goecs/util/roaring"
)


type System struct {
  query    [4]uint64 // 256 kinds of allowed entity types
  entities roaring.RoaringBitset
  OnEntityCreated func(Entity, EntityType) (*any, error)
  CustumOnTick    func(roaring.RoaringBitset)
}

func (s *System) AddEntityTypeToQuery(et EntityType) {
  s.query[et/64] |= 1 << (et % 64)
}

func (s *System) MatchesQuery(e EntityType) bool {
  return (s.query[e/64] & (1 << (e % 64))) != 0
}

func (s *System) AddEntity(e Entity) {
  s.entities.InsertOne(uint64(e))
}

func (s *System) OnTick(ecs *ECS) {
  if (ecs.dirty) {
    // We only ever have to delete entities from the system
    // to keep ourselves from getting out of sync with the ECS
    // since when we add entities we always add them to the matching systems
    ecs.entities.IntersectWith(s.entities)
  }

  if s.OnTick != nil {
    s.CustumOnTick(s.entities)
  }
}
