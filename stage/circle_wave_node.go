package stage

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type circleWaveNode struct {
	canvas *Canvas
	color  ge.ColorScale

	x        float32
	y        float32
	t        float64
	r        float64
	duration float64

	disposed bool
}

func newCircleWaveNode(canvas *Canvas, pos gmath.Vec, clr color.RGBA, duration float64) *circleWaveNode {
	var colorScale ge.ColorScale
	colorScale.SetColor(clr)

	return &circleWaveNode{
		canvas:   canvas,
		duration: duration,
		color:    colorScale,
		x:        float32(pos.X),
		y:        float32(pos.Y),
	}
}

func (c *circleWaveNode) IsDisposed() bool { return c.disposed }

func (c *circleWaveNode) Dispose() { c.disposed = true }

func (c *circleWaveNode) Update(delta float64) {
	c.t += delta
	c.r = math.Sqrt(c.t) * c.canvas.ctx.PlotScale
	if c.t >= c.duration {
		c.Dispose()
		return
	}
}

func (c *circleWaveNode) Draw(screen *ebiten.Image) {
	var p vector.Path
	p.Arc(c.x, c.y, float32(c.r), 0, 2*math.Pi, vector.Clockwise)
	c.canvas.DrawPath(screen, p, 1, c.color)
}
