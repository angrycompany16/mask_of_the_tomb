package game

import (
	"fmt"
	"mask_of_the_tomb/internal/engine/advertisers"
	"mask_of_the_tomb/internal/engine/events"
	ui "mask_of_the_tomb/internal/entities/UI"
	pubui "mask_of_the_tomb/internal/entities/UI/pub"
	"mask_of_the_tomb/internal/entities/camera"
	pubgame "mask_of_the_tomb/internal/entities/game/pub"
	"mask_of_the_tomb/internal/entities/player"
	"mask_of_the_tomb/internal/entities/world"
	"mask_of_the_tomb/internal/libraries/rendering"
	save "mask_of_the_tomb/internal/libraries/savesystem"
)

type Game struct {
	state          pubgame.GameState
	gameAdvertiser pubgame.GameAdvertiser
	uiListener     *events.EventListener
}

func NewGame() *Game {
	_game := Game{
		uiListener: events.NewEventListener(pubui.UISelected),
		state:      pubgame.StateMainMenu,
	}

	advertisers.RegisterAdvertiser(&_game.gameAdvertiser, pubgame.GameEntityName)

	_, initLevelInfo := world.NewWorld()
	fmt.Println("Initialized at", initLevelInfo.SpawnX, initLevelInfo.SpawnY)

	save.GlobalSave.LoadGame()

	_, initPlayerInfo := player.New(initLevelInfo.SpawnX, initLevelInfo.SpawnY)

	ui.NewUI()

	camera.New(
		initLevelInfo.LevelWidth,
		initLevelInfo.LevelHeight,
		(rendering.GameWidth-initPlayerInfo.Width)/2,
		(rendering.GameHeight-initPlayerInfo.Height)/2,
	)
	return &_game
}

func (g *Game) Update() {
	info, raised := g.uiListener.Poll()
	if !raised {
		return
	}

	uiSelect := info.Data.(pubui.UISelect)

	switch uiSelect {
	case pubui.SelectPlay:
		g.state = pubgame.StatePlaying
	case pubui.SelectOpts:
	case pubui.SelectMainMenu:
		g.state = pubgame.StateMainMenu
	case pubui.SelectQuit:
	}
}

func (g *Game) PostUpdate() {
	g.gameAdvertiser.State = g.state
}

// func (g *Game) updateGameplay() error {
// TODO: Rewrite with events
// playerMove := g.player.InputBuffer.Read()
// if slamming {
// 	select {
// 	case <-slamFinishChan:
// 		slamming = false
// 	default:
// 	}
// }

// 	// if playerMove != maths.DirNone && g.player.CanMove() && !g.player.Disabled {
// 	// 	g.player.InputBuffer.Clear()
// 	// 	slambox := g.world.ActiveLevel.GetSlamboxHit(g.player.Hitbox, playerMove)
// 	// 	if slambox != nil {
// 	// 		g.player.StartSlamming(playerMove)
// 	// 		if !slamming {
// 	// 			slamming = true
// 	// 			go g.DoSlam(slambox, playerMove)
// 	// 		}
// 	// 	} else {
// 	// 		newRect, _ := g.world.ActiveLevel.TilemapCollider.ProjectRect(g.player.Hitbox, playerMove, g.world.ActiveLevel.GetSlamboxColliders())
// 	// 		if newRect != *g.player.Hitbox {
// 	// 			g.player.EnterDashAnim()
// 	// 			g.player.SetRot(playerMove)
// 	// 			g.player.SetTarget(newRect.Left(), newRect.Top())
// 	// 			g.player.State = player.Moving
// 	// 		}
// 	// 	}
// 	// }

// 	if g.player.GetLevelSwapInput() {
// 		hit, levelIid, entityIid := g.world.ActiveLevel.GetDoorHit(g.player.Hitbox)

// 		if hit {
// 			err := world.ChangeActiveLevel(g.world, levelIid)
// 			if err != nil {
// 				fmt.Println("Error occured when swapping to level with iid: ", levelIid)
// 				return err
// 			}

// 			camera.GlobalCamera.SetBorders(g.world.ActiveLevel.GetLevelBounds())

