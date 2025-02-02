package game

import (
	"image/color"
	"mask_of_the_tomb/commons"
	. "mask_of_the_tomb/ebitenRenderUtil"
	ui "mask_of_the_tomb/game/UI"
	"mask_of_the_tomb/game/camera"
	"mask_of_the_tomb/game/player"
	"mask_of_the_tomb/game/world"
	"mask_of_the_tomb/save"
	. "mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	worldSurf  *ebiten.Image
	uiSurf     *ebiten.Image
	screenSurf *ebiten.Image
	player     *player.Player
	world      *world.World
	camera     *camera.Camera
	ui         *ui.UI
}

func (g *Game) Init() {
	save.GlobalSave.LoadGame()
	g.world.Init()
	g.player.Init(g.world.ActiveLevel.GetSpawnPoint())
	width, height := g.world.ActiveLevel.GetLevelBounds()
	playerWidth, playerHeight := g.player.GetSize()
	g.camera.Init(
		width,
		height,
		(commons.GameWidth-playerWidth)/2,
		(commons.GameHeight-playerHeight)/2,
	)
	g.ui.SetText(g.ui.GenerateScoreMessage(0))
}

func (g *Game) Update() error {
	playerMove := g.player.GetMoveInput()
	playerX, playerY := g.player.GetPos()
	if playerMove != player.DirNone && !g.player.IsMoving() {
		targetX, targetY := g.world.ActiveLevel.GetCollision(playerMove, playerX, playerY)
		g.player.SetTarget(targetX, targetY)
	}

	if g.player.GetLevelSwapInput() {
		validSwapPosition, entityInstance := g.world.ActiveLevel.TryDoorOverlap(g.player.GetPos())

		if validSwapPosition {
			newPosX, newPosY := g.world.ExitByDoor(entityInstance)

			g.camera.SetBorders(g.world.ActiveLevel.GetLevelBounds())
			g.player.SetPos(F64(newPosX), F64(newPosY))
		}
	}

	dx, dy := g.player.GetMovementSize()
	collectibleOverlapCount := g.world.ActiveLevel.TryCollectibleOverlap(playerX, playerY, dx, dy)

	if collectibleOverlapCount > 0 {
		g.player.SetScore(g.player.GetScore() + collectibleOverlapCount)
		g.ui.SetText(g.ui.GenerateScoreMessage(g.player.GetScore()))
	}

	g.player.Update()

	playerX, playerY = g.player.GetPos()
	g.camera.SetPos(playerX, playerY)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return commons.Terminated
	}

	return nil
}

func (g *Game) Draw() *ebiten.Image {
	camX, camY := g.camera.GetPos()
	g.worldSurf.Fill(color.RGBA{255, 255, 0, 255})
	g.uiSurf.Clear()

	g.world.Draw(g.worldSurf, camX, camY)
	g.player.Draw(g.worldSurf, camX, camY)
	g.ui.Draw(g.uiSurf)

	DrawAtScaled(g.worldSurf, g.screenSurf, 0, 0, commons.PixelScale, commons.PixelScale)
	DrawAt(g.uiSurf, g.screenSurf, 0, 0)

	return g.screenSurf
}

func NewGame() *Game {
	return &Game{
		worldSurf:  ebiten.NewImage(commons.GameWidth, commons.GameHeight),
		uiSurf:     ebiten.NewImage(commons.GameWidth*commons.PixelScale, commons.GameHeight*commons.PixelScale),
		screenSurf: ebiten.NewImage(commons.GameWidth*commons.PixelScale, commons.GameHeight*commons.PixelScale),
		player:     player.NewPlayer(),
		world:      &world.World{},
		camera:     camera.NewCamera(),
		ui:         ui.NewUi(),
	}
}
