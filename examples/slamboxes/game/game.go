package game

import (
	"errors"
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/audiocontext"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/camera"
	"mask_of_the_tomb/internal/libraries/gamestate"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"
	"mask_of_the_tomb/internal/plugins/musicplayer"
	"mask_of_the_tomb/internal/plugins/player"
	"mask_of_the_tomb/internal/plugins/world"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	gameEntryDirection = maths.DirDown
)

var (
	ErrTerminated = errors.New("Terminated")
	InitLevelName = "Level_0"
	SaveProfile   int
)

var (
	loadFinishedChan  = make(chan int)
	mainMenuPath      = filepath.Join("assets", "menus", "game", "mainmenu.yaml")
	pauseMenuPath     = filepath.Join("assets", "menus", "game", "pausemenu.yaml")
	hudPath           = filepath.Join("assets", "menus", "game", "hud.yaml")
	loadingScreenPath = filepath.Join("assets", "menus", "game", "loadingscreen.yaml")
	optionsMenuPath   = filepath.Join("assets", "menus", "game", "options.yaml")
	introScreenPath   = filepath.Join("assets", "menus", "game", "intro.yaml")
	LDTKMapPath       = filepath.Join("assets", "LDTK", "slambox.ldtk")
)

type Game struct {
	State       gamestate.GameState
	player      *player.Player
	world       *world.World
	gameUI      *ui.UI
	musicPlayer *musicplayer.MusicPlayer
	// Events
	// Listeners
	deathEffectEnterListener *events.EventListener
	playerMoveListener       *events.EventListener
}

func (g *Game) Load() {
	g.gameUI.LoadPreamble(loadingScreenPath)
	g.world.Load(LDTKMapPath)
	g.player.CreateAssets()
	g.gameUI.Load(mainMenuPath, pauseMenuPath, hudPath, optionsMenuPath, introScreenPath)

	go assetloader.LoadAll(loadFinishedChan)
}

func (g *Game) Init() {
	g.gameUI.SwitchActiveDisplay("mainmenu")
	g.musicPlayer = musicplayer.NewMusicPlayer(audiocontext.Current().Context)
}

