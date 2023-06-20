package stage

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/sinecord/gamedata"
)

type waveNode struct {
	canvas *Canvas
	color  ge.ColorScale

	shape    gamedata.Shape
	x        float32
	y        float32
	t        float64
	r        float64
	duration float64

	EventFinished gsignal.Event[float64]

	disposed bool
}

func newWaveNode(canvas *Canvas, shape gamedata.Shape, pos gmath.Vec, clr color.RGBA, duration float64) *waveNode {
	var colorScale ge.ColorScale
	colorScale.SetColor(clr)

	return &waveNode{
		shape:    shape,
		canvas:   canvas,
		duration: duration,
		color:    colorScale,
		x:        float32(pos.X),
		y:        float32(pos.Y),
	}
}

func (n *waveNode) IsDisposed() bool { return n.disposed }

func (n *waveNode) Dispose() {
	n.disposed = true
}

func (n *waveNode) Update(delta float64) {
	n.t = gmath.ClampMax(n.t+delta, n.duration)
	n.r = math.Sqrt(n.t*0.9) * n.canvas.ctx.Scaler.Factor
	if n.t == n.duration {
		n.Dispose()
		n.EventFinished.Emit(n.r)
		return
	}
}

func (n *waveNode) Draw(screen *ebiten.Image) {
	var angle gmath.Rad
	if n.shape != gamedata.ShapeCircle {
		angle = gmath.Rad(1.5 * (n.t / n.duration))
	}
	r := float32(n.r)
	n.canvas.drawShape(screen, n.shape, n.x, n.y, r, angle, n.color)
}
