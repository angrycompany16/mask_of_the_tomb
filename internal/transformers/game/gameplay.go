package game

import (
	"fmt"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/camera"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"
	"mask_of_the_tomb/internal/plugins/world"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Initializes player, world, etc
func (g *Game) InitGameplayStage(gameData save.SaveData, enter bool) {
	g.world.Init(InitLevelName, gameData)

	resetX, resetY := g.world.ActiveLevel.GetResetPoint()
	gameEntryX, gameEntryY := g.world.ActiveLevel.GetGameEntryPos()

	if enter {
		g.player.Init(gameEntryX, gameEntryY, maths.DirNone)
	} else {
		g.player.Init(resetX, resetY, maths.DirNone)
	}
	playerWidth, playerHeight := g.player.GetSize()

	w, h := g.world.ActiveLevel.GetBounds()

	camera.Init(
		w, h,
		(rendering.GAME_WIDTH-playerWidth)/2,
		(rendering.GAME_HEIGHT-playerHeight)/2,
	)
	g.gameplayUI.SwitchActiveDisplay("hud", nil)
	g.mainUI.SwitchActiveDisplay("empty", nil)
}

func (g *Game) GameplayStageUpdate() {
	fmt.Println("Update gameplay")
	g.IntroStageUpdate()

	camera.Update()
	g.musicPlayer.PlayGameMusic(g.world.ActiveLevel.GetBiome())

	velX, velY := g.player.GetMovementSize()
	posX, posY := g.player.GetPosCentered()

	if g.introDashTimer == nil {
		g.introDashTimer = time.NewTimer(time.Second)
		g.introDashTimer.Stop()
	}

	if _, timedout := threads.Poll(g.introDashTimer.C); timedout {
		newRect, _ := g.world.ActiveLevel.TilemapCollider.ProjectRect(g.player.GetHitbox(), gameEntryDirection, g.world.ActiveLevel.GetSlamboxRects())
		if newRect != *g.player.GetHitbox() {
			g.player.Dash(gameEntryDirection, newRect.Left(), newRect.Top())
		}
	}

	if resources.State != resources.Playing {
		return
	}

	g.world.Update(posX, posY, velX, velY)
	if eventInfo, ok := g.playerMoveListener.Poll(); ok {
		moveDir := eventInfo.Data.(maths.Direction)
		slambox, hit := g.world.ActiveLevel.GetSlamboxHit(g.player.GetHitbox(), moveDir)
		// also check if we can dash into a catcher
		if hit {
			g.player.StartSlamming(moveDir)
			slambox.StartSlam(moveDir, &g.world.ActiveLevel.TilemapCollider, g.world.ActiveLevel.GetDisconnectedColliders(slambox))
		} else {
			newRect, _ := g.world.ActiveLevel.TilemapCollider.ProjectRect(g.player.GetHitbox(), moveDir, g.world.ActiveLevel.GetSlamboxTargetRects())
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
			titleCard, _ := titleCardOverlay.OverlayContent.(*ui.TitleCard)
			titleCard.ChangeText(newBiome)
			titleCardOverlay.StartFadeIn()
		}
		camera.SetBorders(g.world.ActiveLevel.GetBounds())
		g.player.SetHitboxPos(g.world.ActiveLevel.GetResetPoint())
		levelCardOverlay := g.gameplayUI.GetOverlay("levelcard")
		levelCard, _ := levelCardOverlay.OverlayContent.(*ui.LevelCard)
		levelCard.ChangeText(g.world.ActiveLevel.GetTitle())
		levelCardOverlay.StartFadeIn()
	}

	restartPrompted := inpututil.IsKeyJustReleased(ebiten.KeyR)
	hitHazard := g.world.ActiveLevel.GetHazardHit(g.player.GetHitbox())
	hitTurret := g.world.ActiveLevel.CheckTurretHit(g.player.GetHitbox())
	if hitTurret && !g.player.Disabled || hitHazard && !g.player.Disabled || restartPrompted {
		g.player.Die()
		screenFade := g.mainUI.GetOverlay("screenfade")
		screenFade.StartFadeIn()
	}

	_, raised := g.deathEffectEnterListener.Poll()
	if raised {
		posX, posY := g.world.ResetActiveLevel()
		g.player.SetPos(posX, posY)
		g.player.Respawn()

		screenFade := g.mainUI.GetOverlay("screenfade")
		screenFade.StartFadeOut()
	}

	g.player.Update()
	camera.SetPos(g.player.GetPosCentered())

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		resources.State = resources.Paused
		g.mainUI.SwitchActiveDisplay("pausemenu", nil)
	}
}

func (g *Game) GameplayStageDraw() {
	g.IntroStageDraw()
	pX, pY := g.player.GetPosCentered()
	cX, cY := camera.GetPos()
	drawCtx := rendering.Ctx{
		CamX:    cX,
		CamY:    cY,
		PlayerX: pX,
		PlayerY: pY,
	}

	g.player.Draw(rendering.WithLayer(drawCtx, rendering.ScreenLayers.Playerspace))
	g.world.ActiveLevel.Draw(drawCtx)

	g.gameplayUI.Draw()
	// UI is HARD-CODED to render at the UI layer...
	// I sck at programming
}
