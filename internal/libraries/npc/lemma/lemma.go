package lemma

import (
	"fmt"
	"image/color"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/particles"
	"time"
)

var DefaultLines = []string{
	"Hi! I'm Lemma, the \nuseless NPC.",
	"For now, all I do is zip \naround and show some \nplaceholder dialogue.",
	"Maybe one day someone \nwill give me some \nuseful features...",
}

// TODO: Add color gradient to vfx body!
// And then, we must implement some kind of actual hint/NPC system
// Also we gotta give the guy some dialogue... Now that will be a bit tricky
// Idk how well it's gonna work but we'll see
// Then we also want to add some sound effects to the NPC itself
// But overall not bad for a days work!

// Some other TODOs
// Redo door animation/visual
// Finish implementing slambox chains
// Redo turret visuals
// Create button prompts
// Create shaders and ambient particle systems for other biomes than just
// basement
// Slambox particles have a draw order defined by the slambox draw order,
// which is not too good of an idea.
// Also, lemma does not in any sense follow the camera position (i guess
// that's kind of fair though)
// Lowkey maybe it would be nice to implement a non-pixel perfect camera..
// (take aarthificial approach)
// IMPROVE PERFORMANCE

type LemmaState int

const (
	IDLE_GONE LemmaState = iota
	APPEARING
	DISAPPEARING
	IDLE_PRESENT
)

const (
	idleEmission      = 8
	appearEmission    = 8
	disappearEmission = 8
	zipTimeMin        = 5.0
	zipTimeMax        = 10.0
)

var dark = []float64{14, 9, 47}
var light = []float64{125, 242, 207}

const (
	// Just to be sure
	appearTime    = 810 * time.Millisecond
	disappearTime = 510 * time.Millisecond
)

type Lemma struct {
	appearTimer          *time.Timer
	disappearTimer       *time.Timer
	zipTimer             *time.Timer
	state                LemmaState
	idleParticleSys      *particles.ParticleSystem
	appearParticleSys    *particles.ParticleSystem
	disappearParticleSys *particles.ParticleSystem
	appearTime           float64
	disappearTime        float64
	vfx                  *vfx
}

func (l *Lemma) Update() {
	dt := 0.0166666666

	l.idleParticleSys.Update()
	l.appearParticleSys.Update()
	l.disappearParticleSys.Update()
	l.idleParticleSys.SetPos(l.vfx.visualX, l.vfx.visualY)
	l.appearParticleSys.SetPos(l.vfx.visualX, l.vfx.visualY)
	l.disappearParticleSys.SetPos(l.vfx.visualX, l.vfx.visualY)

	switch l.state {
	case IDLE_PRESENT:
		l.vfx.Update()

		// if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// 	mouseX, mouseY := ebiten.CursorPosition()
		// 	l.vfx.SeekTarget(float64(mouseX)/rendering.PIXEL_SCALE, float64(mouseY)/rendering.PIXEL_SCALE)
		// }
		if _, ok := threads.Poll(l.zipTimer.C); ok {
			targetX := maths.RandomRange(60, rendering.GAME_WIDTH-60)
			targetY := maths.RandomRange(40, rendering.GAME_HEIGHT-40)
			l.vfx.SeekTarget(targetX, targetY)
			l.zipTimer = time.NewTimer(time.Duration(maths.RandomRange(zipTimeMin, zipTimeMax) * 1e9))
		}

	case IDLE_GONE:

	case APPEARING:
		l.appearTime += dt
		l.vfx.Update()
		t := maths.Clamp(l.appearTime*(5.0/4.0), 0, 1)
		s := maths.CubicIn(t)
		l.vfx.color = color.RGBA{
			uint8(maths.Lerp(dark[0], light[0], s)),
			uint8(maths.Lerp(dark[1], light[1], s)),
			uint8(maths.Lerp(dark[2], light[2], s)),
			255,
		}

		l.vfx.shakeStrength = maths.Lerp(0, 7, t)
		l.vfx.size = maths.Lerp(l.vfx.size, 1, 0.08)

		if _, ok := threads.Poll(l.appearTimer.C); ok {
			l.vfx.color = color.RGBA{255, 253, 240, 255}
			l.state = IDLE_PRESENT
			l.idleParticleSys.Play()
			l.appearParticleSys.Stop()
			l.vfx.shakeStrength = 0
			l.appearTime = 0
		}
	case DISAPPEARING:
		l.disappearTime += dt
		l.vfx.Update()

		t := maths.Clamp(l.disappearTime*2, 0, 1)
		s := maths.CubicIn(t)
		l.vfx.color = color.RGBA{
			uint8(maths.Lerp(light[0], dark[0], t)),
			uint8(maths.Lerp(light[1], dark[1], t)),
			uint8(maths.Lerp(light[2], dark[2], t)),
			255,
		}

		l.vfx.shakeStrength = maths.Lerp(7, 0, t)
		l.vfx.size = maths.Lerp(1, 0, s)

		if _, ok := threads.Poll(l.disappearTimer.C); ok {
			l.state = IDLE_GONE
			l.idleParticleSys.Stop()
			l.disappearParticleSys.Stop()
			l.disappearTime = 0
		}
	}
	l.vfx.Update()
}

