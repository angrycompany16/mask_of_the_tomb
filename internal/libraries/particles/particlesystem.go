package particles

import (
	"image/color"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/threads"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// I'm very strongly considering rewriting this to a multithreaded
// (and almost completely reworked) system cause this is broken AF

// Maybe we can convert the non-rendering part of the particlesystem into
// a struct and then only use the drawing part outside of that
// Kind of ugly tho

// TODO: Implement stop time so it doesn't run forever
// TODO: Make burst count random
// TODO: Make color random
// TODO: Convert angles to degrees
// NOTE: there's an effective cap on emission (it cannot be higher than 60) because
// we only ever add one particle to the system
// TODO: Make it thread safe

type ParticleSystem struct {
	GlobalSpace     bool    `yaml:"GlobalSpace"`
	PosX            float64 `yaml:"PosX"`
	PosY            float64 `yaml:"PosY"`
	Angle           float64 `yaml:"Angle"`
	particles       []*Particle
	Emission        float64             `yaml:"Emission"`
	Bursts          []*ParticleBurst    `yaml:"Bursts"`
	SpawnPosX       maths.RandomFloat64 `yaml:"SpawnPosX"`
	SpawnPosY       maths.RandomFloat64 `yaml:"SpawnPosY"`
	SpawnVelX       maths.RandomFloat64 `yaml:"SpawnVelX"`
	SpawnVelY       maths.RandomFloat64 `yaml:"SpawnVelY"`
	SpawnAngle      maths.RandomFloat64 `yaml:"SpawnAngle"`
	SpawnAngularVel maths.RandomFloat64 `yaml:"SpawnAngularVel"`
	AirFriction     maths.RandomFloat64 `yaml:"AirFriction"`
	StartScale      maths.RandomFloat64 `yaml:"StartScale"`
	EndScale        maths.RandomFloat64 `yaml:"EndScale"`
	Lifetime        maths.RandomFloat64 `yaml:"Lifetime"`
	NoiseFactorX    maths.RandomFloat64 `yaml:"NoiseFactorX"`
	NoiseFactorY    maths.RandomFloat64 `yaml:"NoiseFactorY"`
	t               float64
	// TODO: Use ebiten.color by making embedded struct with custom unmarshaler
	StartColor  [4]uint8      `yaml:"StartColor"`
	EndColor    [4]uint8      `yaml:"EndColor"`
	ImageWidth  int           `yaml:"ImageWidth"`
	ImageHeight int           `yaml:"ImageHeight"`
	SpritePath  string        `yaml:"SpritePath"`
	surf        *ebiten.Image // The image that the particles are drawn onto
	sprite      *ebiten.Image // The sprite for the particles
	layer       *ebiten.Image
	burstTimers []*time.Timer
	isPlaying   bool
}

func (ps *ParticleSystem) Play() {
	ps.t = 0
	ps.isPlaying = true
	// ps.particles = make([]*Particle, 0)
	for i, burst := range ps.Bursts {
		ps.burstTimers[i] = time.NewTimer(time.Duration(burst.Time * 1e9))
	}
}

func (ps *ParticleSystem) Stop() {
	if !ps.isPlaying {
		return
	}
	ps.isPlaying = false
	for i := range ps.Bursts {
		ps.burstTimers[i].Stop()
	}
}

func (ps *ParticleSystem) Update() {
	ps.t += 0.016666666667

	// Hmm... this might not be that efficient
	j := 0
	for i, particle := range ps.particles {
		finished := particle.update()
		if finished {
			ps.particles[i] = ps.particles[len(ps.particles)-1-j]
			j++
		}
	}
	ps.particles = ps.particles[:len(ps.particles)-j]

	if !ps.isPlaying {
		return
	}

	for i, burst := range ps.Bursts {
		if _, ok := threads.Poll(ps.burstTimers[i].C); ok {
			for range burst.Count {
				ps.particles = append(ps.particles, ps.newParticle())
			}
		}
	}

	if ps.Emission > 0 && ps.t > 1/ps.Emission {
		ps.particles = append(ps.particles, ps.newParticle())
		ps.t = 0
	}
}

func (ps *ParticleSystem) Draw(ctx rendering.Ctx) {
	if ps.GlobalSpace {
		for _, particle := range ps.particles {
			particle.draw(ctx.Dst, ctx.CamX, ctx.CamY)
		}
		return
	}

	// Local-space rendering is very cursed ngl
	s := ps.surf.Bounds().Size()
	for _, particle := range ps.particles {
		particle.draw(ps.surf, 0, 0)
	}

	ebitenrenderutil.DrawAtRotated(
		ps.surf, ctx.Dst,
		ps.PosX-ctx.CamX-float64(s.X)/2,
		ps.PosY-ctx.CamY-float64(s.Y)/2,
		ps.Angle, 0.5, 0.5)

	ps.surf.Clear()
}

func (ps *ParticleSystem) SetPos(x, y float64) {
	ps.PosX = x
	ps.PosY = y
}

func (ps *ParticleSystem) newParticle() *Particle {
	var x, y float64
	if ps.GlobalSpace {
		x, y = ps.SpawnPosX.Eval()+ps.PosX-float64(ps.ImageWidth)/2, ps.SpawnPosY.Eval()+ps.PosY-float64(ps.ImageHeight)/2
	} else {
		x, y = ps.SpawnPosX.Eval(), ps.SpawnPosY.Eval()
	}
	startScale := ps.StartScale.Eval()
	startColor := color.RGBA{
		R: ps.StartColor[0],
		G: ps.StartColor[1],
		B: ps.StartColor[2],
		A: ps.StartColor[3],
	}
	endColor := color.RGBA{
		R: ps.EndColor[0],
		G: ps.EndColor[1],
		B: ps.EndColor[2],
		A: ps.EndColor[3],
	}
	return &Particle{
		x, y, ps.SpawnVelX.Eval(), ps.SpawnVelY.Eval(),
		ps.SpawnAngle.Eval(), ps.SpawnAngularVel.Eval(),
		ps.AirFriction.Eval(),
		startScale, startScale, ps.EndScale.Eval(),
		ps.Lifetime.Eval(), 0,
		startColor, startColor, endColor,
		ps.NoiseFactorX.Eval(), ps.NoiseFactorY.Eval(),
		ps.sprite,
	}
}

// TODO: This is a challenging case: We have an asset that needs to load a path to another asset
// Can this be solved? It shouldn't be a gigantic performance loss if not so most likely nothing to
// worry about.
func (ps *ParticleSystem) Init() {
	ps.particles = make([]*Particle, 0)
	ps.surf = ebiten.NewImage(ps.ImageWidth, ps.ImageHeight)
	spritePath := errs.Must(filepath.Localize(ps.SpritePath))
	ps.sprite = errs.MustNewImageFromFile(spritePath)
	ps.burstTimers = make([]*time.Timer, len(ps.Bursts))
}

type ParticleBurst struct {
	Count int     `yaml:"Count"`
	Time  float64 `yaml:"Time"`
}
