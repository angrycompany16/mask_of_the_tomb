package scenes

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/core/sound_v2"
	"mask_of_the_tomb/internal/libraries/node"
	ui "mask_of_the_tomb/internal/plugins/UI"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type OptionsScene struct {
	UI        *ui.UI
	lastScene scene.Scene
}

func (o *OptionsScene) Init() {
	o.UI.Reset(map[string]node.OverWriteInfo{
		"Master_vol": {SliderVal: resources.Settings.MasterVolume},
		"Music_vol":  {SliderVal: resources.Settings.MusicVolume},
		"Sound_vol":  {SliderVal: resources.Settings.SoundVolume},
	})
}

func (o *OptionsScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	o.UI.Update()
	confirmations := o.UI.GetConfirmations()

	if confirm, ok := confirmations["Back"]; ok && confirm.IsConfirmed ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
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
	// Apply volume to sound server (maybe we don't need settings as resource anymore?)
	totalSfx := resources.Settings.MasterVolume * resources.Settings.SoundVolume / 10000
	totalMusic := resources.Settings.MasterVolume * resources.Settings.MusicVolume / 10000

	// AN EPIC system
	sound_v2.EditDSPChannelEffect("sfxMaster", "vol", sound_v2.SetVolumeAction(totalSfx))
	sound_v2.EditDSPChannelEffect("musicMaster", "vol", sound_v2.SetVolumeAction(totalMusic))

	return nil, false
}

func (o *OptionsScene) Draw()           { o.UI.Draw() }
func (o *OptionsScene) GetName() string { return "optionsScene" }
func MakeOptionsScene(lastScene scene.Scene) *OptionsScene {
	return &OptionsScene{
		UI:        errs.Must(assettypes.GetYamlAsset("optionsMenu")).(*ui.UI),
		lastScene: lastScene,
	}
}
