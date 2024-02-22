package util

type SpriteInfo struct {
	Name    string
	Index   int
	Height  int
	YOffset int
}

type SpriteSheet struct {
	SpriteWidth  int
	SpriteHeight int
	Sprites      []SpriteInfo
}
