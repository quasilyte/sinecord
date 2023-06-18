package stage

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/synthdb"
)

type waveShape int

const (
	waveCircle waveShape = iota
	waveSquare
	waveTriangle
	waveHexagon
	waveStar
	waveCross
)

func instrumentWaveShape(kind synthdb.InstrumentKind) waveShape {
	switch kind {
	case synthdb.BassInstrument:
		return waveStar
	case synthdb.KeyboardInstrument:
		return waveSquare
	case synthdb.BrassInstrument:
		return waveTriangle
	case synthdb.StringInstrument:
		return waveHexagon
	case synthdb.DrumInstrument:
		return waveCircle
	default:
		return waveCross
	}
}

type waveNode struct {
	canvas *Canvas
	color  ge.ColorScale

	shape    waveShape
	x        float32
	y        float32
	t        float64
	r        float64
	duration float64

	disposed bool
}

func newWaveNode(canvas *Canvas, shape waveShape, pos gmath.Vec, clr color.RGBA, duration float64) *waveNode {
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

func (n *waveNode) Dispose() { n.disposed = true }

func (n *waveNode) Update(delta float64) {
	n.t += delta
	n.r = math.Sqrt(n.t) * n.canvas.ctx.PlotScale
	if n.t >= n.duration {
		n.Dispose()
		return
	}
}

func (n *waveNode) Draw(screen *ebiten.Image) {
	var angle gmath.Rad
	if n.shape != waveCircle {
		angle = gmath.Rad(1.5 * (n.t / n.duration))
	}
	r := float32(n.r)
	n.canvas.drawShape(screen, n.shape, n.x, n.y, r, angle, n.color)
}
