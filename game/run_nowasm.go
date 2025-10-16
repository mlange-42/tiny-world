//go:build !js

package game

import (
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/ark-repl/repl"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/res"
)

type canvasHelper struct{}

func newCanvasHelper() *canvasHelper {
	return &canvasHelper{}
}

func (c *canvasHelper) isMouseInside(width, height int) bool {
	x, y := ebiten.CursorPosition()
	return x >= 0 && y >= 0 && x < width && y < height
}

func addRepl(app *app.App) {
	startServer := len(os.Args) > 1 && os.Args[1] == "monitor"
	if !startServer {
		return
	}

	callbacks := repl.Callbacks{
		Pause: func(out *strings.Builder) {
			ecs.GetResource[res.GameSpeed](&app.World).Pause = true
		},
		Resume: func(out *strings.Builder) {
			ecs.GetResource[res.GameSpeed](&app.World).Pause = false
		},
		Ticks: func() int {
			return int(ecs.GetResource[res.GameTick](&app.World).Tick)
		},
	}

	repl := repl.NewRepl(&app.World, callbacks)

	app.AddUISystem(repl.System())

	repl.StartServer(":9000")
}
