package stage

import (
	"image/color"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type effectNode struct {
	canvas *Canvas
	pos    gmath.Vec
	image  resource.ImageID
	anim   *ge.Animation
	clr    color.RGBA
}

func newEffectNode(canvas *Canvas, pos gmath.Vec, image resource.ImageID, clr color.RGBA) *effectNode {
	return &effectNode{
		canvas: canvas,
		pos:    pos,
		image:  image,
		clr:    clr,
	}
}

func (e *effectNode) Init(scene *ge.Scene) {
	sprite := scene.NewSprite(e.image)
	sprite.Pos.Base = &e.pos
	var colorScale ge.ColorScale
	colorScale.SetColor(e.clr)
	sprite.SetColorScale(colorScale)
	e.canvas.AddSprite(sprite)
	if e.anim == nil {
		e.anim = ge.NewAnimation(sprite, -1)
	}
}

func (e *effectNode) IsDisposed() bool {
	return e.anim.IsDisposed()
}

func (e *effectNode) Dispose() {
	e.anim.Sprite().Dispose()
}

func (e *effectNode) Update(delta float64) {
	if e.anim.Tick(delta) {
		e.Dispose()
		return
	}
}
