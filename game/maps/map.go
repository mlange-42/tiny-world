package maps

import "image"

type Map struct {
	Terrains              []rune
	Data                  [][]rune
	Achievements          []string
	Center                image.Point
	InitialRandomTerrains int
}
