package scenes

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/core/sound_v2"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"
)

type MenuScene struct {
	UI *ui.UI
}

func (m *MenuScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	m.UI.Update()

	confirmations := m.UI.GetConfirmations()

	if confirm, ok := confirmations["Play"]; ok && confirm.IsConfirmed {
		sound_v2.StopSound("menuTheme")
		gameData := errs.Must(save.GetSaveAsset("saveData"))
		if gameData.SpawnRoomName == "" {
			return &scene.SceneTransition{
				Kind:       scene.Replace,
				Name:       m.GetName(),
				OtherScene: MakeIntroScene(),
			}, true
		} else {
			return &scene.SceneTransition{
				Kind:       scene.Replace,
				Name:       m.GetName(),
				OtherScene: MakeGameplayScene(),
			}, true
		}
	} else if confirm, ok := confirmations["Quit"]; ok && confirm.IsConfirmed {
		save.SaveGame(save.SaveData{resources.PreviousLevelName, resources.Settings}, SaveProfile)
		return &scene.SceneTransition{Kind: scene.Quit}, true
	} else if confirm, ok := confirmations["Options"]; ok && confirm.IsConfirmed {
		return &scene.SceneTransition{
			Kind:       scene.Replace,
			Name:       m.GetName(),
			OtherScene: MakeOptionsScene(m),
		}, true
	}

	return nil, false
}

func (m *MenuScene) Init() {
	m.UI.Reset(nil)
	sound_v2.PlaySound("menuTheme", "musicMaster", 0)
}
func (m *MenuScene) Draw()           { m.UI.Draw() }
func (m *MenuScene) GetName() string { return "menuScene" }
func MakeMenuScene() *MenuScene {
	return &MenuScene{
		UI: errs.Must(assettypes.GetYamlAsset("mainMenu")).(*ui.UI),
	}
}
