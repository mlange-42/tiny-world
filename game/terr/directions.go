package terr

type Direction uint8

const (
	N Direction = iota
	E
	S
	W
	EndDirection
)

func (d Direction) Deltas() (int, int) {
	dir := directionXY[d]
	return dir[0], dir[1]
}

func (d Direction) Opposite() Direction {
	return (d + 4) % EndDirection
}

type Directions uint8

func NewDirections(dirs ...Direction) Directions {
	d := Directions(0)
	for _, dir := range dirs {
		d |= (1 << dir)
	}
	return d
}

func (d *Directions) Set(dir Direction) {
	*d |= (1 << dir)
}

func (d *Directions) Unset(dir Direction) {
	*d &= ^(1 << dir)
}

// Contains checks whether all the argument's bits are contained in this Subscription.
func (d Directions) Contains(dir Direction) bool {
	bits := Directions(1 << dir)
	return (bits & d) == bits
}

var directionXY = [4][2]int{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}
