package ui

import (
	"image/color"
	"mask_of_the_tomb/internal/engine/advertisers"
	"mask_of_the_tomb/internal/engine/entities"
	"mask_of_the_tomb/internal/engine/events"
	pubui "mask_of_the_tomb/internal/entities/UI/pub"
	pubgame "mask_of_the_tomb/internal/entities/game/pub"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	TextColorNormal = ColorPair{
		BrightColor: color.RGBA{205, 247, 226, 255},
		DarkColor:   color.RGBA{199, 176, 139, 255},
	}
	TextColorSelected = ColorPair{
		BrightColor: color.RGBA{255, 255, 255, 255},
		DarkColor:   color.RGBA{0, 0, 0, 255},
	}

	DefaultShadowX, DefaultShadowY = -4.0, 4.0
)

// TODO: Convert into asset files? BASed
var (
	mainMenu = newMenu(
		make([]*textbox, 0),
		[]*selectable{
			newSelectable(
				"Play",
				newTextBoxSimple("Play video game", 48, 0, -100, 10, text.AlignCenter, Centered),
				TextColorNormal,
				TextColorSelected,
			),
			newSelectable(
				"Options",
				newTextBoxSimple("Options", 48, 0, 0, 10, text.AlignCenter, Centered),
				TextColorNormal,
				TextColorSelected,
			),
			newSelectable(
				"Quit",
				newTextBoxSimple("Don't play video game", 48, 0, 100, 10, text.AlignCenter, Centered),
				TextColorNormal,
				TextColorSelected,
			),
		},
	)

	Pausemenu = newMenu(
		make([]*textbox, 0),
		[]*selectable{
			newSelectable(
				"Resume",
				newTextBoxSimple("Resume", 48, 0, -100, 10, text.AlignCenter, Centered),
				TextColorNormal,
				TextColorSelected,
			),
			newSelectable(
				"Options",
				newTextBoxSimple("Options", 48, 0, 0, 10, text.AlignCenter, Centered),
				TextColorNormal,
				TextColorSelected,
			),
			newSelectable(
				"Quit",
				newTextBoxSimple("Quit", 48, 0, 100, 10, text.AlignCenter, Centered),
				TextColorNormal,
				TextColorSelected,
			),
		},
	)
	Hud = newMenu(
		[]*textbox{
			newTextBoxSimple("Text", defaultFontSize, 24, 24, defaultLineSpacing, text.AlignStart, TopLeft),
		},
		make([]*selectable, 0),
	)
)

type UI struct {
	activeMenu *menu
}

type ColorPair struct {
	BrightColor color.Color
	DarkColor   color.Color
}

func NewUI() *UI {
	_ui := UI{
		activeMenu: mainMenu,
	}
	entities.RegisterEntity(&_ui, "UI")

	return &_ui
}

func (ui *UI) Update() {
	adv := advertisers.GetAdvertiser(pubgame.GameEntityName)
	val := adv.Read().(pubgame.GameAdvertiser)

	confirmations := ui.activeMenu.getConfirmed()

	switch val.State {
	case pubgame.StateMainMenu:
		if val, ok := confirmations["Play"]; ok && val {
			pubui.UISelected.Raise(events.EventInfo{Data: pubui.SelectPlay})
			ui.activeMenu = Hud
		} else if val, ok := confirmations["Quit"]; ok && val {
			pubui.UISelected.Raise(events.EventInfo{Data: pubui.SelectQuit})
		}
	case pubgame.StatePlaying:

	case pubgame.StatePaused:
		if val, ok := confirmations["Resume"]; ok && val {
			// pubui.UISelected.Raise(events.EventInfo{Data: pubui.SelectPlay})
			// pubui.UISelected.Raise(events.EventInfo{Data: pubui.SelectQuit})
			// ui.activeMenu = mainMenu
			// ui.activeMenu.selectorPos = 0
		} else if val, ok := confirmations["Quit"]; ok && val {
			pubui.UISelected.Raise(events.EventInfo{Data: pubui.SelectMainMenu})
			// Save data and stuff
			// Loading screens
			// etc
			ui.activeMenu.selectorPos = 0
			ui.activeMenu = mainMenu
		}
	}

	inputDir := 0
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		inputDir = 1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		inputDir = -1
	}

	ui.activeMenu.update(inputDir)
}

func (ui *UI) Draw() {
	ui.activeMenu.draw()
}
