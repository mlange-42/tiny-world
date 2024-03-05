package res

import (
	"fmt"
	stdimage "image"
	"image/color"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/sprites"
	"github.com/mlange-42/tiny-world/game/terr"
	"github.com/mlange-42/tiny-world/game/util"
	"golang.org/x/image/font"
)

const tooltipSpecial = "\n(*) Can be placed over existing tiles."

type randomButton struct {
	Terrain      terr.Terrain
	RandomSprite uint16
	AllowRemove  bool
	Button       *widget.Button
}

type UI struct {
	RandomTerrains []RandomTerrain

	ui              *ebitenui.UI
	sprites         *Sprites
	saveEvent       *SaveEvent
	resourceLabels  []*widget.Text
	populationLabel *widget.Text
	timerLabel      *widget.Text
	terrainButtons  []*widget.Button

	buttonImages           []widget.ButtonImage
	buttonTooltip          []string
	randomButtonsContainer *widget.Container
	randomButtons          map[int]randomButton
	mouseBlockers          []*widget.Container

	markerSprite int

	selection *Selection
	font      font.Face
	idPool    util.IntPool[int]
}

type RandomTerrain struct {
	Terrain     terr.Terrain
	AllowRemove bool
}

func (ui *UI) UI() *ebitenui.UI {
	return ui.ui
}

func (ui *UI) SetResourceLabel(id resource.Resource, text string) {
	ui.resourceLabels[id].Label = text
}

func (ui *UI) SetPopulationLabel(text string) {
	ui.populationLabel.Label = text
}

func (ui *UI) SetTimerLabel(text string) {
	ui.timerLabel.Label = text
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
	for _, w := range ui.mouseBlockers {
		if pt.In(w.GetWidget().Rect) {
			return true
		}
	}

	return false
}

func NewUI(selection *Selection, font font.Face, sprts *Sprites, save *SaveEvent) UI {
	ui := UI{
		randomButtons: map[int]randomButton{},
		selection:     selection,
		font:          font,
		idPool:        util.NewIntPool[int](8),
		sprites:       sprts,
		saveEvent:     save,
		markerSprite:  sprts.GetIndex(sprites.SpecialCardMarker),
	}

	ui.prepareButtons()

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}),
		)),
	)

	hudContainer := ui.createHUD(font)
	rootContainer.AddChild(hudContainer)

	uiContainer := ui.createUI()
	rootContainer.AddChild(uiContainer)

	eui := ebitenui.UI{
		Container: rootContainer,
	}
	ui.ui = &eui

	return ui
}

func (ui *UI) createRandomButton(rules *Rules) {
	t := rules.RandomTerrains[rand.Intn(len(rules.RandomTerrains))]
	randSprite := uint16(rand.Int31n(math.MaxUint16))
	allowRemove := terr.Properties[t].TerrainBits.Contains(terr.IsTerrain) &&
		rand.Float64() < rules.SpecialCardProbability

	button, id := ui.createButton(t, allowRemove, randSprite)
	ui.randomButtonsContainer.AddChild(button)
	ui.randomButtons[id] = randomButton{t, randSprite, allowRemove, button}
}

func (ui *UI) ReplaceButton(stock *Stock, rules *Rules) bool {
	id := ui.selection.ButtonID
	if bt, ok := ui.randomButtons[id]; ok {
		ui.randomButtonsContainer.RemoveChild(bt.Button)
		delete(ui.randomButtons, id)

		ui.createRandomButton(rules)
		ui.updateRandomTerrains()

		ui.selection.Reset()
		for id2, bt2 := range ui.randomButtons {
			if bt2.Terrain == bt.Terrain && bt2.AllowRemove == bt.AllowRemove {
				ui.selection.SetBuild(bt2.Terrain, id2, bt2.RandomSprite, bt2.AllowRemove)
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

func (ui *UI) ReplaceAllButtons(rules *Rules) {
	ui.selection.Reset()
	ids := []int{}
	for id := range ui.randomButtons {
		ids = append(ids, id)
	}
	for _, id := range ids {
		bt := ui.randomButtons[id]
		ui.randomButtonsContainer.RemoveChild(bt.Button)
		delete(ui.randomButtons, id)

		ui.createRandomButton(rules)
		ui.updateRandomTerrains()
	}
}

func (ui *UI) createUI() *widget.Container {
	anchor := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionEnd,
			}),
		),
	)

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
				widget.GridLayoutOpts.Columns(3),
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

	ui.terrainButtons = make([]*widget.Button, len(terr.Properties))
	for i := range terr.Properties {
		canBuy := terr.Properties[i].TerrainBits.Contains(terr.CanBuy)
		if !canBuy && i != int(terr.Bulldoze) {
			continue
		}
		button, _ := ui.createButton(terr.Terrain(i), false)
		ui.terrainButtons[i] = button
		buildButtonsContainer.AddChild(button)
	}
	innerContainer.AddChild(buildButtonsContainer)

	ui.randomButtonsContainer = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255})),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(3),
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

	anchor.AddChild(innerContainer)
	ui.mouseBlockers = append(ui.mouseBlockers, innerContainer)

	return anchor
}