func (l *Lemma) Hide() {
	l.disappearTimer = time.NewTimer(disappearTime)
	l.state = DISAPPEARING
	l.disappearParticleSys.Play()
}

func (l *Lemma) Reveal() {
	l.appearTimer = time.NewTimer(appearTime)
	l.state = APPEARING
	l.appearParticleSys.Play()
}

// Returns whether we was turned on or off
func (l *Lemma) ToggleVisibility() bool {
	fmt.Println("Brh?")
	fmt.Println(l.state)
	if l.state == IDLE_GONE {
		l.Reveal()
		fmt.Println("I was hidden")
		return true
	} else if l.state == IDLE_PRESENT {
		fmt.Println("I was present")
		l.Hide()
		return false
	}
	return false
}

func (l *Lemma) Draw(ctx rendering.Ctx) {
	l.idleParticleSys.Draw(ctx)
	l.appearParticleSys.Draw(ctx)
	l.disappearParticleSys.Draw(ctx)
	switch l.state {
	case IDLE_PRESENT:
		l.vfx.Draw(ctx)
	case IDLE_GONE:
	case DISAPPEARING:
		l.vfx.Draw(ctx)
	case APPEARING:
		l.vfx.Draw(ctx)
	}
}

func (l *Lemma) X() float64 {
	return l.vfx.visualX
}

func (l *Lemma) Y() float64 {
	return l.vfx.visualY
}

func NewLemma(x, y float64) *Lemma {
	newLemma := Lemma{
		vfx: newVfx([]float64{
			5.0, 5.0, 4.0, 3.5, 3.0, 2.5, 2.0, 1.5,
		}, x, y, 2.0, 0.55, 0, color.RGBA{uint8(dark[0]), uint8(dark[1]), uint8(dark[2]), 255}),
	}

	idleParticles := errs.Must(assettypes.GetYamlAsset("lemmaIdleParticles")).(*particles.ParticleSystem)
	idleParticles.Init()
	newLemma.idleParticleSys = idleParticles

	appearParticles := errs.Must(assettypes.GetYamlAsset("lemmaAppearParticles")).(*particles.ParticleSystem)
	appearParticles.Init()
	newLemma.appearParticleSys = appearParticles

	disappearParticles := errs.Must(assettypes.GetYamlAsset("lemmaDisappearParticles")).(*particles.ParticleSystem)
	disappearParticles.Init()
	newLemma.disappearParticleSys = disappearParticles

	newLemma.zipTimer = time.NewTimer(time.Duration(maths.RandomRange(zipTimeMin, zipTimeMax) * 1e9))

	return &newLemma
}
