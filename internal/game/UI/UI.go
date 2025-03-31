package ui

import (
	"fmt"
	"mask_of_the_tomb/internal/game/UI/fonts"
	"mask_of_the_tomb/internal/game/UI/menu"
	"mask_of_the_tomb/internal/game/UI/screenfade"
	"mask_of_the_tomb/internal/game/core/assetloader"
	"mask_of_the_tomb/internal/game/core/assetloader/assettypes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TODO: Make menu select into an event (With event info!!!)
type UI struct {
	activeMenu  *menu.Menu
	menus       []*menu.Menu
	DeathEffect *screenfade.DeathEffect
}

// Loads a single menu file and sets it as the active menu
func (ui *UI) LoadPreamble(path string) {
	fonts.FontRegistry.LoadPreamble()
	loadingscreen, err := menu.FromFile(path)
	if err != nil {
		panic(err)
	}

	ui.activeMenu = loadingscreen
}

func (ui *UI) Load(menuPaths ...string) {
	for _, menuPath := range menuPaths {
		menuAsset := assettypes.NewMenuAsset(menuPath)
		assetloader.AddAsset(menuAsset)
		// menuAsset.Menu.Name is not loaded yet (idiot)
		ui.menus = append(ui.menus, &menuAsset.Menu)
	}
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

	ui.activeMenu.UpdateSelectables(inputDir)
	ui.activeMenu.UpdateInputboxes()
	ui.DeathEffect.Update()
}

func (ui *UI) SwitchActiveMenu(name string) {
	for _, menu := range ui.menus {
		if menu.Name != name {
			continue
		}
		ui.activeMenu = menu
		ui.activeMenu.SelectorPos = 0
		return
	}

	panic(fmt.Sprintln("Failed to switch to menu with name", name))
}

func (ui *UI) Draw() {
	ui.activeMenu.Draw()
	ui.DeathEffect.Draw()
}

func (ui *UI) GetConfirmations() map[string]bool {
	return ui.activeMenu.GetConfirmed()
}

func (ui *UI) GetSubmits() map[string]string {
	return ui.activeMenu.GetSubmitted()
}

// TODO?: replace this?
func NewUI() *UI {
	return &UI{
		DeathEffect: screenfade.NewDeathEffect(),
		menus:       make([]*menu.Menu, 0),
	}
}
