package menu

import (
	"fmt"
	stdimage "image"
	"image/color"
	"io/fs"
	"log"
	"math"
	"math/rand"
	"slices"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/res/achievements"
	"github.com/mlange-42/tiny-world/game/save"
	"github.com/mlange-42/tiny-world/game/sprites"
	"github.com/mlange-42/tiny-world/game/terr"
	"github.com/mlange-42/tiny-world/game/util"
)

const panelWidth = 500
const panelHeight = 400

type startFunction = func(name string, mapLoc save.MapLocation, loadType save.LoadType, isEditor bool)
type menuFunction = func(tab int)

const editorModeText = "Shift+click for scenario editor mode."

type UI struct {
	fs         fs.FS
	saveFolder string
	mapsFolder string

	ui *ebitenui.UI

	sprites         *res.Sprites
	infoLabel       *widget.Text
	tabContainer    *widget.TabBook
	tabs            []*widget.TabBookTab
	scenarioButtons []*widget.Button

	empty             *image.NineSlice
	background        *image.NineSlice
	backgroundHover   *image.NineSlice
	backgroundPressed *image.NineSlice
}

func (ui *UI) UI() *ebitenui.UI {
	return ui.ui
}

func (ui *UI) UnlockAll() {
	for _, b := range ui.scenarioButtons {
		b.GetWidget().Disabled = false
	}
}

func NewUI(f fs.FS, folder, mapsFolder string, selectedTab int, sprts *res.Sprites, fonts *res.Fonts,
	achievements *achievements.Achievements,
	start startFunction, restart menuFunction) UI {
	ui := UI{
		fs:         f,
		saveFolder: folder,
		mapsFolder: mapsFolder,
		sprites:    sprts,
	}

	sp := ui.sprites.Get(ui.sprites.GetIndex(sprites.UiPanel))
	w := sp.Bounds().Dx()
	ui.background = image.NewNineSliceSimple(sp, w/4, w/2)

	sp = ui.sprites.Get(ui.sprites.GetIndex(sprites.UiPanelHover))
	w = sp.Bounds().Dx()
	ui.backgroundHover = image.NewNineSliceSimple(sp, w/4, w/2)

	sp = ui.sprites.Get(ui.sprites.GetIndex(sprites.UiPanelPressed))
	w = sp.Bounds().Dx()
	ui.backgroundPressed = image.NewNineSliceSimple(sp, w/4, w/2)

	ui.empty = image.NewNineSliceColor(color.Transparent)

	games, err := save.ListSaveGames(folder)
	if err != nil {
		panic(err)
	}

	ui.infoLabel = widget.NewText(
		widget.TextOpts.Text("   ", fonts.Default, ui.sprites.TextHighlightColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(10, 32),
		),
	)

	rootGrid := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{true, false, true}, []bool{true}),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(6)),
			widget.GridLayoutOpts.Spacing(6, 6),
		)),
	)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(12)),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
		),
	)

	newWorldTab := widget.NewTabBookTab("New World",
		widget.ContainerOpts.BackgroundImage(ui.background),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(6)),
		)),
	)
	newWorldTab.AddChild(ui.createNewWorldPanel(games, fonts, start))

	scenariosTab := widget.NewTabBookTab("Scenarios",
		widget.ContainerOpts.BackgroundImage(ui.background),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(6)),
		)),
	)
	scenariosTab.AddChild(ui.createScenariosPanel(games, achievements, fonts, start))

	loadWorldTab := widget.NewTabBookTab("Load World",
		widget.ContainerOpts.BackgroundImage(ui.background),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(6)),
		)),
	)
	loadWorldTab.AddChild(ui.createLoadPanel(games, fonts, start, restart))

	achievementTab := widget.NewTabBookTab("Achievements",
		widget.ContainerOpts.BackgroundImage(ui.background),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(6)),
		)),
	)
	achievementTab.AddChild(ui.createAchievementsPanel(achievements, fonts))

	ui.tabs = append(ui.tabs, newWorldTab, scenariosTab, loadWorldTab, achievementTab)

	img := ui.defaultButtonImage()
	ui.tabContainer = widget.NewTabBook(
		widget.TabBookOpts.TabButtonImage(img),
		widget.TabBookOpts.TabButtonText(fonts.Default, &widget.ButtonTextColor{Idle: ui.sprites.TextColor, Disabled: ui.sprites.TextColor}),
		widget.TabBookOpts.TabButtonSpacing(0),
		widget.TabBookOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
					Stretch:  true,
				}),
				widget.WidgetOpts.MinSize(panelWidth, panelHeight),
			),
		),
		widget.TabBookOpts.TabButtonOpts(
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(60, 0)),
		),
		widget.TabBookOpts.Tabs(ui.tabs...),
		widget.TabBookOpts.InitialTab(ui.tabs[selectedTab]),
	)

	menuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
			widget.RowLayoutOpts.Padding(widget.Insets{Top: 24}),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
		),
	)

	titleLabel := widget.NewText(
		widget.TextOpts.Text("Tiny World", fonts.Title, ui.sprites.TextColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
	)

	menuContainer.AddChild(titleLabel)
	menuContainer.AddChild(ui.infoLabel)
	menuContainer.AddChild(ui.tabContainer)

	rootContainer.AddChild(menuContainer)

	t1, t2 := ui.drawRandomSprites()

	rootGrid.AddChild(ui.createIconContainer(t1))
	rootGrid.AddChild(rootContainer)
	rootGrid.AddChild(ui.createIconContainer(t2))

	eui := ebitenui.UI{
		Container: rootGrid,
	}
	ui.ui = &eui

	return ui
}

