package actors

import (
	"fmt"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/engine"

	"github.com/ebitengine/debugui"
)

type Transform2D struct {
	Node
	local  *transform
	global *transform
}

func (t *Transform2D) Init() {}

func (t *Transform2D) Update(tree *engine.NodeTree) {
	node, found := tree.GetNode(t.treeID)
	if !found {
		fmt.Errorf("could not find tree node for Transform2D with id %s", t.treeID)
		return
	}

	parentNode := node.GetParent()
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
	ctx.SetGridLayout(make([]int, 1), make([]int, 1))
	ctx.Text("Transform")
	// Write some nice rendering for transforms
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

func (t *Transform2D) SetAngle(angle float64) {
	t.local.setAngle(angle)
}

func (t *Transform2D) SetScale(scaleX, scaleY float64) {
	t.local.setScale(scaleX, scaleY)
}

func NewTransform2D(x, y, angle, scaleX, scaleY float64) *Transform2D {
	return &Transform2D{
		local:  newTransform(x, y, angle, scaleX, scaleY),
		global: newTransform(x, y, angle, scaleX, scaleY),
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
		transformMat: t.transformMat.TimesMat(&other.transformMat),
		origin:       t.origin.Plus(&other.origin),
	}
}

func (t *transform) generateTransformationMatrix() maths.Mat2x2 {
	rot := maths.Mat2x2FromRot(t.angle)
	scale := maths.NewMatrix(t.scaleX, 0, t.scaleY, 0)
	return scale.TimesMat(&rot)
}

func (t *transform) getPos() (float64, float64) {
	return t.origin.X, t.origin.Y
}

func (t *transform) setPos(x, y float64) {
	t.origin = maths.NewVec2(x, y)
}

func (t *transform) setAngle(angle float64) {
	t.angle = angle
	t.transformMat = t.generateTransformationMatrix()
}

func (t *transform) setScale(scaleX, scaleY float64) {
	t.scaleX = scaleX
	t.scaleY = scaleY
	t.transformMat = t.generateTransformationMatrix()
}

func newTransform(x, y, angle, scaleX, scaleY float64) *transform {
	rot := maths.Mat2x2FromRot(angle)
	scale := maths.NewMatrix(scaleX, 0, scaleY, 0)
	transformMat := scale.TimesMat(&rot)
	origin := maths.NewVec2(x, y)
	return &transform{
		transformMat: transformMat,
		origin:       origin,
	}
}
