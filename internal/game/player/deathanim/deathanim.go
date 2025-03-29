package deathanim

import (
	"image/color"
	"mask_of_the_tomb/internal/ebitenrenderutil"
	ui "mask_of_the_tomb/internal/game/UI"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/core/rendering/camera"
	"mask_of_the_tomb/internal/maths"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	smallRadius    = 6.0
	largeRadius    = 30.0
	ballCount      = 6
	speed          = 0.012
	border         = 2.0
	surfacePadding = 20.0
)

var (
	ballColors = ui.ColorPair{
		BrightColor: color.RGBA{255, 253, 240, 255},
		DarkColor:   color.RGBA{21, 10, 31, 255},
	}
)

type DeathAnim struct {
	circleSurfaces []*ebiten.Image
	playing        bool
	x, y           float64
	t              float64
	radius         float64
	angleOffset    float64
}

func (d *DeathAnim) Update() {
	if d.playing {
		d.t += speed
		d.radius = getRadius(d.t)
		d.angleOffset += speed

		if d.t >= 1 {
			d.playing = false
		}
	} else {
	}
}

func (d *DeathAnim) Draw() {
	if !d.playing {
		return
	}
	for i := range ballCount {
		camX, camY := camera.GetPos()
		xPos := d.x + d.radius*math.Cos(2*math.Pi/float64(ballCount)*float64(i)+d.angleOffset)
		yPos := d.y - d.radius*math.Sin(2*math.Pi/float64(ballCount)*float64(i)+d.angleOffset)
		vector.DrawFilledCircle(
			d.circleSurfaces[i],
			smallRadius+surfacePadding/2,
			smallRadius+surfacePadding/2,
			smallRadius,
			ballColors.DarkColor,
			false,
		)
		vector.StrokeCircle(
			d.circleSurfaces[i],
			smallRadius+surfacePadding/2,
			smallRadius+surfacePadding/2,
			smallRadius,
			border,
			ballColors.BrightColor,
			false,
		)

		ebitenrenderutil.DrawAt(d.circleSurfaces[i], rendering.RenderLayers.Overlay, xPos-camX, yPos-camY)
	}
}

func (d *DeathAnim) Play() {
	// Reset values
	d.angleOffset = 0
	d.radius = 0
	d.t = 0
	d.playing = true
}

func (d *DeathAnim) SetPos(x, y float64) {
	d.x = x
	d.y = y
}

func NewDeathAnim() *DeathAnim {
	deathAnim := &DeathAnim{
		playing: false,
	}

	for i := range ballCount {
		func(i int) {}(i)
		deathAnim.circleSurfaces = append(deathAnim.circleSurfaces, ebiten.NewImage(smallRadius*2+surfacePadding, smallRadius*2+surfacePadding))
	}
	return deathAnim
}

func getRadius(t float64) float64 {
	if t < 0.5 {
		return largeRadius * maths.SineInOut(t*2)
	}
	return largeRadius * (1 - maths.QuadInOut(t*2-1))
}
