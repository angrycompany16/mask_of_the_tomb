package scenes

// Now we change all the ui layers into different scenes...

import (
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/libraries/camera"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	ui "mask_of_the_tomb/internal/plugins/UI"
	"mask_of_the_tomb/internal/plugins/player"
	"mask_of_the_tomb/internal/plugins/world"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameplayScene struct {
	UI                       *ui.UI
	world                    *world.World
	player                   *player.Player
	deathEffectEnterListener *events.EventListener
	titleCardTimeoutListener *events.EventListener
	levelCardTimeoutListener *events.EventListener
	playerMoveListener       *events.EventListener
	// TODO: introduce gametime/globaltime?
	// this would make it easier to pause the game
}

func (g *GameplayScene) Init() {
	gameData := errs.Must(save.GetSaveAsset("saveData"))
	g.world.Init(InitLevelName, gameData)

	resetX, resetY := g.world.ActiveLevel.GetResetPoint()

	g.player.Init(resetX, resetY, maths.DirNone)
	playerWidth, playerHeight := g.player.GetSize()

	w, h := g.world.ActiveLevel.GetBounds()

	camera.Init(
		w, h,
		(rendering.GAME_WIDTH-playerWidth)/2,
		(rendering.GAME_HEIGHT-playerHeight)/2,
	)

	screenFade := ui.NewOverlay(ui.NewScreenFade(), time.Second*2)
	g.deathEffectEnterListener = events.NewEventListener(screenFade.OnFinishEnter)
	g.UI.AddOverlay("screenfade", screenFade)

	titleCard := ui.NewOverlay(ui.NewTitleCard(), time.Second*2)
	g.titleCardTimeoutListener = events.NewEventListener(titleCard.OnIdleTimeout)
	g.UI.AddOverlay("titlecard", titleCard)

	levelCard := ui.NewOverlay(ui.NewLevelCard(), time.Second)
	g.levelCardTimeoutListener = events.NewEventListener(levelCard.OnIdleTimeout)
	g.UI.AddOverlay("levelcard", levelCard)

	g.playerMoveListener = events.NewEventListener(g.player.OnMove)
}

func (g *GameplayScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	// TODO: Find a way to pause the game while in the menu screen
	if musicScene, _, ok := sceneStack.GetScene("musicScene"); ok {
		musicScene.(*BaseScene).musicPlayer.PlayGameMusic(g.world.ActiveLevel.GetBiome())
	} else {
		fmt.Println("Music player was not found in game")
	}

	camera.Update()

	g.UI.Update()
	titlecard := g.UI.GetOverlay("titlecard")
	if _, raised := g.titleCardTimeoutListener.Poll(); raised {
		titlecard.StartFadeOut()
	}

	levelcard := g.UI.GetOverlay("levelcard")
	if _, raised := g.levelCardTimeoutListener.Poll(); raised {
		levelcard.StartFadeOut()
	}

	velX, velY := g.player.GetMovementSize()
	posX, posY := g.player.GetPosCentered()

	// TODO: Move even more logic out of here pleeease
	g.world.Update(posX, posY, velX, velY)
	if eventInfo, ok := g.playerMoveListener.Poll(); ok {
		moveDir := eventInfo.Data.(maths.Direction)
		slambox, hitSlambox := g.world.ActiveLevel.GetSlamboxHit(g.player.GetHitbox(), moveDir)

		catcherHitboxes := g.world.ActiveLevel.GetCatcherRects()
		hitCatcher, catcherX, catcherY := false, 0.0, 0.0
		for _, catcherHitbox := range catcherHitboxes {
			hit, x, y := catcherHitbox.RaycastDirectional(posX, posY, moveDir)
			hitCatcher = hitCatcher || hit
			if hit {
				catcherX = x
				catcherY = y
			}
		}

		platformHitboxes := make([]*maths.Rect, 0)
		if moveDir == maths.DirUp {
			platformHitboxes = g.world.ActiveLevel.GetPlatformHitboxes(false)
		} else if moveDir == maths.DirDown {
			platformHitboxes = g.world.ActiveLevel.GetPlatformHitboxes(true)
		}

		if hitSlambox {
			g.player.StartSlamming(moveDir)
			slambox.StartSlam(moveDir, &g.world.ActiveLevel.TilemapCollider, g.world.ActiveLevel.GetDisconnectedColliders(slambox))
		} else if hitCatcher {
			g.player.Dash(moveDir, catcherX, catcherY)
		} else {
			newRect, _ := g.world.ActiveLevel.TilemapCollider.ProjectRect(
				g.player.GetHitbox(),
				moveDir,
				append(g.world.ActiveLevel.GetSlamboxTargetRects(), platformHitboxes...),
			)
			if newRect != *g.player.GetHitbox() {
				g.player.Dash(moveDir, newRect.Left(), newRect.Top())
			}
		}
	}

	doorOverlap, levelIid, doorEntityIid := g.world.ActiveLevel.CheckDoorOverlap(g.player.GetHitbox())
	if g.player.GetLevelSwapInput() && doorOverlap && !g.player.Disabled {
		newBiome := errs.Must(world.ChangeActiveLevel(g.world, levelIid, doorEntityIid))
		if newBiome != "" {
			titleCardOverlay := g.UI.GetOverlay("titlecard")
			titleCard, _ := titleCardOverlay.OverlayContent.(*ui.TitleCard)
			titleCard.ChangeText(newBiome)
			titleCardOverlay.StartFadeIn()
		}
		camera.SetBorders(g.world.ActiveLevel.GetBounds())
		g.player.SetHitboxPos(g.world.ActiveLevel.GetResetPoint())
		levelCardOverlay := g.UI.GetOverlay("levelcard")
		levelCard, _ := levelCardOverlay.OverlayContent.(*ui.LevelCard)
		levelCard.ChangeText(g.world.ActiveLevel.GetTitle())
		levelCardOverlay.StartFadeIn()
	}

	restartPrompted := inpututil.IsKeyJustReleased(ebiten.KeyR)
	hitHazard := g.world.ActiveLevel.GetHazardHit(g.player.GetHitbox())
	hitTurret := g.world.ActiveLevel.CheckTurretHit(g.player.GetHitbox())
	if hitTurret && !g.player.Disabled || hitHazard && !g.player.Disabled || restartPrompted {
		g.player.Die()
		screenFade := g.UI.GetOverlay("screenfade")
		screenFade.StartFadeIn()
	}

	_, raised := g.deathEffectEnterListener.Poll()
	if raised {
		posX, posY := g.world.ResetActiveLevel()
		g.player.SetPos(posX, posY)
		g.player.Respawn()

		screenFade := g.UI.GetOverlay("screenfade")
		screenFade.StartFadeOut()
	}

	g.player.Update()
	camera.SetPos(g.player.GetPosCentered())

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return &scene.SceneTransition{
			Kind:       scene.Push,
			OtherScene: MakePauseScene(),
		}, true
	}

	return nil, false
}

func (g *GameplayScene) Draw() {
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

	g.UI.Draw()
	// UI is HARD-CODED to render at the UI layer...
	// I sck at programming
}

func (g *GameplayScene) GetName() string { return "gameplayScene" }
func MakeGameplayScene() *GameplayScene {
	return &GameplayScene{
		UI:     errs.Must(assettypes.GetYamlAsset("hud")).(*ui.UI),
		player: player.NewPlayer(),
		world:  world.NewWorld(),
	}
}
