package ui

import (
	"fmt"
	"maps"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/libraries/node"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

type UI struct {
	Root          node.Root `yaml:"Root"`
	overlays      map[string]*Overlay
	confirmations map[string]node.ConfirmInfo
}

// Create asset groups that can be loaded with a single function call?
// Loads a single menu file and sets it as the active menu
// TODO: Some problems here. First, the font loading is a bit dumb (why is it the way it is?)
// Second, we should probably just use the standard asset loading for this too
func LoadPreambleUI(path string) *UI {
	assetloader.LoadPreamble()

	ui := UI{}
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&ui)
	if err != nil {
		return nil
	}
	return &ui
}

func (ui *UI) AddOverlay(name string, overlay *Overlay) {
	if ui.overlays == nil {
		ui.overlays = make(map[string]*Overlay)
	}
	ui.overlays[name] = overlay
}

func (ui *UI) Update() {
	ui.confirmations = make(map[string]node.ConfirmInfo)
	ui.Root.Update(ui.confirmations)
	for _, overlay := range ui.overlays {
		overlay.Update()
	}
}

func (ui *UI) GetOverlay(name string) *Overlay {
	if ui.overlays[name] == nil {
		fmt.Println("-------------")
		fmt.Println("ERROR: Could not get screen overlay with name", name)
		fmt.Println("List of names:")
		names := slices.Collect(maps.Keys(ui.overlays))
		for _, name := range names {
			fmt.Println(name)
		}
		fmt.Println("-------------")
	}
	return ui.overlays[name]
}

func (ui *UI) Draw() {
	ui.Root.Draw(0, 0, rendering.GAME_WIDTH*rendering.PIXEL_SCALE, rendering.GAME_HEIGHT*rendering.PIXEL_SCALE)

	for _, overlay := range ui.overlays {
		overlay.Draw()
	}
}

func (ui *UI) GetConfirmations() map[string]node.ConfirmInfo {
	return ui.confirmations
}

func (ui *UI) GetSubmits() map[string]string {
	// TODO: Implement??
	return nil
}

func (ui *UI) Reset(overWriteInfo map[string]node.OverWriteInfo) {
	ui.Root.Reset(overWriteInfo)
}
