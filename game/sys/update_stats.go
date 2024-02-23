package sys

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
)

// UpdateStats system.
type UpdateStats struct {
	production generic.Resource[res.Production]
	stock      generic.Resource[res.Stock]
	ui         generic.Resource[res.HUD]
	prodFilter generic.Filter1[comp.Production]
	consFilter generic.Filter1[comp.Consumption]
}

// Initialize the system
func (s *UpdateStats) Initialize(world *ecs.World) {
	s.production = generic.NewResource[res.Production](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.ui = generic.NewResource[res.HUD](world)

	s.prodFilter = *generic.NewFilter1[comp.Production]()
	s.consFilter = *generic.NewFilter1[comp.Consumption]()
}

// Update the system
func (s *UpdateStats) Update(world *ecs.World) {
	ui := s.ui.Get()
	production := s.production.Get()
	stock := s.stock.Get()
	production.Reset()

	prodQuery := s.prodFilter.Query(world)
	for prodQuery.Next() {
		prod := prodQuery.Get()
		production.Prod[prod.Type] += prod.Amount
	}
	consQuery := s.consFilter.Query(world)
	for consQuery.Next() {
		cons := consQuery.Get()
		production.Cons[resource.Food] += cons.Amount
	}

	for i := resource.Resource(0); i < resource.EndResources; i++ {
		if i == resource.Food {
			ui.ResourceLabels[i].Label = fmt.Sprintf("+%d-%d (%d)", production.Prod[i], production.Cons[i], stock.Res[i])
		} else {
			ui.ResourceLabels[i].Label = fmt.Sprintf("+%d (%d)", production.Prod[i], stock.Res[i])
		}
	}
}

// Finalize the system
func (s *UpdateStats) Finalize(world *ecs.World) {}
