package game

import (
	"image/color"
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/game/camera"
	"mask_of_the_tomb/game/player"
	"mask_of_the_tomb/game/world"
	. "mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameWidth, GameHeight = 480, 270
	PixelScale            = 4
)

type Game struct {
	worldSurf  *ebiten.Image
	screenSurf *ebiten.Image
	player     *player.Player
	world      *world.World
	camera     *camera.Camera
}

func (g *Game) Init() {
	g.world.Init()
	g.player.Init(g.world.GetSpawnPoint())
	g.camera.Init(g.world.GetLevelBorders(), GameWidth/2, GameHeight/2-1)
}

func (g *Game) Update() error {
	// Every entity should manage its own data, signals are passed to create communication
	// between entities

	playerMove := g.player.GetMoveInput()
	if playerMove != player.DirNone && !g.player.IsMoving() {
		playerX, playerY := g.player.GetPos()
		targetX, targetY := g.world.GetCollision(playerMove, playerX, playerY)
		g.player.SetTarget(targetX, targetY)
	}

	// TODO: Finish level swapping. Current problems:
	//  - Camera movement is not implemented
	//  - Level tile data needs to be regenerated when switching levels due to different level sizes

	if g.player.GetLevelSwapInput() {
		validSwapPosition, entityInstance := g.world.TryDoorOverlap(g.player.GetPos())
		// Check position
		if validSwapPosition {
			newPosX, newPosY := g.world.ExitByDoor(entityInstance)

			g.player.SetPos(F64(newPosX), F64(newPosY))
			// println("Swap level")
		}
		// Swap to correct level
		// Set player position
	}

	// Get player position

	// Input player position to camera controls

	g.player.Update()

	return nil
}

func (g *Game) Draw() *ebiten.Image {
	g.worldSurf.Fill(color.RGBA{120, 120, 255, 255})

	g.world.Draw(g.worldSurf)
	g.player.Draw(g.worldSurf)

	DrawAtScaled(g.worldSurf, g.screenSurf, 0, 0, PixelScale, PixelScale)

	return g.screenSurf
}

func NewGame() *Game {
	return &Game{
		worldSurf:  ebiten.NewImage(GameWidth, GameHeight),
		screenSurf: ebiten.NewImage(GameWidth*PixelScale, GameHeight*PixelScale),
		player:     player.NewPlayer(),
		world:      &world.World{},
		camera:     camera.NewCamera(),
	}
}
