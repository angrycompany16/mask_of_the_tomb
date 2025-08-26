package scenes

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/libraries/node"
	ui "mask_of_the_tomb/internal/plugins/UI"
)

type OptionsScene struct {
	UI        *ui.UI
	lastScene scene.Scene
}

func (o *OptionsScene) Init() {
	o.UI.SwitchActiveDisplay("options", map[string]node.OverWriteInfo{
		"Master_vol": {SliderVal: resources.Settings.MasterVolume},
		"Music_vol":  {SliderVal: resources.Settings.MusicVolume},
		"Sound_vol":  {SliderVal: resources.Settings.SoundVolume},
	})
}

func (o *OptionsScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	o.UI.Update()
	confirmations := o.UI.GetConfirmations()

	if confirm, ok := confirmations["Back"]; ok && confirm.IsConfirmed {
		return &scene.SceneTransition{
			Kind:       scene.Replace,
			Name:       o.GetName(),
			OtherScene: o.lastScene,
		}, true
	} else if confirm, ok := confirmations["Master_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.MasterVolume = confirm.SliderVal
	} else if confirm, ok := confirmations["Music_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.MusicVolume = confirm.SliderVal
	} else if confirm, ok := confirmations["Sound_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.SoundVolume = confirm.SliderVal
	}

	return nil, false
}

func (o *OptionsScene) Draw()           { o.UI.Draw() }
func (o *OptionsScene) GetName() string { return "optionsScene" }
func MakeOptionsScene(lastScene scene.Scene) *OptionsScene {
	optionsMenu := errs.Must(assettypes.GetYamlAsset("optionsMenu")).(*ui.Layer)
	// pauseMenuLayer := errs.Must(assettypes.GetYamlAsset("pauseMenu")).(*ui.Layer)

	return &OptionsScene{
		UI:        ui.NewUI([]*ui.Layer{optionsMenu}, make(map[string]*ui.Overlay)),
		lastScene: lastScene,
	}
}
