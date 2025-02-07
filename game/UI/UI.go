package ui

import (
	"fmt"
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TODO: Convert into asset files? BASed
var (
	Mainmenu = newMenu(
		make([]*textbox, 0),
		[]*selectable{
			newSelectable(
				"Play",
				newTextBoxSimple("Play video game", 48, 0, -100, 10, text.AlignCenter, Centered),
				utils.TextColorNormal,
				utils.TextColorSelected,
			),
			newSelectable(
				"Options",
				newTextBoxSimple("Options", 48, 0, 0, 10, text.AlignCenter, Centered),
				utils.TextColorNormal,
				utils.TextColorSelected,
			),
			newSelectable(
				"Quit",
				newTextBoxSimple("Don't play video game", 48, 0, 100, 10, text.AlignCenter, Centered),
				utils.TextColorNormal,
				utils.TextColorSelected,
			),
		},
	)

	Pausemenu = newMenu(
		make([]*textbox, 0),
		[]*selectable{
			newSelectable(
				"Resume",
				newTextBoxSimple("Resume", 48, 0, -100, 10, text.AlignCenter, Centered),
				utils.TextColorNormal,
				utils.TextColorSelected,
			),
			newSelectable(
				"Options",
				newTextBoxSimple("Options", 48, 0, 0, 10, text.AlignCenter, Centered),
				utils.TextColorNormal,
				utils.TextColorSelected,
			),
			newSelectable(
				"Quit",
				newTextBoxSimple("Quit", 48, 0, 100, 10, text.AlignCenter, Centered),
				utils.TextColorNormal,
				utils.TextColorSelected,
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

func (ui *UI) Init() {

}

func (ui *UI) Update() {
	inputDir := 0
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		inputDir = 1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		inputDir = -1
	}

	ui.activeMenu.update(inputDir)
}

func (ui *UI) SwitchActiveMenu(menu *menu) {
	ui.activeMenu = menu
	ui.activeMenu.selectorPos = 0
}

func (ui *UI) Draw() {
	ui.activeMenu.draw()
}

func (ui *UI) GetConfirmations() map[string]bool {
	return ui.activeMenu.getConfirmed()
}

// Not great, really not great
func (ui *UI) SetScore(score int) {
	Hud.textboxes[0].text = fmt.Sprintf("YOUR SCORE IS: %d", score)
}

// TODO?: replace this?
func NewUI() *UI {
	return &UI{
		activeMenu: Mainmenu,
	}
}
