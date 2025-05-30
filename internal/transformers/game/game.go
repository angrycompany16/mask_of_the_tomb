package game

import (
	"errors"
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/sound"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/assettypes"
	"mask_of_the_tomb/internal/libraries/camera"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"
	"mask_of_the_tomb/internal/plugins/musicplayer"
	"mask_of_the_tomb/internal/plugins/player"
	"mask_of_the_tomb/internal/plugins/world"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
	player     *player.Player
	world      *world.World
	menuUI     *ui.UI
	gameplayUI *ui.UI

	musicPlayer *musicplayer.MusicPlayer
	// Events
	// Listeners
	deathEffectEnterListener *events.EventListener
	playerMoveListener       *events.EventListener
}

func (g *Game) InitLoad() {
	g.menuUI.LoadPreamble(loadingScreenPath)
	assetloader.Load("any", &delayAsset)
	g.world.Load(LDTKMapPath)
	g.player.CreateAssets()
	g.menuUI.Load(mainMenuPath, pauseMenuPath, optionsMenuPath, introScreenPath, emptyMenuPath)
	g.gameplayUI.Load(hudPath)

	go assetloader.LoadAll(loadFinishedChan)
}

func (g *Game) InitMenu() {
	initTime = time.Now()
	g.menuUI.SwitchActiveDisplay("mainmenu")
	g.gameplayUI.SwitchActiveDisplay("hud")
	g.musicPlayer = musicplayer.NewMusicPlayer(sound.GetCurrentAudioContext().Context)
}

func (g *Game) PreloadUpdate() {
	if _, done := threads.Poll(loadFinishedChan); done {
		fmt.Println("Finished loading stage")
		g.InitMenu()
		resources.State = resources.MainMenu
	}
}

