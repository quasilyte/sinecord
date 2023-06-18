package stage

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/styles"
	"github.com/quasilyte/sinecord/synthdb"
)

type targetNode struct {
	board      *Board
	pos        gmath.Vec
	instrument synthdb.InstrumentKind
	disposed   bool
	r          float32
	color      ge.ColorScale
}

func newTargetNode(b *Board, t Target) *targetNode {
	var size float64
	switch t.Size {
	case SmallTarget:
		size = 28.0
	}

	var colorScale ge.ColorScale
	colorScale.SetColor(styles.TargetColor)
	return &targetNode{
		board:      b,
		pos:        b.canvas.scalePos(t.Pos),
		instrument: t.Instrument,
		r:          float32(size / 2),
		color:      colorScale,
	}
}

func (n *targetNode) IsDisposed() bool { return n.disposed }

func (n *targetNode) Dispose() { n.disposed = true }

func (n *targetNode) Update(delta float64) {
}

func (n *targetNode) Draw(screen *ebiten.Image) {
	shape := instrumentWaveShape(n.instrument)
	r := float32(n.r)
	n.board.canvas.drawFilledShape(screen, shape, float32(n.pos.X), float32(n.pos.Y), r, 0, n.color)
}
