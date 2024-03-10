package util

import "image"

func Clamp(r image.Rectangle, p image.Point) image.Point {
	if p.X < r.Min.X {
		p.X = r.Min.X
	}
	if p.Y < r.Min.Y {
		p.Y = r.Min.Y
	}
	if p.X > r.Max.X {
		p.X = r.Max.X
	}
	if p.Y > r.Max.Y {
		p.Y = r.Max.Y
	}
	return p
}
