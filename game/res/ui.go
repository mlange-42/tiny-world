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
	return stdimage.Pt(x, y).In(ui.UI.Container.Children()[0].GetWidget().Rect)
}

func NewUI(selection *Selection, font font.Face, sprites *Sprites) UI {
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

	buttons := []widget.RadioGroupElement{}
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
		buttons = append(buttons, button)
	}
	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(buttons...),

		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {

		}),
	)

	rootContainer.AddChild(innerContainer)

	eui := ebitenui.UI{
		Container: rootContainer,
	}
	ui.UI = &eui

	return ui
}

func createButton(selection *Selection, terrain terr.Terrain, buttonImage *widget.ButtonImage, face font.Face) *widget.Button {
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
