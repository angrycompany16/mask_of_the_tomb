package game

import (
	"fmt"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/libraries/assettypes"
	"mask_of_the_tomb/internal/libraries/node"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Initializes menu, music settings and such by retrieveing loaded assets
func (g *Game) InitMenuStage() {
	initTime = time.Now()
	gameData := errs.Must(save.GetSaveAsset("saveData"))
	resources.Settings = gameData.Settings
	g.mainUI.SwitchActiveDisplay("mainmenu", nil)
	g.gameplayUI.SwitchActiveDisplay("empty", nil)
	g.musicPlayer.Init()
	node.SelectSound = errs.Must(assettypes.GetEffectPlayerAsset("selectSound"))
	node.DialogueSound = errs.Must(assettypes.GetEffectPlayerAsset("dialogueSound"))
}

func (g *Game) MenuStageUpdate() error {
	g.LoadingStageUpdate()
	confirmations := g.mainUI.GetConfirmations()
	resources.Time = time.Since(initTime).Seconds()

	g.gameplayUI.Update()

	titlecard := g.gameplayUI.GetOverlay("titlecard")
	if _, raised := g.titleCardTimeoutListener.Poll(); raised {
		titlecard.StartFadeOut()
	}

	ebitenutil.DebugPrint(rendering.ScreenLayers.Overlay, fmt.Sprintf("TPS: %0.2f \nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

	if resources.State != resources.MainMenu {
		return nil
	}
	g.musicPlayer.ResetMusic()

	// --- MAIN MENU SPECIFIC ---
	g.musicPlayer.PlayMenuMusic()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mainUI.SwitchActiveDisplay("mainmenu", nil)
	}
	if confirm, ok := confirmations["Play"]; ok && confirm.IsConfirmed {
		gameData := errs.Must(save.GetSaveAsset("saveData"))
		if gameData.SpawnRoomName == "" {
			resources.State = resources.Intro
			g.InitIntroStage()
		} else {
			resources.State = resources.Playing
			g.InitGameplayStage(gameData, false)
		}
	} else if confirm, ok := confirmations["Quit"]; ok && confirm.IsConfirmed {
		save.SaveGame(save.SaveData{resources.PreviousLevelName, resources.Settings}, SaveProfile)
		return ErrTerminated
	} else if confirm, ok := confirmations["Options"]; ok && confirm.IsConfirmed {
		g.mainUI.SwitchActiveDisplay("options", map[string]node.OverWriteInfo{
			"Master_vol": {SliderVal: resources.Settings.MasterVolume},
			"Music_vol":  {SliderVal: resources.Settings.MusicVolume},
			"Sound_vol":  {SliderVal: resources.Settings.SoundVolume},
		})
	}
	updateOptions(confirmations)
	return nil
}

func (g *Game) MenuStageDraw() {
	g.LoadingStageDraw()
}

func updateOptions(confirmations map[string]node.ConfirmInfo) {
	if confirm, ok := confirmations["Master_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.MasterVolume = confirm.SliderVal
	} else if confirm, ok := confirmations["Music_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.MusicVolume = confirm.SliderVal
	} else if confirm, ok := confirmations["Sound_vol"]; ok && confirm.IsConfirmed {
		resources.Settings.SoundVolume = confirm.SliderVal
	}
}
