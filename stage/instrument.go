package stage

import "github.com/quasilyte/sinecord/synthdb"

type instrument struct {
	fx    string
	oldFx string

	compiledFx func(x float64) float64

	instrumentIndex int

	period  float64
	enabled bool
	volume  int32
	kind    synthdb.InstrumentKind
}

func (inst *instrument) SetFx(fx string) {
	inst.oldFx = inst.fx
	inst.fx = fx
}
