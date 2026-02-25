package lemma

import (
	"image/color"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/vector64"
	"math"
)

const (
	bobFreqX     = 0.9
	bobFreqY     = 1.5
	length       = 10
	jointSpacing = 2
)

// Both position parameters are modelled as 2nd order ODEs tracking a
// reference. The equations are:
//
//	d(dx) + bX * dx + cX * x = cX * targetX
//	d(dy) + bY * dy + cY * y = cY * targetY
type vfxPosition struct {
	x, y             float64
	targetX, targetY float64
	dx, dy           float64
	gainX            float64
	gainY            float64
	bX               float64
	bY               float64
	cX               float64
	cY               float64
}

func (vp *vfxPosition) Update() {
	dt := 0.016666666666

	vp.dx += (vp.cX*vp.targetX - vp.cX*vp.x - vp.bX*vp.dx) * dt * vp.gainX
	vp.dy += (vp.cY*vp.targetY - vp.cY*vp.y - vp.bY*vp.dy) * dt * vp.gainY

	vp.x += vp.dx * dt
	vp.y += vp.dy * dt
}

func (vp *vfxPosition) SetTarget(x, y float64) {
	vp.targetX = x
	vp.targetY = y
	distX := vp.targetX - vp.x
	distY := vp.targetY - vp.y
	dir := maths.Normalize(distX, distY)
	vp.dx = dir[1] * 9
	vp.dy = -dir[0] * 9
}

func newVFXPosition(x, y float64, nat_freq, damping_coeff float64) vfxPosition {
	c := math.Pow(nat_freq, 2)
	b := 2 * nat_freq * damping_coeff
	return vfxPosition{
		x:       x,
		y:       y,
		targetX: x,
		targetY: y,
		bX:      b,
		bY:      b * 0.8,
		cX:      c,
		cY:      c * 2.0,
		gainX:   1,
		gainY:   1.3,
	}
}

type vfx struct {
	radii            []float64
	jointPositionsX  []float64
	jointPositionsY  []float64
	position         vfxPosition
	visualX, visualY float64
	shakeX, shakeY   float64
	bobX, bobY       float64
	shakeStrength    float64
	t                float64
	color            color.Color
	size             float64
}

func (vfx *vfx) Update() {
	vfx.t += 0.016666666
	vfx.position.Update()

	// Note: This is a purely visual offset, so it doesn't affect calculations
	vfx.bobX = 2 * math.Sin(vfx.t*bobFreqX+0.3)
	vfx.bobY = 1.4 * math.Cos(vfx.t*bobFreqY)

	// Same for this
	vfx.shakeX = maths.RandomRange(-1, 1) * vfx.shakeStrength
	vfx.shakeY = maths.RandomRange(-1, 1) * vfx.shakeStrength

	vfx.visualX = vfx.position.x + vfx.bobX
	vfx.visualY = vfx.position.y + vfx.bobY

	vfx.jointPositionsX[0] = vfx.visualX
	vfx.jointPositionsY[0] = vfx.visualY

	for i := 1; i < len(vfx.radii); i++ {
		vfx.jointPositionsX[i], vfx.jointPositionsY[i] = ConstrainDist(
			vfx.jointPositionsX[i], vfx.jointPositionsY[i],
			vfx.jointPositionsX[i-1], vfx.jointPositionsY[i-1],
			jointSpacing,
		)
	}
}

func ConstrainDist(x, y float64, anchorX, anchorY float64, distance float64) (float64, float64) {
	dx := x - anchorX
	dy := y - anchorY
	if dx == 0 && dy == 0 {
		return x, y
	}
	dir := maths.Normalize(dx, dy)
	return anchorX + dir[0]*distance, anchorY + dir[1]*distance
}

func (vfx *vfx) SeekTarget(x, y float64) {
	vfx.position.SetTarget(x, y)
}

func (vfx *vfx) Draw(ctx rendering.Ctx) {
	for i := range vfx.radii {
		vector64.FillCircle(
			ctx.Dst,
			vfx.jointPositionsX[i]+vfx.shakeX,
			vfx.jointPositionsY[i]+vfx.shakeY,
			vfx.radii[i]*vfx.size, vfx.color, false)
	}

	for i := 0; i < len(vfx.radii)-1; i++ {
		vector64.StrokeLine(
			ctx.Dst,
			vfx.jointPositionsX[i]+vfx.shakeX,
			vfx.jointPositionsY[i]+vfx.shakeY,
			vfx.jointPositionsX[i+1]+vfx.shakeX,
			vfx.jointPositionsY[i+1]+vfx.shakeY,
			vfx.radii[i+1]*1.9, vfx.color, false,
		)
	}
}

func newVfx(radii []float64, headX, headY float64, natural_freq, damping_coeff float64, size float64, color color.Color) *vfx {
	jointPositionsX := make([]float64, len(radii))
	jointPositionsY := make([]float64, len(radii))
	for i := range radii {
		jointPositionsX[i] = headX
		jointPositionsY[i] = headY + float64(i)*jointSpacing
	}
	return &vfx{
		radii:           radii,
		jointPositionsX: jointPositionsX,
		jointPositionsY: jointPositionsY,
		position:        newVFXPosition(headX, headY, natural_freq, damping_coeff),
		size:            size,
		color:           color,
	}
}
