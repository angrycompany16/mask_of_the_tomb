package ui

import (
	"fmt"
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type UI struct {
	pausemenu *menu
	mainmenu  *menu
	hud       *textbox
}

func (ui *UI) Init() {

}

func (ui *UI) Update() {
	inputDir := 0
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		inputDir = -1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		inputDir = 1
	}

	confirmInput := inpututil.IsKeyJustReleased(ebiten.KeySpace) || inpututil.IsKeyJustReleased(ebiten.KeyEnter)

	switch utils.GlobalState {
	case utils.StateMainMenu:
		ui.mainmenu.update(inputDir, confirmInput)
	case utils.StatePlaying:
	case utils.StatePaused:
	}
	// Take input
	// Update active UI elements
}

func (ui *UI) Draw() {
	switch utils.GlobalState {
	case utils.StateMainMenu:
		ui.mainmenu.draw()
	case utils.StatePlaying:
		ui.hud.draw()
	case utils.StatePaused:
		ui.pausemenu.draw()
	}
}

func (ui *UI) SetScore(score int) {
	ui.hud.text = fmt.Sprintf("YOUR SCORE IS: %d", score)
}

func NewUI() *UI {
	return &UI{
		pausemenu: newMenu(
			[]*selectable{
				newSelectable(
					newTextBoxSimple("Return", 48, 24, 24, 10, text.AlignCenter),
					utils.TextColorNormal,
					utils.TextColorSelected,
				),
			},
		),
		mainmenu: newMenu(
			[]*selectable{
				newSelectable(
					newTextBoxSimple("Play video game", 48, 24, 24, 10, text.AlignCenter),
					utils.TextColorNormal,
					utils.TextColorSelected,
				),
				newSelectable(
					newTextBoxSimple("Don't play video game", 48, 24, 124, 10, text.AlignCenter),
					utils.TextColorNormal,
					utils.TextColorSelected,
				),
			},
		),
		hud: newTextBoxSimple("Text", defaultFontSize, 24, 24, defaultLineSpacing, text.AlignCenter),
	}
}
