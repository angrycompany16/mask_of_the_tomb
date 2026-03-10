package transform2D

import (
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/node"
	"mask_of_the_tomb/internal/engine/servers"

	"github.com/ebitengine/debugui"
)

type Option func(*Transform2D)

type Transform2D struct {
	node.Node // Actually, should this be a pointer?
	local     *transform
	global    *transform
}

func (t *Transform2D) Init() {
	t.Node.Init()
}

func (t *Transform2D) Update(servers *servers.Servers) {
	t.Node.Update(servers) // best practice

	parentNode := t.Node.GetNode().GetParent()
	if parentNode == nil {
		return
	}

	if parentTransform, ok := engine.GetActor[*Transform2D](*parentNode.GetValue()); ok {
		t.global = t.local.times(parentTransform.global)
	} else {
		t.global = t.local
	}
}

func (t *Transform2D) DrawInspector(ctx *debugui.Context) {
	ctx.SetGridLayout([]int{0}, []int{0})

	ctx.Header("Transform", false, func() {
		ctx.SetGridLayout([]int{-1, -1, -1}, []int{0, 0, 0})
		ctx.Text("Position")
		ctx.NumberFieldF(&t.local.origin.X, 10, 0)
		ctx.NumberFieldF(&t.local.origin.Y, 10, 0)

		ctx.Text("Scale")
		ctx.NumberFieldF(&t.local.scaleX, 0.1, 2).On(t.local.recompute)
		ctx.NumberFieldF(&t.local.scaleY, 0.1, 2).On(t.local.recompute)

		ctx.SetGridLayout([]int{-1, -2}, []int{0, 0})
		ctx.Text("Angle")
		ctx.NumberFieldF(&t.local.angle, 0.01, 2).On(t.local.recompute)
	})

	t.Node.DrawInspector(ctx)
}

func (t *Transform2D) GetPos(local bool) (float64, float64) {
	if local {
		return t.local.getPos()
	}
	return t.global.getPos()
}

func (t *Transform2D) SetPos(x, y float64) {
	t.local.setPos(x, y)
}

func (t *Transform2D) GetAngle(local bool) float64 {
	if local {
		return t.local.getAngle()
	}
	return t.global.getAngle()
}

func (t *Transform2D) SetAngle(angle float64) {
	t.local.setAngle(angle)
}

func (t *Transform2D) GetScale(local bool) (float64, float64) {
	if local {
		return t.local.getScale()
	}
	return t.global.getScale()
}

func (t *Transform2D) SetScale(scaleX, scaleY float64) {
	t.local.setScale(scaleX, scaleY)
}

// Right now we are instantiating the node as a zero object, which
// is fine-ish (we don't get any nil references), but not great
// (this could break if we change node, and is not very clean)
func NewTransform2D(node node.Node, options ...Option) *Transform2D {
	t := defaultTransform2D(node)

	for _, option := range options {
		option(t)
	}

	return t
}

func defaultTransform2D(node node.Node) *Transform2D {
	return &Transform2D{
		Node:   node,
		local:  newTransform(0, 0, 0, 1, 1),
		global: newTransform(0, 0, 0, 1, 1),
	}
}

func WithPos(x, y float64) Option {
	return func(t *Transform2D) {
		t.local.origin.X = x
		t.local.origin.Y = y
		t.global.origin.X = x
		t.global.origin.Y = y
	}
}

func WithAngle(angle float64) Option {
	return func(t *Transform2D) {
		t.local.angle = angle
		t.global.angle = angle
	}
}

func WithScale(scaleX, scaleY float64) Option {
	return func(t *Transform2D) {
		t.local.scaleX = scaleX
		t.global.scaleY = scaleY
	}
}

// Could cache transformation multiplication results for
// possibly slight performance boost

type transform struct {
	transformMat   maths.Mat2x2
	angle          float64
	scaleX, scaleY float64
	origin         maths.Vec2
}

func (t *transform) times(other *transform) *transform {
	return &transform{
		angle:        t.angle + other.angle,
		scaleX:       t.scaleX * other.scaleX,
		scaleY:       t.scaleY * other.scaleY,
		transformMat: t.transformMat.TimesMat(other.transformMat),
		origin:       other.origin.Plus(other.transformMat.TimesVec(t.origin)),
	}
}

func (t *transform) recompute() {
	t.transformMat = t.generateTransformationMatrix()
}

func (t *transform) generateTransformationMatrix() maths.Mat2x2 {
	rot := maths.Mat2x2FromRot(t.angle)
	scale := maths.NewMatrix(t.scaleX, 0, 0, t.scaleY)
	return scale.TimesMat(rot)
}

func (t *transform) getPos() (float64, float64) {
	return t.origin.X, t.origin.Y
}

func (t *transform) setPos(x, y float64) {
	t.origin = maths.NewVec2(x, y)
}

func (t *transform) getAngle() float64 {
	return t.angle
}

func (t *transform) setAngle(angle float64) {
	t.angle = angle
	t.transformMat = t.generateTransformationMatrix()
}

func (t *transform) getScale() (float64, float64) {
	return t.scaleX, t.scaleY
}

func (t *transform) setScale(scaleX, scaleY float64) {
	t.scaleX = scaleX
	t.scaleY = scaleY
	t.transformMat = t.generateTransformationMatrix()
}

func newTransform(x, y, angle, scaleX, scaleY float64) *transform {
	newTransform := transform{
		angle:  angle,
		scaleX: scaleX,
		scaleY: scaleY,
		origin: maths.NewVec2(x, y),
	}

	newTransform.recompute()
	return &newTransform
}
