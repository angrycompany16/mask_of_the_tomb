package game

import (
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/assettypes"
	save "mask_of_the_tomb/internal/libraries/savesystem"
)

// Loads preamble, sets up other assets
func (g *Game) InitLoadingStage() {
	g.mainUI.LoadPreamble(loadingScreenPath)

	// --- Load all game assets ---
	assetloader.Load("any", &delayAsset)
	g.world.Load(LDTKMapPath)
	g.player.CreateAssets()
	g.mainUI.Load(mainMenuPath, pauseMenuPath, optionsMenuPath, introScreenPath, emptyMenuPath)
	g.gameplayUI.Load(hudPath, emptyMenuPath)
	g.musicPlayer.Load()
	assetloader.Load("selectSound", assettypes.MakeSoundAsset(assets.Select_ogg, assettypes.Ogg))
	assetloader.Load("dialogueSound", assettypes.MakeSoundAsset(assets.Text_scroll_ogg, assettypes.Ogg))
	assetloader.Load("saveData", save.MakeSaveAsset(SaveProfile))

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
