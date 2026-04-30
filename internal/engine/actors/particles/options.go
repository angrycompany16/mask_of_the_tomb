package particles

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/utils"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewParticleSystem(graphic *graphic.Graphic, options ...utils.Option[ParticleSystem]) *ParticleSystem {
	particleSys := defaultParticleSystem(graphic)

	for _, option := range options {
		option(particleSys)
	}

	return particleSys
}

func defaultParticleSystem(
	graphic *graphic.Graphic,
) *ParticleSystem {
	return &ParticleSystem{
		Graphic:   graphic,
		particles: make([]*Particle, 0),
		Bursts: []*Burst{
			&Burst{50, 1},
		},
		GlobalSpace:     true,
		Emission:        5,
		Angle:           0,
		SpawnPosX:       maths.RandomFloat64{Min: -10, Max: 10},
		SpawnPosY:       maths.RandomFloat64{Min: -10, Max: 10},
		SpawnVelX:       maths.RandomFloat64{Min: -6, Max: 5},
		SpawnVelY:       maths.RandomFloat64{Min: -5, Max: 5},
		SpawnAngle:      maths.RandomFloat64{Min: 0, Max: 0},
		SpawnAngularVel: maths.RandomFloat64{Min: 0, Max: 0},
		AirFriction:     maths.RandomFloat64{Min: 0, Max: 0},
		StartScale:      maths.RandomFloat64{Min: 1, Max: 1},
		EndScale:        maths.RandomFloat64{Min: 0, Max: 0},
		Lifetime:        maths.RandomFloat64{Min: 0.5, Max: 1},
		NoiseFactorX:    maths.RandomFloat64{Min: 0, Max: 0},
		NoiseFactorY:    maths.RandomFloat64{Min: 0, Max: 0},
		t:               0,
		StartColor:      [4]uint8{255, 255, 255, 255},
		EndColor:        [4]uint8{0, 0, 0, 255},
		ImageWidth:      1,
		ImageHeight:     1,
		SpritePath:      "sprites/icons/circle-15x15.png",
		surf:            ebiten.NewImage(1, 1),
		gizmosImage:     ebiten.NewImage(1, 1),
		burstTimers:     make([]*time.Timer, 1),
		isPlaying:       true,
		target:          renderer.RenderTarget{Type: renderer.SCREEN, Name: "Playerspace"},
		drawOrder:       0,
	}
}

func WithBursts(bursts ...*Burst) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.Bursts = bursts
	}
}

func WithGlobalSpace(global bool) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.GlobalSpace = global
	}
}

func WithEmission(rate float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.Emission = rate
	}
}

func WithAngle(angle float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.Angle = angle
	}
}

func WithSpawnPos(minX, maxX, minY, maxY float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.SpawnPosX = maths.RandomFloat64{Min: minX, Max: maxX}
		ps.SpawnPosY = maths.RandomFloat64{Min: minY, Max: maxY}
	}
}

func WithSpawnVel(minX, maxX, minY, maxY float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.SpawnVelX = maths.RandomFloat64{Min: minX, Max: maxX}
		ps.SpawnVelY = maths.RandomFloat64{Min: minY, Max: maxY}
	}
}

func WithSpawnAngle(min, max float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.SpawnAngle = maths.RandomFloat64{Min: min, Max: max}
	}
}

func WithSpawnAngularVel(min, max float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.SpawnAngularVel = maths.RandomFloat64{Min: min, Max: max}
	}
}

func WithAirFriction(min, max float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.AirFriction = maths.RandomFloat64{Min: min, Max: max}
	}
}

func WithScale(startMin, startMax, endMin, endMax float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.StartScale = maths.RandomFloat64{Min: startMin, Max: startMax}
		ps.EndScale = maths.RandomFloat64{Min: endMin, Max: endMax}
	}
}

func WithLifetime(min, max float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.Lifetime = maths.RandomFloat64{Min: min, Max: max}
	}
}

func WithNoiseFactor(minX, maxX, minY, maxY float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.NoiseFactorX = maths.RandomFloat64{Min: minX, Max: maxX}
		ps.NoiseFactorY = maths.RandomFloat64{Min: minY, Max: maxY}
	}
}

func WithColors(start, end [4]uint8) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.StartColor = start
		ps.EndColor = end
	}
}

func WithImageSize(width, height int) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.ImageWidth = width
		ps.ImageHeight = height
		ps.surf = ebiten.NewImage(width, height)
		ps.gizmosImage = ebiten.NewImage(width, height)
	}
}

func WithSprite(path string) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.SpritePath = path
	}
}

func WithRenderInfo(target renderer.RenderTarget, drawOrder int) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.target = target
		ps.drawOrder = drawOrder
	}
}

func WithPlaying(playing bool) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.isPlaying = playing
	}
}
