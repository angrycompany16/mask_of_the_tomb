package entities

import (
	"image"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/resources"

	"math"
	"math/rand/v2"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/aquilax/go-perlin"
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Entities seriously need to have their own libraries...
const (
	bladesPertile = 8.0
)

type Grass struct {
	dst      *ebiten.Image
	extents  maths.Rect
	gridSize float64
	blades   []*grassBlade
	perlin   *perlin.Perlin
}

// TODO: Optimize with grass groups
func (g *Grass) Update(playerX, playerY, playerVelX, playerVelY float64) {
	for _, blade := range g.blades {
		blade.Update(playerX, playerY, playerVelX, playerVelY, g.perlin)
	}
}

func (g *Grass) Draw(ctx rendering.Ctx) {
	for _, blade := range g.blades {
		ebitenrenderutil.DrawAtRotated(blade.sprite, ctx.Dst, blade.posX, blade.posY, blade.angle+blade.angleOffset, 0.5, 1.0)
	}
}

func NewGrass(
	entity *ebitenLDTK.Entity,
	gridSize float64,
	grassTilemap *ebiten.Image,
	dst *ebiten.Image,
) Grass {
	newGrass := Grass{
		extents: *maths.NewRect(
			entity.Px[0],
			entity.Px[1],
			entity.Width,
			entity.Height,
		),
		gridSize: gridSize,
		dst:      dst,
		perlin:   perlin.NewPerlin(2.0, 0.05, 4, resources.GrassWindSeed),
	}

	for i := 0; i < int(bladesPertile*entity.Width/gridSize); i++ {
		x := entity.Px[0] + float64(i)/bladesPertile*gridSize
		newGrass.blades = append(newGrass.blades, newGrassBlade(grassTilemap, x, entity.Px[1]-(SPRITE_HEIGHT-entity.Height)))
	}

	return newGrass
}

const GRASS_FREQUENCY = 0.4
const DAMPING_RATIO = 0.3
const WIND_SPEED = 0.45
const WIND_BIAS = 0.3
const WIND_STRENGTH = 0.03
const SPRITE_WIDTH, SPRITE_HEIGHT = 4, 16

type grassBlade struct {
	sprite      *ebiten.Image
	angle       float64
	angularVel  float64
	angleOffset float64
	posX, posY  float64
}

func (g *grassBlade) Update(playerX, playerY, playerVelX, playerVelY float64, perlin *perlin.Perlin) {
	// Interaction forces
	// Add interactions from slamboxes and other entities
	playerDist := maths.Length(playerX-g.posX, playerY-g.posY)
	velocityEffect := 1 - maths.SmoothStep(0, 20, playerDist)

	velocityStrengthX := 0.0
	if playerVelX < 0 {
		velocityStrengthX = maths.SmoothStep(-50, 0, playerVelX) - 1
	} else {
		velocityStrengthX = maths.SmoothStep(0, 50, playerVelX)
	}

	velocityStrengthY := maths.SmoothStep(0, 40, math.Abs(playerVelY))
	velocityStrengthY *= math.Copysign(1, g.posX-playerX)
	velocityStrength := velocityStrengthX + velocityStrengthY

	g.angularVel += velocityEffect * velocityStrength * 20

	// Wind forces
	noiseFactor := perlin.Noise3D(g.posX, g.posY, resources.Time*WIND_SPEED) + WIND_BIAS
	g.angularVel += noiseFactor * WIND_STRENGTH

	g.angularVel += (-2*DAMPING_RATIO*GRASS_FREQUENCY*g.angularVel - math.Pow(GRASS_FREQUENCY, 2)*g.angle) * 0.2
	g.angle += g.angularVel * 0.2
}

func newGrassBlade(grassTilemap *ebiten.Image, posX, posY float64) *grassBlade {
	tileIndex := rand.IntN(8)
	grassSprite := grassTilemap.SubImage(image.Rect(tileIndex*SPRITE_WIDTH, 0, (tileIndex+1)*SPRITE_WIDTH, SPRITE_HEIGHT)).(*ebiten.Image)

	return &grassBlade{
		posX: posX, posY: posY, sprite: grassSprite,
	}
}
