package scenes

import (
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/scene"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	UI *ui.UI
}

func (p *PauseScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	p.UI.Update()
	confirmations := p.UI.GetConfirmations()

	if confirm, ok := confirmations["Resume"]; ok && confirm.IsConfirmed ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return &scene.SceneTransition{
			Kind: scene.Pop,
		}, true
	} else if confirm, ok := confirmations["Quit"]; ok && confirm.IsConfirmed {
		save.SaveGame(save.SaveData{resources.PreviousLevelName, resources.Settings}, SaveProfile)

		if gameplayScene, _, ok := sceneStack.GetScene("gameplayScene"); ok {
			InitLevelName = gameplayScene.(*GameplayScene).world.ActiveLevel.GetName()
		} else {
			fmt.Println("Could not find gameplay scene")
		}

		return &scene.SceneTransition{
			Kind:       scene.Replace,
			Name:       "gameplayScene",
			OtherScene: MakeMenuScene(),
		}, true
	} else if confirm, ok := confirmations["Options"]; ok && confirm.IsConfirmed {
		return &scene.SceneTransition{
			Kind:       scene.Replace,
			Name:       p.GetName(),
			OtherScene: MakeOptionsScene(p),
		}, true
	}

	return nil, false
}

func (p *PauseScene) Init()           { p.UI.Reset(nil) }
func (p *PauseScene) Draw()           { p.UI.Draw() }
func (p *PauseScene) GetName() string { return "pauseScene" }
func MakePauseScene() *PauseScene {
	return &PauseScene{
		UI: errs.Must(assettypes.GetYamlAsset("pauseMenu")).(*ui.UI),
	}
}
