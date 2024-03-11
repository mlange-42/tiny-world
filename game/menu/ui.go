package menu

import (
	"fmt"
	"image/color"
	"slices"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/save"
)

type UI struct {
	ui *ebitenui.UI

	infoLabel *widget.Text
}

func (ui *UI) UI() *ebitenui.UI {
	return ui.ui
}

func NewUI(folder string, fonts *res.Fonts, start func(string, bool)) UI {
	games, err := save.ListSaveGames(folder)
	if err != nil {
		panic(err)
	}

	ui := UI{}
	ui.infoLabel = widget.NewText(
		widget.TextOpts.Text("   ", fonts.Default, color.RGBA{R: 255, G: 255, B: 150, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(12)),
		)),
	)

	newWorldTab := widget.NewTabBookTab("New World",
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 0, 255})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)
	newWorldTab.AddChild(ui.createNewWorldPanel(games, fonts, start))

	loadWorldTab := widget.NewTabBookTab("Load World",
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 0, 255})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)
	loadWorldTab.AddChild(ui.createLoadPanel(folder, games, fonts, start))

	img := loadButtonImage()
	tabContainer := widget.NewTabBook(
		widget.TabBookOpts.TabButtonImage(img),
		widget.TabBookOpts.TabButtonText(fonts.Default, &widget.ButtonTextColor{Idle: color.White, Disabled: color.White}),
		widget.TabBookOpts.TabButtonSpacing(0),
		widget.TabBookOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
				widget.WidgetOpts.MinSize(260, 20),
			),
		),
		widget.TabBookOpts.TabButtonOpts(
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(98, 0)),
		),
		widget.TabBookOpts.Tabs(newWorldTab, loadWorldTab),
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
			widget.WidgetOpts.MinSize(260, 360),
		),
	)

	titleLabel := widget.NewText(
		widget.TextOpts.Text("Tiny World", fonts.Title, color.RGBA{R: 255, G: 255, B: 255, A: 255}),
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
	menuContainer.AddChild(tabContainer)

	rootContainer.AddChild(menuContainer)

	eui := ebitenui.UI{
		Container: rootContainer,
	}
	ui.ui = &eui

	return ui
}

func (ui *UI) createNewWorldPanel(games []string, fonts *res.Fonts, start func(string, bool)) *widget.Container {
	menuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(260, 20),
		),
	)

	img := loadButtonImage()

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
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 40, G: 40, B: 40, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 80, G: 80, B: 80, A: 255}),
		}),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{254, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(fonts.Default, 2),
		),
	)

	newButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
		widget.ButtonOpts.Image(img),
		widget.ButtonOpts.Text("New World", fonts.Default, &widget.ButtonTextColor{
			Idle: color.NRGBA{255, 255, 255, 255},
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
			start(name, false)
		}),
	)

	menuContainer.AddChild(newName)
	menuContainer.AddChild(newButton)

	return menuContainer
}

func (ui *UI) createLoadPanel(folder string, games []string, fonts *res.Fonts, start func(string, bool)) *widget.Container {
	menuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(260, 20),
		),
	)

	if len(games) > 0 {
		worldsLabel := widget.NewText(
			widget.TextOpts.Text("Load world:", fonts.Default, color.White),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		)
		menuContainer.AddChild(worldsLabel)
	}

	img := loadButtonImage()

	for _, game := range games {
		contextMenu := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(4)),
			)),
			widget.ContainerOpts.BackgroundImage(
				image.NewNineSliceColor(color.NRGBA{20, 20, 20, 255}),
			),
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
				Idle:     color.NRGBA{255, 255, 255, 255},
				Disabled: color.NRGBA{180, 180, 180, 255},
			}),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				start(game, true)
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
				Idle:     color.NRGBA{255, 255, 255, 255},
				Disabled: color.NRGBA{180, 180, 180, 255},
			}),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				if err := deleteGame(folder, game); err != nil {
					ui.infoLabel.Label = err.Error()
					return
				}
				menuContainer.RemoveChild(gameButton)
			}),
		)
		contextMenu.AddChild(deleteButton)

	}

	return menuContainer
}

func loadButtonImage() *widget.ButtonImage {
	idle := image.NewNineSliceColor(color.NRGBA{60, 60, 60, 255})

	hover := image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255})

	pressed := image.NewNineSliceColor(color.NRGBA{20, 20, 20, 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}

func deleteGame(folder, game string) error {
	return save.DeleteGame(folder, game)
}
