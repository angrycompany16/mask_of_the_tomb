package scenes

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/libraries/node"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	UI *ui.UI
}

func (p *PauseScene) Init() {
	pauseMenuLayer := errs.Must(assettypes.GetYamlAsset("pauseMenu")).(*ui.Layer)
	optionsMenuLayer := errs.Must(assettypes.GetYamlAsset("optionsMenu")).(*ui.Layer)

	p.UI = ui.NewUI([]*ui.Layer{pauseMenuLayer, optionsMenuLayer}, make(map[string]*ui.Overlay))
	p.UI.SwitchActiveDisplay("pausemenu", nil)
}

func (p *PauseScene) Update() {
	confirmations := p.UI.GetConfirmations()

	if confirm, ok := confirmations["Resume"]; ok && confirm.IsConfirmed || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Exit back to play stage
	} else if confirm, ok := confirmations["Quit"]; ok && confirm.IsConfirmed {
		save.SaveGame(save.SaveData{resources.PreviousLevelName, resources.Settings}, SaveProfile)
		// Exit back to main menu

		// Solve this somehow
		InitLevelName = g.world.ActiveLevel.GetName()
	} else if confirm, ok := confirmations["Options"]; ok && confirm.IsConfirmed {
		p.UI.SwitchActiveDisplay("options", map[string]node.OverWriteInfo{
			"Master_vol": {SliderVal: resources.Settings.MasterVolume},
			"Music_vol":  {SliderVal: resources.Settings.MusicVolume},
			"Sound_vol":  {SliderVal: resources.Settings.SoundVolume},
		})
	} else if confirm, ok := confirmations["Back"]; ok && confirm.IsConfirmed {
		p.UI.SwitchActiveDisplay("pausemenu", nil)
	}

	if confirm, ok := confirmations["Master_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.MasterVolume = confirm.SliderVal
	} else if confirm, ok := confirmations["Music_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.MusicVolume = confirm.SliderVal
	} else if confirm, ok := confirmations["Sound_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.SoundVolume = confirm.SliderVal
	}
}

func (p *PauseScene) Draw() {
	p.UI.Draw()
}

func (p *PauseScene) Exit() bool {
	return false
}
