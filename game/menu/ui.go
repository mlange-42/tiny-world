package menu

import (
	"fmt"
	stdimage "image"
	"image/color"
	"io/fs"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
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
	tabContainer    *widget.FlipBook
	selectedTab     int
	tabs            []*widget.Container
	scenarioButtons []*widget.Button

	loadButtonsGroup     *widget.RadioGroup
	scenarioButtonsGroup *widget.RadioGroup

	empty             *image.NineSlice
	background        *image.NineSlice
	backgroundHover   *image.NineSlice
	backgroundPressed *image.NineSlice
	textHighlightHex  string
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
		fs:               f,
		saveFolder:       folder,
		mapsFolder:       mapsFolder,
		sprites:          sprts,
		textHighlightHex: util.ColorToBB(sprts.TextHighlightColor),
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

	mainTab := ui.createMainMenuPanel(games, fonts, start)
	newWorldTab := ui.createNewWorldPanel(games, fonts, start)
	scenariosTab := ui.createScenariosPanel(games, achievements, fonts, start)
	loadWorldTab := ui.createLoadPanel(games, fonts, start, restart)
	achievementTab := ui.createAchievementsPanel(achievements, fonts)
	ui.tabs = append(ui.tabs, mainTab, newWorldTab, scenariosTab, loadWorldTab, achievementTab)

	ui.tabContainer = widget.NewFlipBook(
		widget.FlipBookOpts.ContainerOpts(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
					Stretch:  true,
				}),
				widget.WidgetOpts.MinSize(panelWidth, panelHeight),
			),
			widget.ContainerOpts.BackgroundImage(ui.background),
		),
	)
	ui.selectPage(selectedTab)

	menuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
			widget.RowLayoutOpts.Padding(&widget.Insets{Top: 24}),
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

func (ui *UI) selectPage(idx int) {
	ui.tabContainer.SetPage(ui.tabs[idx])
	ui.selectedTab = idx
	ui.infoLabel.Label = "   "
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

func (ui *UI) createMainMenuPanel(games []save.SaveGame, fonts *res.Fonts, start startFunction) *widget.Container {
	menuContainer := ui.createTabPanel()

	mainLabel := ui.createMainMenuLabel("Main menu", fonts)
	menuContainer.AddChild(mainLabel)

	enabled := false
	text := "Continue"
	if len(games) > 0 {
		enabled = true
		text = fmt.Sprintf("Continue %s", games[0].Name)
	}
	continueButton := ui.createMainMenuButton(text, fonts,
		func(args *widget.ButtonClickedEventArgs) {
			if enabled {
				start(games[0].Name, save.MapLocation{}, save.LoadTypeGame, false)
			}
		})
	continueButton.GetWidget().Disabled = !enabled
	menuContainer.AddChild(continueButton)

	newButton := ui.createMainMenuButton("New World", fonts,
		func(args *widget.ButtonClickedEventArgs) { ui.selectPage(1) })
	menuContainer.AddChild(newButton)

	scenariosButton := ui.createMainMenuButton("Scenarios", fonts,
		func(args *widget.ButtonClickedEventArgs) { ui.selectPage(2) })
	menuContainer.AddChild(scenariosButton)

	loadButton := ui.createMainMenuButton("Load World", fonts,
		func(args *widget.ButtonClickedEventArgs) { ui.selectPage(3) })
	menuContainer.AddChild(loadButton)
	if len(games) == 0 {
		loadButton.GetWidget().Disabled = true
	}

	achievementsButton := ui.createMainMenuButton("Achievements", fonts,
		func(args *widget.ButtonClickedEventArgs) { ui.selectPage(4) })
	menuContainer.AddChild(achievementsButton)

	if runtime.GOOS != "js" {
		quitButton := ui.createMainMenuButton("Quit", fonts,
			func(args *widget.ButtonClickedEventArgs) { os.Exit(0) })
		menuContainer.AddChild(quitButton)
	}

	return menuContainer
}

func (ui *UI) createMainMenuButton(text string, fonts *res.Fonts, click func(args *widget.ButtonClickedEventArgs)) *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
			widget.WidgetOpts.MinSize(240, 0),
		),
		widget.ButtonOpts.Image(ui.defaultButtonImage()),
		widget.ButtonOpts.Text(text, fonts.Default, &widget.ButtonTextColor{
			Idle:     ui.sprites.TextColor,
			Disabled: ui.sprites.TextColor,
		}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		widget.ButtonOpts.ClickedHandler(click),
	)

	return button
}

