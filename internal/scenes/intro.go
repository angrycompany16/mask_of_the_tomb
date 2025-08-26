package scenes

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/scene"
	ui "mask_of_the_tomb/internal/plugins/UI"
)

type IntroScene struct {
	UI *ui.UI
}

func (i *IntroScene) Init() {
	i.UI.SwitchActiveDisplay("intro", nil)
}

func (i *IntroScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	i.UI.Update()
	confirmations := i.UI.GetConfirmations()
	if confirm, ok := confirmations["Introtext"]; ok && confirm.IsConfirmed {
		return &scene.SceneTransition{
			Kind:       scene.Replace,
			Name:       i.GetName(),
			OtherScene: MakeGameplayScene(),
		}, true
	}
	return nil, false
}

func (i *IntroScene) Draw()           { i.UI.Draw() }
func (i *IntroScene) GetName() string { return "introScene" }
func MakeIntroScene() *IntroScene {
	introScene := IntroScene{}
	introScreenLayer := errs.Must(assettypes.GetYamlAsset("introScreen")).(*ui.Layer)

	introScene.UI = ui.NewUI([]*ui.Layer{introScreenLayer}, make(map[string]*ui.Overlay))
	return &introScene
}
