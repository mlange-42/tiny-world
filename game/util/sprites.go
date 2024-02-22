package util

type SpriteInfo struct {
	Name      string
	Index     int
	Height    int
	YOffset   int
	MultiTile bool
}

type SpriteSheet struct {
	SpriteWidth  int
	SpriteHeight int
	Sprites      []SpriteInfo
}
