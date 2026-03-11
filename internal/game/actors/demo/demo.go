package demo

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"math"

	"github.com/ebitengine/debugui"
)

type Option func(*Demo)

type Demo struct {
	sprite.Sprite
	t              float64
	onlyRotate     bool
	scaleX, scaleY float64
	x, y           float64
	angle          float64
}

func (d *Demo) Update(servers *engine.Servers) {
	d.t += 0.016666
	d.angle = 0.5 * d.t
	d.Sprite.SetAngle(d.angle)

	if d.onlyRotate {
	} else {
		d.scaleX = 2 * math.Sin(d.t)
		d.scaleY = 1 + 2*math.Cos(d.t)
		d.Sprite.SetScale(d.scaleX, d.scaleY)
	}

	d.Sprite.Update(servers)
}

func (d *Demo) DrawInspector(ctx *debugui.Context) {
	d.Sprite.DrawInspector(ctx)
}

func NewDemo(sprite sprite.Sprite, options ...Option) *Demo {
	d := defaultDemo(sprite)

	for _, option := range options {
		option(d)
	}

	return d
}

func defaultDemo(sprite sprite.Sprite) *Demo {
	return &Demo{
		Sprite:     sprite,
		onlyRotate: false,
	}
}

// At some point it'll get tedious to write these...
func WithOnlyRotate(onlyRotate bool) Option {
	return func(d *Demo) {
		d.onlyRotate = onlyRotate
	}
}
