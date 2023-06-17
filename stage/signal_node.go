package stage

import (
	"image/color"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
)

type signalNode struct {
	pos gmath.Vec

	canvas *Canvas

	sprite *ge.Sprite

	clr color.RGBA
}

func newSignalNode(canvas *Canvas, clr color.RGBA) *signalNode {
	return &signalNode{
		canvas: canvas,
		clr:    clr,
	}
}

func (n *signalNode) Init(scene *ge.Scene) {
	n.sprite = scene.NewSprite(assets.ImageSignal)
	var cs ge.ColorScale
	cs.SetColor(n.clr)
	n.sprite.SetColorScale(cs)
	n.sprite.Pos.Base = &n.pos
	n.canvas.AddSprite(n.sprite)
}

func (n *signalNode) IsDisposed() bool {
	return n.sprite.IsDisposed()
}

func (n *signalNode) Dispose() {
	n.sprite.Dispose()
}

func (n *signalNode) Update(delta float64) {}
