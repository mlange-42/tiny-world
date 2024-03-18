package sys

import (
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// GameControls system.
type GameControls struct {
	PauseKey      ebiten.Key
	SlowerKey     rune
	FasterKey     rune
	FullscreenKey ebiten.Key

	speed     generic.Resource[res.GameSpeed]
	update    generic.Resource[res.UpdateInterval]
	prevSpeed int8

	inputChars []rune
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

	s.inputChars = ebiten.AppendInputChars(s.inputChars)

	if speed.Speed > speed.MinSpeed && slices.Contains(s.inputChars, s.SlowerKey) {
		speed.Speed -= 1
	}
	if speed.Speed < speed.MaxSpeed && slices.Contains(s.inputChars, s.FasterKey) {
		speed.Speed += 1
	}

	s.inputChars = s.inputChars[:0]

	if s.prevSpeed != speed.Speed {
		ebiten.SetTPS(int(math.Pow(2, float64(speed.Speed)) * float64(update.Interval)))
		s.prevSpeed = speed.Speed
	}
}

// Finalize the system
func (s *GameControls) Finalize(world *ecs.World) {}
