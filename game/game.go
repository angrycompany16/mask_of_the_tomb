package game

import (
	"image/color"
	"mask_of_the_tomb/player"
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameWidth, GameHeight = 480, 270
	PixelScale            = 4
)

type World struct {
	worldSurf  *ebiten.Image
	screenSurf *ebiten.Image
	player     *player.Player
}

func (w *World) Update() {
	// Every entity should manage its own data, signals are passed to create communication
	// between entities

	w.player.Update()
}

func (w *World) Draw() *ebiten.Image {
	w.worldSurf.Fill(color.RGBA{120, 120, 255, 255})

	w.player.Draw(w.worldSurf)

	// Pixel scaling
	op := utils.Op()
	utils.OpScale(op, PixelScale, PixelScale)
	w.screenSurf.DrawImage(w.worldSurf, op)
	return w.screenSurf
}

func MakeWorld() *World {
	return &World{
		worldSurf:  ebiten.NewImage(GameWidth, GameHeight),
		screenSurf: ebiten.NewImage(GameWidth*PixelScale, GameHeight*PixelScale),
		player:     player.NewPlayer(),
	}
}
