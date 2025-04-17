package game

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	ui "mask_of_the_tomb/internal/game/UI"
	"mask_of_the_tomb/internal/game/UI/fonts"
	"mask_of_the_tomb/internal/game/core/assetloader"
	"mask_of_the_tomb/internal/game/core/assetloader/delayasset"
	"mask_of_the_tomb/internal/game/core/events"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/core/rendering/camera"
	save "mask_of_the_tomb/internal/game/core/savesystem"
	"mask_of_the_tomb/internal/game/gamestate"
	"mask_of_the_tomb/internal/game/player"
	"mask_of_the_tomb/internal/game/sound"
	"mask_of_the_tomb/internal/game/world"
	"mask_of_the_tomb/internal/maths"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	ErrTerminated = errors.New("Terminated")
)

var (
	loadFinishedChan  = make(chan int)
	delayAsset        = delayasset.NewDelayAsset(time.Second)
	mainMenuPath      = filepath.Join("assets", "menus", "game", "mainmenu.yaml")
	pauseMenuPath     = filepath.Join("assets", "menus", "game", "pausemenu.yaml")
	hudPath           = filepath.Join("assets", "menus", "game", "hud.yaml")
	loadingScreenPath = filepath.Join("assets", "menus", "game", "loadingscreen.yaml")
	optionsMenuPath   = filepath.Join("assets", "menus", "game", "options.yaml")
)

type Game struct {
	State       gamestate.State
	player      *player.Player
	world       *world.World
	gameUI      *ui.UI
	musicPlayer *sound.MusicPlayer
	// Events
	// Listeners
	deathEffectEnterListener *events.EventListener
	playerDeathListener      *events.EventListener
	playerMoveListener       *events.EventListener
}

func (g *Game) Load() {
	g.gameUI.LoadPreamble(loadingScreenPath)
	assetloader.AddAsset(&delayAsset)
	save.GlobalSave.LoadGame()
	g.player.CreateAssets()
	g.world.Load()
	g.gameUI.Load(mainMenuPath, pauseMenuPath, hudPath, optionsMenuPath)
	fonts.Load()

	go assetloader.LoadAll(loadFinishedChan)
}

func (g *Game) Init() {
	g.world.Init()
	g.player.Init(g.world.ActiveLevel.GetSpawnPoint())
	width, height := g.world.ActiveLevel.GetLevelBounds()
	playerWidth, playerHeight := g.player.GetSize()

	camera.Init(
		width,
		height,
		(rendering.GameWidth-playerWidth)/2,
		(rendering.GameHeight-playerHeight)/2,
	)

	g.State = gamestate.MainMenu
	g.gameUI.SwitchActiveDisplay("mainmenu")

	g.musicPlayer = sound.NewMusicPlayer(audio.CurrentContext())
}

// Design goal: switching on the global state should not be needed in every update
// function, as this (død)
// TODO: Create an UpdateUI() function
func (g *Game) Update() error {
	events.Update()
	confirmations := g.gameUI.GetConfirmations()
	g.gameUI.Update()

	biome := ""
	if g.world.ActiveLevel != nil {
		biome = g.world.ActiveLevel.GetBiome()
	}

	g.musicPlayer.Update(g.State, biome)
	ebitenutil.DebugPrint(rendering.RenderLayers.Overlay, fmt.Sprintf("TPS: %0.2f \nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

	var err error
	switch g.State {
	case gamestate.Loading:
		select {
		case <-loadFinishedChan:
			fmt.Println("Finished loading")
			g.Init()
		default:
		}
	case gamestate.MainMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.gameUI.SwitchActiveDisplay("mainmenu")
		}

		if val, ok := confirmations["Play"]; ok && val {
			g.EnterPlayMode()
		} else if val, ok := confirmations["Quit"]; ok && val {
			return ErrTerminated
		} else if val, ok := confirmations["Options"]; ok && val {
			g.gameUI.SwitchActiveDisplay("options")
		}
	case gamestate.Playing:
		g.world.Update()
		err = g.updateGameplay()
		if err != nil {
			return err
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State = gamestate.Paused
			g.gameUI.SwitchActiveDisplay("pausemenu")
		}

		g.player.Update()
	case gamestate.Paused:
		// TODO: Make esc go back to playing state, and implement back button
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.gameUI.SwitchActiveDisplay("pausemenu")
		}

		if val, ok := confirmations["Resume"]; ok && val {
			g.EnterPlayMode()
		} else if val, ok := confirmations["Quit"]; ok && val {
			// Save data and stuff
			// Loading screens
			// etc
			g.State = gamestate.MainMenu
			g.gameUI.SwitchActiveDisplay("mainmenu")
		} else if val, ok := confirmations["Options"]; ok && val {
			g.gameUI.SwitchActiveDisplay("options")
		}
	}

	return nil
}

