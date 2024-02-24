package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/tiny-world/game/res"
)

// Game container
type Game struct {
	Model  *model.Model
	Screen res.EbitenImage
}

// NewGame returns a new game
func NewGame(mod *model.Model) Game {
	return Game{
		Model:  mod,
		Screen: res.EbitenImage{Image: nil, Width: 0, Height: 0},
	}
}

// Initialize the game.
func (g *Game) Initialize() {
	//ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Tiny World")
	g.Model.Initialize()
}

// Run the game.
func (g *Game) Run() error {
	if err := ebiten.RunGame(g); err != nil {
		return err
	}
	return nil
}

// Layout the game.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := 1.0 //ebiten.DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}

// Update the game.
func (g *Game) Update() error {
	g.Model.Update()
	return nil
}

// Draw the game.
func (g *Game) Draw(screen *ebiten.Image) {
	g.Screen.Image = screen
	g.Screen.Width = screen.Bounds().Dx()
	g.Screen.Height = screen.Bounds().Dy()
	g.Model.UpdateUI()
}
