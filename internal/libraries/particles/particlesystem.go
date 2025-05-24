package particles

// Q: Should fileio be a core component? Not completely sure
import (
	"image/color"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/fileio"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Maybe we can convert the non-rendering part of the particlesystem into
// a struct and then only use the drawing part outside of that
// Kind of ugly tho

// TODO: Implement stop time so it doesn't run forever
// TODO: Make burst count random
// TODO: Make color random
// TODO: Some kind of air friction coefficient?
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
	surf        *ebiten.Image // The surf that the particles are drawn onto
	sprite      *ebiten.Image // The sprite for the particles
	layer       *ebiten.Image
}

func (ps *ParticleSystem) Play() {
	ps.particles = nil
	// Not good, really not good...
	// This is literally just not thread safe at all...
	go func() {
		for _, burst := range ps.Bursts {
			time.Sleep(time.Duration(burst.Time) * time.Second)
			for i := 0; i < burst.Count; i++ {
				ps.particles = append(ps.particles, ps.newParticle())
			}
		}
	}()
}

func (ps *ParticleSystem) Update() {
	ps.t += 0.016666666667

	j := 0
	for i, particle := range ps.particles {
		finished := particle.update()
		if finished {
			ps.particles[i] = ps.particles[len(ps.particles)-1-j]
			j++
		}
	}
	ps.particles = ps.particles[:len(ps.particles)-j]

	if ps.Emission == 0 {
		return
	}

	if ps.t > 1/ps.Emission {
		ps.particles = append(ps.particles, ps.newParticle())
		ps.t = 0
	}
}

// TODO: Maybe take in the layer as a function parameter?
func (ps *ParticleSystem) Draw(ctx rendering.Ctx) {
	if ps.GlobalSpace {
		for _, particle := range ps.particles {
			particle.draw(ctx.Dst, ctx.CamX, ctx.CamY)
		}
		return
	}

	s := ps.surf.Bounds().Size()
	for _, particle := range ps.particles {
		particle.draw(ps.surf, -float64(s.X)/2, -float64(s.Y)/2)
	}
	ebitenrenderutil.DrawAtRotated(
		ps.surf, ctx.Dst,
		ps.PosX-float64(s.X)/2-ctx.CamX,
		ps.PosY-float64(s.Y)/2-ctx.CamY,
		ps.Angle, 0.5, 0.5)

	ps.surf.Clear()
}

func (ps *ParticleSystem) newParticle() *Particle {
	var x, y float64
	if ps.GlobalSpace {
		x, y = ps.SpawnPosX.Eval()+ps.PosX, ps.SpawnPosY.Eval()+ps.PosY
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
		startScale, startScale, ps.EndScale.Eval(),
		ps.Lifetime.Eval(), 0,
		startColor, startColor, endColor,
		ps.NoiseFactorX.Eval(), ps.NoiseFactorY.Eval(),
		ps.sprite,
	}
}

func FromFile(path string, dest *ebiten.Image) (*ParticleSystem, error) {
	particleSystem := &ParticleSystem{}
	errs.MustSingle(fileio.UnmarshalStruct(path, particleSystem))
	particleSystem.particles = make([]*Particle, 0)
	particleSystem.surf = ebiten.NewImage(particleSystem.ImageWidth, particleSystem.ImageHeight)
	spritePath := errs.Must(filepath.Localize(particleSystem.SpritePath))
	particleSystem.sprite = errs.MustNewImageFromFile(spritePath)
	particleSystem.layer = dest
	return particleSystem, nil
}

type ParticleBurst struct {
	Count int     `yaml:"Count"`
	Time  float64 `yaml:"Time"`
}
