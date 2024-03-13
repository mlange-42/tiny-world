package res

import (
	"io/fs"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/cmd/util"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/terr"
)

type Achievement struct {
	Name        string
	Description string
	Conditions  []Condition
}

type Condition struct {
	Type   string
	IDs    uint32
	Number int
}

type Achievements struct {
	world *ecs.World

	terrainFilter generic.Filter1[comp.Terrain]

	checks       map[string]func([]int, int) bool
	achievements []Achievement
}

func NewAchievements(world *ecs.World, f fs.FS, file string) *Achievements {
	a := Achievements{
		world:         world,
		terrainFilter: *generic.NewFilter1[comp.Terrain](),
	}

	a.checks = map[string]func([]int, int) bool{
		"terrain": a.checkTerrain,
	}
	parse := map[string]func(...string) uint32{
		"terrain": a.parseTerrains,
	}

	ach := []achievementJs{}
	err := util.FromJsonFs(f, file, &ach)
	if err != nil {
		panic(err)
	}

	for _, achieve := range ach {
		conditions := []Condition{}

		for _, c := range achieve.Conditions {
			conditions = append(conditions,
				Condition{
					Type:   c.Type,
					IDs:    parse[c.Type](c.IDs...),
					Number: c.Number,
				},
			)
		}

		a.achievements = append(a.achievements,
			Achievement{
				Name:        achieve.Name,
				Description: achieve.Description,
				Conditions:  conditions,
			},
		)
	}

	return &a
}

func (a *Achievements) checkTerrain(ids []int, num int) bool {

	return false
}

func (a *Achievements) parseTerrains(ids ...string) uint32 {
	return uint32(terr.ToTerrains(ids...))
}

type achievementJs struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Conditions  []conditionJs `json:"conditions"`
}

type conditionJs struct {
	Type   string   `json:"type"`
	IDs    []string `json:"ids"`
	Number int      `json:"number"`
}
