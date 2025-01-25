package game

import (
	"image/color"
	"mask_of_the_tomb/player"
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
	world      *World
}

func (g *Game) Init() {
	g.world.Init()
	g.player.Init(g.world.GetSpawnPoint())
}

func (g *Game) Update() error {
	// Every entity should manage its own data, signals are passed to create communication
	// between entities

	playerMove := g.player.GetInput()

	if playerMove != player.DirNone && !g.player.IsMoving() {
		playerX, playerY := g.player.GetPos()
		targetX, targetY := g.world.getCollision(playerMove, playerX, playerY)
		g.player.SetTarget(targetX, targetY)
	}

	g.player.Update()

	return nil
}

func (g *Game) Draw() *ebiten.Image {
	g.worldSurf.Fill(color.RGBA{120, 120, 255, 255})

	g.world.Draw(g.worldSurf)
	g.player.Draw(g.worldSurf)

	// Pixel scaling
	op := Op()
	OpScale(op, PixelScale, PixelScale)
	g.screenSurf.DrawImage(g.worldSurf, op)
	return g.screenSurf
}

func NewGame() *Game {
	return &Game{
		worldSurf:  ebiten.NewImage(GameWidth, GameHeight),
		screenSurf: ebiten.NewImage(GameWidth*PixelScale, GameHeight*PixelScale),
		player:     player.NewPlayer(),
		world:      NewWorld(),
	}
}
