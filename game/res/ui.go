package res

import (
	stdimage "image"
	"image/color"
	"math/rand"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
	"github.com/mlange-42/tiny-world/game/util"
	"golang.org/x/image/font"
)

type UI struct {
	ui             *ebitenui.UI
	resourceLabels [resource.EndResources]*widget.Text
	terrainButtons [terr.EndTerrain]*widget.Button

	buttonImages           [terr.EndTerrain]widget.ButtonImage
	randomButtonsContainer *widget.Container
	randomButtons          map[int]*widget.Button

	selection *Selection
	font      font.Face
	idPool    util.IntPool[int]
}

func (ui *UI) UI() *ebitenui.UI {
	return ui.ui
}

func (ui *UI) SetResourceLabel(id resource.Resource, text string) {
	ui.resourceLabels[id].Label = text
}

func (ui *UI) SetButtonEnabled(id terr.Terrain, enabled bool) {
	button := ui.terrainButtons[id]
	if button == nil {
		return
	}
	button.GetWidget().Disabled = !enabled
}

func (ui *UI) MouseInside(x, y int) bool {
	pt := stdimage.Pt(x, y)
	for _, w := range ui.ui.Container.Children() {
		if pt.In(w.GetWidget().Rect) {
			return true
		}
	}

	return false
}

func NewUI(selection *Selection, font font.Face, sprites *Sprites, randomTerrains int, tileWidth int) UI {
	ui := UI{
		randomButtons: map[int]*widget.Button{},
		selection:     selection,
		font:          font,
		idPool:        util.NewIntPool[int](8),
	}

	ui.createImages(sprites, tileWidth)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	uiContainer := ui.createUI(sprites, randomTerrains)
	hudContainer := ui.createHUD(font)
	rootContainer.AddChild(uiContainer)
	rootContainer.AddChild(hudContainer)

	eui := ebitenui.UI{
		Container: rootContainer,
	}
	ui.ui = &eui

	return ui
}

func (ui *UI) CreateRandomButton() {
	t := terr.RandomTerrain[rand.Intn(len(terr.RandomTerrain))]
	button, id := ui.createButton(t)
	ui.randomButtonsContainer.AddChild(button)
	ui.randomButtons[id] = button
}

func (ui *UI) RemoveButton(id int) bool {
	if bt, ok := ui.randomButtons[id]; ok {
		ui.randomButtonsContainer.RemoveChild(bt)
		delete(ui.randomButtons, id)
		return true
	}
	return false
}

func (ui *UI) createUI(sprites *Sprites, randomTerrains int) *widget.Container {
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
			widget.WidgetOpts.MinSize(40, 10),
		),
	)

	for i := terr.Terrain(0); i < terr.EndTerrain; i++ {
		if !terr.Properties[i].CanBuy {
			continue
		}
		button, _ := ui.createButton(i)
		innerContainer.AddChild(button)
		ui.terrainButtons[i] = button
	}

	ui.randomButtonsContainer = widget.NewContainer(
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
			widget.WidgetOpts.MinSize(40, 10),
		),
	)
	for i := 0; i < randomTerrains; i++ {
		ui.CreateRandomButton()
	}

	innerContainer.AddChild(ui.randomButtonsContainer)

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
		ui.resourceLabels[i] = counter
	}

	return innerContainer
}

func (ui *UI) createImages(sprites *Sprites, tileWidth int) {
	for i := terr.Terrain(0); i < terr.EndTerrain; i++ {
		idx := sprites.GetTerrainIndex(i)
		img, _ := sprites.Get(idx)
		slice := image.NewNineSliceSimple(img, 0, tileWidth)
		pressed := ebiten.NewImageFromImage(img)
		vector.DrawFilledRect(pressed, 0, 0,
			float32(img.Bounds().Dx()), float32(img.Bounds().Dy()),
			color.RGBA{0, 0, 0, 80}, false)
		slicePressed := image.NewNineSliceSimple(pressed, 0, tileWidth)

		ui.buttonImages[i] = widget.ButtonImage{
			Idle:     slice,
			Hover:    slicePressed,
			Pressed:  slicePressed,
			Disabled: slicePressed,
		}
	}
}

func (ui *UI) createButton(terrain terr.Terrain) (*widget.Button, int) {
	id := ui.idPool.Get()
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(&ui.buttonImages[terrain]),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			p := &terr.Properties[terrain]
			println("paint", p.Name)
			ui.selection.SetBuild(terrain, id)
		}),
	)

	return button, id
}
