package sys

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

// UpdateStats system.
type UpdateStats struct {
	production  generic.Resource[res.Production]
	stock       generic.Resource[res.Stock]
	ui          generic.Resource[res.UI]
	prodFilter  generic.Filter1[comp.Production]
	consFilter  generic.Filter1[comp.Consumption]
	stockFilter generic.Filter1[comp.Terrain]
}

// Initialize the system
func (s *UpdateStats) Initialize(world *ecs.World) {
	s.production = generic.NewResource[res.Production](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.ui = generic.NewResource[res.UI](world)

	s.prodFilter = *generic.NewFilter1[comp.Production]()
	s.consFilter = *generic.NewFilter1[comp.Consumption]()
	s.stockFilter = *generic.NewFilter1[comp.Terrain]().With(generic.T[comp.Warehouse]())
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
		production.Prod[prod.Resource] += prod.Amount
	}
	consQuery := s.consFilter.Query(world)
	for consQuery.Next() {
		cons := consQuery.Get()
		production.Cons[cons.Resource] += cons.Amount
	}

	for i := range resource.Properties {
		stock.Cap[i] = 0
	}
	stockQuery := s.stockFilter.Query(world)
	for stockQuery.Next() {
		tp := stockQuery.Get()
		st := terr.Properties[tp.Terrain].Storage
		for i := range resource.Properties {
			stock.Cap[i] += st[i]
		}
	}

	for i := range resource.Properties {
		if stock.Res[i] > stock.Cap[i] {
			stock.Res[i] = stock.Cap[i]
		}
		if production.Cons[i] > 0 {
			ui.SetResourceLabel(resource.Resource(i), fmt.Sprintf("+%d-%d (%d/%d)", production.Prod[i], production.Cons[i], stock.Res[i], stock.Cap[i]))
		} else {
			ui.SetResourceLabel(resource.Resource(i), fmt.Sprintf("+%d (%d/%d)", production.Prod[i], stock.Res[i], stock.Cap[i]))
		}
	}

	for i := range terr.Properties {
		props := &terr.Properties[i]
		if !props.CanBuy {
			continue
		}
		ui.SetButtonEnabled(terr.Terrain(i), stock.CanPay(props.BuildCost))
	}
}

// Finalize the system
func (s *UpdateStats) Finalize(world *ecs.World) {}
