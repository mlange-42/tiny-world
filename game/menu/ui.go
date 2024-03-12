package menu

import (
	"fmt"
	"image/color"
	"io/fs"
	"math/rand"
	"slices"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/save"
	"github.com/mlange-42/tiny-world/game/sprites"
	"github.com/mlange-42/tiny-world/game/terr"
)

type UI struct {
	fs         fs.FS
	saveFolder string
	mapsFolder string

	ui *ebitenui.UI

	sprites      *res.Sprites
	infoLabel    *widget.Text
	tabContainer *widget.TabBook
	tabs         []*widget.TabBookTab

	background        *image.NineSlice
	backgroundHover   *image.NineSlice
	backgroundPressed *image.NineSlice
}

func (ui *UI) UI() *ebitenui.UI {
	return ui.ui
}

func NewUI(f fs.FS, folder, mapsFolder string, selectedTab int, sprts *res.Sprites, fonts *res.Fonts,
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
	ui.tabs = append(ui.tabs, newWorldTab, scenariosTab, loadWorldTab)

	img := ui.defaultButtonImage()
	ui.tabContainer = widget.NewTabBook(
		widget.TabBookOpts.TabButtonImage(img),
		widget.TabBookOpts.TabButtonText(fonts.Default, &widget.ButtonTextColor{Idle: ui.sprites.TextColor, Disabled: ui.sprites.TextColor}),
		widget.TabBookOpts.TabButtonSpacing(0),
		widget.TabBookOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
				widget.WidgetOpts.MinSize(360, 20),
			),
		),
		widget.TabBookOpts.TabButtonOpts(
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(60, 0)),
		),
		widget.TabBookOpts.Tabs(
			newWorldTab,
			scenariosTab,
			loadWorldTab,
		),
		widget.TabBookOpts.InitialTab(ui.tabs[selectedTab]),
	)

	menuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(360, 360),
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
			}),
			widget.WidgetOpts.MinSize(360, 20),
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
				Stretch:  true,
			}),
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
			}),
			widget.WidgetOpts.MinSize(360, 20),
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
		menuContainer.AddChild(gameButton)

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
			}),
			widget.WidgetOpts.MinSize(360, 20),
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
		menuContainer.AddChild(newButton)

	}

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
	scale := 2

	props := &terr.Properties[t]

	bx, by := ui.sprites.TileWidth+4, ui.sprites.TileWidth+4

	xOff := (bx - ui.sprites.TileWidth) / 2
	yOff := (by - ui.sprites.TileWidth) / 2
	img := ebiten.NewImage(bx*scale, by*scale)

	ui.background.Draw(img, bx, by, func(opts *ebiten.DrawImageOptions) {
		opts.GeoM.Scale(float64(scale), float64(scale))
	})

	idx := ui.sprites.GetTerrainIndex(t)
	height := 0

	if props.TerrainBelow != terr.Air {
		idx2 := ui.sprites.GetTerrainIndex(props.TerrainBelow)
		info2 := ui.sprites.GetInfo(idx2)

		sp2 := ui.sprites.Get(idx2)
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(xOff), float64(ui.sprites.TileWidth+yOff-sp2.Bounds().Dy()))
		op.GeoM.Scale(float64(scale), float64(scale))
		img.DrawImage(sp2, &op)

		height = info2.Height
	}

	sp1 := ui.sprites.GetRand(idx, 0, 0)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xOff), float64(ui.sprites.TileWidth+yOff-sp1.Bounds().Dy()-height))
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
