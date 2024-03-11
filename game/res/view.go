package res

import (
	"image"
)

// View resource, holding the game's view state.
type View struct {
	// Tile width and height, for convenient access.
	TileWidth, TileHeight int
	// Current global coordinates of the top-left corner of the screen, in pixels.
	X, Y int
	// Y offset used to determine the mouse position.
	MouseOffset int
	// Current zoom factor.
	Zoom float64
}

func NewView(tileWidth, tileHeight int) View {
	return View{
		TileWidth:   tileWidth,
		TileHeight:  tileHeight,
		Zoom:        1,
		MouseOffset: tileHeight,
	}
}

func (v *View) Center(cell image.Point, screenWidth, screenHeight int) {
	pos := v.TileToGlobal(cell.X, cell.Y)
	v.Zoom = 1
	v.X = pos.X - screenWidth/2
	v.Y = pos.Y - screenHeight/2
}

func (v *View) Offset() image.Point {
	return image.Pt(int(float64(v.X)*v.Zoom), int(float64(v.Y)*v.Zoom))
}

func (v *View) Bounds(w, h int) image.Rectangle {
	vw, vh := int(float64(w)/v.Zoom), int(float64(h)/v.Zoom)

	return image.Rect(
		v.X-v.TileWidth, v.Y-2*v.TileHeight,
		v.X+vw+v.TileWidth, v.Y+vh+5*v.TileHeight,
	)
}

func (v *View) TileToGlobal(x, y int) image.Point {
	return image.Pt((x-y)*v.TileWidth/2,
		(x+y)*v.TileHeight/2)
}

func (v *View) SubtileToGlobal(x, y float64) image.Point {
	return image.Pt(int((x-y)*float64(v.TileWidth)/2),
		int((x+y)*float64(v.TileHeight/2)))
}

func (v *View) GlobalToTile(x, y int) image.Point {
	y += v.MouseOffset

	w, h := float64(v.TileWidth), float64(v.TileHeight)
	xx, yy := float64(x), float64(y)
	i := xx/w + yy/h
	j := yy/h - xx/w

	return image.Pt(int(i), int(j))
}

func (v *View) ScreenToGlobal(x, y int) (int, int) {
	return v.X + int(float64(x)/v.Zoom),
		v.Y + int(float64(y)/v.Zoom)
}

func (v *View) MapBounds(screenWidth, screenHeight int) image.Rectangle {
	p := v.GlobalToTile(v.ScreenToGlobal(0, 0))
	xMin := p.X
	p = v.GlobalToTile(v.ScreenToGlobal(screenWidth+v.TileWidth, screenHeight+5*v.TileHeight))
	xMax := p.X
	p = v.GlobalToTile(v.ScreenToGlobal(screenWidth+v.TileWidth, 0))
	yMin := p.Y
	p = v.GlobalToTile(v.ScreenToGlobal(0, screenHeight+5*v.TileHeight))
	yMax := p.Y

	return image.Rect(xMin, yMin, xMax, yMax)
}

func (v *View) BoundsToGlobal(b *WorldBounds) image.Rectangle {
	p := v.TileToGlobal(b.Min.X, b.Min.Y)
	yMin := p.Y
	p = v.TileToGlobal(b.Max.X, b.Max.Y)
	yMax := p.Y
	p = v.TileToGlobal(b.Max.X, b.Min.Y)
	xMax := p.X
	p = v.TileToGlobal(b.Min.X, b.Max.Y)
	xMin := p.X

	return image.Rect(xMin, yMin, xMax, yMax)
}
