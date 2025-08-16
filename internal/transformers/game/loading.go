package game

import (
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/threads"
	save "mask_of_the_tomb/internal/libraries/savesystem"
)

// Loads preamble, sets up other assets
func (g *Game) InitLoadingStage() {
	g.mainUI.LoadPreamble(loadingScreenPath)

	// --- Load all game assets ---
	assetloader.Add("any", &delayAsset)
	g.world.Load(LDTKMapPath)
	g.player.CreateAssets()
	g.mainUI.Load(mainMenuPath, pauseMenuPath, optionsMenuPath, introScreenPath, emptyMenuPath)
	g.gameplayUI.Load(hudPath, emptyMenuPath)
	g.musicPlayer.Load()
	assetloader.Add("transitionShader", assettypes.MakeShaderAsset(assets.Transition_kage))
	assetloader.Add("selectSound", assettypes.MakeAudioStreamAsset(assets.Select_ogg, assettypes.Ogg))
	assetloader.Add("dialogueSound", assettypes.MakeAudioStreamAsset(assets.Text_scroll_ogg, assettypes.Ogg))
	assetloader.Add("saveData", save.MakeSaveAsset(SaveProfile))
	assetloader.Add("titleCard", assettypes.MakeImageAsset(assets.Level_titlecard_sprite))

	go assetloader.LoadAll(loadFinishedChan)
}

func (g *Game) LoadingStageUpdate() {
	events.Update()
	g.mainUI.Update()

	if _, done := threads.Poll(loadFinishedChan); done {
		fmt.Println("Finished loading stage")
		fmt.Println("Loaded assets:")
		assetloader.PrintAssetRegistry()

		resources.State = resources.MainMenu
		g.InitMenuStage()
	}
}

func (g *Game) LoadingStageDraw() {
	g.mainUI.Draw()
}
