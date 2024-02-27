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
