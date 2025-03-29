package game

import (
	"errors"
	"fmt"
	"time"

	ui "mask_of_the_tomb/internal/game/UI"
	"mask_of_the_tomb/internal/game/core/assetloader"
	"mask_of_the_tomb/internal/game/core/events"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/core/rendering/camera"
	save "mask_of_the_tomb/internal/game/core/savesystem"
	"mask_of_the_tomb/internal/game/player"
	"mask_of_the_tomb/internal/game/world"
	"mask_of_the_tomb/internal/maths"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	ErrTerminated = errors.New("Terminated")
)

type GameState int

const (
	StateLoading GameState = iota
	StateMainMenu
	StatePlaying
	StatePaused
)

var (
	State            GameState
	loadFinishedChan = make(chan int)
	delayAsset       = assetloader.NewDelayAsset(time.Second)
)

type Game struct {
	player *player.Player
	world  *world.World
	gameUI *ui.UI
	// Events
	// Listeners
	deathEffectEnterListener *events.EventListener
	playerDeathListener      *events.EventListener
}

func (g *Game) Load() {
	assetloader.AddAsset(&delayAsset)
	save.GlobalSave.LoadGame()
	g.player.Load()
	g.world.Load()

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
	g.gameUI.SetScore(0)

	State = StateMainMenu
	g.gameUI.SwitchActiveMenu(ui.Mainmenu)
}

// Design goal: switching on the global state should not be needed in every update
// function, as this (død)
func (g *Game) Update() error {
	events.Update()
	confirmations := g.gameUI.GetConfirmations()
	g.gameUI.Update()
	var err error
	switch State {
	case StateLoading:
		select {
		case <-loadFinishedChan:
			fmt.Println("Finished loading")
			g.Init()
		default:
		}
	case StateMainMenu:
		if val, ok := confirmations["Play"]; ok && val {
			g.EnterPlayMode()
		} else if val, ok := confirmations["Quit"]; ok && val {
			return ErrTerminated
		}
	case StatePlaying:
		g.world.Update()
		err = g.updateGameplay()
		if err != nil {
			return err
		}
		g.player.Update()
	case StatePaused:
		if val, ok := confirmations["Resume"]; ok && val {
			State = StatePlaying
			g.gameUI.SwitchActiveMenu(ui.Hud)
		} else if val, ok := confirmations["Quit"]; ok && val {
			// Save data and stuff
			// Loading screens
			// etc
			State = StateMainMenu
			g.gameUI.SwitchActiveMenu(ui.Mainmenu)
		}
	}
	return nil
}

func (g *Game) updateGameplay() error {
	playerMove := g.player.InputBuffer.Read()

	// TODO: Change this by listening to player move event?
	if playerMove != maths.DirNone && g.player.CanMove() && !g.player.Disabled {
		g.player.InputBuffer.Clear()
		slambox := g.world.ActiveLevel.GetSlamboxHit(g.player.Hitbox, playerMove)
		if slambox != nil {
			g.player.StartSlamming(playerMove)
			slambox.DoSlam(playerMove, &g.world.ActiveLevel.TilemapCollider, g.world.ActiveLevel.DisconnectedColliders(slambox))
		} else {
			newRect, _ := g.world.ActiveLevel.TilemapCollider.ProjectRect(g.player.Hitbox, playerMove, g.world.ActiveLevel.GetSlamboxColliders())
			if newRect != *g.player.Hitbox {
				g.player.EnterDashAnim()
				g.player.SetRot(playerMove)
				g.player.SetTarget(newRect.Left(), newRect.Top())
				g.player.State = player.Moving
			}
		}
	}

	if g.player.GetLevelSwapInput() {
		hit, levelIid, entityIid := g.world.ActiveLevel.GetDoorHit(g.player.Hitbox)

		if hit {
			err := world.ChangeActiveLevel(g.world, levelIid)
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

	damage := g.world.ActiveLevel.GetHazardHit(g.player.Hitbox)
	if damage && !g.player.Disabled {
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

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		State = StatePaused
		g.gameUI.SwitchActiveMenu(ui.Pausemenu)
	}

	return nil
}

func (g *Game) Draw() {
	g.gameUI.Draw()
	switch State {
	case StateLoading:
	case StateMainMenu:
	case StatePlaying:
		g.world.ActiveLevel.Draw()
		g.player.Draw()
	case StatePaused:
		g.world.ActiveLevel.Draw()
		g.player.Draw()
	}
}

func (g *Game) EnterPlayMode() {
	State = StatePlaying
	g.gameUI.SwitchActiveMenu(ui.Hud)
}

func NewGame() *Game {
	game := &Game{
		player: player.NewPlayer(),
		world:  &world.World{},
		gameUI: ui.NewUI(),
	}

	game.playerDeathListener = events.NewEventListener(game.player.OnDeath)
	game.deathEffectEnterListener = events.NewEventListener(game.gameUI.DeathEffect.OnFinishEnter)
	return game
}
