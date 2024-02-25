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
	rules       generic.Resource[res.Rules]
	production  generic.Resource[res.Production]
	stock       generic.Resource[res.Stock]
	ui          generic.Resource[res.UI]
	prodFilter  generic.Filter1[comp.Production]
	consFilter  generic.Filter1[comp.Consumption]
	stockFilter generic.Filter1[comp.Warehouse]
}

// Initialize the system
func (s *UpdateStats) Initialize(world *ecs.World) {
	s.rules = generic.NewResource[res.Rules](world)
	s.production = generic.NewResource[res.Production](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.ui = generic.NewResource[res.UI](world)

	s.prodFilter = *generic.NewFilter1[comp.Production]()
	s.consFilter = *generic.NewFilter1[comp.Consumption]()
	s.stockFilter = *generic.NewFilter1[comp.Warehouse]()
}

// Update the system
func (s *UpdateStats) Update(world *ecs.World) {
	rules := s.rules.Get()
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

	stockQuery := s.stockFilter.Query(world)
	numWarehouses := stockQuery.Count()
	stockQuery.Close()

	for i := resource.Resource(0); i < resource.EndResources; i++ {
		stock.Cap[i] = numWarehouses * rules.StockPerWarehouse[i]
		if stock.Res[i] > stock.Cap[i] {
			stock.Res[i] = stock.Cap[i]
		}

		if i == resource.Food {
			ui.SetResourceLabel(i, fmt.Sprintf("+%d-%d (%d/%d)", production.Prod[i], production.Cons[i], stock.Res[i], stock.Cap[i]))
		} else {
			ui.SetResourceLabel(i, fmt.Sprintf("+%d (%d/%d)", production.Prod[i], stock.Res[i], stock.Cap[i]))
		}
	}

	for i := terr.Terrain(0); i < terr.EndTerrain; i++ {
		props := &terr.Properties[i]
		if !props.CanBuy {
			continue
		}
		ui.SetButtonEnabled(i, stock.CanPay(props.BuildCost))
	}
}

// Finalize the system
func (s *UpdateStats) Finalize(world *ecs.World) {}
