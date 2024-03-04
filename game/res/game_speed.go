package res

type GameSpeed struct {
	Pause    bool
	Speed    int8
	MinSpeed int8
	MaxSpeed int8
}

type GameTick struct {
	Tick int64
}
