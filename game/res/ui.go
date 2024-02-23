package res

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
	"golang.org/x/image/font"
)

type UI struct {
	UI             *ebitenui.UI
	ResourceLabels [resource.EndResources]*widget.Text
}

func NewUI(font font.Face, sprites *Sprites) UI {
	ui := UI{}
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

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

		buttonImage := widget.ButtonImage{
			Idle:    slice,
			Hover:   slice,
			Pressed: slice,
		}
		button := createButton(&buttonImage, font)
		innerContainer.AddChild(button)
	}

	rootContainer.AddChild(innerContainer)

	eui := ebitenui.UI{
		Container: rootContainer,
	}
	ui.UI = &eui

	return ui
}

func createButton(buttonImage *widget.ButtonImage, face font.Face) *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImage),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			fmt.Println("button clicked")
		}),
	)

	return button
}
