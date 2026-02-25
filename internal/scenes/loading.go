package scenes

import (
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/core/sound_v2"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/animation"
	"mask_of_the_tomb/internal/libraries/particles"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"
	"path/filepath"
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
	// LDTK
	assetloader.Add("LDTKAsset", assettypes.NewLDTKAsset(LDTKMapPath))

	// Save asset
	assetloader.Add("saveData", save.MakeSaveAsset(SaveProfile))

	// Images
	assetloader.Add("slamboxTilemap", assettypes.MakeImageAsset(assets.Slambox_tilemap_png))
	assetloader.Add("grassTilemap", assettypes.MakeImageAsset(assets.Grass_png))
	assetloader.Add("turretSprite", assettypes.MakeImageAsset(assets.Turret_png))
	assetloader.Add("catcherSprite", assettypes.MakeImageAsset(assets.Catcher_png))
	assetloader.Add("slamboxGemBlue", assettypes.MakeImageAsset(assets.Slambox_gem_blue_png))
	assetloader.Add("slamboxGemRed", assettypes.MakeImageAsset(assets.Slambox_gem_red_png))
	assetloader.Add("lanternSprite", assettypes.MakeImageAsset(assets.Lantern_png))
	assetloader.Add("textBoxTickSprite", assettypes.MakeImageAsset(assets.Textbox_tick_png))
	assetloader.Add("playerSprite", assettypes.MakeImageAsset(assets.Player_png))
	assetloader.Add("titleCard", assettypes.MakeImageAsset(assets.Level_titlecard_sprite))

	// Shaders
	assetloader.Add("fogShader", assettypes.MakeShaderAsset(assets.Fog_kage))
	assetloader.Add("vignetteShader", assettypes.MakeShaderAsset(assets.Vignette_kage))
	assetloader.Add("pixelLightsShader", assettypes.MakeShaderAsset(assets.Pixel_lights_kage))
	assetloader.Add("deathTransitionShader", assettypes.MakeShaderAsset(assets.Death_transition_kage))
	assetloader.Add("levelTransitionEnterShader", assettypes.MakeShaderAsset(assets.Level_transition_enter_kage))
	assetloader.Add("levelTransitionExitShader", assettypes.MakeShaderAsset(assets.Level_transition_exit_kage))

	// YAML assets
	//   Particles
	assetloader.Add("ambientParticles", assettypes.MakeYamlAsset(assets.Basement_yaml, &particles.ParticleSystem{}))
	assetloader.Add("jumpParticlesBroad", assettypes.MakeYamlAsset(assets.Jump_broad_yaml, &particles.ParticleSystem{}))
	assetloader.Add("jumpParticlesTight", assettypes.MakeYamlAsset(assets.Jump_tight_yaml, &particles.ParticleSystem{}))
	assetloader.Add("slamboxParticles", assettypes.MakeYamlAsset(assets.Slambox_particles_yaml, &particles.ParticleSystem{}))
	//   Animations
	assetloader.Add("teleporterAnimation", assettypes.MakeYamlAsset(assets.Teleporter_yaml, &animation.AnimationInfo{}))
	assetloader.Add("dashInitAnim", assettypes.MakeYamlAsset(assets.Dash_init_yaml, &animation.AnimationInfo{}))
	assetloader.Add("dashLoopAnim", assettypes.MakeYamlAsset(assets.Dash_loop_yaml, &animation.AnimationInfo{}))
	assetloader.Add("playerIdleAnim", assettypes.MakeYamlAsset(assets.Player_idle_yaml, &animation.AnimationInfo{}))
	assetloader.Add("playerSlamAnim", assettypes.MakeYamlAsset(assets.Player_slam_yaml, &animation.AnimationInfo{}))
	//   Menus
	assetloader.Add("mainMenu", assettypes.MakeYamlAsset(assets.Main_menu_yaml, &ui.UI{}))
	assetloader.Add("pauseMenu", assettypes.MakeYamlAsset(assets.Pause_menu_yaml, &ui.UI{}))
	assetloader.Add("optionsMenu", assettypes.MakeYamlAsset(assets.Options_yaml, &ui.UI{}))
	assetloader.Add("introScreen", assettypes.MakeYamlAsset(assets.Intro_yaml, &ui.UI{}))
	assetloader.Add("emptyMenu", assettypes.MakeYamlAsset(assets.Empty_yaml, &ui.UI{}))
	assetloader.Add("hud", assettypes.MakeYamlAsset(assets.Hud_yaml, &ui.UI{}))

	// Load Audio
	assetloader.Add("dashSound", assettypes.MakeAudioStreamAsset(assets.Dash_wav, assettypes.Wav))
	assetloader.Add("slamSound", assettypes.MakeAudioStreamAsset(assets.Slam_wav, assettypes.Wav))
	assetloader.Add("deathSound", assettypes.MakeAudioStreamAsset(assets.Death_mp3, assettypes.Mp3))
	assetloader.Add("slamboxLandSound", assettypes.MakeAudioStreamAsset(assets.Slambox_land_wav, assettypes.Wav))

	soundCatalogue := map[string]sound_v2.SoundData{
		"menuTheme":       sound_v2.Loop("music/menu.ogg"),
		"basementTheme":   sound_v2.Loop("music/basement.ogg"),
		"libraryTheme":    sound_v2.Loop("music/library.ogg"),
		"grasslandsTheme": sound_v2.Loop("music/grasslands.ogg"),
		"vowelA":          sound_v2.Oneshot("sfx/speech/vowel_A.ogg", 5),
		"vowelE":          sound_v2.Oneshot("sfx/speech/vowel_E.ogg", 5),
		"vowelI":          sound_v2.Oneshot("sfx/speech/vowel_I.ogg", 5),
		"vowelO":          sound_v2.Oneshot("sfx/speech/vowel_O.ogg", 5),
		"vowelU":          sound_v2.Oneshot("sfx/speech/vowel_U.ogg", 5),
		"selectUI":        sound_v2.Oneshot("sfx/select.ogg", 2),
		"dialogueUI":      sound_v2.Oneshot("sfx/text-scroll.ogg", 5),
		"playerDash":      sound_v2.Oneshot("sfx/dash.ogg", 2),
		"playerSlam":      sound_v2.Oneshot("sfx/slam.ogg", 2),
		"playerDeath":     sound_v2.Oneshot("sfx/death.ogg", 2),
		"slamboxLand":     sound_v2.Oneshot("sfx/slambox-land.ogg", 2),
	}

	DSPChannelNames := []string{
		// Volume-only channels
		"sfxMaster",
		"musicMaster",
	}

	go assetloader.LoadAll(l.loadFinishedChan)
	go sound_v2.SoundServer(soundCatalogue, DSPChannelNames)
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
