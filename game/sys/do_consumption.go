package sys

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
)

// DoConsumption system.
type DoConsumption struct {
	speed  generic.Resource[res.GameSpeed]
	time   generic.Resource[res.GameTick]
	update generic.Resource[res.UpdateInterval]
	stock  generic.Resource[res.Stock]

	filter generic.Filter2[comp.UpdateTick, comp.Consumption]
}

// Initialize the system
func (s *DoConsumption) Initialize(world *ecs.World) {
	s.speed = generic.NewResource[res.GameSpeed](world)
	s.time = generic.NewResource[res.GameTick](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.stock = generic.NewResource[res.Stock](world)

	s.filter = *generic.NewFilter2[comp.UpdateTick, comp.Consumption]()
}

// Update the system
func (s *DoConsumption) Update(world *ecs.World) {
	if s.speed.Get().Pause {
		return
	}

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
		cons.Countdown -= int(cons.Amount)
		if cons.Countdown < 0 {
			cons.Countdown += update.Countdown
			stock.Res[cons.Resource]--
		}
	}
}

// Finalize the system
func (s *DoConsumption) Finalize(world *ecs.World) {}
