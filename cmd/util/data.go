package util

type RawSprite struct {
	Id        string     `json:"id"`
	File      []string   `json:"file"`
	Height    int        `json:"height,omitempty"`
	YOffset   int        `json:"y_offset,omitempty"`
	Multitile [][]string `json:"multitile,omitempty"`
}

type Sprite struct {
	Id        string  `json:"id"`
	Index     []int   `json:"index"`
	Height    int     `json:"height,omitempty"`
	YOffset   int     `json:"y_offset,omitempty"`
	Multitile [][]int `json:"multitile,omitempty"`
}

func (s *Sprite) IsMultitile() bool {
	return len(s.Multitile) > 0
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
