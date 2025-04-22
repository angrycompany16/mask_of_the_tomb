package ui

import (
	"fmt"
	"mask_of_the_tomb/internal/game/UI/display"
	"mask_of_the_tomb/internal/game/UI/fonts"
	"mask_of_the_tomb/internal/game/UI/screenfade"
	"mask_of_the_tomb/internal/game/core/assetloader/assettypes"
)

// TODO: Reset selector position on active menu change
// TODO: allow for optional borders on elements (shouldn't be too hard)

// Everything can be rewritten with events...
// TODO: Make menu select into an event (With event info!!!)
type UI struct {
	activeDisplay *display.Display
	displays      []*display.Display
	DeathEffect   *screenfade.DeathEffect
}

// Loads a single menu file and sets it as the active menu
func (ui *UI) LoadPreamble(path string) {
	fonts.LoadPreamble()
	loadingscreen, err := display.FromFile(path)
	if err != nil {
		panic(err)
	}

	ui.activeDisplay = loadingscreen
}

func (ui *UI) Load(menuPaths ...string) {
	for _, menuPath := range menuPaths {
		ui.displays = append(ui.displays, assettypes.NewMenuAsset(menuPath))
	}
}

func (ui *UI) AddDisplayManual(display *display.Display) {
	ui.displays = append(ui.displays, display)
	ui.activeDisplay = display
}

func (ui *UI) Init() {
}

func (ui *UI) Update() {
	ui.activeDisplay.Update()
	// ui.DeathEffect.Update()
}

// TODO: Try to enable switching active menu with enum instead of string
func (ui *UI) SwitchActiveDisplay(name string) {
	ui.activeDisplay.Root.Reset()

	for _, menu := range ui.displays {
		if menu.Name != name {
			continue
		}
		ui.activeDisplay = menu
		return
	}

	panic(fmt.Sprintln("Failed to switch to menu with name", name))
}

func (ui *UI) Draw() {
	ui.activeDisplay.Draw()
	ui.DeathEffect.Draw()
}

func (ui *UI) GetConfirmations() map[string]bool {
	return ui.activeDisplay.GetConfirmed()
}

func (ui *UI) GetSubmits() map[string]string {
	return ui.activeDisplay.GetSubmitted()
}

func (ui *UI) GetFileSearchValue() string {
	// if ui.activeMenu.FileSearch != nil {
	// 	return ui.activeMenu.FileSearch.Searchfield.Textbox.Text
	// }
	return ""
}

func (ui *UI) ResetFileSearch() {
	// ui.activeMenu.FileSearch.Searchfield.Textbox.Text = ""
}

// TODO?: replace this?
func NewUI() *UI {
	return &UI{
		DeathEffect: screenfade.NewDeathEffect(),
		displays:    make([]*display.Display, 0),
	}
}