func (ui *UI) createIconContainer(t terr.Terrain) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(24)),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
		),
	)

	graphic := ui.createTerrainGraphic(t)
	container.AddChild(graphic)

	return container
}

func (ui *UI) createNewWorldPanel(games []string, fonts *res.Fonts, start startFunction) *widget.Container {
	menuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				Stretch:  true,
			}),
		),
	)

	img := ui.defaultButtonImage()

	newName := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
		widget.TextInputOpts.Placeholder("World name"),
		widget.TextInputOpts.Face(fonts.Default),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     ui.background,
			Disabled: ui.backgroundHover,
		}),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          ui.sprites.TextColor,
			Disabled:      ui.sprites.TextColor,
			Caret:         ui.sprites.TextColor,
			DisabledCaret: ui.sprites.TextColor,
		}),
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(fonts.Default, 2),
		),
	)

	menuContainer.AddChild(newName)

	newButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
			widget.WidgetOpts.MinSize(160, 0),
		),
		widget.ButtonOpts.Image(img),
		widget.ButtonOpts.Text("New World", fonts.Default, &widget.ButtonTextColor{
			Idle:     ui.sprites.TextColor,
			Disabled: ui.sprites.TextColor,
		}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			name := newName.GetText()
			if len(name) == 0 {
				ui.infoLabel.Label = "Give your world a name!"
				return
			}
			if slices.Contains(games, name) {
				ui.infoLabel.Label = "World already exists!"
				return
			}
			if !save.IsValidName(name) {
				ui.infoLabel.Label = "Use only letters, numbers,\nspaces, '-' and '_'!"
				return
			}
			isEditor := ebiten.IsKeyPressed(ebiten.KeyShift)
			start(name, save.MapLocation{}, save.LoadTypeNone, isEditor)
		}),
	)
	menuContainer.AddChild(newButton)

	editorLabel := widget.NewText(
		widget.TextOpts.Text(editorModeText, fonts.Default, ui.sprites.TextColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionEnd),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(10, 48),
		),
	)
	menuContainer.AddChild(editorLabel)

	return menuContainer
}

