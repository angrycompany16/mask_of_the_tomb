package game

import (
	"errors"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/resources"
	ui "mask_of_the_tomb/internal/plugins/UI"
	"mask_of_the_tomb/internal/plugins/musicplayer"
	"mask_of_the_tomb/internal/plugins/player"
	"mask_of_the_tomb/internal/plugins/world"
	"path/filepath"
	"time"
)

// IMPORTANT NOTE: If there is a missing entity reference, the game world will not be able to load !!!
// Instead it becomes a nil pointer...

// Some problems:
// - When exiting to main menu, music volume remains low
//   - This follows from the game not playing the proper song on pressing play
// - When opening menus for the first time, the select sound plays
// - Somehow it seems that the game itself gets darker if we don't draw the UI???

const (
	gameEntryDirection = maths.DirDown
)

var (
	ErrTerminated = errors.New("Terminatednow")
	InitLevelName string
	SaveProfile   int
	initTime      time.Time
)

var (
	loadFinishedChan  = make(chan int)
	delayAsset        = assettypes.NewDelayAsset(time.Second)
	mainMenuPath      = filepath.Join("assets", "menus", "game", "mainmenu.yaml")
	pauseMenuPath     = filepath.Join("assets", "menus", "game", "pausemenu.yaml")
	hudPath           = filepath.Join("assets", "menus", "game", "hud.yaml")
	loadingScreenPath = filepath.Join("assets", "menus", "game", "loadingscreen.yaml")
	optionsMenuPath   = filepath.Join("assets", "menus", "game", "options.yaml")
	emptyMenuPath     = filepath.Join("assets", "menus", "game", "empty.yaml")
	introScreenPath   = filepath.Join("assets", "menus", "game", "intro.yaml")
	LDTKMapPath       = filepath.Join("assets", "LDTK", "world.ldtk")
)

type Game struct {
	player      *player.Player
	world       *world.World
	mainUI      *ui.UI
	gameplayUI  *ui.UI
	musicPlayer *musicplayer.MusicPlayer

	introDashTimer *time.Timer

	deathEffectEnterListener *events.EventListener
	playerMoveListener       *events.EventListener
	titleCardTimeoutListener *events.EventListener
	levelCardTimeoutListener *events.EventListener
}

func (g *Game) Update() error {
	var err error
	switch resources.State {
	case resources.Loading:
		g.LoadingStageUpdate()
	case resources.MainMenu:
		err = g.MenuStageUpdate()
	case resources.Intro:
		g.IntroStageUpdate()
	case resources.Playing:
		g.GameplayStageUpdate()
	case resources.Paused:
		g.PausedStageUpdate()
	}
	return err
}

func (g *Game) Draw() error {
	var err error
	switch resources.State {
	case resources.Loading:
		g.LoadingStageDraw()
	case resources.MainMenu:
		g.MenuStageDraw()
	case resources.Intro:
		g.IntroStageDraw()
	case resources.Playing:
		g.GameplayStageDraw()
	case resources.Paused:
		g.PausedStageDraw()
	}
	return err
}

func NewGame() *Game {
	game := &Game{
		player:      player.NewPlayer(),
		world:       world.NewWorld(),
		mainUI:      ui.NewUI(make(map[string]*ui.Overlay)),
		gameplayUI:  ui.NewUI(make(map[string]*ui.Overlay)),
		musicPlayer: musicplayer.NewMusicPlayer(),
	}

	game.playerMoveListener = events.NewEventListener(game.player.OnMove)

	return game
}
