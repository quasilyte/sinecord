package stage

import (
	"github.com/quasilyte/sinecord/exprc"
	"github.com/quasilyte/sinecord/gamedata"
)

type instrument struct {
	fx    string
	oldFx string

	periodFunc string

	compiledFx *exprc.FuncRunner

	instrumentIndex int
	patchNumber     int32

	period       float64
	oldPeriod    float64
	enabled      bool
	mappedVolume int32
	volume       float64
	kind         gamedata.InstrumentKind
}

func (inst *instrument) SetPeriod(src string, period float64) {
	inst.oldPeriod = inst.period
	inst.period = period
	inst.periodFunc = src
}

func (inst *instrument) SetFx(fx string) {
	inst.oldFx = inst.fx
	inst.fx = fx
}
