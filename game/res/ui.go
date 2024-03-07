package res

import (
	"fmt"
	stdimage "image"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
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
	Index        int
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

	animMapper generic.Map1[comp.CardAnimation]

	buttonImages           []widget.ButtonImage
	buttonTooltip          []string
	randomButtonsContainer *widget.Container
	randomContainers       []*widget.Container
	randomButtons          map[int]randomButton
	mouseBlockers          []*widget.Container

	specialCardSprite    int
	buttonIdleSprite     int
	buttonHoverSprite    int
	buttonPressedSprite  int
	buttonDisabledSprite int

	background        *image.NineSlice
	backgroundHover   *image.NineSlice
	backgroundPressed *image.NineSlice

	selection *Selection
	font      font.Face
	idPool    util.IntPool[int]

	buttonSize stdimage.Point
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

func NewUI(world *ecs.World, selection *Selection, font font.Face, sprts *Sprites, save *SaveEvent) UI {
	ui := UI{
		randomButtons: map[int]randomButton{},
		selection:     selection,
		font:          font,
		idPool:        util.NewIntPool[int](8),
		sprites:       sprts,
		saveEvent:     save,

		specialCardSprite:    sprts.GetIndex(sprites.SpecialCardMarker),
		buttonIdleSprite:     sprts.GetIndex(sprites.Button),
		buttonHoverSprite:    sprts.GetIndex(sprites.ButtonHover),
		buttonPressedSprite:  sprts.GetIndex(sprites.ButtonPressed),
		buttonDisabledSprite: sprts.GetIndex(sprites.ButtonDisabled),

		animMapper: generic.NewMap1[comp.CardAnimation](world),
	}
	sp := ui.sprites.Get(ui.buttonIdleSprite)
	ui.buttonSize = sp.Bounds().Max

	sp = ui.sprites.Get(ui.sprites.GetIndex(sprites.UiPanel))
	w := sp.Bounds().Dx()
	ui.background = image.NewNineSliceSimple(sp, w/4, w/2)

	sp = ui.sprites.Get(ui.sprites.GetIndex(sprites.UiPanelHover))
	w = sp.Bounds().Dx()
	ui.backgroundHover = image.NewNineSliceSimple(sp, w/4, w/2)

	sp = ui.sprites.Get(ui.sprites.GetIndex(sprites.UiPanelPressed))
	w = sp.Bounds().Dx()
	ui.backgroundPressed = image.NewNineSliceSimple(sp, w/4, w/2)

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

func (ui *UI) createRandomButton(rules *Rules, index int) {
	t := rules.RandomTerrains[rand.Intn(len(rules.RandomTerrains))]
	randSprite := uint16(rand.Int31n(math.MaxUint16))
	allowRemove := terr.Properties[t].TerrainBits.Contains(terr.IsTerrain) &&
		rand.Float64() < rules.SpecialCardProbability

	button, id := ui.createButton(t, allowRemove, randSprite)
	ui.randomContainers[index].AddChild(button)
	ui.randomButtons[id] = randomButton{t, randSprite, allowRemove, button, index}
}

func (ui *UI) ReplaceButton(stock *Stock, rules *Rules, tick int64, target stdimage.Point) bool {
	id := ui.selection.ButtonID
	if bt, ok := ui.randomButtons[id]; ok {
		ui.animMapper.NewWith(&comp.CardAnimation{
			Point:      bt.Button.GetWidget().Rect.Min,
			Target:     target,
			Terrain:    bt.Terrain,
			RandSprite: bt.RandomSprite,
			StartTick:  tick,
		})

		ui.randomContainers[bt.Index].RemoveChild(bt.Button)
		delete(ui.randomButtons, id)

		ui.createRandomButton(rules, bt.Index)
		ui.updateRandomTerrains()

		ui.ClearSelection()
		// Try at the same index first
		for id2, bt2 := range ui.randomButtons {
			if bt2.Index == bt.Index && bt2.Terrain == bt.Terrain && bt2.AllowRemove == bt.AllowRemove {
				ui.selectTerrain(bt2.Button, bt2.Terrain, id2, bt2.RandomSprite, bt2.AllowRemove)
				return true
			}
		}
		// Try to find any
		for id2, bt2 := range ui.randomButtons {
			if bt2.Terrain == bt.Terrain && bt2.AllowRemove == bt.AllowRemove {
				ui.selectTerrain(bt2.Button, bt2.Terrain, id2, bt2.RandomSprite, bt2.AllowRemove)
				return true
			}
		}
		return true
	}
	if !stock.CanPay(terr.Properties[ui.selection.BuildType].BuildCost) {
		ui.ClearSelection()
	}
	return false
}

func (ui *UI) ReplaceAllButtons(rules *Rules) {
	ui.ClearSelection()
	ids := []int{}
	for id := range ui.randomButtons {
		ids = append(ids, id)
	}
	for _, id := range ids {
		bt := ui.randomButtons[id]
		ui.randomContainers[bt.Index].RemoveChild(bt.Button)
		delete(ui.randomButtons, id)

		ui.createRandomButton(rules, bt.Index)
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
		widget.ContainerOpts.BackgroundImage(ui.background),
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
	ui.randomContainers = make([]*widget.Container, randomTerrains)
	if len(ui.RandomTerrains) == 0 {
		for i := 0; i < randomTerrains; i++ {
			randSprite := uint16(rand.Int31n(math.MaxUint16))
			button, id := ui.createButton(terr.Default, false, randSprite)

			container := widget.NewContainer(widget.ContainerOpts.Layout(
				widget.NewGridLayout(widget.GridLayoutOpts.Columns(1))))
			container.AddChild(button)
			ui.randomContainers[i] = container
			ui.randomButtonsContainer.AddChild(container)
			ui.randomButtons[id] = randomButton{terr.Default, randSprite, false, button, i}
		}
		ui.updateRandomTerrains()
	} else {
		for i, t := range ui.RandomTerrains {
			randSprite := uint16(rand.Int31n(math.MaxUint16))
			button, id := ui.createButton(t.Terrain, t.AllowRemove, randSprite)

			container := widget.NewContainer(widget.ContainerOpts.Layout(
				widget.NewGridLayout(widget.GridLayoutOpts.Columns(1))))
			container.AddChild(button)
			ui.randomContainers[i] = container
			ui.randomButtonsContainer.AddChild(container)
			ui.randomButtons[id] = randomButton{t.Terrain, randSprite, t.AllowRemove, button, i}
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
		widget.ContainerOpts.BackgroundImage(ui.background),
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
		widget.ButtonOpts.Image(ui.simpleButtonImage()),
		widget.ButtonOpts.Text("Save", font, &widget.ButtonTextColor{
			Idle: ui.sprites.TextColor,
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
		widget.ContainerOpts.BackgroundImage(ui.background),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 4, Bottom: 4, Left: 12, Right: 12}),
				widget.RowLayoutOpts.Spacing(12),
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
		cont, lab := ui.createLabel(resource.Properties[i].Short,
			fmt.Sprintf("%s:\n   +production -consumption\n   (stock / max)", util.Capitalize(resource.Properties[i].Name)), 130)
		infoContainer.AddChild(cont)
		ui.resourceLabels[i] = lab
	}
	cont, lab := ui.createLabel("Pop", "Population: current/max", 50)
	infoContainer.AddChild(cont)
	ui.populationLabel = lab

	cont, lab = ui.createLabel("", "Total game time", 40)
	infoContainer.AddChild(cont)
	ui.timerLabel = lab

	topBar.AddChild(menuContainer)
	topBar.AddChild(innerAnchor)
	innerAnchor.AddChild(infoContainer)
	anchor.AddChild(topBar)

	ui.mouseBlockers = append(ui.mouseBlockers, infoContainer, menuContainer)

	return anchor
}

func (ui *UI) createLabel(text, tooltip string, width int) (*widget.Container, *widget.Text) {
	tooltipContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{Top: 6, Bottom: 6, Left: 12, Right: 12}),
		)),
		widget.ContainerOpts.AutoDisableChildren(),
		widget.ContainerOpts.BackgroundImage(ui.background),
	)
	label := widget.NewText(
		widget.TextOpts.Text(tooltip, ui.font, ui.sprites.TextColor),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
		widget.TextOpts.MaxWidth(360),
	)
	tooltipContainer.AddChild(label)

	cont := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(4),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.ToolTip(widget.NewToolTip(
				widget.ToolTipOpts.Content(tooltipContainer),
				widget.ToolTipOpts.Offset(stdimage.Point{-5, 5}),
				widget.ToolTipOpts.Position(widget.TOOLTIP_POS_WIDGET),
				widget.ToolTipOpts.Delay(time.Millisecond*300),
			)),
		),
	)

	if len(text) > 0 {
		label := widget.NewText(
			widget.TextOpts.Text(text, ui.font, ui.sprites.TextColor),
			widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
		)
		cont.AddChild(label)
	}
	counter := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(width, 0),
		),
		widget.TextOpts.Text("", ui.font, ui.sprites.TextColor),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	)
	cont.AddChild(counter)

	return cont, counter
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

	xOff := (ui.buttonSize.X - ui.sprites.TileWidth) / 2
	yOff := (ui.buttonSize.Y - ui.sprites.TileWidth) / 2
	img := ebiten.NewImage(ui.buttonSize.X, ui.buttonSize.Y)
	idx := ui.sprites.GetTerrainIndex(terr.Terrain(t))

	height := 0

	if props.TerrainBelow != terr.Air {
		idx2 := ui.sprites.GetTerrainIndex(props.TerrainBelow)
		info2 := ui.sprites.GetInfo(idx2)

		sp2 := ui.sprites.Get(idx2)
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(xOff), float64(ui.buttonSize.X-sp2.Bounds().Dy()-yOff))
		img.DrawImage(sp2, &op)

		height = info2.Height
	}

	sp1 := ui.sprites.GetRand(idx, 0, int(randSprite))
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xOff), float64(ui.buttonSize.X-sp1.Bounds().Dy()-height-yOff))
	img.DrawImage(sp1, &op)

	if allowRemove {
		marker := ui.sprites.Get(ui.specialCardSprite)
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(xOff), float64(ui.buttonSize.X-marker.Bounds().Dy()-yOff))
		img.DrawImage(marker, &op)
	}

	op = ebiten.DrawImageOptions{}
	idle := ebiten.NewImageFromImage(img)
	button := ui.sprites.Get(ui.buttonIdleSprite)
	idle.DrawImage(button, &op)
	sliceIdle := image.NewNineSlice(idle, [3]int{ui.buttonSize.X, 0, 0}, [3]int{ui.buttonSize.Y, 0, 0})

	hover := ebiten.NewImageFromImage(img)
	button = ui.sprites.Get(ui.buttonPressedSprite)
	hover.DrawImage(button, &op)
	sliceHover := image.NewNineSlice(hover, [3]int{ui.buttonSize.X, 0, 0}, [3]int{ui.buttonSize.Y, 0, 0})

	pressed := ebiten.NewImageFromImage(img)
	button = ui.sprites.Get(ui.buttonPressedSprite)
	pressed.DrawImage(button, &op)
	slicePressed := image.NewNineSlice(pressed, [3]int{ui.buttonSize.X, 0, 0}, [3]int{ui.buttonSize.Y, 0, 0})

	disabled := ebiten.NewImageFromImage(img)
	button = ui.sprites.Get(ui.buttonDisabledSprite)
	disabled.DrawImage(button, &op)
	sliceDisabled := image.NewNineSlice(disabled, [3]int{ui.buttonSize.X, 0, 0}, [3]int{ui.buttonSize.Y, 0, 0})

	return widget.ButtonImage{
		Idle:     sliceIdle,
		Hover:    sliceHover,
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
		widget.ContainerOpts.BackgroundImage(ui.background),
	)

	text := ui.buttonTooltip[terrain]
	if allowRemove {
		text += tooltipSpecial
	}
	label := widget.NewText(
		widget.TextOpts.Text(text, ui.font, ui.sprites.TextColor),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
		widget.TextOpts.MaxWidth(360),
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
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxWidth:  ui.sprites.TileWidth,
				MaxHeight: ui.sprites.TileWidth,
			}),
			widget.WidgetOpts.MinSize(ui.sprites.TileWidth, ui.sprites.TileWidth),
			widget.WidgetOpts.ToolTip(widget.NewToolTip(
				widget.ToolTipOpts.Content(tooltipContainer),
				widget.ToolTipOpts.Offset(stdimage.Point{-5, 5}),
				widget.ToolTipOpts.Position(widget.TOOLTIP_POS_WIDGET),
				widget.ToolTipOpts.Delay(time.Millisecond*300),
			)),
		),
		widget.ButtonOpts.ToggleMode(),
		widget.ButtonOpts.Image(&bImage),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.selectTerrain(args.Button, terrain, id, randSpriteVal, allowRemove)
		}),
	)

	return button, id
}

func (ui *UI) selectTerrain(button *widget.Button, terrain terr.Terrain, id int, randSprite uint16, allowRemove bool) {
	for _, bt := range ui.terrainButtons {
		if bt != nil {
			bt.SetState(widget.WidgetUnchecked)
		}
	}
	for _, bt := range ui.randomButtons {
		bt.Button.SetState(widget.WidgetUnchecked)
	}

	ui.selection.SetBuild(terrain, id, randSprite, allowRemove)
	button.SetState(widget.WidgetChecked)
}

func (ui *UI) ClearSelection() {
	for _, bt := range ui.terrainButtons {
		if bt != nil {
			bt.SetState(widget.WidgetUnchecked)
		}
	}
	for _, bt := range ui.randomButtons {
		bt.Button.SetState(widget.WidgetUnchecked)
	}
	ui.selection.Reset()
}

func (ui *UI) simpleButtonImage() *widget.ButtonImage {
	return &widget.ButtonImage{
		Idle:    ui.background,
		Hover:   ui.backgroundHover,
		Pressed: ui.backgroundPressed,
	}
}
