package sys

import (
	ares "github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
)

// DoProduction system.
type DoProduction struct {
	time   generic.Resource[ares.Tick]
	update generic.Resource[res.UpdateInterval]
	stock  generic.Resource[res.Stock]

	filter generic.Filter2[comp.UpdateTick, comp.Production]
}

// Initialize the system
func (s *DoProduction) Initialize(world *ecs.World) {
	s.time = generic.NewResource[ares.Tick](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.stock = generic.NewResource[res.Stock](world)

	s.filter = *generic.NewFilter2[comp.UpdateTick, comp.Production]()
}

// Update the system
func (s *DoProduction) Update(world *ecs.World) {
	stock := s.stock.Get()
	tick := s.time.Get().Tick
	interval := s.update.Get().Interval
	tickMod := tick % interval

	query := s.filter.Query(world)
	for query.Next() {
		up, pr := query.Get()

		if up.Tick != tickMod {
			continue
		}
		pr.Countdown -= pr.Amount
		if pr.Countdown < 0 {
			pr.Countdown += 100
			stock.Res[pr.Type]++
		}
	}
}

// Finalize the system
func (s *DoProduction) Finalize(world *ecs.World) {}
