package scenes

import (
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
	UI              *ui.UI
	sceneTransition scene.SceneTransition
	exit            bool
}

// We *need* a more sophisticated exit system
// Need to ideally specify which scene we switch to

func (m *MenuScene) Init() {
	gameData := errs.Must(save.GetSaveAsset("saveData"))
	resources.Settings = gameData.Settings

	mainMenuLayer := errs.Must(assettypes.GetYamlAsset("mainMenu")).(*ui.Layer)
	optionsMenuLayer := errs.Must(assettypes.GetYamlAsset("optionsMenu")).(*ui.Layer)

	m.UI = ui.NewUI([]*ui.Layer{mainMenuLayer, optionsMenuLayer}, make(map[string]*ui.Overlay))
	m.UI.SwitchActiveDisplay("mainmenu", nil)
}

func (m *MenuScene) Update() {
	confirmations := m.UI.GetConfirmations()

	// How can we achieve this result using the scene system?
	// Using some kind of event
	// Find node by name
	// We may then access that node's member functions and stuff
	// profit

	// IREMEMBER THE IDEA
	// - Instead of a scene with a scenebehaviour, just have a scene (with behaviour) which contains
	//   sceneMetadata as a child object! That way it should be possible to include scene switching
	//   behaviour easily into the scene's interface methods.
	// YEah
	m.musicPlayer.PlayMenuMusic()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // TODO: Check if we are already in the main menu
		m.UI.SwitchActiveDisplay("mainmenu", nil)
	}
	if confirm, ok := confirmations["Play"]; ok && confirm.IsConfirmed {
		gameData := errs.Must(save.GetSaveAsset("saveData"))
		if gameData.SpawnRoomName == "" {
			m.sceneTransition = scene.SceneTransition{
				Kind: scene.Replace,
				Name: "introScene",
			}
			m.exit = true
		} else {
			m.sceneTransition = scene.SceneTransition{
				Kind: scene.Replace,
				Name: "gameplayScene",
			}
			m.exit = true
		}
	} else if confirm, ok := confirmations["Quit"]; ok && confirm.IsConfirmed {
		save.SaveGame(save.SaveData{resources.PreviousLevelName, resources.Settings}, SaveProfile)
		// Spawn some kind of termination scene?
		return // ErrTerminated
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
}

func (m *MenuScene) Draw() {
	m.UI.Draw()
}

func (m *MenuScene) Exit() bool {
	return m.exit
}
