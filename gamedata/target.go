package gamedata

import (
	"github.com/quasilyte/gmath"
)

type TargetSize int

const (
	SmallTarget TargetSize = iota
	NormalTarget
	BigTarget
)

type Target struct {
	Pos gmath.Vec

	Instrument InstrumentKind

	Size TargetSize
}
