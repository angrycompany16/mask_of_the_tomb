package particles

import (
	"image/color"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/utils"
	"math"
	"time"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
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

type ParticleSystem struct {
	*graphic.Graphic
	particles       []*Particle
	Bursts          []*Burst
	GlobalSpace     bool                `debug:"auto"`
	Emission        float64             `debug:"auto"`
	Angle           float64             `debug:"auto"`
	SpawnPosX       maths.RandomFloat64 `debug:"auto"`
	SpawnPosY       maths.RandomFloat64 `debug:"auto"`
	SpawnVelX       maths.RandomFloat64 `debug:"auto"`
	SpawnVelY       maths.RandomFloat64 `debug:"auto"`
	SpawnAngle      maths.RandomFloat64 `debug:"auto"`
	SpawnAngularVel maths.RandomFloat64 `debug:"auto"`
	AirFriction     maths.RandomFloat64 `debug:"auto"`
	StartScale      maths.RandomFloat64 `debug:"auto"`
	EndScale        maths.RandomFloat64 `debug:"auto"`
	Lifetime        maths.RandomFloat64 `debug:"auto"`
	NoiseFactorX    maths.RandomFloat64 `debug:"auto"`
	NoiseFactorY    maths.RandomFloat64 `debug:"auto"`
	t               float64             `debug:"auto"`
	// TODO: Use ebiten.color by making embedded struct with custom unmarshaler
	StartColor  [4]uint8                            `debug:"auto"`
	EndColor    [4]uint8                            `debug:"auto"`
	ImageWidth  int                                 `debug:"auto"`
	ImageHeight int                                 `debug:"auto"`
	SpritePath  string                              `debug:"auto"`
	surf        *ebiten.Image                       // The image that the particles are drawn onto
	imageAsset  *assetloader.AssetRef[ebiten.Image] // The sprite for the particles
	// layer       *ebiten.Image
	gizmosImage *ebiten.Image
	burstTimers []*time.Timer
	isPlaying   bool   `debug:"auto"`
	layer       string `debug:"auto"`
	drawOrder   int    `debug:"auto"`
}

// This could maybe just be part of the constructor
func (ps *ParticleSystem) OnTreeAdd(node *engine.Node, cmd *engine.Commands) {
	ps.Graphic.OnTreeAdd(node, cmd)
	ps.imageAsset = assetloader.StageAsset[ebiten.Image](
		cmd.AssetLoader(),
		ps.SpritePath,
		assettypes.NewImageAsset(ps.SpritePath),
	)
}

func (ps *ParticleSystem) Init(cmd *engine.Commands) {
	ps.Graphic.Init(cmd)
	ps.gizmosImage.Fill(color.RGBA{12, 123, 222, 123})
	ps.Play()
}

func (ps *ParticleSystem) DrawGizmo(cmd *engine.Commands) {
	ps.Graphic.DrawGizmo(cmd)

	worldX, worldY := ps.Transform2D.GetPos(false)
	camX, camY := ps.GetCamera().WorldToCam(worldX, worldY, false)
	cmd.Renderer().Request(opgen.Pos(ps.gizmosImage, camX, camY, 0.5, 0.5), ps.gizmosImage, "Overlay", ps.drawOrder+1)
}

func (ps *ParticleSystem) Update(cmd *engine.Commands) {
	ps.Graphic.Update(cmd)
	ps.t += 0.016666666667
	ps.surf.Clear()

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
		if _, ok := utils.PollThread(ps.burstTimers[i].C); ok {
			for range burst.Count {
				ps.particles = append(ps.particles, ps.newParticle())
			}
		}
	}

	if ps.Emission > 0 && ps.t > 1/ps.Emission {
		nParticles := int(math.Floor(ps.t * ps.Emission))
		for range nParticles {
			ps.particles = append(ps.particles, ps.newParticle())
			ps.t = 0
		}
	}

	if ps.GlobalSpace {
		for _, particle := range ps.particles {
			camX, camY := ps.GetCamera().WorldToCam(particle.posX, particle.posY, true)
			c, op := particle.makeOp(camX, camY)
			cmd.Renderer().RequestColorM(c, op, particle.sprite, ps.layer, ps.drawOrder)
		}
		return
	}

	// Local-space rendering is very cursed ngl
	for _, particle := range ps.particles {
		c, op := particle.makeOp(particle.posX, particle.posY)
		colorm.DrawImage(ps.surf, particle.sprite, c, op)
	}

	worldX, worldY := ps.Transform2D.GetPos(false)
	angle := ps.Transform2D.GetAngle(false)
	scaleX, scaleY := ps.Transform2D.GetScale(false)
	camX, camY := ps.GetCamera().WorldToCam(worldX, worldY, true)
	cmd.Renderer().Request(opgen.PosRotScale(ps.surf, camX, camY, angle, scaleX, scaleY, 0.5, 0.5), ps.surf, ps.layer, ps.drawOrder)
}

