package scenes

import (
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	"time"
)

func (g *Game) InitIntroStage() {
	g.mainUI.SwitchActiveDisplay("intro", nil)
}

func (g *Game) IntroStageUpdate() {
	g.MenuStageUpdate()

	if resources.State != resources.Intro {
		return
	}

	confirmations := g.mainUI.GetConfirmations()
	if confirm, ok := confirmations["Introtext"]; ok && confirm.IsConfirmed {
		gameData := errs.Must(save.GetSaveAsset("saveData"))
		resources.State = resources.Playing
		g.introDashTimer = time.NewTimer(time.Second)
		g.InitGameplayStage(gameData, true)
	}
}

func (g *Game) IntroStageDraw() {
	g.MenuStageDraw()
}