func (g *Game) Update() error {
	events.Update()
	g.gameUI.Update()
	camera.Update()

	ebitenutil.DebugPrint(rendering.ScreenLayers.Overlay, fmt.Sprintf("TPS: %0.2f \nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

	var err error
	switch g.State.S {
	case gamestate.Loading:
		if _, done := threads.Poll(loadFinishedChan); done {
			fmt.Println("Finished loading stage")
			g.Init()
			g.world.Init(InitLevelName, save.GameData{""})
			spawnX, spawnY := g.world.ActiveLevel.GetResetPoint()
			g.player.Init(spawnX, spawnY, maths.DirNone)
			playerWidth, playerHeight := g.player.GetSize()

			w, h := g.world.ActiveLevel.GetBounds()
			camera.Init(
				w, h,
				(rendering.GAME_WIDTH-playerWidth)/2,
				(rendering.GAME_HEIGHT-playerHeight)/2,
			)
			g.gameUI.SwitchActiveDisplay("hud")
			g.State.S = gamestate.Playing
		}
	case gamestate.Playing:
		resources.Time += 0.016666
		// TODO: Center everything in some updateContext?
		velX, velY := g.player.GetMovementSize()
		posX, posY := g.player.GetPosCentered()
		g.world.Update(posX, posY, velX, velY)
		err = g.updateGameplay()
		if err != nil {
			return err
		}

		g.player.Update()
	}

	return nil
}

// TODO: Revamp and shorten this function, move more logic into player

// Q: How could this logic be expressed if entities were more separated?
func (g *Game) updateGameplay() error {
	if eventInfo, ok := g.playerMoveListener.Poll(); ok {
		moveDir := eventInfo.Data.(maths.Direction)
		slambox := g.world.ActiveLevel.GetSlamboxHit(g.player.GetHitbox(), moveDir)
		if slambox != nil {
			g.player.StartSlamming(moveDir)
			slambox.StartSlam(moveDir, &g.world.ActiveLevel.TilemapCollider, g.world.ActiveLevel.GetDisconnectedColliders(slambox))
		} else {
			newRect, _ := g.world.ActiveLevel.TilemapCollider.ProjectRect(g.player.GetHitbox(), moveDir, g.world.ActiveLevel.GetSlamboxColliders())
			if newRect != *g.player.GetHitbox() {
				g.player.Dash(moveDir, newRect.Left(), newRect.Top())
			}
		}
	}

	doorOverlap, levelIid, doorEntityIid := g.world.ActiveLevel.CheckDoorOverlap(g.player.GetHitbox())
	if g.player.GetLevelSwapInput() && doorOverlap && !g.player.Disabled {
		newBiome := errs.Must(world.ChangeActiveLevel(g.world, levelIid, doorEntityIid))
		if newBiome != "" {
			titleCardOverlay := g.gameUI.GetOverlay("titlecard")
			titleCard, ok := titleCardOverlay.OverlayContent.(*ui.TitleCard)
			if !ok {
				panic("Shit and piss")
			}
			titleCard.ChangeText(newBiome)
			titleCardOverlay.StartFadeIn()
		}
		camera.SetBorders(g.world.ActiveLevel.GetBounds())
		g.player.SetHitboxPos(g.world.ActiveLevel.GetResetPoint())
	}

	restartPrompted := inpututil.IsKeyJustReleased(ebiten.KeyR)
	wasHit := g.world.ActiveLevel.GetHazardHit(g.player.GetHitbox())
	if wasHit && !g.player.Disabled || restartPrompted {
		g.player.Die()
		screenFade := g.gameUI.GetOverlay("screenfade")
		screenFade.StartFadeIn()
	}

	_, raised := g.deathEffectEnterListener.Poll()
	if raised {
		posX, posY := g.world.ResetActiveLevel()
		g.player.SetPos(posX, posY)
		g.player.Respawn()

		screenFade := g.gameUI.GetOverlay("screenfade")
		screenFade.StartFadeOut()
	}

	titlecard := g.gameUI.GetOverlay("titlecard")
	if titlecard.IdleTime > 2 {
		fmt.Println(titlecard.IdleTime)
		titlecard.StartFadeOut()
	}

	g.player.Update()

	camera.SetPos(g.player.GetPosCentered())

	return nil
}

func (g *Game) Draw() {
	g.gameUI.Draw()
	switch g.State.S {
	case gamestate.Loading:
	case gamestate.MainMenu:
	case gamestate.Playing:
		pX, pY := g.player.GetPosCentered()
		cX, cY := camera.GetPos()
		drawCtx := rendering.Ctx{
			CamX:    cX,
			CamY:    cY,
			PlayerX: pX,
			PlayerY: pY,
		}

		g.world.ActiveLevel.Draw(drawCtx)
		g.player.Draw(rendering.WithLayer(drawCtx, rendering.ScreenLayers.Playerspace))
	case gamestate.Paused:
		pX, pY := g.player.GetPosCentered()
		cX, cY := camera.GetPos()
		drawCtx := rendering.Ctx{
			CamX:    cX,
			CamY:    cY,
			PlayerX: pX,
			PlayerY: pY,
		}

		g.world.ActiveLevel.Draw(drawCtx)
		g.player.Draw(rendering.WithLayer(drawCtx, rendering.ScreenLayers.Playerspace))
	}
}

func NewGame() *Game {
	game := &Game{
		player: player.NewPlayer(),
		world:  world.NewWorld(),
		gameUI: ui.NewUI(),
	}

	screenFade := game.gameUI.GetOverlay("screenfade")
	game.deathEffectEnterListener = events.NewEventListener(screenFade.OnFinishEnter)
	game.playerMoveListener = events.NewEventListener(game.player.OnMove)
	return game
}
