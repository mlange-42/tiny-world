package sys

import (
	ares "github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
)

// DoConsumption system.
type DoConsumption struct {
	time   generic.Resource[ares.Tick]
	update generic.Resource[res.UpdateInterval]
	stock  generic.Resource[res.Stock]

	filter generic.Filter2[comp.UpdateTick, comp.Consumption]
}

// Initialize the system
func (s *DoConsumption) Initialize(world *ecs.World) {
	s.time = generic.NewResource[ares.Tick](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.stock = generic.NewResource[res.Stock](world)

	s.filter = *generic.NewFilter2[comp.UpdateTick, comp.Consumption]()
}

// Update the system
func (s *DoConsumption) Update(world *ecs.World) {
	stock := s.stock.Get()
	tick := s.time.Get().Tick
	update := s.update.Get()
	tickMod := tick % update.Interval

	query := s.filter.Query(world)
	for query.Next() {
		up, cons := query.Get()

		if up.Tick != tickMod {
			continue
		}
		cons.Countdown -= cons.Amount
		if cons.Countdown < 0 {
			cons.Countdown += update.Countdown
			stock.Res[resource.Food]--
		}
	}
}

// Finalize the system
func (s *DoConsumption) Finalize(world *ecs.World) {}
