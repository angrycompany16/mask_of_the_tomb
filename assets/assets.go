package assets

import (
	"embed"
	_ "embed"
	"path/filepath"
)

var (
	//go:embed *
	FS embed.FS
)

var (
	EnvironmentFolder = filepath.Join("assets", "sprites", "environment")
	PlayerFolder      = filepath.Join("assets", "sprites", "player")

	// Images
	//go:embed sprites/environment/slambox_tilemap.png
	Slambox_tilemap_png []byte
	//go:embed sprites/player/player.png
	Player_png []byte
	//go:embed sprites/environment/grass.png
	Grass_png []byte
	//go:embed sprites/environment/turret.png
	Turret_png []byte
	//go:embed sprites/environment/catcher.png
	Catcher_png []byte
	//go:embed sprites/environment/lantern.png
	Lantern_png []byte
	//go:embed sprites/icons/textbox-tick.png
	Textbox_tick_png []byte
	//go:embed sprites/environment/slambox-gem-blue.png
	Slambox_gem_blue_png []byte
	//go:embed sprites/environment/slambox-gem-red.png
	Slambox_gem_red_png []byte
	//go:embed sprites/UI/level-titlecard.png
	Level_titlecard_sprite []byte

	// Shaders
	//go:embed shaders/fog.kage
	Fog_kage []byte
	//go:embed shaders/vignette.kage
	Vignette_kage []byte
	//go:embed shaders/pixel_lights.kage
	Pixel_lights_kage []byte
	//go:embed shaders/death_transition.kage
	Death_transition_kage []byte
	//go:embed shaders/level_transition_enter.kage
	Level_transition_enter_kage []byte
	//go:embed shaders/level_transition_exit.kage
	Level_transition_exit_kage []byte

	// YAML assets
	// Particles
	//go:embed particlesystems/jump-tight.yaml
	Jump_tight_yaml []byte
	//go:embed particlesystems/jump-broad.yaml
	Jump_broad_yaml []byte
	//go:embed particlesystems/slambox.yaml
	Slambox_particles_yaml []byte
	//go:embed particlesystems/lemma-idle.yaml
	Lemma_idle_particles_yaml []byte
	//go:embed particlesystems/lemma-appear.yaml
	Lemma_appear_particles_yaml []byte
	//go:embed particlesystems/lemma-disappear.yaml
	Lemma_disappear_particles_yaml []byte
	//go:embed particlesystems/basement.yaml
	Basement_yaml []byte
	// Animations
	//go:embed animations/teleporter.yaml
	Teleporter_yaml []byte
	//go:embed animations/dash-init.yaml
	Dash_init_yaml []byte
	//go:embed animations/dash-loop.yaml
	Dash_loop_yaml []byte
	//go:embed animations/player-idle.yaml
	Player_idle_yaml []byte
	//go:embed animations/player-slam.yaml
	Player_slam_yaml []byte
	// Menus
	//go:embed menus/game/mainmenu.yaml
	Main_menu_yaml []byte
	//go:embed menus/game/pausemenu.yaml
	Pause_menu_yaml []byte
	//go:embed menus/game/options.yaml
	Options_yaml []byte
	//go:embed menus/game/intro.yaml
	Intro_yaml []byte
	//go:embed menus/game/empty.yaml
	Empty_yaml []byte
	//go:embed menus/game/hud.yaml
	Hud_yaml []byte

	// Fonts
	//go:embed fonts/JSE_AmigaAMOS.ttf
	JSE_AmigaAMOS_ttf []byte
	//go:embed fonts/JSE_ZXSpectrum.ttf
	JSE_ZXSpectrum_ttf []byte
	//go:embed fonts/dude.ttf
	C_AND_C_Red_Alert_ttf []byte
)
