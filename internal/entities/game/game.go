package game

import (
	"errors"
	"fmt"
	"math"
	"time"

	"mask_of_the_tomb/internal/engine/advertisers"
	"mask_of_the_tomb/internal/engine/entities"
	ui "mask_of_the_tomb/internal/game/UI"
	"mask_of_the_tomb/internal/game/camera"
	"mask_of_the_tomb/internal/game/player"
	"mask_of_the_tomb/internal/game/rendering"
	save "mask_of_the_tomb/internal/game/savesystem"
	"mask_of_the_tomb/internal/game/world"
	"mask_of_the_tomb/internal/maths"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	GameEntityName = "Game"
)

var (
	ErrTerminated = errors.New("Terminated")
	_game         = Game{}
)

type Game struct {
	advertiser gameAdvertiser
	state      GameState
}

func Init() {
	_game.state = StateMainMenu
	entities.RegisterEntity(&_game, GameEntityName)
	advertisers.RegisterAdvertiser(&_game.advertiser, GameEntityName)

	// Q: Should we init a bunch of stuff from here?
	// A: Maybe, this could be an easy way to get parent-child hierarchies
	save.GlobalSave.LoadGame()

	worldInit := world.Init()

	playerInit := player.Init(worldInit.SpawnX, worldInit.SpawnY)
	ui.Init()

	camera.Init(
		worldInit.LevelWidth,
		worldInit.LevelHeight,
		(rendering.GameWidth-playerInit.Width)/2,
		(rendering.GameHeight-playerInit.Height)/2,
	)
}

var (
	slamming       = false
	slamFinishChan = make(chan int, 1)
)

// Design goal: switching on the global state should not be needed in every update
// function, as this (død)
func (g *Game) Update() {
	entities.UpdateEntities()

	var err error
	switch State {
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
	// TODO: Rewrite with events
	playerMove := g.player.InputBuffer.Read()
	if slamming {
		select {
		case <-slamFinishChan:
			slamming = false
		default:
		}
	}

	if playerMove != maths.DirNone && g.player.CanMove() && !g.player.Disabled {
		g.player.InputBuffer.Clear()
		slambox := g.world.ActiveLevel.GetSlamboxHit(g.player.Hitbox, playerMove)
		if slambox != nil {
			g.player.StartSlamming(playerMove)
			if !slamming {
				slamming = true
				go g.DoSlam(slambox, playerMove)
			}
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

			camera.GlobalCamera.SetBorders(g.world.ActiveLevel.GetLevelBounds())

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
	if damage > 0 && !g.player.Invincible && !g.player.Disabled {
		g.player.TakeDamage(damage)
	}

	g.player.Update()

	camera.GlobalCamera.SetPos(g.player.PosX, g.player.PosY)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		State = StatePaused
		g.gameUI.SwitchActiveMenu(ui.Pausemenu)
	}

	return nil
}

func (g *Game) DoSlam(slambox *world.Slambox, playerMove maths.Direction) {
	// holy fuCKNING SHIT IT WORKSsss!!!!!!!!!!!!!!!
	// YEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAH
	// FUCK
	// YEAH
	// TODO: rewrite a little bit as this is not very beautiful (function is
	// over 100 lines long)
	// Some small bugs hehe!

	// TODO: event
	time.Sleep(500 * time.Millisecond)

	projectedSlamboxRect, dist := g.world.ActiveLevel.TilemapCollider.ProjectRect(
		&slambox.GetCollider().Rect,
		playerMove,
		g.world.ActiveLevel.DisconnectedColliders(slambox),
	)
	shortestDist := dist

	for _, otherSlambox := range slambox.ConnectedBoxes {
		_, otherDist := g.world.ActiveLevel.TilemapCollider.ProjectRect(
			&otherSlambox.GetCollider().Rect,
			playerMove,
			g.world.ActiveLevel.DisconnectedColliders(otherSlambox),
		)

		if math.Abs(otherDist) < math.Abs(dist) {
			shortestDist = otherDist
		}
	}

	for _, otherSlambox := range slambox.ConnectedBoxes {
		otherProjRect, _dist := g.world.ActiveLevel.TilemapCollider.ProjectRect(
			&otherSlambox.GetCollider().Rect,
			playerMove,
			g.world.ActiveLevel.DisconnectedColliders(otherSlambox),
		)

		offset := _dist - shortestDist

		switch playerMove {
		case maths.DirUp:
			otherProjRect.SetPos(otherSlambox.Collider.Left(), otherProjRect.Top()+offset)
		case maths.DirDown:
			otherProjRect.SetPos(otherSlambox.Collider.Left(), otherProjRect.Top()-offset)
		case maths.DirRight:
			otherProjRect.SetPos(otherProjRect.Left()-offset, otherSlambox.Collider.Top())
		case maths.DirLeft:
			otherProjRect.SetPos(otherProjRect.Left()+offset, otherSlambox.Collider.Top())
		}
		otherSlambox.SetTarget(otherProjRect.Left(), otherProjRect.Top())
	}

	offset := math.Abs(dist - shortestDist)

	switch playerMove {
	case maths.DirUp:
		projectedSlamboxRect.SetPos(slambox.Collider.Left(), projectedSlamboxRect.Top()+offset)
	case maths.DirDown:
		projectedSlamboxRect.SetPos(slambox.Collider.Left(), projectedSlamboxRect.Top()-offset)
	case maths.DirRight:
		projectedSlamboxRect.SetPos(projectedSlamboxRect.Left()-offset, slambox.Collider.Top())
	case maths.DirLeft:
		projectedSlamboxRect.SetPos(projectedSlamboxRect.Left()+offset, slambox.Collider.Top())
	}

	slambox.SetTarget(projectedSlamboxRect.Left(), projectedSlamboxRect.Top())
	slamFinishChan <- 1
}

func (g *Game) Draw() {
	g.gameUI.Draw()
	switch State {
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
	g.Init()
	State = StatePlaying
	g.gameUI.SwitchActiveMenu(ui.Hud)
}

func NewGame() *Game {
	return &Game{
		player: player.NewPlayer(),
		// world:  &world.World{},
		gameUI: ui.NewUI(),
	}
}
