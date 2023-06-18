package stage

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/synthdb"
)

type TargetSize int

const (
	SmallTarget TargetSize = iota
)

type Target struct {
	Pos gmath.Vec

	Instrument synthdb.InstrumentKind

	Size TargetSize
}
