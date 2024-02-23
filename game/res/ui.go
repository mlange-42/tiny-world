package res

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/mlange-42/tiny-world/game/resource"
	"golang.org/x/image/font"
)

type UserInterface struct {
	UI             *ebitenui.UI
	ResourceLabels [resource.EndResources]*widget.Text
}

func NewUserInterface(font font.Face) UserInterface {
	ui := UserInterface{}
	rootContainer := widget.NewContainer(
		//widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	innerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 0, 0, 255})),
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
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
			widget.TextOpts.Text("0", font, color.White),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
			widget.TextOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		)
		innerContainer.AddChild(label)
		ui.ResourceLabels[i] = label
	}

	rootContainer.AddChild(innerContainer)

	eui := ebitenui.UI{
		Container: rootContainer,
	}
	ui.UI = &eui

	return ui
}
