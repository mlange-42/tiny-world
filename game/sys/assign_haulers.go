package sys

import (
	"cmp"
	"slices"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
)

// AssignHaulers system.
type AssignHaulers struct {
	haulerFilter generic.Filter1[comp.Hauler]
	pathFilter   generic.Filter1[comp.Path]

	pathMapper generic.Map1[comp.Path]

	update   generic.Resource[res.UpdateInterval]
	landUseE generic.Resource[res.LandUseEntities]
}

// Initialize the system
func (s *AssignHaulers) Initialize(world *ecs.World) {
	s.haulerFilter = *generic.NewFilter1[comp.Hauler]()
	s.pathFilter = *generic.NewFilter1[comp.Path]()

	s.pathMapper = generic.NewMap1[comp.Path](world)

	s.update = generic.NewResource[res.UpdateInterval](world)
	s.landUseE = generic.NewResource[res.LandUseEntities](world)
}

// Update the system
func (s *AssignHaulers) Update(world *ecs.World) {
	update := s.update.Get()
	landUseE := s.landUseE.Get()

	pathQuery := s.pathFilter.Query(world)
	for pathQuery.Next() {
		path := pathQuery.Get()
		path.Haulers = path.Haulers[:0]
	}

	haulerQuery := s.haulerFilter.Query(world)
	for haulerQuery.Next() {
		haul := haulerQuery.Get()

		frac := float64(haul.PathFraction) / float64(update.Interval)
		p1 := haul.Path[haul.Index]
		p2 := haul.Path[haul.Index-1]

		x := float64(p1.X)*(1-frac) + float64(p2.X)*frac
		y := float64(p1.Y)*(1-frac) + float64(p2.Y)*frac
		yPos := y - x/2

		xx, yy := int(x+0.5), int(y+0.5)

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
}

// Finalize the system
func (s *AssignHaulers) Finalize(world *ecs.World) {}
