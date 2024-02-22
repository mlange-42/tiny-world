package util

type SpriteInfo struct {
	Name  string
	Index int
}

type SpriteSheet struct {
	SpriteWidth  int
	SpriteHeight int
	Sprites      []SpriteInfo
}
