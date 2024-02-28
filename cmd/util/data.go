package util

import (
	"fmt"

	"github.com/mlange-42/tiny-world/game/terr"
)

type TileSet struct {
	TileWidth  int `json:"tile_width"`
	TileHeight int `json:"tile_height"`
}

type RawSprite struct {
	Id         string     `json:"id"`
	File       []string   `json:"file"`
	Height     int        `json:"height,omitempty"`
	Below      string     `json:"below,omitempty"`
	YOffset    int        `json:"y_offset,omitempty"`
	AnimFrames int        `json:"anim_frames,omitempty"`
	AnimSpeed  int        `json:"anim_speed,omitempty"`
	Multitile  [][]string `json:"multitile,omitempty"`
}

type JsSprite struct {
	Id         string  `json:"id"`
	Index      []int   `json:"index"`
	Height     int     `json:"height,omitempty"`
	Below      string  `json:"below,omitempty"`
	YOffset    int     `json:"y_offset,omitempty"`
	AnimFrames int     `json:"anim_frames,omitempty"`
	AnimSpeed  int     `json:"anim_speed,omitempty"`
	Multitile  [][]int `json:"multitile,omitempty"`
}

type Sprite struct {
	Id         string
	Index      []int
	Height     int
	YOffset    int
	Below      terr.Terrain
	AnimFrames int
	AnimSpeed  int
	Multitile  [][]int
}

func NewSprite(s *JsSprite, lookup func(string) (terr.Terrain, bool)) Sprite {
	var below terr.Terrain

	if s.Below != "" {
		var ok bool
		below, ok = lookup(s.Below)
		if !ok {
			panic(fmt.Sprintf("unknown terrain '%s' in sprite %s", s.Below, s.Id))
		}
	}

	return Sprite{
		Id:         s.Id,
		Index:      s.Index,
		Height:     s.Height,
		YOffset:    s.YOffset,
		Below:      below,
		AnimFrames: s.AnimFrames,
		AnimSpeed:  s.AnimSpeed,
		Multitile:  s.Multitile,
	}
}

func (s *Sprite) IsMultitile() bool {
	return len(s.Multitile) > 0
}

func (s *Sprite) IsAnimated() bool {
	return s.AnimFrames > 1
}

type RawSpriteSheet struct {
	Directory string
	Width     int
	Height    int
}

type SpriteSheet struct {
	SpriteWidth  int        `json:"sprite_width"`
	SpriteHeight int        `json:"sprite_height"`
	Sprites      []JsSprite `json:"sprites"`
	TotalSprites int        `json:"total_sprites"`
}

type Directory struct {
	Dir         string
	HasJson     bool
	Files       []File
	Directories []Directory
}

type File struct {
	Name   string
	IsJson bool
}
