package res

import (
	"fmt"
	stdimage "image"
	"image/color"
	"math/rand"
	"strings"
	"time"

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

type randomButton struct {
	Terrain terr.Terrain
	Button  *widget.Button
}

type UI struct {
	RandomTerrains []terr.Terrain

	ui             *ebitenui.UI
	resourceLabels [resource.EndResources]*widget.Text
	terrainButtons [terr.EndTerrain]*widget.Button

	buttonImages           [terr.EndTerrain]widget.ButtonImage
	buttonTooltip          [terr.EndTerrain]string
	randomButtonsContainer *widget.Container
	randomButtons          map[int]randomButton

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

func NewUI(selection *Selection, font font.Face, sprites *Sprites, tileWidth int) UI {
	ui := UI{
		randomButtons: map[int]randomButton{},
		selection:     selection,
		font:          font,
		idPool:        util.NewIntPool[int](8),
	}

	ui.prepareButtons(sprites, tileWidth)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	uiContainer := ui.createUI()
	hudContainer := ui.createHUD(font)
	rootContainer.AddChild(uiContainer)
	rootContainer.AddChild(hudContainer)

	eui := ebitenui.UI{
		Container: rootContainer,
	}
	ui.ui = &eui

	return ui
}

func (ui *UI) createRandomButton() {
	t := terr.RandomTerrain[rand.Intn(len(terr.RandomTerrain))]
	button, id := ui.createButton(t)
	ui.randomButtonsContainer.AddChild(button)
	ui.randomButtons[id] = randomButton{t, button}
}

func (ui *UI) ReplaceButton(stock *Stock) bool {
	id := ui.selection.ButtonID
	if bt, ok := ui.randomButtons[id]; ok {
		ui.randomButtonsContainer.RemoveChild(bt.Button)
		delete(ui.randomButtons, id)
		ui.createRandomButton()
		ui.updateRandomTerrains()

		ui.selection.Reset()
		for id2, bt2 := range ui.randomButtons {
			if bt2.Terrain == bt.Terrain {
				ui.selection.SetBuild(bt2.Terrain, id2)
				break
			}
		}
		return true
	}
	if !stock.CanPay(terr.Properties[ui.selection.BuildType].BuildCost) {
		ui.selection.Reset()
	}
	return false
}

func (ui *UI) createUI() *widget.Container {
	innerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(4)),
				widget.RowLayoutOpts.Spacing(24),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
			widget.WidgetOpts.MinSize(40, 10),
		),
	)

	buildButtonsContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255})),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(4)),
				widget.GridLayoutOpts.Spacing(4, 4),
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
		ui.terrainButtons[i] = button
		buildButtonsContainer.AddChild(button)
	}
	innerContainer.AddChild(buildButtonsContainer)

	ui.randomButtonsContainer = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255})),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(4)),
				widget.GridLayoutOpts.Spacing(4, 4),
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
	innerContainer.AddChild(ui.randomButtonsContainer)

	return innerContainer
}

func (ui *UI) CreateRandomButtons(randomTerrains int) {
	if len(ui.RandomTerrains) == 0 {
		for i := 0; i < randomTerrains; i++ {
			button, id := ui.createButton(terr.Grass)
			ui.randomButtonsContainer.AddChild(button)
			ui.randomButtons[id] = randomButton{terr.Grass, button}
		}
		ui.updateRandomTerrains()
	} else {
		for _, t := range ui.RandomTerrains {
			button, id := ui.createButton(t)
			ui.randomButtonsContainer.AddChild(button)
			ui.randomButtons[id] = randomButton{t, button}
		}
	}
}

func (ui *UI) updateRandomTerrains() {
	ui.RandomTerrains = ui.RandomTerrains[:0]
	for _, bt := range ui.randomButtons {
		ui.RandomTerrains = append(ui.RandomTerrains, bt.Terrain)
	}
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

func (ui *UI) prepareButtons(sprites *Sprites, tileWidth int) {
	for i := terr.Terrain(0); i < terr.EndTerrain; i++ {
		idx := sprites.GetTerrainIndex(i)
		img := sprites.Get(idx)
		slice := image.NewNineSliceSimple(img, 0, tileWidth)

		pressed := ebiten.NewImageFromImage(img)
		vector.DrawFilledRect(pressed, 0, 0,
			float32(img.Bounds().Dx()), float32(img.Bounds().Dy()),
			color.RGBA{0, 0, 0, 80}, false)
		slicePressed := image.NewNineSliceSimple(pressed, 0, tileWidth)

		disabled := ebiten.NewImageFromImage(img)
		vector.StrokeLine(disabled, 0, 0, float32(img.Bounds().Dx()), float32(img.Bounds().Dy()), 6, color.RGBA{120, 0, 0, 160}, false)
		sliceDisabled := image.NewNineSliceSimple(disabled, 0, tileWidth)

		ui.buttonImages[i] = widget.ButtonImage{
			Idle:     slice,
			Hover:    slicePressed,
			Pressed:  slicePressed,
			Disabled: sliceDisabled,
		}

		props := &terr.Properties[i]
		costs := ""
		if len(props.BuildCost) > 0 {
			costs = "Cost: "
			for _, cost := range props.BuildCost {
				costs += fmt.Sprintf("%d %s, ", cost.Amount, resource.Properties[cost.Type].Short)
			}
			costs += "\n"
		}
		requires := ""
		if props.Production.ConsumesFood > 0 {
			requires = fmt.Sprintf("Requires: %d F/min\n", props.Production.ConsumesFood)
		}
		maxProd := ""
		if props.Production.MaxProduction > 0 {
			maxProd = fmt.Sprintf(" (max %d)", props.Production.MaxProduction)
		}
		ui.buttonTooltip[i] = fmt.Sprintf("%s\n%s%s%s%s.", strings.ToUpper(props.Name), costs, requires, terr.Descriptions[i], maxProd)
	}
}

func (ui *UI) createButton(terrain terr.Terrain) (*widget.Button, int) {
	id := ui.idPool.Get()

	tooltipContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{Top: 6, Bottom: 6, Left: 12, Right: 12}),
		)),
		widget.ContainerOpts.AutoDisableChildren(),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{20, 20, 20, 255})),
	)
	label := widget.NewText(
		widget.TextOpts.Text(ui.buttonTooltip[terrain], ui.font, color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
		widget.TextOpts.MaxWidth(250),
	)
	tooltipContainer.AddChild(label)

	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),

			widget.WidgetOpts.ToolTip(widget.NewToolTip(
				widget.ToolTipOpts.Content(tooltipContainer),
				widget.ToolTipOpts.Offset(stdimage.Point{-5, 5}),
				widget.ToolTipOpts.Position(widget.TOOLTIP_POS_WIDGET),
				widget.ToolTipOpts.Delay(time.Millisecond*300),
			)),
		),
		widget.ButtonOpts.Image(&ui.buttonImages[terrain]),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.selection.SetBuild(terrain, id)
		}),
	)

	return button, id
}
