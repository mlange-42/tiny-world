package res

import "github.com/mlange-42/tiny-world/game/terr"

type Selection struct {
	BuildType terr.Terrain
	ButtonID  int
}

func (s *Selection) SetBuild(build terr.Terrain, button int) {
	s.BuildType = build
	s.ButtonID = button
}

func (s *Selection) Reset() {
	s.BuildType = terr.Air
	s.ButtonID = -1
}
