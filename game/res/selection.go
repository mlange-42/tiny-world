package res

import "github.com/mlange-42/tiny-world/game/terr"

// Selection from the build toolbar.
type Selection struct {
	// Selected terrain type.
	BuildType terr.Terrain
	// ID of the selected button.
	ButtonID int
	// Random sprite index of the selected button.
	RandSprite uint16
	// Whether a special tile has been selected.
	AllowRemove bool
}

// SetBuild sets all fields to the given selection.
func (s *Selection) SetBuild(build terr.Terrain, button int, randSprite uint16, allowRemove bool) {
	s.BuildType = build
	s.ButtonID = button
	s.RandSprite = randSprite
	s.AllowRemove = allowRemove
}

// Reset to select nothing.
func (s *Selection) Reset() {
	s.BuildType = terr.Air
	s.ButtonID = -1
}
