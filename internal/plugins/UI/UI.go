package ui

import (
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/libraries/node"
)

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
		// fmt.Println("Loaded layer", path)
		ui.layers = append(ui.layers, NewLayerAsset(path))
	}
}

func (ui *UI) AddDisplayManual(display *Layer) {
	ui.layers = append(ui.layers, display)
	ui.activeLayer = display
}

func (ui *UI) Update() {
	// fmt.Println(ui)
	ui.activeLayer.Update()
	for _, overlay := range ui.overlays {
		overlay.Update()
	}
}

func (ui *UI) SwitchActiveDisplay(name string, overWriteInfo map[string]node.OverWriteInfo) {
	// if ui.activeLayer != nil {
	// }

	for _, menu := range ui.layers {
		if menu.Name != name {
			continue
		}
		ui.activeLayer = menu
		ui.activeLayer.Root.Reset(overWriteInfo)
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

func (ui *UI) GetConfirmations() map[string]node.ConfirmInfo {
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
