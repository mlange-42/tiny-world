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
	ui         generic.Resource[res.UserInterface]
	filter     generic.Filter1[comp.Production]
}

// Initialize the system
func (s *UpdateStats) Initialize(world *ecs.World) {
	s.production = generic.NewResource[res.Production](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.ui = generic.NewResource[res.UserInterface](world)
	s.filter = *generic.NewFilter1[comp.Production]()
}

// Update the system
func (s *UpdateStats) Update(world *ecs.World) {
	ui := s.ui.Get()
	production := s.production.Get()
	stock := s.stock.Get()
	production.Reset()

	query := s.filter.Query(world)
	for query.Next() {
		prod := query.Get()
		production.Res[prod.Type] += prod.Amount
	}

	for i := resource.Resource(0); i < resource.EndResources; i++ {
		ui.ResourceLabels[i].Label = fmt.Sprintf("%d/%d", production.Res[i], stock.Res[i])
	}
}

// Finalize the system
func (s *UpdateStats) Finalize(world *ecs.World) {}
