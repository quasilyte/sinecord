package stage

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/synthdb"
)

type waveShape int

const (
	waveCircle waveShape = iota
	waveSquare
	waveTriangle
	waveOctagon
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
		return waveCross
	case synthdb.DrumInstrument:
		return waveCircle
	default:
		return waveOctagon
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
	var p vector.Path
	r := float32(n.r)
	switch n.shape {
	case waveCircle:
		p.Arc(n.x, n.y, r, 0, 2*math.Pi, vector.Clockwise)
	case waveSquare:
		angle := gmath.Rad(1.5 * (n.t / n.duration))
		p.MoveTo(n.translate(rotate(-r, -r, angle)))
		p.LineTo(n.translate(rotate(-r, +r, angle)))
		p.LineTo(n.translate(rotate(+r, +r, angle)))
		p.LineTo(n.translate(rotate(+r, -r, angle)))
		p.Close()
	case waveTriangle:
		angle := gmath.Rad(1.5 * (n.t / n.duration))
		p.MoveTo(n.translate(rotate(-r, -r, angle)))
		p.LineTo(n.translate(rotate(0, r, angle)))
		p.LineTo(n.translate(rotate(r, -r, angle)))
		p.Close()
	case waveOctagon:
		angle := gmath.Rad(1.5 * (n.t / n.duration))
		r2 := r / 2
		p.MoveTo(n.translate(rotate(-r, -r2, angle)))
		p.LineTo(n.translate(rotate(-r, r2, angle)))
		p.LineTo(n.translate(rotate(-r2, r, angle)))
		p.LineTo(n.translate(rotate(r2, r, angle)))
		p.LineTo(n.translate(rotate(r, r2, angle)))
		p.LineTo(n.translate(rotate(r, -r2, angle)))
		p.LineTo(n.translate(rotate(r2, -r, angle)))
		p.LineTo(n.translate(rotate(-r2, -r, angle)))
		p.Close()
	case waveStar:
		angle := gmath.Rad(1.5 * (n.t / n.duration))
		r3 := r / 3
		p.MoveTo(n.translate(rotate(-r3, -r3, angle)))
		p.LineTo(n.translate(rotate(-r, 0, angle)))
		p.LineTo(n.translate(rotate(-r3, r3, angle)))
		p.LineTo(n.translate(rotate(0, r, angle)))
		p.LineTo(n.translate(rotate(r3, r3, angle)))
		p.LineTo(n.translate(rotate(r, 0, angle)))
		p.LineTo(n.translate(rotate(r3, -r3, angle)))
		p.LineTo(n.translate(rotate(0, -r, angle)))
		p.Close()
	case waveCross:
		angle := gmath.Rad(1.5 * (n.t / n.duration))
		r2 := r / 2
		p.MoveTo(n.translate(rotate(-r, -r2, angle)))
		p.LineTo(n.translate(rotate(-r, r2, angle)))
		p.LineTo(n.translate(rotate(-r2, r2, angle)))
		p.LineTo(n.translate(rotate(-r2, r, angle)))
		p.LineTo(n.translate(rotate(r2, r, angle)))
		p.LineTo(n.translate(rotate(r2, r2, angle)))
		p.LineTo(n.translate(rotate(r, r2, angle)))
		p.LineTo(n.translate(rotate(r, -r2, angle)))
		p.LineTo(n.translate(rotate(r2, -r2, angle)))
		p.LineTo(n.translate(rotate(r2, -r, angle)))
		p.LineTo(n.translate(rotate(-r2, -r, angle)))
		p.LineTo(n.translate(rotate(-r2, -r2, angle)))
		p.Close()
	}
	n.canvas.DrawPath(screen, p, 1, n.color)
}

func (n *waveNode) translate(x, y float32) (float32, float32) {
	return x + n.x, y + n.y
}

func rotate(x, y float32, angle gmath.Rad) (float32, float32) {
	sine := float32(angle.Sin())
	cosi := float32(angle.Cos())
	return x*cosi - y*sine, x*sine + y*cosi
}