// 			otherSideDoor, err := g.world.ActiveLevel.GetEntityByIid(entityIid)
// 			if err != nil {
// 				fmt.Println("Didn't find the other side door, iid ", entityIid)
// 				return err
// 			}
// 			posX, posY := otherSideDoor.Px[0], otherSideDoor.Px[1]
// 			g.player.SetPos(posX, posY)
// 		}
// 	}

// 	damage := g.world.ActiveLevel.GetHazardHit(g.player.Hitbox)
// 	if damage > 0 && !g.player.Invincible && !g.player.Disabled {
// 		g.player.TakeDamage(damage)
// 	}

// 	g.player.Update()

// 	camera.GlobalCamera.SetPos(g.player.PosX, g.player.PosY)

// 	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
// 		State = StatePaused
// 		g.gameUI.SwitchActiveMenu(ui.Pausemenu)
// 	}

// 	return nil
// }

// func (g *Game) DoSlam(slambox *world.Slambox, playerMove maths.Direction) {
// 	// holy fuCKNING SHIT IT WORKSsss!!!!!!!!!!!!!!!
// 	// YEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAH
// 	// FUCK
// 	// YEAH
// 	// TODO: rewrite a little bit as this is not very beautiful (function is
// 	// over 100 lines long)
// 	// Some small bugs hehe!

// 	// TODO: event
// 	time.Sleep(500 * time.Millisecond)

// 	projectedSlamboxRect, dist := g.world.ActiveLevel.TilemapCollider.ProjectRect(
// 		&slambox.GetCollider().Rect,
// 		playerMove,
// 		g.world.ActiveLevel.DisconnectedColliders(slambox),
// 	)
// 	shortestDist := dist

// 	for _, otherSlambox := range slambox.ConnectedBoxes {
// 		_, otherDist := g.world.ActiveLevel.TilemapCollider.ProjectRect(
// 			&otherSlambox.GetCollider().Rect,
// 			playerMove,
// 			g.world.ActiveLevel.DisconnectedColliders(otherSlambox),
// 		)

// 		if math.Abs(otherDist) < math.Abs(dist) {
// 			shortestDist = otherDist
// 		}
// 	}

// 	for _, otherSlambox := range slambox.ConnectedBoxes {
// 		otherProjRect, _dist := g.world.ActiveLevel.TilemapCollider.ProjectRect(
// 			&otherSlambox.GetCollider().Rect,
// 			playerMove,
// 			g.world.ActiveLevel.DisconnectedColliders(otherSlambox),
// 		)

// 		offset := _dist - shortestDist

// 		switch playerMove {
// 		case maths.DirUp:
// 			otherProjRect.SetPos(otherSlambox.Collider.Left(), otherProjRect.Top()+offset)
// 		case maths.DirDown:
// 			otherProjRect.SetPos(otherSlambox.Collider.Left(), otherProjRect.Top()-offset)
// 		case maths.DirRight:
// 			otherProjRect.SetPos(otherProjRect.Left()-offset, otherSlambox.Collider.Top())
// 		case maths.DirLeft:
// 			otherProjRect.SetPos(otherProjRect.Left()+offset, otherSlambox.Collider.Top())
// 		}
// 		otherSlambox.SetTarget(otherProjRect.Left(), otherProjRect.Top())
// 	}

// 	offset := math.Abs(dist - shortestDist)

// 	switch playerMove {
// 	case maths.DirUp:
// 		projectedSlamboxRect.SetPos(slambox.Collider.Left(), projectedSlamboxRect.Top()+offset)
// 	case maths.DirDown:
// 		projectedSlamboxRect.SetPos(slambox.Collider.Left(), projectedSlamboxRect.Top()-offset)
// 	case maths.DirRight:
// 		projectedSlamboxRect.SetPos(projectedSlamboxRect.Left()-offset, slambox.Collider.Top())
// 	case maths.DirLeft:
// 		projectedSlamboxRect.SetPos(projectedSlamboxRect.Left()+offset, slambox.Collider.Top())
// 	}

// 	slambox.SetTarget(projectedSlamboxRect.Left(), projectedSlamboxRect.Top())
// 	slamFinishChan <- 1
// }
