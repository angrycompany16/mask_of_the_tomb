package entities

import (
	"image/color"
	"mask_of_the_tomb/internal/core/arrays"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/shaders"
	"math"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	dt             = 0.016666666666666667
	tol            = 0.001
	ropeResolution = 8
)

type RopeSegment struct {
	VerletObject
	next     *RopeSegment
	length   float64
	anchored bool
}

func (r *RopeSegment) Update(forceX, forceY float64) {
	// TODO: Maybe add some smoothing/damping force?
	GRAVITY := 10.0

	velX, velY := r.GetVelocity()

	if !r.anchored {
		r.VerletObject.Update(forceX-velX/2, forceY+GRAVITY-velY/2)
	}

	r.ApplyConstraints()
}

func (r *RopeSegment) ApplyConstraints() {
	if r.next == nil {
		return
	}

	dist := maths.Length(r.x-r.next.x, r.y-r.next.y)
	signedResidue := (dist - r.length) / 2.0

	// Create normalized direction vector
	dir := arrays.MapSlice(
		[]float64{r.next.x - r.x, r.next.y - r.y},
		func(x float64) float64 { return x / dist },
	)

	if !r.anchored {
		if r.next.anchored {
			r.x += signedResidue * dir[0] * 2.0
			r.y += signedResidue * dir[1] * 2.0
		} else {
			r.x += signedResidue * dir[0]
			r.y += signedResidue * dir[1]
		}
	}

	if !r.next.anchored {
		if r.anchored {
			r.next.x -= signedResidue * dir[0] * 2.0
			r.next.y -= signedResidue * dir[1] * 2.0
		} else {
			r.next.x -= signedResidue * dir[0]
			r.next.y -= signedResidue * dir[1]
		}
	}
}

func (r *RopeSegment) GetAngle() float64 {
	if r.next == nil {
		return 0
	}

	dx := r.next.x - r.x
	dy := r.next.y - r.y

	return math.Atan2(dx, dy)
}

type VerletObject struct {
	x, y         float64
	prevX, prevY float64
}

func (v *VerletObject) Update(accelX, accelY float64) {
	nextX := 2*v.x - v.prevX + accelX*dt
	nextY := 2*v.y - v.prevY + accelY*dt

	v.prevX = v.x
	v.prevY = v.y

	v.x = nextX
	v.y = nextY
}

func (v *VerletObject) GetVelocity() (float64, float64) {
	return v.x - v.prevX, v.y - v.prevY
}

type Lantern struct {
	sprite *ebiten.Image
	x, y   float64
	rope   []*RopeSegment
	Light  *shaders.Light
}

func (l *Lantern) Update(playerX, playerY, playerVelX, playerVelY float64) {
	for _, ropeSeg := range l.rope {
		dist := maths.Clamp(maths.Length(playerX-ropeSeg.x, playerY-ropeSeg.y), 1, 1000)
		forceX := playerVelX / math.Pow(dist, 2) * 400
		forceY := playerVelY / math.Pow(dist, 2) * 400
		ropeSeg.Update(forceX, forceY)
	}
	l.Light.X = l.rope[len(l.rope)-1].x - float64(l.sprite.Bounds().Dx()/2)
	l.Light.Y = l.rope[len(l.rope)-1].y - float64(l.sprite.Bounds().Dy()/2)
}

func (l *Lantern) Draw(ctx rendering.Ctx) {
	for _, ropeSeg := range l.rope {
		if ropeSeg.next == nil {
			continue
		}
		vector.StrokeLine(
			ctx.Dst,
			// centered. Maybe instead apply the centering directly to the rope? probably better
			float32(ropeSeg.x),
			float32(ropeSeg.y),
			float32(ropeSeg.next.x),
			float32(ropeSeg.next.y),
			2,
			color.RGBA{128, 128, 128, 255},
			false,
		)
	}

	endPointX := l.rope[len(l.rope)-1].x - float64(l.sprite.Bounds().Dx()/2)
	endPointY := l.rope[len(l.rope)-1].y - float64(l.sprite.Bounds().Dy()/2)
	angle := l.GetRopeEndAngle()
	ebitenrenderutil.DrawAtRotated(l.sprite, ctx.Dst, endPointX, endPointY, -angle, 0.5, 0.5)
}

func (l *Lantern) GetRopeEndAngle() float64 {
	nPoints := 1
	totalAngle := 0.0
	for i := len(l.rope) - 1; i >= len(l.rope)-nPoints-1; i-- {
		totalAngle += l.rope[i].GetAngle()
	}
	return totalAngle / float64(nPoints)
}

func NewLantern(
	entity *ebitenLDTK.Entity,
	tileSize float64,
) *Lantern {
	newLantern := Lantern{}
	newLantern.x, newLantern.y = entity.Px[0], entity.Px[1]
	newLantern.sprite = errs.Must(assettypes.GetImageAsset("lanternSprite"))

	anchorPointField := errs.Must(entity.GetFieldByName("Anchor"))

	anchorX := anchorPointField.Point.X * tileSize
	anchorY := anchorPointField.Point.Y * tileSize

	centerX := newLantern.sprite.Bounds().Dx() / 2
	centerY := newLantern.sprite.Bounds().Dy() / 2

	newLantern.Light = &shaders.Light{
		X:           entity.Px[0],
		Y:           entity.Px[1],
		InnerRadius: 0,
		OuterRadius: 50,
		ZOffset:     0,
		Intensity:   0.5,
		R:           1.0,
		G:           1.0,
		B:           1.0,
	}

	rope := []*RopeSegment{
		{
			VerletObject: VerletObject{
				x:     anchorX + float64(centerX),
				y:     anchorY + float64(centerY),
				prevX: anchorX + float64(centerX),
				prevY: anchorY + float64(centerY),
			},
			length:   ropeResolution,
			anchored: true,
		},
	}

	numRopeSegments := int(math.Ceil(maths.Length(newLantern.x-anchorX, newLantern.y-anchorY) / ropeResolution))

	for i := 1; i < numRopeSegments; i++ {
		posX := maths.Lerp(anchorX, newLantern.x, float64(i)/float64(numRopeSegments)) + float64(centerX)
		posY := maths.Lerp(anchorY, newLantern.y, float64(i)/float64(numRopeSegments)) + float64(centerY)
		rope = append(rope, &RopeSegment{
			VerletObject: VerletObject{
				x: posX, y: posY, prevX: posX, prevY: posY,
			},
			length:   ropeResolution,
			anchored: false,
		})
	}

	for i := 0; i < numRopeSegments-1; i++ {
		rope[i].next = rope[i+1]
	}

	newLantern.rope = rope
	return &newLantern
}
