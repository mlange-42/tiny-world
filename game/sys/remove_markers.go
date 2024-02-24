package sys

import (
	ares "github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
)

// RemoveMarkers system.
type RemoveMarkers struct {
	MaxTime int64

	time   generic.Resource[ares.Tick]
	filter generic.Filter1[comp.ProductionMarker]

	toRemove []ecs.Entity
}

// Initialize the system
func (s *RemoveMarkers) Initialize(world *ecs.World) {
	s.time = generic.NewResource[ares.Tick](world)

	s.filter = *generic.NewFilter1[comp.ProductionMarker]()
}

// Update the system
func (s *RemoveMarkers) Update(world *ecs.World) {
	tick := s.time.Get().Tick

	query := s.filter.Query(world)
	for query.Next() {
		mark := query.Get()
		if tick > mark.StartTick+s.MaxTime {
			s.toRemove = append(s.toRemove, query.Entity())
		}
	}

	for _, e := range s.toRemove {
		world.RemoveEntity(e)
	}
	s.toRemove = s.toRemove[:0]
}

// Finalize the system
func (s *RemoveMarkers) Finalize(world *ecs.World) {}
