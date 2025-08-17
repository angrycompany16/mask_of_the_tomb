package scenes

import (
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/animation"
	"mask_of_the_tomb/internal/libraries/particles"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"
)

type LoadingSceneBehaviour struct {
	UI              *ui.UI
	SceneTransition scene.SceneTransition
	exit            bool
}

func (l *LoadingSceneBehaviour) Init() {
	l.UI.LoadPreamble(loadingScreenPath)

	assetloader.Add("any", &delayAsset)

	assetloader.Add("LDTKAsset", assettypes.NewLDTKAsset(LDTKMapPath))
	assetloader.Add("slamboxTilemap", assettypes.MakeImageAsset(assets.Slambox_tilemap))
	assetloader.Add("grassTilemap", assettypes.MakeImageAsset(assets.Grass_tiles))
	assetloader.Add("turretSprite", assettypes.MakeImageAsset(assets.Turret_sprite))
	assetloader.Add("fogShader", assettypes.MakeShaderAsset(assets.Fog_kage))
	assetloader.Add("vignetteShader", assettypes.MakeShaderAsset(assets.Vignette_kage))
	assetloader.Add("pixelLightsShader", assettypes.MakeShaderAsset(assets.Pixel_lights_kage))
	assetloader.Add("ambientParticles", assettypes.MakeYamlAsset(assets.Basement_yaml, &particles.ParticleSystem{}))
	assetloader.Add("teleporterAnimation", assettypes.MakeYamlAsset(assets.Teleporter_yaml, &animation.AnimationInfo{}))

	assetloader.Add("playerSprite", assettypes.MakeImageAsset(assets.Player_sprite))
	assetloader.Add("dashSound", assettypes.MakeAudioStreamAsset(assets.Dash_wav, assettypes.Wav))
	assetloader.Add("slamSound", assettypes.MakeAudioStreamAsset(assets.Slam_wav, assettypes.Wav))
	assetloader.Add("deathSound", assettypes.MakeAudioStreamAsset(assets.Death_mp3, assettypes.Mp3))
	assetloader.Add("jumpParticlesBroad", assettypes.MakeYamlAsset(assets.Jump_broad_yaml, &particles.ParticleSystem{}))
	assetloader.Add("jumpParticlesTight", assettypes.MakeYamlAsset(assets.Jump_tight_yaml, &particles.ParticleSystem{}))
	assetloader.Add("dashInitAnim", assettypes.MakeYamlAsset(assets.Dash_init_yaml, &animation.AnimationInfo{}))
	assetloader.Add("dashLoopAnim", assettypes.MakeYamlAsset(assets.Dash_loop_yaml, &animation.AnimationInfo{}))
	assetloader.Add("playerIdleAnim", assettypes.MakeYamlAsset(assets.Player_idle_yaml, &animation.AnimationInfo{}))
	assetloader.Add("playerSlamAnim", assettypes.MakeYamlAsset(assets.Player_slam_yaml, &animation.AnimationInfo{}))

	assetloader.Add("mainMenu", assettypes.MakeYamlAsset(assets.Main_menu_yaml, &ui.Layer{}))
	assetloader.Add("pauseMenu", assettypes.MakeYamlAsset(assets.Pause_menu_yaml, &ui.Layer{}))
	assetloader.Add("optionsMenu", assettypes.MakeYamlAsset(assets.Options_yaml, &ui.Layer{}))
	assetloader.Add("introScreen", assettypes.MakeYamlAsset(assets.Intro_yaml, &ui.Layer{}))
	assetloader.Add("emptyMenu", assettypes.MakeYamlAsset(assets.Empty_yaml, &ui.Layer{}))
	assetloader.Add("hud", assettypes.MakeYamlAsset(assets.Hud_yaml, &ui.Layer{}))

	assetloader.Add("menuTheme", assettypes.MakeAudioStreamAsset(assets.Menu_mp3, assettypes.Mp3))
	assetloader.Add("basementTheme", assettypes.MakeAudioStreamAsset(assets.Basement_wav, assettypes.Wav))
	assetloader.Add("libraryTheme", assettypes.MakeAudioStreamAsset(assets.Library_mp3, assettypes.Mp3))

	assetloader.Add("transitionShader", assettypes.MakeShaderAsset(assets.Transition_kage))
	assetloader.Add("selectSound", assettypes.MakeAudioStreamAsset(assets.Select_ogg, assettypes.Ogg))
	assetloader.Add("dialogueSound", assettypes.MakeAudioStreamAsset(assets.Text_scroll_ogg, assettypes.Ogg))
	assetloader.Add("saveData", save.MakeSaveAsset(SaveProfile))
	assetloader.Add("titleCard", assettypes.MakeImageAsset(assets.Level_titlecard_sprite))

	go assetloader.LoadAll(loadFinishedChan)
}

func (l *LoadingSceneBehaviour) Update() {
	events.Update()
	l.UI.Update()

	if _, done := threads.Poll(loadFinishedChan); done {
		fmt.Println("Finished loading stage")
		assetloader.PrintAssetRegistry()
		l.exit = true
	}
}

func (l *LoadingSceneBehaviour) Draw() { l.UI.Draw() }
func (l *LoadingSceneBehaviour) Transition() (scene.SceneTransition, bool) {
	return l.SceneTransition, l.exit
}
