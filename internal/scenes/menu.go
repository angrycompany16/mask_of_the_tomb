package scenes

import (
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/libraries/node"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MenuScene struct {
	UI *ui.UI
}

func (m *MenuScene) Init() {
	m.UI.SwitchActiveDisplay("mainmenu", nil)
}

func (m *MenuScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	m.UI.Update()

	confirmations := m.UI.GetConfirmations()

	if musicScene, ok := sceneStack.GetScene("musicScene"); ok {
		musicScene.(*BaseScene).musicPlayer.PlayMenuMusic()
	} else {
		fmt.Println("Music player was not found in main menu")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // TODO: Check if we are already in the main menu
		m.UI.SwitchActiveDisplay("mainmenu", nil)
	}
	if confirm, ok := confirmations["Play"]; ok && confirm.IsConfirmed {
		gameData := errs.Must(save.GetSaveAsset("saveData"))
		if gameData.SpawnRoomName == "" {
			return &scene.SceneTransition{
				Kind:       scene.Replace,
				OtherScene: MakeIntroScene(),
			}, true
		} else {
			return &scene.SceneTransition{
				Kind:       scene.Replace,
				OtherScene: MakeGameplayScene(),
			}, true
		}
	} else if confirm, ok := confirmations["Quit"]; ok && confirm.IsConfirmed {
		save.SaveGame(save.SaveData{resources.PreviousLevelName, resources.Settings}, SaveProfile)
		return &scene.SceneTransition{Kind: scene.Quit}, true
	} else if confirm, ok := confirmations["Options"]; ok && confirm.IsConfirmed {
		m.UI.SwitchActiveDisplay("options", map[string]node.OverWriteInfo{
			"Master_vol": {SliderVal: resources.Settings.MasterVolume},
			"Music_vol":  {SliderVal: resources.Settings.MusicVolume},
			"Sound_vol":  {SliderVal: resources.Settings.SoundVolume},
		})
	} else if confirm, ok := confirmations["Back"]; ok && confirm.IsConfirmed {
		m.UI.SwitchActiveDisplay("mainmenu", nil)
	}

	if confirm, ok := confirmations["Master_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.MasterVolume = confirm.SliderVal
	} else if confirm, ok := confirmations["Music_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.MusicVolume = confirm.SliderVal
	} else if confirm, ok := confirmations["Sound_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.SoundVolume = confirm.SliderVal
	}

	return nil, false
}

func (m *MenuScene) Draw()           { m.UI.Draw() }
func (m *MenuScene) GetName() string { return "menuScene" }
func MakeMenuScene() *MenuScene {
	mainMenuLayer := errs.Must(assettypes.GetYamlAsset("mainMenu")).(*ui.Layer)
	optionsMenuLayer := errs.Must(assettypes.GetYamlAsset("optionsMenu")).(*ui.Layer)

	return &MenuScene{
		UI: ui.NewUI([]*ui.Layer{mainMenuLayer, optionsMenuLayer}, make(map[string]*ui.Overlay)),
	}
}
