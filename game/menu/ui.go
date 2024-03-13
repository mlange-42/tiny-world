package menu

import (
	"fmt"
	"image/color"
	"io/fs"
	"math"
	"math/rand"
	"slices"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/res/achievements"
	"github.com/mlange-42/tiny-world/game/save"
	"github.com/mlange-42/tiny-world/game/sprites"
	"github.com/mlange-42/tiny-world/game/terr"
)

const panelHeight = 400

type UI struct {
	fs         fs.FS
	saveFolder string
	mapsFolder string

	ui *ebitenui.UI

	sprites      *res.Sprites
	infoLabel    *widget.Text
	tabContainer *widget.TabBook
	tabs         []*widget.TabBookTab

	empty             *image.NineSlice
	background        *image.NineSlice
	backgroundHover   *image.NineSlice
	backgroundPressed *image.NineSlice
}

func (ui *UI) UI() *ebitenui.UI {
	return ui.ui
}

func NewUI(f fs.FS, folder, mapsFolder string, selectedTab int, sprts *res.Sprites, fonts *res.Fonts,
	achievements *achievements.Achievements,
	start func(string, save.MapLocation, save.LoadType),
	restart func(tab int)) UI {
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
		widget.TextOpts.Text("   ", fonts.Default, ui.sprites.TextColor),
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
	scenariosTab.AddChild(ui.createScenariosPanel(games, fonts, start))

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
				widget.WidgetOpts.MinSize(480, panelHeight),
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

func (ui *UI) createNewWorldPanel(games []string, fonts *res.Fonts, start func(string, save.MapLocation, save.LoadType)) *widget.Container {
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
			start(name, save.MapLocation{}, save.LoadTypeNone)
		}),
	)
	menuContainer.AddChild(newButton)

	return menuContainer
}

func (ui *UI) createLoadPanel(games []string, fonts *res.Fonts,
	start func(string, save.MapLocation, save.LoadType),
	restart func(tab int)) *widget.Container {
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
				start(game, save.MapLocation{}, save.LoadTypeGame)
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

func (ui *UI) createScenariosPanel(games []string, fonts *res.Fonts, start func(string, save.MapLocation, save.LoadType)) *widget.Container {
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

	for _, m := range maps {
		newButton := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				}),
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
				start(name, m, save.LoadTypeMap)
			}),
		)
		content.AddChild(newButton)

	}

	menuContainer.AddChild(scroll)

	return menuContainer
}

func (ui *UI) createAchievementsPanel(achievements *achievements.Achievements, fonts *res.Fonts) *widget.Container {
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

	label := widget.NewText(
		widget.TextOpts.Text("Achievements:", fonts.Default, ui.sprites.TextColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)
	menuContainer.AddChild(label)

	scroll, content := ui.createScrollPanel(panelHeight - 110)

	for i := range achievements.Achievements {
		ach := achievements.Achievements[i]

		achButton := widget.NewButton(
			widget.ButtonOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				}),
			),
			widget.ButtonOpts.Image(img),
			widget.ButtonOpts.Text(ach.Name, fonts.Default, &widget.ButtonTextColor{
				Idle:     ui.sprites.TextColor,
				Disabled: color.NRGBA{180, 180, 180, 255},
			}),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		)
		if !ach.Completed {
			achButton.GetWidget().Disabled = true
		}

		content.AddChild(achButton)
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

func (ui *UI) createTerrainGraphic(terrain terr.Terrain) *widget.Graphic {
	img := ui.createTerrainImage(terrain)

	button := widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(ui.sprites.TileWidth, ui.sprites.TileWidth),
		),
		widget.GraphicOpts.Image(img),
	)

	return button
}

func (ui *UI) createTerrainImage(t terr.Terrain) *ebiten.Image {
	scale := 3

	props := &terr.Properties[t]

	bx, by := ui.sprites.TileWidth, ui.sprites.TileWidth

	img := ebiten.NewImage(bx*scale, by*scale)

	idx := ui.sprites.GetTerrainIndex(t)
	height := 0

	if props.TerrainBelow != terr.Air {
		idx2 := ui.sprites.GetTerrainIndex(props.TerrainBelow)
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
		vSlider.Current -= int(math.Round(a.Y * float64(p)))
	})

	rootContainer.AddChild(vSlider)

	return rootContainer, content
}