func (ui *UI) CreateRandomButtons(randomTerrains int) {
	if len(ui.RandomTerrains) == 0 {
		for i := 0; i < randomTerrains; i++ {
			randSprite := uint16(rand.Int31n(math.MaxUint16))
			button, id := ui.createButton(terr.Default, false, randSprite)
			ui.randomButtonsContainer.AddChild(button)
			ui.randomButtons[id] = randomButton{terr.Default, randSprite, false, button}
		}
		ui.updateRandomTerrains()
	} else {
		for _, t := range ui.RandomTerrains {
			randSprite := uint16(rand.Int31n(math.MaxUint16))
			button, id := ui.createButton(t.Terrain, t.AllowRemove, randSprite)
			ui.randomButtonsContainer.AddChild(button)
			ui.randomButtons[id] = randomButton{t.Terrain, randSprite, t.AllowRemove, button}
		}
	}
}

func (ui *UI) updateRandomTerrains() {
	ui.RandomTerrains = ui.RandomTerrains[:0]
	for _, bt := range ui.randomButtons {
		ui.RandomTerrains = append(ui.RandomTerrains, RandomTerrain{bt.Terrain, bt.AllowRemove})
	}
}

func (ui *UI) createHUD(font font.Face) *widget.Container {
	anchor := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
		),
	)

	topBar := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
			widget.WidgetOpts.MinSize(30, 30),
		),
	)

	menuContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(4)),
				widget.RowLayoutOpts.Spacing(4),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionStart,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
		),
	)

	saveButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				Stretch:  false,
			}),
		),
		widget.ButtonOpts.Image(simpleButtonImage()),
		widget.ButtonOpts.Text("Save", font, &widget.ButtonTextColor{
			Idle: color.NRGBA{255, 255, 255, 255},
		}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.saveEvent.ShouldSave = true
		}),
	)
	menuContainer.AddChild(saveButton)

	innerAnchor := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
		),
	)

	infoContainer := widget.NewContainer(
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
			widget.WidgetOpts.MinSize(30, 30),
		),
	)

	ui.resourceLabels = make([]*widget.Text, len(resource.Properties))
	for i := range resource.Properties {
		label := widget.NewText(
			widget.TextOpts.Text("  "+resource.Properties[i].Short, font, color.White),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		)
		infoContainer.AddChild(label)
		counter := widget.NewText(
			widget.TextOpts.Text("0", font, color.White),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		)
		infoContainer.AddChild(counter)
		ui.resourceLabels[i] = counter
	}
	label := widget.NewText(
		widget.TextOpts.Text("  Pop", font, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)
	infoContainer.AddChild(label)
	counter := widget.NewText(
		widget.TextOpts.Text("0", font, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)
	infoContainer.AddChild(counter)
	ui.populationLabel = counter

	labelTimer := widget.NewText(
		widget.TextOpts.Text("   ", font, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)
	infoContainer.AddChild(labelTimer)
	counterTimer := widget.NewText(
		widget.TextOpts.Text("0", font, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)
	infoContainer.AddChild(counterTimer)
	ui.timerLabel = counterTimer

	topBar.AddChild(menuContainer)
	topBar.AddChild(innerAnchor)
	innerAnchor.AddChild(infoContainer)
	anchor.AddChild(topBar)

	ui.mouseBlockers = append(ui.mouseBlockers, infoContainer, menuContainer)

	return anchor
}

func (ui *UI) prepareButtons() {
	ui.buttonImages = make([]widget.ButtonImage, len(terr.Properties))
	ui.buttonTooltip = make([]string, len(terr.Properties))

	for i := range terr.Properties {
		props := &terr.Properties[i]

		ui.buttonImages[i] = ui.createButtonImage(terr.Terrain(i), 0, false)

		costs := ""
		if len(props.BuildCost) > 0 {
			costs = "Cost: "
			for i, cost := range props.BuildCost {
				if i > 0 {
					costs += ", "
				}
				costs += fmt.Sprintf("%d %s", cost.Amount, resource.Properties[cost.Resource].Short)
			}
			costs += "\n"
		}
		maxProd := ""
		if props.Production.MaxProduction > 0 {
			maxProd = fmt.Sprintf(" (max %d)", props.Production.MaxProduction)
		}
		radius := ""
		if props.BuildRadius > 0 {
			radius = fmt.Sprintf("Radius: %d\n", props.BuildRadius)
		}
		pop := ""
		if props.Population > 0 {
			pop = fmt.Sprintf("Population: %d\n", props.Population)
		}

		requires := ""
		requiresTemp := ui.resourcesToString(props.Consumption)
		if len(requiresTemp) > 0 {
			requires = fmt.Sprintf("Requires: %s /min\n", requiresTemp)
		}

		storage := ""
		if props.TerrainBits.Contains(terr.IsWarehouse) {
			storage = fmt.Sprintf("Stores: %s\n", ui.resourcesToString(props.Storage))
		}

		ui.buttonTooltip[i] = fmt.Sprintf("%s\n%s%s%s%s%s%s%s.",
			strings.ToUpper(props.Name), costs, requires, pop, radius, storage, props.Description, maxProd)
	}
}

func (ui *UI) resourcesToString(res []uint8) string {
	out := ""
	cnt := 0
	for i, st := range res {
		if st == 0 {
			continue
		}
		if cnt > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%d %s", st, resource.Properties[i].Short)
		cnt++
	}
	return out
}

func (ui *UI) createButtonImage(t terr.Terrain, randSprite uint16, allowRemove bool) widget.ButtonImage {
	props := &terr.Properties[t]

	tileWidth := ui.sprites.TileWidth
	img := ebiten.NewImage(tileWidth, tileWidth)
	idx := ui.sprites.GetTerrainIndex(terr.Terrain(t))

	height := 0

	if props.TerrainBelow != terr.Air {
		idx2 := ui.sprites.GetTerrainIndex(props.TerrainBelow)
		info2 := ui.sprites.GetInfo(idx2)

		sp2 := ui.sprites.Get(idx2)
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(tileWidth-sp2.Bounds().Dy()))
		img.DrawImage(sp2, &op)

		height = info2.Height
	}

	sp1 := ui.sprites.GetRand(idx, 0, int(randSprite))
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(tileWidth-sp1.Bounds().Dy()-height))
	img.DrawImage(sp1, &op)

	if allowRemove {
		marker := ui.sprites.Get(ui.markerSprite)
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(tileWidth-marker.Bounds().Dy()))
		img.DrawImage(marker, &op)
	}

	slice := image.NewNineSliceSimple(img, 0, tileWidth)

	pressed := ebiten.NewImageFromImage(img)
	vector.DrawFilledRect(pressed, 0, 0,
		float32(img.Bounds().Dx()), float32(img.Bounds().Dy()),
		color.RGBA{0, 0, 0, 80}, false)
	slicePressed := image.NewNineSliceSimple(pressed, 0, tileWidth)

	disabled := ebiten.NewImageFromImage(img)
	vector.StrokeLine(disabled, 0, 0, float32(img.Bounds().Dx()), float32(img.Bounds().Dy()), 6, color.RGBA{120, 0, 0, 160}, false)
	sliceDisabled := image.NewNineSliceSimple(disabled, 0, tileWidth)

	return widget.ButtonImage{
		Idle:     slice,
		Hover:    slicePressed,
		Pressed:  slicePressed,
		Disabled: sliceDisabled,
	}
}

