package res

import "image"

type WorldBounds struct {
	image.Rectangle
}

func (b *WorldBounds) Contains(p image.Point) bool {
	return b.Min.X <= p.X && p.X <= b.Max.X &&
		b.Min.Y <= p.Y && p.Y <= b.Max.Y
}

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
