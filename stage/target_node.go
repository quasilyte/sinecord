package stage

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/styles"
)

type targetNode struct {
	board      *Board
	pos        gmath.Vec
	instrument gamedata.InstrumentKind
	disposed   bool
	outline    bool
	r          float32
	color      ge.ColorScale
	hp         int
	prevHit    int
}

func newTargetNode(b *Board, t gamedata.Target) *targetNode {
	var size float64
	hp := 1
	switch t.Size {
	case gamedata.SmallTarget:
		size = 15
	case gamedata.NormalTarget:
		size = 30
	case gamedata.BigTarget:
		size = 60
		hp = 2
	default:
		panic("unexpected target size")
	}

	var colorScale ge.ColorScale
	colorScale.SetColor(styles.TargetColor)
	return &targetNode{
		board:      b,
		pos:        t.Pos,
		instrument: t.Instrument,
		r:          float32(size / 2),
		color:      colorScale,
		hp:         hp,
		outline:    t.Outline,
		prevHit:    -1,
	}
}

func (n *targetNode) IsDisposed() bool { return n.disposed }

func (n *targetNode) Dispose() { n.disposed = true }

func (n *targetNode) Update(delta float64) {
}

func (n *targetNode) Draw(screen *ebiten.Image) {
	shape := gamedata.InstrumentShape(n.instrument)
	isOutline := n.outline
	clr := n.color
	if n.instrument == gamedata.AnyInstrument {
		isOutline = true
		clr.SetColor(styles.TargetColorBonus)
	}
	r := float32(n.r)
	if isOutline {
		n.board.canvas.drawShape(screen, shape, float32(n.pos.X), float32(n.pos.Y), r, 0, clr)
	} else {
		n.board.canvas.drawFilledShape(screen, shape, float32(n.pos.X), float32(n.pos.Y), r, 0, clr)
	}
}

func (n *targetNode) OnDamage(instrumentID int) bool {
	if n.prevHit == instrumentID {
		return false
	}
	n.hp--
	n.prevHit = instrumentID
	destroyed := n.hp <= 0
	if destroyed {
		n.Dispose()
	}
	return destroyed
}