func (ui *UI) createLoadPanel(games []string, fonts *res.Fonts,
	start startFunction,
	restart menuFunction) *widget.Container {
	menuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				Stretch:  true,
			}),
		),
	)

	if len(games) > 0 {
		worldsLabel := widget.NewText(
			widget.TextOpts.Text("Load world:", fonts.Default, ui.sprites.TextColor),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		)
		menuContainer.AddChild(worldsLabel)
	}

	img := ui.defaultButtonImage()

	scroll, content := ui.createScrollPanel(panelHeight - 72)

	for _, game := range games {
		contextMenu := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(4)),
			)),
			widget.ContainerOpts.BackgroundImage(ui.background),
		)

		gameButton := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				}),
				widget.WidgetOpts.ContextMenu(contextMenu),
			),
			widget.ButtonOpts.Image(img),
			widget.ButtonOpts.Text(game, fonts.Default, &widget.ButtonTextColor{
				Idle:     ui.sprites.TextColor,
				Disabled: ui.sprites.TextColor,
			}),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				start(game, save.MapLocation{}, save.LoadTypeGame, false)
			}),
		)
		content.AddChild(gameButton)

		deleteButton := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				}),
				widget.WidgetOpts.ContextMenu(contextMenu),
			),
			widget.ButtonOpts.Image(img),
			widget.ButtonOpts.Text(fmt.Sprintf("Delete '%s'", game), fonts.Default, &widget.ButtonTextColor{
				Idle:     ui.sprites.TextColor,
				Disabled: ui.sprites.TextColor,
			}),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				if err := deleteGame(ui.saveFolder, game); err != nil {
					ui.infoLabel.Label = err.Error()
					return
				}
				menuContainer.RemoveChild(gameButton)

				currTab := ui.tabContainer.Tab()
				idx := slices.Index(ui.tabs, currTab)
				restart(idx)
			}),
		)
		contextMenu.AddChild(deleteButton)

	}

	menuContainer.AddChild(scroll)

	return menuContainer
}

func (ui *UI) createScenariosPanel(games []string, achievements *achievements.Achievements, fonts *res.Fonts, start startFunction) *widget.Container {
	maps, err := save.ListMaps(ui.fs, ui.mapsFolder)
	if err != nil {
		panic(err)
	}

	menuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				Stretch:  true,
			}),
		),
	)

	img := ui.defaultButtonImage()

	newName := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
		widget.TextInputOpts.Placeholder("World name"),
		widget.TextInputOpts.Face(fonts.Default),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     ui.background,
			Disabled: ui.background,
		}),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          ui.sprites.TextColor,
			Disabled:      ui.sprites.TextColor,
			Caret:         ui.sprites.TextColor,
			DisabledCaret: ui.sprites.TextColor,
		}),
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(fonts.Default, 2),
		),
	)

	mapsLabel := widget.NewText(
		widget.TextOpts.Text("Scenarios:", fonts.Default, ui.sprites.TextColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)
	menuContainer.AddChild(newName)
	menuContainer.AddChild(mapsLabel)

	scroll, content := ui.createScrollPanel(panelHeight - 110)

	mapsUnlocked := []save.MapLocation{}
	mapsLocked := []save.MapLocation{}
	achUnlocked := [][]string{}
	achLocked := [][]string{}
	cntEnabled := 0

	for _, m := range maps {
		ach, err := save.LoadMapAchievements(ui.fs, ui.mapsFolder, m)
		if err != nil {
			log.Fatalf("error loading achievements for map %s: %s", m.Name, err.Error())
		}

		enabled := true
		for _, a := range ach {
			a2, ok := achievements.IdMap[a]
			if !ok {
				log.Printf("WARNING: Achievement '%s' in map '%s' not found", a, m.Name)
				continue
			}
			if !a2.Completed {
				enabled = false
				break
			}
		}
		if enabled {
			mapsUnlocked = append(mapsUnlocked, m)
			achUnlocked = append(achUnlocked, ach)
			cntEnabled++
		} else {
			mapsLocked = append(mapsLocked, m)
			achLocked = append(achLocked, ach)
		}
	}
	mapsUnlocked = append(mapsUnlocked, mapsLocked...)
	achUnlocked = append(achUnlocked, achLocked...)

	for i, m := range mapsUnlocked {
		ach := achUnlocked[i]
		enabled := i < cntEnabled

		tooltipContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 6, Bottom: 6, Left: 12, Right: 12}),
			)),
			widget.ContainerOpts.AutoDisableChildren(),
			widget.ContainerOpts.BackgroundImage(ui.background),
		)

		achieve := ""
		if len(ach) > 0 {
			for _, a := range ach {
				if name, ok := achievements.IdMap[a]; ok {
					achieve += "\n - " + name.Name
				}
			}
		} else {
			achieve = "\n - none"
		}

		label := widget.NewText(
			widget.TextOpts.Text(fmt.Sprintf("%s\n\nRequired achievements:\n%s", m.Name, achieve), fonts.Default, ui.sprites.TextColor),
			widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
			widget.TextOpts.MaxWidth(360),
		)
		tooltipContainer.AddChild(label)

		newButton := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				}),
				widget.WidgetOpts.ToolTip(widget.NewToolTip(
					widget.ToolTipOpts.Content(tooltipContainer),
					widget.ToolTipOpts.Offset(stdimage.Point{-5, 5}),
					widget.ToolTipOpts.Position(widget.TOOLTIP_POS_WIDGET),
					widget.ToolTipOpts.Delay(time.Millisecond*300),
				)),
			),
			widget.ButtonOpts.Image(img),
			widget.ButtonOpts.Text(m.Name, fonts.Default, &widget.ButtonTextColor{
				Idle:     ui.sprites.TextColor,
				Disabled: color.NRGBA{180, 180, 180, 255},
			}),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				name := newName.GetText()
				if len(name) == 0 {
					ui.infoLabel.Label = "Give your world a name!"
					return
				}
				if slices.Contains(games, name) {
					ui.infoLabel.Label = "World already exists!"
					return
				}
				if !save.IsValidName(name) {
					ui.infoLabel.Label = "Use only letters, numbers,\nspaces, '-' and '_'!"
					return
				}
				isEditor := ebiten.IsKeyPressed(ebiten.KeyShift)
				start(name, m, save.LoadTypeMap, isEditor)
			}),
		)
		newButton.GetWidget().Disabled = !enabled
		content.AddChild(newButton)
		ui.scenarioButtons = append(ui.scenarioButtons, newButton)
	}

	menuContainer.AddChild(scroll)

	return menuContainer
}

