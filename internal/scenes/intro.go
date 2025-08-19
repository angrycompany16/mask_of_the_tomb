package scenes

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/scene"
	ui "mask_of_the_tomb/internal/plugins/UI"
)

// "introScene"
type IntroSceneBehaviour struct {
	UI *ui.UI
}

func (i *IntroSceneBehaviour) Init() {
	introScreenLayer := errs.Must(assettypes.GetYamlAsset("introScreen")).(*ui.Layer)

	i.UI = ui.NewUI([]*ui.Layer{introScreenLayer}, make(map[string]*ui.Overlay))
	i.UI.SwitchActiveDisplay("intro", nil)
}

func (i *IntroSceneBehaviour) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	confirmations := i.UI.GetConfirmations()
	if confirm, ok := confirmations["Introtext"]; ok && confirm.IsConfirmed {
		return &scene.SceneTransition{
			Kind:       scene.Replace,
			OtherScene: &GameplayScene{},
		}, true
	}
	return nil, false
}

func (i *IntroSceneBehaviour) Draw()           { i.UI.Draw() }
func (i *IntroSceneBehaviour) GetName() string { return "introScene" }