func (ui *UI) createBackStartButtons(text string, fonts *res.Fonts,
	click func(args *widget.ButtonClickedEventArgs)) (*widget.Container, *widget.Button) {
	cols := 2
	if len(text) == 0 {
		cols = 1
	}
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(cols),
			widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{false}),
			widget.GridLayoutOpts.Spacing(48, 12),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				},
			),
		),
	)

	backButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionStart,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
		),
		widget.ButtonOpts.Image(ui.defaultButtonImage()),
		widget.ButtonOpts.Text("Back", fonts.Default, &widget.ButtonTextColor{
			Idle:     ui.sprites.TextColor,
			Disabled: ui.sprites.TextColor,
		}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.selectPage(0)
		}),
	)
	container.AddChild(backButton)

	if len(text) == 0 {
		return container, nil
	}

	startButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionEnd,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
		),
		widget.ButtonOpts.Image(ui.defaultButtonImage()),
		widget.ButtonOpts.Text(text, fonts.Default, &widget.ButtonTextColor{
			Idle:     ui.sprites.TextColor,
			Disabled: ui.sprites.TextColor,
		}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		widget.ButtonOpts.ClickedHandler(click),
	)
	container.AddChild(startButton)

	return container, startButton
}

func (ui *UI) createMainMenuLabel(text string, fonts *res.Fonts) *widget.Text {
	label := widget.NewText(
		widget.TextOpts.Text(text, fonts.Default, ui.sprites.TextColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				},
			),
		),
	)

	return label
}

func (ui *UI) createNewWorldPanel(games []save.SaveGame, fonts *res.Fonts, start startFunction) *widget.Container {
	menuContainer := ui.createTabPanel()

	newLabel := ui.createMainMenuLabel("New World", fonts)
	menuContainer.AddChild(newLabel)

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

	click := func(args *widget.ButtonClickedEventArgs) {
		name := newName.GetText()
		if len(name) == 0 {
			ui.infoLabel.Label = "Give your world a name!"
			return
		}
		for _, g := range games {
			if g.Name == name {
				ui.infoLabel.Label = "World already exists!"
				return
			}
		}
		if !save.IsValidName(name) {
			ui.infoLabel.Label = "Use only letters, numbers,\nspaces, '-' and '_'!"
			return
		}
		isEditor := ebiten.IsKeyPressed(ebiten.KeyShift)
		start(name, save.MapLocation{}, save.LoadTypeNone, isEditor)
	}

	buttons, _ := ui.createBackStartButtons("New World", fonts, click)
	menuContainer.AddChild(buttons)

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

func (ui *UI) createLoadPanel(games []save.SaveGame, fonts *res.Fonts,
	start startFunction,
	restart menuFunction) *widget.Container {
	menuContainer := ui.createTabPanel()

	worldsLabel := ui.createMainMenuLabel("Load World", fonts)
	menuContainer.AddChild(worldsLabel)

	img := ui.defaultButtonImage()

	scroll, content := ui.createScrollPanel(panelHeight - 76)

	buttons := make([]widget.RadioGroupElement, len(games))
	for i, game := range games {
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
			widget.ButtonOpts.Text(game.Name, fonts.Default, &widget.ButtonTextColor{
				Idle:     ui.sprites.TextColor,
				Disabled: ui.sprites.TextColor,
			}),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.ToggleMode(),
		)
		content.AddChild(gameButton)
		buttons[i] = gameButton

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
				if err := deleteGame(ui.saveFolder, game.Name); err != nil {
					ui.infoLabel.Label = err.Error()
					return
				}
				menuContainer.RemoveChild(gameButton)
				restart(ui.selectedTab)
			}),
		)
		contextMenu.AddChild(deleteButton)
	}

	ui.loadButtonsGroup = widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(buttons...),
	)

	menuContainer.AddChild(scroll)

	btn, _ := ui.createBackStartButtons("Load World", fonts,
		func(args *widget.ButtonClickedEventArgs) {
			idx := slices.Index(buttons, ui.loadButtonsGroup.Active())
			start(games[idx].Name, save.MapLocation{}, save.LoadTypeGame, false)
		},
	)
	menuContainer.AddChild(btn)

	return menuContainer
}

