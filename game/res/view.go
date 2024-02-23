package res

import "image"

type View struct {
	TileWidth, TileHeight int
	X, Y                  int
	Zoom                  float64
	MouseOffset           int
}

func (v *View) Offset() image.Point {
	return image.Pt(int(float64(v.X)*v.Zoom), int(float64(v.Y)*v.Zoom))
}

func (v *View) Bounds(w, h int) image.Rectangle {
	vw, vh := int(float64(w)/v.Zoom), int(float64(h)/v.Zoom)

	return image.Rect(
		v.X-v.TileWidth, v.Y-2*v.TileHeight,
		v.X+vw, v.Y+vh+2*v.TileHeight,
	)
}

func (v View) TileToGlobal(x, y int) image.Point {
	return image.Pt((x-y)*v.TileWidth/2,
		(x+y)*v.TileHeight/2)
}

func (v View) GlobalToTile(x, y int) image.Point {
	y += v.MouseOffset

	// TODO: fix the integer version!
	//twh := v.TileWidth / 2
	//thh := v.TileHeight / 2

	//i := (x/twh + y/thh) / 2
	//j := (y/thh - x/twh) / 2

	w, h := float64(v.TileWidth), float64(v.TileHeight)
	xx, yy := float64(x), float64(y)
	i := xx/w + yy/h
	j := yy/h - xx/w

	return image.Pt(int(i), int(j))
}

func (v View) ScreenToGlobal(x, y int) (int, int) {
	return v.X + int(float64(x)/v.Zoom),
		v.Y + int(float64(y)/v.Zoom)
}
