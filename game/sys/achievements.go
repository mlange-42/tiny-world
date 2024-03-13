package sys

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/res/achievements"
	"github.com/mlange-42/tiny-world/game/save"
)

// Achievements system.
type Achievements struct {
	PlayerFile string

	time         generic.Resource[res.GameTick]
	update       generic.Resource[res.UpdateInterval]
	achievements generic.Resource[achievements.Achievements]
}

// Initialize the system
func (s *Achievements) Initialize(world *ecs.World) {
	s.time = generic.NewResource[res.GameTick](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.achievements = generic.NewResource[achievements.Achievements](world)
}

// Update the system
func (s *Achievements) Update(world *ecs.World) {
	tick := s.time.Get().Tick

	if tick%s.update.Get().Interval != 0 {
		return
	}
	achievements := s.achievements.Get()

	for i := range achievements.Achievements {
		ach := &achievements.Achievements[i]
		if ach.Completed {
			continue
		}
		achievements.Check(ach)
		if ach.Completed {
			achievements.Completed = append(achievements.Completed, ach.ID)
			save.SaveAchievements(s.PlayerFile, achievements.Completed)
			println(fmt.Sprintf("Achievement completed: %s", ach.Name))
		}
	}
}

// Finalize the system
func (s *Achievements) Finalize(world *ecs.World) {}
