package scenes

import (
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/libraries/node"
	save "mask_of_the_tomb/internal/libraries/savesystem"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) InitPausedStage() {
	g.mainUI.SwitchActiveDisplay("pausemenu", nil)
}

func (g *Game) PausedStageUpdate() {
	g.GameplayStageUpdate()

	confirmations := g.mainUI.GetConfirmations()

	if confirm, ok := confirmations["Resume"]; ok && confirm.IsConfirmed || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		resources.State = resources.Playing
		g.mainUI.SwitchActiveDisplay("empty", nil)
	} else if confirm, ok := confirmations["Quit"]; ok && confirm.IsConfirmed {
		// g.world.SaveLevel(g.world.ActiveLevel)
		save.SaveGame(save.SaveData{resources.PreviousLevelName, resources.Settings}, SaveProfile)
		resources.State = resources.MainMenu
		g.mainUI.SwitchActiveDisplay("mainmenu", nil)
		InitLevelName = g.world.ActiveLevel.GetName()
	} else if confirm, ok := confirmations["Options"]; ok && confirm.IsConfirmed {
		g.mainUI.SwitchActiveDisplay("options", map[string]node.OverWriteInfo{
			"Master_vol": {SliderVal: resources.Settings.MasterVolume},
			"Music_vol":  {SliderVal: resources.Settings.MusicVolume},
			"Sound_vol":  {SliderVal: resources.Settings.SoundVolume},
		})
	} else if confirm, ok := confirmations["Back"]; ok && confirm.IsConfirmed {
		g.mainUI.SwitchActiveDisplay("pausemenu", nil)
	}
	updateOptions(confirmations)
}

func (g *Game) PausedStageDraw() {
	g.GameplayStageDraw()
}
