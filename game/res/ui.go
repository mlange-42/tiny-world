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
	"github.com/mlange-42/tiny-world/game/terr"
	"github.com/mlange-42/tiny-world/game/util"
	"golang.org/x/image/font"
)

type randomButton struct {
	Terrain      terr.Terrain
	RandomSprite uint16
	Button       *widget.Button
}

type UI struct {
	RandomTerrains []terr.Terrain

	ui             *ebitenui.UI
	sprites        *Sprites
	saveEvent      *SaveEvent
	resourceLabels []*widget.Text
	terrainButtons []*widget.Button

	buttonImages           []widget.ButtonImage
	buttonTooltip          []string
	randomButtonsContainer *widget.Container
	randomButtons          map[int]randomButton
	mouseBlockers          []*widget.Container

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
	for _, w := range ui.mouseBlockers {
		if pt.In(w.GetWidget().Rect) {
			return true
		}
	}

	return false
}

func NewUI(selection *Selection, font font.Face, sprites *Sprites, save *SaveEvent) UI {
	ui := UI{
		randomButtons: map[int]randomButton{},
		selection:     selection,
		font:          font,
		idPool:        util.NewIntPool[int](8),
		sprites:       sprites,
		saveEvent:     save,
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
	button, id := ui.createButton(t, randSprite)
	ui.randomButtonsContainer.AddChild(button)
	ui.randomButtons[id] = randomButton{t, randSprite, button}
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
			if bt2.Terrain == bt.Terrain {
				ui.selection.SetBuild(bt2.Terrain, id2, bt2.RandomSprite)
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

	ui.terrainButtons = make([]*widget.Button, len(terr.Properties))
	for i := range terr.Properties {
		if !terr.Properties[i].TerrainBits.Contains(terr.CanBuy) {
			continue
		}
		button, _ := ui.createButton(terr.Terrain(i))
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

	anchor.AddChild(innerContainer)
	ui.mouseBlockers = append(ui.mouseBlockers, innerContainer)

	return anchor
}

func (ui *UI) CreateRandomButtons(randomTerrains int) {
	if len(ui.RandomTerrains) == 0 {
		for i := 0; i < randomTerrains; i++ {
			randSprite := uint16(rand.Int31n(math.MaxUint16))
			button, id := ui.createButton(terr.Default, randSprite)
			ui.randomButtonsContainer.AddChild(button)
			ui.randomButtons[id] = randomButton{terr.Default, randSprite, button}
		}
		ui.updateRandomTerrains()
	} else {
		for _, t := range ui.RandomTerrains {
			randSprite := uint16(rand.Int31n(math.MaxUint16))
			button, id := ui.createButton(t, randSprite)
			ui.randomButtonsContainer.AddChild(button)
			ui.randomButtons[id] = randomButton{t, randSprite, button}
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

		ui.buttonImages[i] = ui.createButtonImage(terr.Terrain(i), 0)

		costs := ""
		if len(props.BuildCost) > 0 {
			costs = "Cost: "
			for _, cost := range props.BuildCost {
				costs += fmt.Sprintf("%d %s, ", cost.Amount, resource.Properties[cost.Resource].Short)
			}
			costs += "\n"
		}
		requires := ""
		if props.Consumption.Amount > 0 {
			requires = fmt.Sprintf("Requires: %d F/min\n", props.Consumption.Amount)
		}
		maxProd := ""
		if props.Production.MaxProduction > 0 {
			maxProd = fmt.Sprintf(" (max %d)", props.Production.MaxProduction)
		}
		ui.buttonTooltip[i] = fmt.Sprintf("%s\n%s%s%s%s.", strings.ToUpper(props.Name), costs, requires, props.Description, maxProd)
	}
}

func (ui *UI) createButtonImage(t terr.Terrain, randSprite uint16) widget.ButtonImage {
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

func (ui *UI) createButton(terrain terr.Terrain, randSprite ...uint16) (*widget.Button, int) {
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

	bImage := ui.buttonImages[terrain]
	var randSpriteVal uint16 = 0
	if len(randSprite) > 0 {
		bImage = ui.createButtonImage(terrain, randSprite[0])
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
			ui.selection.SetBuild(terrain, id, randSpriteVal)
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
