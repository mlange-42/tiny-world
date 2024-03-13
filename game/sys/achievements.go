package sys

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/save"
)

// Achievements system.
type Achievements struct {
	PlayerFile string

	time         generic.Resource[res.GameTick]
	update       generic.Resource[res.UpdateInterval]
	achievements generic.Resource[res.Achievements]

	completed []string
}

// Initialize the system
func (s *Achievements) Initialize(world *ecs.World) {
	s.time = generic.NewResource[res.GameTick](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.achievements = generic.NewResource[res.Achievements](world)

	err := save.LoadAchievements(s.PlayerFile, &s.completed)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return
		}
		log.Fatal("error parsing achievement: ", err)
	}

	ach := s.achievements.Get()
	for i := range ach.Achievements {
		if slices.Contains(s.completed, ach.Achievements[i].Name) {
			ach.Achievements[i].Completed = true
		}
	}
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
			s.completed = append(s.completed, ach.Name)
			save.SaveAchievements(s.PlayerFile, s.completed)
			println(fmt.Sprintf("Achievement completed: %s", ach.Name))
		}
	}
}

// Finalize the system
func (s *Achievements) Finalize(world *ecs.World) {}
