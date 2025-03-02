package game

import (
	"fmt"

	ui "mask_of_the_tomb/game/UI"
	"mask_of_the_tomb/game/camera"
	"mask_of_the_tomb/game/physics"
	"mask_of_the_tomb/game/player"
	"mask_of_the_tomb/game/save"
	"mask_of_the_tomb/game/world"
	"mask_of_the_tomb/rendering"
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	player *player.Player
	world  *world.World
	gameUI *ui.UI
}

func (g *Game) Init() {
	save.GlobalSave.LoadGame()
	g.world.Init()
	g.player.Init(g.world.ActiveLevel.GetSpawnPoint())
	width, height := g.world.ActiveLevel.GetLevelBounds()
	playerWidth, playerHeight := g.player.GetSize()

	camera.GlobalCamera.Init(
		width,
		height,
		(rendering.GameWidth-playerWidth)/2,
		(rendering.GameHeight-playerHeight)/2,
	)
	g.gameUI.SetScore(0)
}

// Design goal: switching on the global state should not be needed in every update
// function, as this (død)
func (g *Game) Update() error {
	confirmations := g.gameUI.GetConfirmations()
	g.gameUI.Update()
	var err error
	switch utils.GlobalState {
	case utils.StateMainMenu:
		if val, ok := confirmations["Play"]; ok && val {
			g.EnterPlayMode()
		} else if val, ok := confirmations["Quit"]; ok && val {
			return utils.Terminated
		}
	case utils.StatePlaying:
		g.world.Update()
		err = g.updateGameplay()
		if err != nil {
			return err
		}
	case utils.StatePaused:
		if val, ok := confirmations["Resume"]; ok && val {
			utils.GlobalState = utils.StatePlaying
			g.gameUI.SwitchActiveMenu(ui.Hud)
		} else if val, ok := confirmations["Quit"]; ok && val {
			// Save data and stuff
			// Loading screens
			// etc
			utils.GlobalState = utils.StateMainMenu
			g.gameUI.SwitchActiveMenu(ui.Mainmenu)
		}
	}
	return nil
}

// TODO: Rewrite the physics system with somem kind of interface or something
// Basically allow the individual objects to define how collisions should be handled.
// TODO: Make something like list comprehension?
func (g *Game) updateGameplay() error {
	playerMove := g.player.GetBufferedInput()
	playerX, playerY := g.player.GetPos()
	if playerMove != utils.DirNone && !g.player.IsMoving() && !g.player.IsDisabled() {
		g.player.DeleteBuffer()
		slambox := g.world.ActiveLevel.GetSlamboxHit(g.player.GetHitbox(), playerMove)

		if slambox != nil {
			otherHitboxes := make([]*physics.RectCollider, 0)
			for _, slambox := range g.world.ActiveLevel.Without(slambox) {
				otherHitboxes = append(otherHitboxes, slambox.GetCollider())
			}

			newSlamboxRect := g.world.ActiveLevel.TilemapCollider.ProjectRect(&slambox.GetCollider().Rect, playerMove, otherHitboxes)

			slambox.SetTarget(newSlamboxRect.Left(), newSlamboxRect.Top())
		}

		newRect := g.world.ActiveLevel.TilemapCollider.ProjectRect(g.player.GetHitbox(), playerMove, g.world.ActiveLevel.GetSlamboxColliders())
		g.player.SetTarget(newRect.Left(), newRect.Top())
	}

	if g.player.GetLevelSwapInput() {
		hit, levelIid, entityIid := g.world.ActiveLevel.GetDoorHit(g.player.GetHitbox())

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

	dx, dy := g.player.GetMovementSize()
	collectibleOverlapCount := g.world.ActiveLevel.GetCollectibleHit(playerX, playerY, dx, dy)

	if collectibleOverlapCount > 0 {
		g.player.SetScore(g.player.GetScore() + collectibleOverlapCount)
		g.gameUI.SetScore(g.player.GetScore())
	}

	damage := g.world.ActiveLevel.GetHazardHit(playerX, playerY)
	if damage > 0 && !g.player.IsInvincible() && !g.player.IsDisabled() {
		g.player.TakeDamage(damage)
	}

	g.player.Update()

	playerX, playerY = g.player.GetPos()
	camera.GlobalCamera.SetPos(playerX, playerY)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		utils.GlobalState = utils.StatePaused
		g.gameUI.SwitchActiveMenu(ui.Pausemenu)
	}

	return nil
}

func (g *Game) Draw() {
	g.gameUI.Draw()
	switch utils.GlobalState {
	case utils.StateMainMenu:
	case utils.StatePlaying:
		g.world.ActiveLevel.Draw()
		g.player.Draw()
	case utils.StatePaused:
		g.world.ActiveLevel.Draw()
		g.player.Draw()
	}
}

func (g *Game) EnterPlayMode() {
	g.Init()
	utils.GlobalState = utils.StatePlaying
	g.gameUI.SwitchActiveMenu(ui.Hud)
}

func NewGame() *Game {
	return &Game{
		player: player.NewPlayer(),
		world:  &world.World{},
		gameUI: ui.NewUI(),
	}
}