func (g *Game) Update() error {
	events.Update()
	confirmations := g.menuUI.GetConfirmations()
	g.menuUI.Update()
	g.gameplayUI.Update()

	titlecard := g.gameplayUI.GetOverlay("titlecard")
	if titlecard.IdleTime > 2 {
		// TODO: Rewrite with timer
		fmt.Println(titlecard.IdleTime)
		titlecard.StartFadeOut()
	}

	camera.Update()

	biome := ""
	if g.world.ActiveLevel != nil {
		biome = g.world.ActiveLevel.GetBiome()
	}

	g.musicPlayer.Update(biome)
	ebitenutil.DebugPrint(rendering.ScreenLayers.Overlay, fmt.Sprintf("TPS: %0.2f \nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

	// Problem: The lifetime of the objects is not obvious with the current setup
	// Instead of using switch to represent multiple stages, maybe we could make the functions
	// layered sort of like in a flame graph?
	var err error
	switch resources.State {
	case resources.MainMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.menuUI.SwitchActiveDisplay("mainmenu")
		}

		if val, ok := confirmations["Play"]; ok && val {
			gameData := save.LoadGame(SaveProfile)
			g.world.Init(InitLevelName, gameData)
			if gameData.SpawnRoomName == "" {
				resources.State = resources.Intro
				return nil
			}
			// TODO: Convert to single function
			spawnX, spawnY := g.world.ActiveLevel.GetResetPoint()
			g.player.Init(spawnX, spawnY, maths.DirNone)
			playerWidth, playerHeight := g.player.GetSize()

			w, h := g.world.ActiveLevel.GetBounds()
			camera.Init(
				w, h,
				(rendering.GAME_WIDTH-playerWidth)/2,
				(rendering.GAME_HEIGHT-playerHeight)/2,
			)
			g.menuUI.SwitchActiveDisplay("empty")
			resources.State = resources.Playing
		} else if val, ok := confirmations["Quit"]; ok && val {
			return ErrTerminated
		} else if val, ok := confirmations["Options"]; ok && val {
			g.menuUI.SwitchActiveDisplay("options")
		}
	case resources.Intro:
		g.menuUI.SwitchActiveDisplay("intro")
		if val, ok := confirmations["Introtext"]; ok && val {
			spawnX, spawnY := g.world.ActiveLevel.GetGameEntryPos()
			g.player.Init(spawnX, spawnY, gameEntryDirection)
			playerWidth, playerHeight := g.player.GetSize()

			w, h := g.world.ActiveLevel.GetBounds()
			camera.Init(
				w, h,
				(rendering.GAME_WIDTH-playerWidth)/2,
				(rendering.GAME_HEIGHT-playerHeight)/2,
			)
			g.menuUI.SwitchActiveDisplay("empty")
			resources.State = resources.Playing

			newRect, _ := g.world.ActiveLevel.TilemapCollider.ProjectRect(g.player.GetHitbox(), gameEntryDirection, g.world.ActiveLevel.GetSlamboxColliders())
			if newRect != *g.player.GetHitbox() {
				g.player.Dash(gameEntryDirection, newRect.Left(), newRect.Top())
			}
		}
	case resources.Playing:
		resources.Time = time.Since(initTime).Seconds()
		velX, velY := g.player.GetMovementSize()
		posX, posY := g.player.GetPosCentered()
		g.world.Update(posX, posY, velX, velY)
		err = g.updateGameplay()
		if err != nil {
			return err
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			resources.State = resources.Paused
			g.menuUI.SwitchActiveDisplay("pausemenu")
		}

		g.player.Update()
	case resources.Paused:
		// TODO: Make esc go back to playing state, and implement back button
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.menuUI.SwitchActiveDisplay("pausemenu")
		}

		if val, ok := confirmations["Resume"]; ok && val {
			resources.State = resources.Playing
			g.menuUI.SwitchActiveDisplay("empty")
		} else if val, ok := confirmations["Quit"]; ok && val {
			g.world.SaveLevel(g.world.ActiveLevel)
			save.SaveGame(save.GameData{
				SpawnRoomName: g.world.ActiveLevel.GetName(),
			}, SaveProfile)
			resources.State = resources.MainMenu
			g.menuUI.SwitchActiveDisplay("mainmenu")
		} else if val, ok := confirmations["Options"]; ok && val {
			g.menuUI.SwitchActiveDisplay("options")
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
			titleCardOverlay := g.gameplayUI.GetOverlay("titlecard")
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
		screenFade := g.menuUI.GetOverlay("screenfade")
		screenFade.StartFadeIn()
	}

	_, raised := g.deathEffectEnterListener.Poll()
	if raised {
		posX, posY := g.world.ResetActiveLevel()
		g.player.SetPos(posX, posY)
		g.player.Respawn()

		screenFade := g.menuUI.GetOverlay("screenfade")
		screenFade.StartFadeOut()
	}

	g.player.Update()

	camera.SetPos(g.player.GetPosCentered())

	return nil
}

func (g *Game) PreloadDraw() {
	g.menuUI.Draw()
}

func (g *Game) Draw() {
	g.menuUI.Draw()
	if resources.State == resources.MainMenu || resources.State == resources.Intro {
		return
	}

	pX, pY := g.player.GetPosCentered()
	cX, cY := camera.GetPos()
	drawCtx := rendering.Ctx{
		CamX:    cX,
		CamY:    cY,
		PlayerX: pX,
		PlayerY: pY,
	}

	// TODO: Add dim and blur filter on pausing the game
	g.player.Draw(rendering.WithLayer(drawCtx, rendering.ScreenLayers.Playerspace))
	g.world.ActiveLevel.Draw(drawCtx)

	g.gameplayUI.Draw()
	// UI is HARD-CODED to render at the UI layer...
	// I sck at programming
}

func NewGame() *Game {
	game := &Game{
		player: player.NewPlayer(),
		world:  world.NewWorld(),
		menuUI: ui.NewUI(map[string]*ui.Overlay{
			"screenfade": ui.NewOverlay(ui.NewScreenFade()),
		}),
		gameplayUI: ui.NewUI(map[string]*ui.Overlay{
			"titlecard": ui.NewOverlay(ui.NewTitleCard()),
		}),
	}

	screenFade := game.menuUI.GetOverlay("screenfade")
	game.deathEffectEnterListener = events.NewEventListener(screenFade.OnFinishEnter)
	game.playerMoveListener = events.NewEventListener(game.player.OnMove)
	return game
}
