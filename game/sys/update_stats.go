package sys

import (
	"fmt"
	"math"
	"time"

	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
	"github.com/mlange-42/tiny-world/game/util"
)

// UpdateStats system.
type UpdateStats struct {
	rules          ecs.Resource[res.Rules]
	production     ecs.Resource[res.Production]
	stock          ecs.Resource[res.Stock]
	ui             ecs.Resource[res.UI]
	tick           ecs.Resource[res.GameTick]
	speed          ecs.Resource[res.GameSpeed]
	interval       ecs.Resource[res.UpdateInterval]
	editor         ecs.Resource[res.EditorMode]
	randomTerrains ecs.Resource[res.RandomTerrains]

	prodFilter              *ecs.Filter1[comp.Production]
	consFilter              *ecs.Filter1[comp.Consumption]
	populationFilter        *ecs.Filter1[comp.Population]
	populationSupportFilter *ecs.Filter1[comp.PopulationSupport]
	stockFilter             *ecs.Filter1[comp.Terrain]
	unlockFilter            *ecs.Filter1[comp.Terrain]
}

// Initialize the system
func (s *UpdateStats) Initialize(world *ecs.World) {
	s.rules = ecs.NewResource[res.Rules](world)
	s.production = ecs.NewResource[res.Production](world)
	s.stock = ecs.NewResource[res.Stock](world)
	s.ui = ecs.NewResource[res.UI](world)
	s.tick = ecs.NewResource[res.GameTick](world)
	s.speed = ecs.NewResource[res.GameSpeed](world)
	s.interval = ecs.NewResource[res.UpdateInterval](world)
	s.editor = ecs.NewResource[res.EditorMode](world)
	s.randomTerrains = ecs.NewResource[res.RandomTerrains](world)

	s.prodFilter = ecs.NewFilter1[comp.Production](world)
	s.consFilter = ecs.NewFilter1[comp.Consumption](world)
	s.populationFilter = ecs.NewFilter1[comp.Population](world)
	s.populationSupportFilter = ecs.NewFilter1[comp.PopulationSupport](world)

	s.stockFilter = ecs.NewFilter1[comp.Terrain](world).With(ecs.C[comp.Warehouse]())
	s.unlockFilter = ecs.NewFilter1[comp.Terrain](world).With(ecs.C[comp.UnlocksTerrain]())
}

// Update the system
func (s *UpdateStats) Update(world *ecs.World) {
	rules := s.rules.Get()
	ui := s.ui.Get()
	production := s.production.Get()
	stock := s.stock.Get()
	production.Reset()
	tick := s.tick.Get().Tick
	speed := s.speed.Get()
	interval := s.interval.Get().Interval
	randomTerrains := s.randomTerrains.Get()

	isEditor := s.editor.Get().IsEditor

	prodQuery := s.prodFilter.Query()
	for prodQuery.Next() {
		prod := prodQuery.Get()
		production.Prod[prod.Resource] += int(prod.Amount)
	}
	consQuery := s.consFilter.Query()
	for consQuery.Next() {
		cons := consQuery.Get()
		for i, c := range cons.Amount {
			production.Cons[i] += int(c)
		}
	}

	for i := range resource.Properties {
		stock.Cap[i] = 0
	}

	stockQuery := s.stockFilter.Query()
	for stockQuery.Next() {
		tp := stockQuery.Get()
		prop := &terr.Properties[tp.Terrain]
		for i := range resource.Properties {
			stock.Cap[i] += int(prop.Storage[i])
		}
	}

	randomTerrains.TotalAvailable = rules.InitialRandomTerrains
	unlockQuery := s.unlockFilter.Query()
	for unlockQuery.Next() {
		tp := unlockQuery.Get()
		prop := &terr.Properties[tp.Terrain]
		randomTerrains.TotalAvailable += int(prop.UnlocksTerrains)
	}

	stock.Population = 0
	popQuery := s.populationFilter.Query()
	for popQuery.Next() {
		stock.Population += int(popQuery.Get().Pop)
	}
	stock.MaxPopulation = rules.InitialPopulation
	suppQuery := s.populationSupportFilter.Query()
	for suppQuery.Next() {
		stock.MaxPopulation += int(suppQuery.Get().Pop)
	}

	for i := range resource.Properties {
		if stock.Res[i] > stock.Cap[i] {
			stock.Res[i] = stock.Cap[i]
		}
		if production.Cons[i] > 0 {
			ui.SetResourceLabel(resource.Resource(i),
				fmt.Sprintf("+%d-%d (%d/%d)", production.Prod[i], production.Cons[i], stock.Res[i], stock.Cap[i]),
				production.Cons[i] >= production.Prod[i],
			)
		} else {
			ui.SetResourceLabel(resource.Resource(i),
				fmt.Sprintf("+%d (%d/%d)", production.Prod[i], stock.Res[i], stock.Cap[i]),
				false)
		}
	}
	ui.SetPopulationLabel(fmt.Sprintf("%d/%d", stock.Population, stock.MaxPopulation), stock.Population >= stock.MaxPopulation)

	secs := tick / interval
	duration := time.Duration(secs) * time.Second
	ui.SetTimerLabel(util.FormatDuration(duration))
	speedStr := "P"
	if !speed.Pause {
		if speed.Speed >= 0 {
			speedStr = fmt.Sprintf("x%d", int(math.Pow(2, float64(speed.Speed))))
		} else {
			speedStr = fmt.Sprintf("x1/%d", int(1/math.Pow(2, float64(speed.Speed))))
		}
	}
	ui.SetSpeedLabel(speedStr)

	ui.SetRandomTilesLabel(fmt.Sprintf("%d/%d",
		randomTerrains.TotalAvailable-randomTerrains.TotalPlaced,
		randomTerrains.TotalAvailable))

	// Do the rest only 3x per second
	if tick%(interval/3) != 0 {
		return
	}

	if isEditor {
		return
	}
	for i := range terr.Properties {
		props := &terr.Properties[i]
		if !props.TerrainBits.Contains(terr.CanBuy) {
			continue
		}

		canBuild := true
		message := ""
		if !stock.CanPay(props.BuildCost) {
			message = "Not enough "
			cnt := 0
			for _, cost := range props.BuildCost {
				if stock.Res[cost.Resource] < int(cost.Amount) {
					if cnt > 0 {
						message += ", "
					}
					message += resource.Properties[cost.Resource].Name
					cnt++
				}
			}
			message += "."
			canBuild = false
		}
		if props.Population > 0 && stock.Population+int(props.Population) > stock.MaxPopulation {
			if len(message) > 0 {
				message += "\n"
			}
			message += "Population limit reached."
			canBuild = false
		}
		if canBuild {
			ui.EnableButton(terr.Terrain(i))
		} else {
			ui.DisableButton(terr.Terrain(i), message)
		}
	}
}

// Finalize the system
func (s *UpdateStats) Finalize(world *ecs.World) {}
