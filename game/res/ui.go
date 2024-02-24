package res

import (
	stdimage "image"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
	"golang.org/x/image/font"
)

type UI struct {
	UI             *ebitenui.UI
	ResourceLabels [resource.EndResources]*widget.Text
}

func (ui *UI) MouseInside(x, y int) bool {
	pt := stdimage.Pt(x, y)
	for _, w := range ui.UI.Container.Children() {
		if pt.In(w.GetWidget().Rect) {
			return true
		}
	}

	return false
}

func NewUI(selection *Selection, font font.Face, sprites *Sprites) UI {
	ui := UI{}
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	uiContainer := ui.createUI(sprites, selection, font)
	hudContainer := ui.createHUD(font)
	rootContainer.AddChild(uiContainer)
	rootContainer.AddChild(hudContainer)

	eui := ebitenui.UI{
		Container: rootContainer,
	}
	ui.UI = &eui

	return ui
}

func (ui *UI) createUI(sprites *Sprites, selection *Selection, font font.Face) *widget.Container {
	innerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(4)),
				widget.RowLayoutOpts.Spacing(4),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
			widget.WidgetOpts.MinSize(40, 200),
		),
	)

	for i := terr.Terrain(0); i < terr.EndTerrain; i++ {
		if !terr.Properties[i].CanBuild {
			continue
		}
		idx := sprites.GetTerrainIndex(i)
		img, _ := sprites.Get(idx)
		slice := image.NewNineSliceSimple(img, 0, 48)
		pressed := ebiten.NewImageFromImage(img)
		vector.DrawFilledRect(pressed, 0, 0,
			float32(img.Bounds().Dx()), float32(img.Bounds().Dy()),
			color.RGBA{0, 0, 0, 80}, false)
		slicePressed := image.NewNineSliceSimple(pressed, 0, 48)

		buttonImage := widget.ButtonImage{
			Idle:    slice,
			Hover:   slicePressed,
			Pressed: slicePressed,
		}
		button := createButton(selection, i, &buttonImage, font)
		innerContainer.AddChild(button)
	}

	return innerContainer
}

func (ui *UI) createHUD(font font.Face) *widget.Container {
	innerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(4)),
				widget.RowLayoutOpts.Spacing(4),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
			widget.WidgetOpts.MinSize(200, 30),
		),
	)

	for i := resource.Resource(0); i < resource.EndResources; i++ {
		label := widget.NewText(
			widget.TextOpts.Text("  "+resource.Properties[i].Short, font, color.White),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		)
		innerContainer.AddChild(label)
		counter := widget.NewText(
			widget.TextOpts.Text("0", font, color.White),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		)
		innerContainer.AddChild(counter)
		ui.ResourceLabels[i] = counter
	}

	return innerContainer
}

func createButton(selection *Selection, terrain terr.Terrain, buttonImage *widget.ButtonImage, font font.Face) *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImage),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			p := &terr.Properties[terrain]
			println("paint", p.Name)
			selection.Build = terrain
		}),
	)

	return button
}