// TODO: Revamp and shorten this function, move more logic into player
func (g *Game) updateGameplay() error {
	if eventInfo, ok := g.playerMoveListener.Poll(); ok {
		moveDir := eventInfo.Data.(maths.Direction)
		slambox := g.world.ActiveLevel.GetSlamboxHit(g.player.GetHitbox(), moveDir)
		if slambox != nil {
			g.player.StartSlamming(moveDir)
			slambox.DoSlam(moveDir, &g.world.ActiveLevel.TilemapCollider, g.world.ActiveLevel.GetDisconnectedColliders(slambox))
		} else {
			newRect, _ := g.world.ActiveLevel.TilemapCollider.ProjectRect(g.player.GetHitbox(), moveDir, g.world.ActiveLevel.GetSlamboxColliders())
			if newRect != *g.player.GetHitbox() {
				g.player.Dash(moveDir, newRect.Left(), newRect.Top())
			}
		}
	}

	if g.player.GetLevelSwapInput() {
		hit, levelIid, entityIid := g.world.ActiveLevel.GetDoorHit(g.player.GetHitbox())

		if hit {
			err := world.ChangeActiveLevel(g.world, levelIid)
			// TODO: On changing active level, store the slambox positions in some kind of short-term memory
			if err != nil {
				fmt.Println("Error occured when swapping to level with iid: ", levelIid)
				return err
			}

			camera.SetBorders(g.world.ActiveLevel.GetLevelBounds())

			otherSideDoor, err := g.world.ActiveLevel.GetEntityByIid(entityIid)
			if err != nil {
				fmt.Println("Didn't find the other side door, iid ", entityIid)
				return err
			}
			posX, posY := otherSideDoor.Px[0], otherSideDoor.Px[1]
			g.player.SetPos(posX, posY)
		}
	}

	restartPrompted := inpututil.IsKeyJustReleased(ebiten.KeyR)
	wasHit := g.world.ActiveLevel.GetHazardHit(g.player.GetHitbox())
	if wasHit && !g.player.Disabled || restartPrompted {
		fmt.Println("Dead")
		g.player.Die()
		g.gameUI.DeathEffect.StartEnter()
	}

	_, raised := g.deathEffectEnterListener.Poll()
	if raised {
		posX, posY := g.world.ResetActiveLevel()
		g.player.SetPos(posX, posY)
		g.player.Respawn()

		g.gameUI.DeathEffect.StartExit()
	}

	g.player.Update()

	camera.SetPos(g.player.GetPosCentered())

	return nil
}

func (g *Game) Draw() {
	g.gameUI.Draw()
	switch g.State {
	case gamestate.Loading:
	case gamestate.MainMenu:
	case gamestate.Playing:
		g.world.ActiveLevel.Draw()
		g.player.Draw()
	case gamestate.Paused:
		// TODO: Add dim and blur filter on pausing the game
		g.world.ActiveLevel.Draw()
		g.player.Draw()
	}
}

func (g *Game) EnterPlayMode() {
	g.State = gamestate.Playing
	g.gameUI.SwitchActiveDisplay("hud")
}

func NewGame() *Game {
	game := &Game{
		player: player.NewPlayer(),
		world:  world.NewWorld(),
		gameUI: ui.NewUI(),
	}

	game.playerDeathListener = events.NewEventListener(game.player.OnDeath)
	game.deathEffectEnterListener = events.NewEventListener(game.gameUI.DeathEffect.OnFinishEnter)
	game.playerMoveListener = events.NewEventListener(game.player.OnMove)
	return game
}
