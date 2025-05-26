package ui

import (
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader"
)

// NOw:
// Fade out title card bug fix (Ideally create a new layer for "game" UI and "menu" UI)
// Fix the bug where the button select sound plays upon opening a menu
// I think we might have to separate this into smalles libraries because this is just horrendous
// Allow options menu to actually set master, sfx, song volume

// TODO: Reset selector position on active menu change
// TODO: allow for optional borders on elements (shouldn't be too hard)

// Everything can be rewritten with events...
// TODO: Make menu select into an event (With event info!!!)
type UI struct {
	activeLayer *Layer
	layers      []*Layer
	overlays    map[string]*Overlay
}

// Create asset groups that can be loaded with a single function call?
// Loads a single menu file and sets it as the active menu
func (ui *UI) LoadPreamble(path string) {
	assetloader.LoadPreamble()
	loadingscreen, err := FromFile(path)
	if err != nil {
		panic(err)
	}

	ui.activeLayer = loadingscreen
}

func (ui *UI) Load(paths ...string) {
	for _, path := range paths {
		ui.layers = append(ui.layers, NewLayerAsset(path))
	}
}

func (ui *UI) AddDisplayManual(display *Layer) {
	ui.layers = append(ui.layers, display)
	ui.activeLayer = display
}

func (ui *UI) Update() {
	ui.activeLayer.Update()
	for _, overlay := range ui.overlays {
		overlay.Update()
	}
}

// TODO: Try to enable switching active menu with enum instead of string
func (ui *UI) SwitchActiveDisplay(name string) {
	if ui.activeLayer != nil {
		ui.activeLayer.Root.Reset()
	}

	for _, menu := range ui.layers {
		if menu.Name != name {
			continue
		}
		ui.activeLayer = menu
		return
	}

	panic(fmt.Sprintln("Failed to switch to menu with name", name))
}

func (ui *UI) GetOverlay(name string) *Overlay {
	return ui.overlays[name]
}

func (ui *UI) Draw() {
	ui.activeLayer.Draw()
	for _, overlay := range ui.overlays {
		overlay.Draw()
	}
}

func (ui *UI) GetConfirmations() map[string]bool {
	return ui.activeLayer.GetConfirmed()
}

func (ui *UI) GetSubmits() map[string]string {
	return ui.activeLayer.GetSubmitted()
}

func NewUI(overlays map[string]*Overlay) *UI {
	return &UI{
		layers:   make([]*Layer, 0),
		overlays: overlays,
	}
}
