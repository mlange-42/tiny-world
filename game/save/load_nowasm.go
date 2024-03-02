//go:build !js

package save

import (
	"os"

	serde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
)

func loadWorld(world *ecs.World, path string) error {
	jsData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return serde.Deserialize(jsData, world)
}
