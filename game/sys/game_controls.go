package sys

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// GameControls system.
type GameControls struct {
	PauseKey      ebiten.Key
	SlowerKey     ebiten.Key
	FasterKey     ebiten.Key
	FullscreenKey ebiten.Key

	speed     generic.Resource[res.GameSpeed]
	update    generic.Resource[res.UpdateInterval]
	prevSpeed int8
}

// Initialize the system
func (s *GameControls) Initialize(world *ecs.World) {
	s.speed = generic.NewResource[res.GameSpeed](world)
	s.update = generic.NewResource[res.UpdateInterval](world)

	speed := s.speed.Get()
	update := s.update.Get()

	ebiten.SetTPS(int(math.Pow(2, float64(speed.Speed)) * float64(update.Interval)))
}

// Update the system
func (s *GameControls) Update(world *ecs.World) {
	speed := s.speed.Get()
	update := s.update.Get()

	if inpututil.IsKeyJustPressed(s.FullscreenKey) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if inpututil.IsKeyJustPressed(s.PauseKey) {
		speed.Pause = !speed.Pause
	}
	if inpututil.IsKeyJustPressed(s.SlowerKey) && speed.Speed > speed.MinSpeed {
		speed.Speed -= 1
	}
	if inpututil.IsKeyJustPressed(s.FasterKey) && speed.Speed < speed.MaxSpeed {
		speed.Speed += 1
	}

	if s.prevSpeed != speed.Speed {
		ebiten.SetTPS(int(math.Pow(2, float64(speed.Speed)) * float64(update.Interval)))
		s.prevSpeed = speed.Speed
	}
}

// Finalize the system
func (s *GameControls) Finalize(world *ecs.World) {}