func (ui *UI) createButton(terrain terr.Terrain, allowRemove bool, randSprite ...uint16) (*widget.Button, int) {
	id := ui.idPool.Get()

	tooltipContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{Top: 6, Bottom: 6, Left: 12, Right: 12}),
		)),
		widget.ContainerOpts.AutoDisableChildren(),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{20, 20, 20, 255})),
	)

	text := ui.buttonTooltip[terrain]
	if allowRemove {
		text += tooltipSpecial
	}
	label := widget.NewText(
		widget.TextOpts.Text(text, ui.font, color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
		widget.TextOpts.MaxWidth(250),
	)
	tooltipContainer.AddChild(label)

	bImage := ui.buttonImages[terrain]
	var randSpriteVal uint16 = 0
	if len(randSprite) > 0 {
		bImage = ui.createButtonImage(terrain, randSprite[0], allowRemove)
		randSpriteVal = randSprite[0]
	}

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
		widget.ButtonOpts.Image(&bImage),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.selection.SetBuild(terrain, id, randSpriteVal, allowRemove)
		}),
	)

	return button, id
}

func simpleButtonImage() *widget.ButtonImage {
	idle := image.NewNineSliceColor(color.NRGBA{60, 60, 60, 255})

	hover := image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255})

	pressed := image.NewNineSliceColor(color.NRGBA{20, 20, 20, 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}