func (ui *UI) createAchievementsPanel(achieves *achievements.Achievements, fonts *res.Fonts) *widget.Container {
	menuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				Stretch:  true,
			}),
		),
	)

	img := ui.emptyImage()

	label := widget.NewText(
		widget.TextOpts.Text("Achievements:", fonts.Default, ui.sprites.TextColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)
	menuContainer.AddChild(label)

	scroll, content := ui.createScrollPanel(panelHeight - 72)

	achUnlocked := []*achievements.Achievement{}
	achLocked := []*achievements.Achievement{}
	for i := range achieves.Achievements {
		ach := &achieves.Achievements[i]
		if ach.Completed {
			achUnlocked = append(achUnlocked, ach)
		} else {
			achLocked = append(achLocked, ach)
		}
	}
	achUnlocked = append(achUnlocked, achLocked...)

	for _, ach := range achUnlocked {
		name := util.Capitalize(ach.Name)

		rowContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(6, 6),
				widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(6)),
				widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{false}),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
					Stretch:  true,
				}),
			),
			widget.ContainerOpts.BackgroundImage(ui.background),
		)
		var icon *ebiten.Image
		if terr.IsTerrainName(ach.Icon) {
			t := terr.ToTerrain(ach.Icon)
			icon = ui.createTerrainImage(t, 1)
		} else {
			icon = ui.sprites.Get(ui.sprites.GetIndex(ach.Icon))
		}

		graphic := widget.NewGraphic(
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
				widget.WidgetOpts.MinSize(ui.sprites.TileWidth, ui.sprites.TileWidth),
			),
			widget.GraphicOpts.Image(icon),
		)

		achButton := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.GridLayoutData{}),
			),
			widget.ButtonOpts.Image(img),
			widget.ButtonOpts.Text(name+"\n"+ach.Description, fonts.Default, &widget.ButtonTextColor{
				Idle:     ui.sprites.TextColor,
				Disabled: color.NRGBA{120, 120, 120, 255},
			}),
			widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		)
		if !ach.Completed {
			achButton.GetWidget().Disabled = true
		}

		rowContainer.AddChild(graphic)
		rowContainer.AddChild(achButton)

		content.AddChild(rowContainer)
	}

	menuContainer.AddChild(scroll)

	return menuContainer
}

