package sys

import (
	"cmp"
	"slices"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// AssignHaulers system.
type AssignHaulers struct {
	speed    generic.Resource[res.GameSpeed]
	update   generic.Resource[res.UpdateInterval]
	landUse  generic.Resource[res.LandUse]
	landUseE generic.Resource[res.LandUseEntities]

	haulerFilter generic.Filter1[comp.Hauler]
	pathFilter   generic.Filter1[comp.Path]

	pathMapper   generic.Map1[comp.Path]
	haulerMapper generic.Map1[comp.Hauler]
	prodMapper   generic.Map1[comp.Production]

	toRemove []ecs.Entity
}

// Initialize the system
func (s *AssignHaulers) Initialize(world *ecs.World) {
	s.speed = generic.NewResource[res.GameSpeed](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.landUse = generic.NewResource[res.LandUse](world)
	s.landUseE = generic.NewResource[res.LandUseEntities](world)

	s.haulerFilter = *generic.NewFilter1[comp.Hauler]()
	s.pathFilter = *generic.NewFilter1[comp.Path]()

	s.pathMapper = generic.NewMap1[comp.Path](world)
	s.haulerMapper = generic.NewMap1[comp.Hauler](world)
	s.prodMapper = generic.NewMap1[comp.Production](world)
}

// Update the system
func (s *AssignHaulers) Update(world *ecs.World) {
	if s.speed.Get().Pause {
		return
	}

	update := s.update.Get()
	landUse := s.landUse.Get()
	landUseE := s.landUseE.Get()

	pathQuery := s.pathFilter.Query(world)
	for pathQuery.Next() {
		path := pathQuery.Get()
		path.Haulers = path.Haulers[:0]
	}

	haulerQuery := s.haulerFilter.Query(world)
	for haulerQuery.Next() {
		haul := haulerQuery.Get()

		// First tile of hauler path - not a path tile.
		if haul.Index == len(haul.Path)-1 && haul.PathFraction < uint8(update.Interval/2+1) {
			continue
		}
		// Last tile of hauler path - not a path tile.
		if haul.Index <= 1 && haul.PathFraction >= uint8(update.Interval/2-1) {
			continue
		}

		frac := float64(haul.PathFraction) / float64(update.Interval)
		p1 := haul.Path[haul.Index]
		p2 := haul.Path[haul.Index-1]

		x := float64(p1.X)*(1-frac) + float64(p2.X)*frac
		y := float64(p1.Y)*(1-frac) + float64(p2.Y)*frac
		yPos := x + y

		xx, yy := int(x+0.5), int(y+0.5)

		lu := landUse.Get(xx, yy)
		if !terr.Properties[lu].TerrainBits.Contains(terr.IsPath) {
			s.toRemove = append(s.toRemove, haulerQuery.Entity())
			continue
		}
		pathHere := landUseE.Get(xx, yy)
		path := s.pathMapper.Get(pathHere)

		path.Haulers = append(path.Haulers, comp.HaulerEntry{Entity: haulerQuery.Entity(), YPos: yPos})
	}

	pathQuery = s.pathFilter.Query(world)
	for pathQuery.Next() {
		path := pathQuery.Get()

		slices.SortStableFunc(path.Haulers, func(a, b comp.HaulerEntry) int {
			return cmp.Compare(a.YPos, b.YPos)
		})
	}

	for _, e := range s.toRemove {
		haul := s.haulerMapper.Get(e)

		if world.Alive(haul.Home) {
			prod := s.prodMapper.Get(haul.Home)
			prod.IsHauling = false
		}

		world.RemoveEntity(e)
	}
	s.toRemove = s.toRemove[:0]
}

// Finalize the system
func (s *AssignHaulers) Finalize(world *ecs.World) {}
