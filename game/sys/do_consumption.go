package sys

import (
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
)

// DoConsumption system.
type DoConsumption struct {
	speed  ecs.Resource[res.GameSpeed]
	time   ecs.Resource[res.GameTick]
	update ecs.Resource[res.UpdateInterval]
	stock  ecs.Resource[res.Stock]
	editor ecs.Resource[res.EditorMode]

	filter *ecs.Filter3[comp.UpdateTick, comp.Production, comp.Consumption]
}

// Initialize the system
func (s *DoConsumption) Initialize(world *ecs.World) {
	s.speed = ecs.NewResource[res.GameSpeed](world)
	s.time = ecs.NewResource[res.GameTick](world)
	s.update = ecs.NewResource[res.UpdateInterval](world)
	s.stock = ecs.NewResource[res.Stock](world)
	s.editor = ecs.NewResource[res.EditorMode](world)

	s.filter = ecs.NewFilter3[comp.UpdateTick, comp.Production, comp.Consumption](world)
}

// Update the system
func (s *DoConsumption) Update(world *ecs.World) {
	if s.speed.Get().Pause {
		return
	}
	isEditor := s.editor.Get().IsEditor

	stock := s.stock.Get()
	tick := s.time.Get().Tick
	update := s.update.Get()
	tickMod := tick % update.Interval

	query := s.filter.Query()
	for query.Next() {
		up, prod, cons := query.Get()

		if up.Tick != tickMod {
			continue
		}

		cons.IsSatisfied = true
		if isEditor {
			continue
		}

		for i, c := range cons.Amount {
			cons.Countdown[i] -= int16(c)
			if cons.Countdown[i] < 0 {
				if prod.Resource == resource.Resource(i) && prod.Stock > 0 {
					cons.Countdown[i] += int16(update.Countdown)
					prod.Stock--
				} else if stock.Res[i] > 0 {
					cons.Countdown[i] += int16(update.Countdown)
					stock.Res[i]--
				} else {
					cons.Countdown[i] = 0
					cons.IsSatisfied = false
				}
			}
		}
	}
}

// Finalize the system
func (s *DoConsumption) Finalize(world *ecs.World) {}