func deleteGame(folder, game string) error {
	return save.DeleteGame(folder, game)
}

func (ui *UI) defaultButtonImage() *widget.ButtonImage {
	return &widget.ButtonImage{
		Idle:    ui.background,
		Hover:   ui.backgroundHover,
		Pressed: ui.backgroundPressed,
	}
}

func (ui *UI) emptyImage() *widget.ButtonImage {
	return &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(color.Transparent),
		Hover:   image.NewNineSliceColor(color.Transparent),
		Pressed: image.NewNineSliceColor(color.Transparent),
	}
}

func (ui *UI) createTerrainGraphic(terrain terr.Terrain) *widget.Graphic {
	img := ui.createTerrainImage(terrain, 3)

	graphic := widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(ui.sprites.TileWidth, ui.sprites.TileWidth),
		),
		widget.GraphicOpts.Image(img),
	)

	return graphic
}

func (ui *UI) createTerrainImage(t terr.Terrain, scale int) *ebiten.Image {
	props := &terr.Properties[t]

	bx, by := ui.sprites.TileWidth, ui.sprites.TileWidth

	img := ebiten.NewImage(bx*scale, by*scale)

	idx := ui.sprites.GetTerrainIndex(t)
	height := 0

	for _, tr := range props.TerrainBelow {
		idx2 := ui.sprites.GetTerrainIndex(tr)
		info2 := ui.sprites.GetInfo(idx2)

		sp2 := ui.sprites.Get(idx2)
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(ui.sprites.TileWidth-sp2.Bounds().Dy()))
		op.GeoM.Scale(float64(scale), float64(scale))
		img.DrawImage(sp2, &op)

		height = info2.Height
	}

	sp1 := ui.sprites.GetRand(idx, 0, 0)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(ui.sprites.TileWidth-sp1.Bounds().Dy()-height))
	op.GeoM.Scale(float64(scale), float64(scale))
	img.DrawImage(sp1, &op)

	return img
}

func (ui *UI) drawRandomSprites() (terr.Terrain, terr.Terrain) {
	candidates := []terr.Terrain{}

	for i := range terr.Properties {
		prop := &terr.Properties[i]
		if prop.TerrainBits.Contains(terr.CanBuy) && !prop.TerrainBits.Contains(terr.IsPath) {
			candidates = append(candidates, terr.Terrain(i))
		}
	}
	return candidates[rand.Intn(len(candidates))], candidates[rand.Intn(len(candidates))]
}

func (ui *UI) createScrollPanel(height int) (*widget.Container, *widget.Container) {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(ui.background),
		// the container will use an grid layout to layout its ScrollableContainer and Slider
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Spacing(2, 0),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionStart,
				Stretch:   true,
				MaxHeight: height,
			}),
		),
	)

	content := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
	)

	scrollContainer := widget.NewScrollContainer(
		widget.ScrollContainerOpts.Content(content),
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Padding(widget.NewInsetsSimple(2)),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: ui.background,
			Mask: ui.background,
		}),
	)
	rootContainer.AddChild(scrollContainer)

	//Create a function to return the page size used by the slider
	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ViewRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
	}

	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(
			&widget.SliderTrackImage{
				Idle:  ui.empty,
				Hover: ui.empty,
			},
			&widget.ButtonImage{
				Idle:    ui.backgroundPressed,
				Hover:   ui.backgroundPressed,
				Pressed: ui.backgroundPressed,
			},
		),
	)
	//Set the slider's position if the scrollContainer is scrolled by other means than the slider
	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.WidgetScrolledEventArgs)
		p := pageSizeFunc() / 3
		if p < 1 {
			p = 1
		}
		dy := math.Copysign(1.0, a.Y)
		vSlider.Current -= int(math.Round(dy * float64(p)))
	})

	rootContainer.AddChild(vSlider)

	return rootContainer, content
}
