package entity

import (
	"github.com/RoaringBitmap/roaring/roaring64"
)

func SystemHash(s *System) string {
	return s.name
}

type System struct {
	name            string
	queries         []EntityType
  entities        []roaring64.Bitmap
	queriesMut      []EntityType
  entitiesMut        []roaring64.Bitmap
	OnEntityCreated func(*ECS, Entity, EntityType) (*any, error)
	// NOTE: Here is a list of functions that it is safe to call from CustumOnTick
	// - ecs.DestroyEntity
	// - ecs.CreateEntity
	CustumOnTick func(*ECS, []EntityType, []roaring64.Bitmap, []EntityType, []roaring64.Bitmap)
}

func (s *System) AddIfMatches(e Entity, et EntityType) {
	for i, query := range s.queries {
		if query == et {
			s.entities[i].Add(uint64(e))
		}
	}

	for i, query := range s.queriesMut {
		if query == et {
			s.entities[i].Add(uint64(e))
		}
	}
}

func (s *System) AddEntityTypeToQuery(et EntityType) {
	s.queries = append(s.queries, et)
	s.entities = append(s.entities, roaring64.Bitmap{})
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

func (s *System) MatchesQueryMut(e EntityType) bool {
	for _, et := range s.queriesMut {
		if et == e {
			return true
		}
	}
	return false
}

func (s *System) OnTick(ecs *ECS) {
	for i, et := range s.queries {
		if ecs.dirtyMap[et] {
			s.entities[i].And(&ecs.entities[et])
		}
	}

	s.CustumOnTick(ecs, s.queries, s.entities, s.queriesMut, s.entitiesMut)
}
