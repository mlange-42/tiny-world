package res

// GameSpeed resource.
type GameSpeed struct {
	// Is the game paused?
	Pause bool
	// Game speed as an exponent for base 2. s = 2^Speed
	Speed int8
	// Minimum game speed, as an exponent for base 2.
	MinSpeed int8
	// Maximum game speed, as an exponent for base 2.
	MaxSpeed int8
}

// GameTick resource.
type GameTick struct {
	// Current update tick. Stops when the game is paused.
	Tick int64
	// Current render tick. Does not stop when the game is paused.
	RenderTick int64
}
