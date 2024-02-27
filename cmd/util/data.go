package util

type Sprite struct {
	Id        string     `json:"id"`
	File      []string   `json:"file"`
	Height    int        `json:"height,omitempty"`
	YOffset   int        `json:"y_offset,omitempty"`
	Multitile [][]string `json:"multitile,omitempty"`
}
