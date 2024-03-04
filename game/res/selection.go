package res

import "github.com/mlange-42/tiny-world/game/terr"

type Selection struct {
	BuildType   terr.Terrain
	ButtonID    int
	RandSprite  uint16
	AllowRemove bool
}

func (s *Selection) SetBuild(build terr.Terrain, button int, randSprite uint16, allowRemove bool) {
	s.BuildType = build
	s.ButtonID = button
	s.RandSprite = randSprite
	s.AllowRemove = allowRemove
}

func (s *Selection) Reset() {
	s.BuildType = terr.Air
	s.ButtonID = -1
}