func (ui *UI) createScenariosPanel(games []save.SaveGame, achievements *achievements.Achievements, fonts *res.Fonts, start startFunction) *widget.Container {
	maps, err := save.ListMaps(ui.fs, ui.mapsFolder)
	if err != nil {
		panic(err)
	}

	menuContainer := ui.createTabPanel()

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

	mapsLabel := ui.createMainMenuLabel("Scenarios", fonts)
	menuContainer.AddChild(mapsLabel)
	menuContainer.AddChild(newName)

	scroll, content := ui.createScrollPanel(panelHeight - 113)

	mapsUnlocked := []save.MapLocation{}
	mapsLocked := []save.MapLocation{}
	achUnlocked := []save.MapInfo{}
	achLocked := []save.MapInfo{}
	cntEnabled := 0

	for _, m := range maps {
		ach, err := save.LoadMapData(ui.fs, ui.mapsFolder, m)
		if err != nil {
			log.Fatalf("error loading achievements for map %s: %s", m.Name, err.Error())
		}

		enabled := true
		for _, a := range ach.Achievements {
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

	buttons := make([]widget.RadioGroupElement, len(mapsUnlocked))
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
		if len(ach.Achievements) > 0 {
			for _, a := range ach.Achievements {
				if name, ok := achievements.IdMap[a]; ok {
					if name.Completed {
						achieve += fmt.Sprintf("\n - %s", name.Name)
					} else {
						achieve += fmt.Sprintf("\n[color=%s] - %s[/color]", ui.textHighlightHex, name.Name)
					}
				}
			}
		} else {
			achieve = "\n - none"
		}
		localText := ""
		localMarker := ""
		if !m.IsEmbedded {
			localText = " (*local)"
			localMarker = " (*)"
		}
		description := ""
		if len(ach.Description) > 0 {
			description = ach.Description + "\n\n"
		}

		label := widget.NewText(
			widget.TextOpts.ProcessBBCode(true),
			widget.TextOpts.Text(fmt.Sprintf("%s%s\n\n%sRequired achievements:\n%s", m.Name, localText, description, achieve), fonts.Default, ui.sprites.TextColor),
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
			widget.ButtonOpts.Text(m.Name+localMarker, fonts.Default, &widget.ButtonTextColor{
				Idle:     ui.sprites.TextColor,
				Disabled: color.NRGBA{180, 180, 180, 255},
			}),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.ToggleMode(),
		)
		newButton.GetWidget().Disabled = !enabled
		content.AddChild(newButton)
		ui.scenarioButtons = append(ui.scenarioButtons, newButton)
		buttons[i] = newButton
	}
	ui.scenarioButtonsGroup = widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(buttons...),
	)

	menuContainer.AddChild(scroll)

	btns, btn := ui.createBackStartButtons("Start Scenario", fonts,
		func(args *widget.ButtonClickedEventArgs) {
			name := newName.GetText()
			idx := slices.Index(buttons, ui.scenarioButtonsGroup.Active())
			if ui.scenarioButtons[idx].GetWidget().Disabled {
				ui.infoLabel.Label = "Select an unlocked scenario!"
				return
			}
			if len(name) == 0 {
				ui.infoLabel.Label = "Give your world a name!"
				return
			}
			for _, g := range games {
				if g.Name == name {
					ui.infoLabel.Label = "World already exists!"
					return
				}
			}
			if !save.IsValidName(name) {
				ui.infoLabel.Label = "Use only letters, numbers,\nspaces, '-' and '_'!"
				return
			}
			isEditor := ebiten.IsKeyPressed(ebiten.KeyShift)
			start(name, mapsUnlocked[idx], save.LoadTypeMap, isEditor)
		},
	)
	if cntEnabled == 0 {
		btn.GetWidget().Disabled = true
	}
	menuContainer.AddChild(btns)

	return menuContainer
}

func (ui *UI) createAchievementsPanel(achieves *achievements.Achievements, fonts *res.Fonts) *widget.Container {
	menuContainer := ui.createTabPanel()
	img := ui.emptyImage()

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

	label := ui.createMainMenuLabel(fmt.Sprintf("Achievements (%d/%d)", len(achUnlocked)-len(achLocked), len(achUnlocked)), fonts)
	menuContainer.AddChild(label)

	scroll, content := ui.createScrollPanel(panelHeight - 76)

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
			icon = ui.createTerrainImage(t, ach.IconIndex, 1)
		} else {
			idx := ui.sprites.GetIndex(ach.Icon)
			icon = ui.sprites.GetSprite(ui.sprites.GetMultiTileIndex(idx, terr.Directions(ach.IconIndex), 0, 0))
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

	buttons, _ := ui.createBackStartButtons("", fonts, nil)
	menuContainer.AddChild(buttons)

	return menuContainer
}

func (ui *UI) createTabPanel() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(6)),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				},
			),
			widget.WidgetOpts.MinSize(panelWidth, panelHeight),
		),
	)
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
	img := ui.createTerrainImage(terrain, 0, 3)

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

func (ui *UI) createTerrainImage(t terr.Terrain, tileIdx int, scale int) *ebiten.Image {
	props := &terr.Properties[t]

	bx, by := ui.sprites.TileWidth, ui.sprites.TileWidth

	img := ebiten.NewImage(bx*scale, by*scale)

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

	idx := ui.sprites.GetTerrainIndex(t)
	subIdx := ui.sprites.GetMultiTileIndex(idx, terr.Directions(tileIdx), 0, 0)
	sp1 := ui.sprites.GetSprite(subIdx)
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
