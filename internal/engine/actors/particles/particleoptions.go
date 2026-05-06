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
		SpawnPosX:       maths.RandomFloat{Min: -10, Max: 10},
		SpawnPosY:       maths.RandomFloat{Min: -10, Max: 10},
		SpawnVelX:       maths.RandomFloat{Min: -6, Max: 5},
		SpawnVelY:       maths.RandomFloat{Min: -5, Max: 5},
		SpawnAngle:      maths.RandomFloat{Min: 0, Max: 0},
		SpawnAngularVel: maths.RandomFloat{Min: 0, Max: 0},
		AirFriction:     maths.RandomFloat{Min: 0, Max: 0},
		StartScale:      maths.RandomFloat{Min: 1, Max: 1},
		EndScale:        maths.RandomFloat{Min: 0, Max: 0},
		Lifetime:        maths.RandomFloat{Min: 0.5, Max: 1},
		NoiseFactorX:    maths.RandomFloat{Min: 0, Max: 0},
		NoiseFactorY:    maths.RandomFloat{Min: 0, Max: 0},
		t:               0,
		StartColor: [4]maths.RandomInt{
			maths.NewRandomInt(255, 255),
			maths.NewRandomInt(255, 255),
			maths.NewRandomInt(255, 255),
			maths.NewRandomInt(255, 255),
		},
		EndColor: [4]maths.RandomInt{
			maths.NewRandomInt(0, 0),
			maths.NewRandomInt(0, 0),
			maths.NewRandomInt(0, 0),
			maths.NewRandomInt(0, 0),
		},
		ImageWidth:  1,
		ImageHeight: 1,
		SpritePath:  "sprites/icons/circle-15x15.png",
		surf:        ebiten.NewImage(1, 1),
		gizmosImage: ebiten.NewImage(1, 1),
		burstTimers: make([]*time.Timer, 1),
		isPlaying:   true,
		target:      renderer.RenderTarget{Type: renderer.SCREEN, Name: "Playerspace"},
		drawOrder:   0,
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
		ps.SpawnPosX = maths.RandomFloat{Min: minX, Max: maxX}
		ps.SpawnPosY = maths.RandomFloat{Min: minY, Max: maxY}
	}
}

func WithSpawnVel(minX, maxX, minY, maxY float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.SpawnVelX = maths.RandomFloat{Min: minX, Max: maxX}
		ps.SpawnVelY = maths.RandomFloat{Min: minY, Max: maxY}
	}
}

func WithSpawnAngle(min, max float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.SpawnAngle = maths.RandomFloat{Min: min, Max: max}
	}
}

func WithSpawnAngularVel(min, max float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.SpawnAngularVel = maths.RandomFloat{Min: min, Max: max}
	}
}

func WithAirFriction(min, max float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.AirFriction = maths.RandomFloat{Min: min, Max: max}
	}
}

func WithScale(startMin, startMax, endMin, endMax float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.StartScale = maths.RandomFloat{Min: startMin, Max: startMax}
		ps.EndScale = maths.RandomFloat{Min: endMin, Max: endMax}
	}
}

func WithLifetime(min, max float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.Lifetime = maths.RandomFloat{Min: min, Max: max}
	}
}

func WithNoiseFactor(minX, maxX, minY, maxY float64) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.NoiseFactorX = maths.RandomFloat{Min: minX, Max: maxX}
		ps.NoiseFactorY = maths.RandomFloat{Min: minY, Max: maxY}
	}
}

func WithColors(startMin, startMax, endMin, endMax [4]uint8) utils.Option[ParticleSystem] {
	return func(ps *ParticleSystem) {
		ps.StartColor = [4]maths.RandomInt{
			maths.NewRandomInt(int(startMin[0]), int(startMax[0])),
			maths.NewRandomInt(int(startMin[1]), int(startMax[1])),
			maths.NewRandomInt(int(startMin[2]), int(startMax[2])),
			maths.NewRandomInt(int(startMin[3]), int(startMax[3])),
		}
		ps.EndColor = [4]maths.RandomInt{
			maths.NewRandomInt(int(endMin[0]), int(endMax[0])),
			maths.NewRandomInt(int(endMin[1]), int(endMax[1])),
			maths.NewRandomInt(int(endMin[2]), int(endMax[2])),
			maths.NewRandomInt(int(endMin[3]), int(endMax[3])),
		}
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
