package res

import (
	"fmt"
	"io/fs"

	"github.com/mlange-42/tiny-world/cmd/util"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Rules resource, holding game rules read from JSON.
type Rules struct {
	// World extent in X and Y direction, in number of tiles.
	WorldSize int
	// Initial build radius around the starting position.
	InitialBuildRadius int
	// Initial population limit.
	InitialPopulation int
	// Number of random terrains/cards.
	RandomTerrainsCount int
	// List of terrains to draw cards from.
	RandomTerrains []terr.Terrain
	// Initial resources of the player.
	InitialResources []int
	// Probability of special cards/terrains that can be placed over existing terrain.
	SpecialCardProbability float64
}

// NewRules reads rules from the given file.
func NewRules(f fs.FS, file string) Rules {
	rulesHelper := rulesJs{}
	err := util.FromJsonFs(f, file, &rulesHelper)
	if err != nil {
		panic(err)
	}

	storage := make([]int, len(resource.Properties))
	for _, entry := range rulesHelper.InitialResources {
		res, ok := resource.ResourceID(entry.Resource)
		if !ok {
			panic(fmt.Sprintf("unknown resource %s", entry.Resource))
		}
		storage[res] = entry.Amount
	}

	randTerr := []terr.Terrain{}
	for _, str := range rulesHelper.RandomTerrains {
		randTerr = append(randTerr, toTerrain(str))
	}

	return Rules{
		WorldSize:              rulesHelper.WorldSize,
		InitialBuildRadius:     rulesHelper.InitialBuildRadius,
		InitialPopulation:      rulesHelper.InitialPopulation,
		InitialResources:       storage,
		RandomTerrainsCount:    rulesHelper.RandomTerrainsCount,
		RandomTerrains:         randTerr,
		SpecialCardProbability: rulesHelper.SpecialCardProbability,
	}
}

type rulesJs struct {
	WorldSize           int `json:"world_size"`
	InitialBuildRadius  int `json:"initial_build_radius"`
	InitialPopulation   int `json:"initial_population"`
	RandomTerrainsCount int `json:"random_terrains_count"`

	RandomTerrains         []string           `json:"random_terrains"`
	InitialResources       []resourceAmountJs `json:"initial_resources"`
	SpecialCardProbability float64            `json:"special_card_probability"`
}

type resourceAmountJs struct {
	Resource string `json:"resource"`
	Amount   int    `json:"amount"`
}

func toTerrain(t string) terr.Terrain {
	id, ok := terr.TerrainID(t)
	if !ok {
		panic(fmt.Sprintf("unknown terrain %s", t))
	}
	return id
}
