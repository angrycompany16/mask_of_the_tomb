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
	assetloader.Add("menuTheme", assettypes.MakeAudioStreamAsset(assets.Menu_mp3, assettypes.Mp3))
	assetloader.Add("basementTheme", assettypes.MakeAudioStreamAsset(assets.Basement_wav, assettypes.Wav))
	assetloader.Add("libraryTheme", assettypes.MakeAudioStreamAsset(assets.Library_mp3, assettypes.Mp3))
	assetloader.Add("selectSound", assettypes.MakeAudioStreamAsset(assets.Select_ogg, assettypes.Ogg))
	assetloader.Add("dialogueSound", assettypes.MakeAudioStreamAsset(assets.Text_scroll_ogg, assettypes.Ogg))
	assetloader.Add("vowelA", assettypes.MakeAudioStreamAsset(assets.Vowel_A_wav, assettypes.Wav))
	assetloader.Add("vowelE", assettypes.MakeAudioStreamAsset(assets.Vowel_E_wav, assettypes.Wav))
	assetloader.Add("vowelI", assettypes.MakeAudioStreamAsset(assets.Vowel_I_wav, assettypes.Wav))
	assetloader.Add("vowelO", assettypes.MakeAudioStreamAsset(assets.Vowel_O_wav, assettypes.Wav))
	assetloader.Add("vowelU", assettypes.MakeAudioStreamAsset(assets.Vowel_U_wav, assettypes.Wav))
	assetloader.Add("constD", assettypes.MakeAudioStreamAsset(assets.Const_D_wav, assettypes.Wav))
	assetloader.Add("constF", assettypes.MakeAudioStreamAsset(assets.Const_F_wav, assettypes.Wav))
	assetloader.Add("constG", assettypes.MakeAudioStreamAsset(assets.Const_G_wav, assettypes.Wav))
	assetloader.Add("constK", assettypes.MakeAudioStreamAsset(assets.Const_K_wav, assettypes.Wav))
	assetloader.Add("constL", assettypes.MakeAudioStreamAsset(assets.Const_L_wav, assettypes.Wav))
	assetloader.Add("constM", assettypes.MakeAudioStreamAsset(assets.Const_M_wav, assettypes.Wav))
	assetloader.Add("constP", assettypes.MakeAudioStreamAsset(assets.Const_P_wav, assettypes.Wav))
	assetloader.Add("constS", assettypes.MakeAudioStreamAsset(assets.Const_S_wav, assettypes.Wav))
	assetloader.Add("constT", assettypes.MakeAudioStreamAsset(assets.Const_T_wav, assettypes.Wav))
	assetloader.Add("constX", assettypes.MakeAudioStreamAsset(assets.Const_X_wav, assettypes.Wav))

	// NO
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
