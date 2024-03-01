package util

import "image/color"

type TileSet struct {
	TileWidth       int        `json:"tile_width"`
	TileHeight      int        `json:"tile_height"`
	BackgroundColor color.RGBA `json:"background_color"`
	TextColor       color.RGBA `json:"text_color"`
}

type RawSprite struct {
	Id         string     `json:"id"`
	File       []string   `json:"file"`
	Height     int        `json:"height,omitempty"`
	YOffset    int        `json:"y_offset,omitempty"`
	AnimFrames int        `json:"anim_frames,omitempty"`
	AnimSpeed  int        `json:"anim_speed,omitempty"`
	Multitile  [][]string `json:"multitile,omitempty"`
}

type Sprite struct {
	Id         string  `json:"id"`
	Index      []int   `json:"index"`
	Height     int     `json:"height,omitempty"`
	YOffset    int     `json:"y_offset,omitempty"`
	AnimFrames int     `json:"anim_frames,omitempty"`
	AnimSpeed  int     `json:"anim_speed,omitempty"`
	Multitile  [][]int `json:"multitile,omitempty"`
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
	SpriteWidth  int      `json:"sprite_width"`
	SpriteHeight int      `json:"sprite_height"`
	Sprites      []Sprite `json:"sprites"`
	TotalSprites int      `json:"total_sprites"`
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
