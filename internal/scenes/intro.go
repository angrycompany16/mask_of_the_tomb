package scenes

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	ui "mask_of_the_tomb/internal/plugins/UI"
)

type IntroScene struct {
	UI   *ui.UI
	exit bool
}

func (i *IntroScene) Init() {
	introScreenLayer := errs.Must(assettypes.GetYamlAsset("introScreen")).(*ui.Layer)

	i.UI = ui.NewUI([]*ui.Layer{introScreenLayer}, make(map[string]*ui.Overlay))
	i.UI.SwitchActiveDisplay("intro", nil)
}

func (i *IntroScene) Update() {
	confirmations := i.UI.GetConfirmations()
	if confirm, ok := confirmations["Introtext"]; ok && confirm.IsConfirmed {
		i.exit = true
	}
}

func (i *IntroScene) Draw() {
	i.UI.Draw()
}

func (i *IntroScene) Exit() bool {
	return i.exit
}
