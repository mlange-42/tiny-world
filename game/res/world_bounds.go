package res

import "image"

// WorldBounds contain the bounding box of the currently built world.
// Used to limit scrolling/panning.
type WorldBounds struct {
	image.Rectangle
}

// Contains checks whether a pont is inside the bounds.
func (b *WorldBounds) Contains(p image.Point) bool {
	return b.Min.X <= p.X && p.X <= b.Max.X &&
		b.Min.Y <= p.Y && p.Y <= b.Max.Y
}

// AddPoint extends the bounds sp that they contain the given point.
func (b *WorldBounds) AddPoint(p image.Point) {
	if !b.Contains(p) {
		if p.X < b.Min.X {
			b.Min.X = p.X
		}
		if p.Y < b.Min.Y {
			b.Min.Y = p.Y
		}
		if p.X > b.Max.X {
			b.Max.X = p.X
		}
		if p.Y > b.Max.Y {
			b.Max.Y = p.Y
		}
	}
}
