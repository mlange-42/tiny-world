package res

import "github.com/mlange-42/tiny-world/game/terr"

type Selection struct {
	BuildType  terr.Terrain
	ButtonID   int
	RandSprite uint16
}

func (s *Selection) SetBuild(build terr.Terrain, button int, randSprite uint16) {
	s.BuildType = build
	s.ButtonID = button
	s.RandSprite = randSprite
}

func (s *Selection) Reset() {
	s.BuildType = terr.Air
	s.ButtonID = -1
}
