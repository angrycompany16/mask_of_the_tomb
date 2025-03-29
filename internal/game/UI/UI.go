package ui

import (
	"fmt"
	"mask_of_the_tomb/internal/game/UI/menu"
	"mask_of_the_tomb/internal/game/UI/screenfade"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// var (
// 	TextColorNormal = ColorPair{
// 		BrightColor: color.RGBA{205, 247, 226, 255},
// 		DarkColor:   color.RGBA{199, 176, 139, 255},
// 	}
// 	TextColorSelected = ColorPair{
// 		BrightColor: color.RGBA{255, 255, 255, 255},
// 		DarkColor:   color.RGBA{0, 0, 0, 255},
// 	}

// 	DefaultShadowX, DefaultShadowY = -4.0, 4.0
// )

// var (
// 	LoadingScreen = newMenu(
// 		[]*Textbox{newTextBoxSimple("CURRENTLY LOADING YOUR FUTURE SUFFERING", defaultFontSize, 24, 24, defaultLineSpacing, text.AlignCenter, Centered)},
// 		make([]*Selectable, 0),
// 	)

// 	Mainmenu = newMenu(
// 		make([]*Textbox, 0),
// 		[]*Selectable{
// 			newSelectable(
// 				"Play",
// 				newTextBoxSimple("Play video game", 48, 0, -100, 10, text.AlignCenter, Centered),
// 				TextColorNormal,
// 				TextColorSelected,
// 			),
// 			newSelectable(
// 				"Options",
// 				newTextBoxSimple("Options", 48, 0, 0, 10, text.AlignCenter, Centered),
// 				TextColorNormal,
// 				TextColorSelected,
// 			),
// 			newSelectable(
// 				"Quit",
// 				newTextBoxSimple("Don't play video game", 48, 0, 100, 10, text.AlignCenter, Centered),
// 				TextColorNormal,
// 				TextColorSelected,
// 			),
// 		},
// 	)

// 	Pausemenu = newMenu(
// 		make([]*Textbox, 0),
// 		[]*Selectable{
// 			newSelectable(
// 				"Resume",
// 				newTextBoxSimple("Resume", 48, 0, -100, 10, text.AlignCenter, Centered),
// 				TextColorNormal,
// 				TextColorSelected,
// 			),
// 			newSelectable(
// 				"Options",
// 				newTextBoxSimple("Options", 48, 0, 0, 10, text.AlignCenter, Centered),
// 				TextColorNormal,
// 				TextColorSelected,
// 			),
// 			newSelectable(
// 				"Quit",
// 				newTextBoxSimple("Quit", 48, 0, 100, 10, text.AlignCenter, Centered),
// 				TextColorNormal,
// 				TextColorSelected,
// 			),
// 		},
// 	)
// 	Hud = newMenu(
// 		[]*Textbox{
// 			newTextBoxSimple("Text", defaultFontSize, 24, 24, defaultLineSpacing, text.AlignStart, TopLeft),
// 		},
// 		make([]*Selectable, 0),
// 	)
// )

// TODO: Make menu select into an event (With event info!!!)
type UI struct {
	activeMenu  *menu.Menu
	menus       map[string]*menu.Menu
	DeathEffect *screenfade.DeathEffect
}

func (ui *UI) Load(menuPaths ...string) {
	// for _, menuPath := range menuPaths {
	// assetloader.AddAsset(assetloader.men)
	// Load menu path
	// }
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

	ui.activeMenu.Update(inputDir)
	ui.DeathEffect.Update()
}

func (ui *UI) SwitchActiveMenu(name string) {
	menu, ok := ui.menus[name]
	if !ok {
		fmt.Println("Failed to switch to menu with name", name)
		return
	}
	ui.activeMenu = menu
	ui.activeMenu.SelectorPos = 0
}

func (ui *UI) Draw() {
	ui.activeMenu.Draw()
	ui.DeathEffect.Draw()
}

func (ui *UI) GetConfirmations() map[string]bool {
	return ui.activeMenu.GetConfirmed()
}

// Not great, really not great
// func (ui *UI) SetScore(score int) {
// Hud.Textboxes[0].text = fmt.Sprintf("YOUR SCORE IS: %d", score)
// }

// TODO?: replace this?
// func NewUI() *UI {
// 	return &UI{
// 		activeMenu:  LoadingScreen,
// 		DeathEffect: screenfade.NewDeathEffect(),
// 	}
// }