func (ps *ParticleSystem) DrawInspector(ctx *debugui.Context) {
	ps.Graphic.DrawInspector(ctx)
	utils.RenderFieldsAuto(ctx, ps)
}

func (ps *ParticleSystem) Play() {
	ps.t = 0
	ps.isPlaying = true
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

func (ps *ParticleSystem) newParticle() *Particle {
	var x, y float64
	gposX, gposY := ps.Transform2D.GetPos(false)
	if ps.GlobalSpace {
		x = ps.SpawnPosX.Eval() + gposX
		y = ps.SpawnPosY.Eval() + gposY
	} else {
		x = ps.SpawnPosX.Eval() + float64(ps.ImageWidth)/2
		y = ps.SpawnPosY.Eval() + float64(ps.ImageHeight)/2
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
		ps.imageAsset.Value(),
	}
}

func NewParticleSystem(
	graphic *graphic.Graphic,
	Bursts []*Burst,
	GlobalSpace bool,
	Emission float64,
	Angle float64,
	SpawnPosX maths.RandomFloat64,
	SpawnPosY maths.RandomFloat64,
	SpawnVelX maths.RandomFloat64,
	SpawnVelY maths.RandomFloat64,
	SpawnAngle maths.RandomFloat64,
	SpawnAngularVel maths.RandomFloat64,
	AirFriction maths.RandomFloat64,
	StartScale maths.RandomFloat64,
	EndScale maths.RandomFloat64,
	Lifetime maths.RandomFloat64,
	NoiseFactorX maths.RandomFloat64,
	NoiseFactorY maths.RandomFloat64,
	StartColor [4]uint8,
	EndColor [4]uint8,
	ImageWidth int,
	ImageHeight int,
	SpritePath string,
	layer string,
	drawOrder int,
) *ParticleSystem {
	return &ParticleSystem{
		Graphic:         graphic,
		particles:       make([]*Particle, 0),
		Bursts:          Bursts,
		GlobalSpace:     GlobalSpace,
		Emission:        Emission,
		Angle:           Angle,
		SpawnPosX:       SpawnPosX,
		SpawnPosY:       SpawnPosY,
		SpawnVelX:       SpawnVelX,
		SpawnVelY:       SpawnVelY,
		SpawnAngle:      SpawnAngle,
		SpawnAngularVel: SpawnAngularVel,
		AirFriction:     AirFriction,
		StartScale:      StartScale,
		EndScale:        EndScale,
		Lifetime:        Lifetime,
		NoiseFactorX:    NoiseFactorX,
		NoiseFactorY:    NoiseFactorY,
		t:               0,
		StartColor:      StartColor,
		EndColor:        EndColor,
		ImageWidth:      ImageWidth,
		ImageHeight:     ImageHeight,
		SpritePath:      SpritePath,
		surf:            ebiten.NewImage(ImageWidth, ImageHeight),
		gizmosImage:     ebiten.NewImage(ImageWidth, ImageHeight),
		burstTimers:     make([]*time.Timer, len(Bursts)),
		isPlaying:       false,
		layer:           layer,
		drawOrder:       drawOrder,
	}
}

type Burst struct {
	Count int     `yaml:"Count"`
	Time  float64 `yaml:"Time"`
}
