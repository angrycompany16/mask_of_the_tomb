package scenes

import (
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/animation"
	"mask_of_the_tomb/internal/libraries/particles"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"
	"path/filepath"
	"time"
)

var (
	loadingScreenPath = filepath.Join("assets", "menus", "game", "loadingscreen.yaml")
	LDTKMapPath       = filepath.Join("assets", "LDTK", "world.ldtk")
)

type LoadingScene struct {
	UI               *ui.UI
	loadFinishedChan chan int
}

func (l *LoadingScene) Init() {
	// l.UI.LoadPreamble(loadingScreenPath)

	assetloader.Add("any", assettypes.NewDelayAsset(time.Second))

	assetloader.Add("LDTKAsset", assettypes.NewLDTKAsset(LDTKMapPath))
	assetloader.Add("slamboxTilemap", assettypes.MakeImageAsset(assets.Slambox_tilemap_png))
	assetloader.Add("grassTilemap", assettypes.MakeImageAsset(assets.Grass_png))
	assetloader.Add("turretSprite", assettypes.MakeImageAsset(assets.Turret_png))
	assetloader.Add("catcherSprite", assettypes.MakeImageAsset(assets.Catcher_png))
	assetloader.Add("lanternSprite", assettypes.MakeImageAsset(assets.Lantern_png))
	assetloader.Add("fogShader", assettypes.MakeShaderAsset(assets.Fog_kage))
	assetloader.Add("vignetteShader", assettypes.MakeShaderAsset(assets.Vignette_kage))
	assetloader.Add("pixelLightsShader", assettypes.MakeShaderAsset(assets.Pixel_lights_kage))
	assetloader.Add("ambientParticles", assettypes.MakeYamlAsset(assets.Basement_yaml, &particles.ParticleSystem{}))
	assetloader.Add("teleporterAnimation", assettypes.MakeYamlAsset(assets.Teleporter_yaml, &animation.AnimationInfo{}))

	assetloader.Add("playerSprite", assettypes.MakeImageAsset(assets.Player_png))
	assetloader.Add("dashSound", assettypes.MakeAudioStreamAsset(assets.Dash_wav, assettypes.Wav))
	assetloader.Add("slamSound", assettypes.MakeAudioStreamAsset(assets.Slam_wav, assettypes.Wav))
	assetloader.Add("deathSound", assettypes.MakeAudioStreamAsset(assets.Death_mp3, assettypes.Mp3))
	assetloader.Add("jumpParticlesBroad", assettypes.MakeYamlAsset(assets.Jump_broad_yaml, &particles.ParticleSystem{}))
	assetloader.Add("jumpParticlesTight", assettypes.MakeYamlAsset(assets.Jump_tight_yaml, &particles.ParticleSystem{}))
	assetloader.Add("dashInitAnim", assettypes.MakeYamlAsset(assets.Dash_init_yaml, &animation.AnimationInfo{}))
	assetloader.Add("dashLoopAnim", assettypes.MakeYamlAsset(assets.Dash_loop_yaml, &animation.AnimationInfo{}))
	assetloader.Add("playerIdleAnim", assettypes.MakeYamlAsset(assets.Player_idle_yaml, &animation.AnimationInfo{}))
	assetloader.Add("playerSlamAnim", assettypes.MakeYamlAsset(assets.Player_slam_yaml, &animation.AnimationInfo{}))

	assetloader.Add("mainMenu", assettypes.MakeYamlAsset(assets.Main_menu_yaml, &ui.UI{}))
	assetloader.Add("pauseMenu", assettypes.MakeYamlAsset(assets.Pause_menu_yaml, &ui.UI{}))
	assetloader.Add("optionsMenu", assettypes.MakeYamlAsset(assets.Options_yaml, &ui.UI{}))
	assetloader.Add("introScreen", assettypes.MakeYamlAsset(assets.Intro_yaml, &ui.UI{}))
	assetloader.Add("emptyMenu", assettypes.MakeYamlAsset(assets.Empty_yaml, &ui.UI{}))
	assetloader.Add("hud", assettypes.MakeYamlAsset(assets.Hud_yaml, &ui.UI{}))

	assetloader.Add("menuTheme", assettypes.MakeAudioStreamAsset(assets.Menu_mp3, assettypes.Mp3))
	assetloader.Add("basementTheme", assettypes.MakeAudioStreamAsset(assets.Basement_wav, assettypes.Wav))
	assetloader.Add("libraryTheme", assettypes.MakeAudioStreamAsset(assets.Library_mp3, assettypes.Mp3))

	assetloader.Add("transitionShader", assettypes.MakeShaderAsset(assets.Transition_kage))
	assetloader.Add("selectSound", assettypes.MakeAudioStreamAsset(assets.Select_ogg, assettypes.Ogg))
	assetloader.Add("dialogueSound", assettypes.MakeAudioStreamAsset(assets.Text_scroll_ogg, assettypes.Ogg))
	assetloader.Add("saveData", save.MakeSaveAsset(SaveProfile))
	assetloader.Add("titleCard", assettypes.MakeImageAsset(assets.Level_titlecard_sprite))

	go assetloader.LoadAll(l.loadFinishedChan)
}

func (l *LoadingScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	l.UI.Update()

	if _, done := threads.Poll(l.loadFinishedChan); done {
		fmt.Println("Finished loading stage")
		assetloader.PrintAssetRegistry()

		return &scene.SceneTransition{
			Kind:       scene.Replace,
			Name:       l.GetName(),
			OtherScene: MakeBaseScene(),
		}, true
	}
	return nil, false
}

func (l *LoadingScene) Draw()           { l.UI.Draw() }
func (l *LoadingScene) GetName() string { return "loadingScene" }
func MakeLoadingScene() *LoadingScene {
	return &LoadingScene{
		UI:               ui.LoadPreambleUI(loadingScreenPath),
		loadFinishedChan: make(chan int),
	}
}
