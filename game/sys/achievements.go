package sys

import (
	"fmt"

	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/res/achievements"
	"github.com/mlange-42/tiny-world/game/save"
)

// Achievements system.
type Achievements struct {
	PlayerFile string

	time         ecs.Resource[res.GameTick]
	update       ecs.Resource[res.UpdateInterval]
	editor       ecs.Resource[res.EditorMode]
	ui           ecs.Resource[res.UI]
	achievements ecs.Resource[achievements.Achievements]
}

// Initialize the system
func (s *Achievements) Initialize(world *ecs.World) {
	s.time = ecs.NewResource[res.GameTick](world)
	s.update = ecs.NewResource[res.UpdateInterval](world)
	s.editor = ecs.NewResource[res.EditorMode](world)
	s.ui = ecs.NewResource[res.UI](world)
	s.achievements = ecs.NewResource[achievements.Achievements](world)
}

// Update the system
func (s *Achievements) Update(world *ecs.World) {
	tick := s.time.Get().Tick

	if tick%s.update.Get().Interval != 0 {
		return
	}
	if s.editor.Get().IsEditor {
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
			s.ui.Get().SetStatusLabel(fmt.Sprintf(" \nAchievement completed!\n\"%s\"\n ", ach.Name))
		}
	}
}

// Finalize the system
func (s *Achievements) Finalize(world *ecs.World) {}
