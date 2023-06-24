package stage

import (
	"sort"

	"github.com/quasilyte/sinecord/exprc"
	"github.com/quasilyte/sinecord/gamedata"
)

type SynthProgram struct {
	Length      float64
	Instruments []SynthProgramInstrument
}

type SynthProgramInstrument struct {
	ID     int // The channel is identical
	Index  int
	Func   *exprc.FuncRunner
	Period float64
	Kind   gamedata.InstrumentKind
}

type programRunner struct {
	events []noteActivation
}

type noteActivation struct {
	index int
	id    int
	t     float64
}

func (r *programRunner) RunProgram(prog SynthProgram) []noteActivation {
	r.events = r.events[:0]

	for i, inst := range prog.Instruments {
		for t := inst.Period; t < prog.Length; t += inst.Period {
			r.events = append(r.events, noteActivation{
				index: i,
				id:    inst.ID,
				t:     t,
			})
		}
	}

	// This whole thing can be done without an extra sorting
	// if we would construct the events slice in already sorted form.
	// But since the number of events is usually relatively small
	// and it's a jam game, I don't want to bother optimizing this.
	// The current solution at least provides the precise timings
	// using a constant t step for every instrument.
	if len(prog.Instruments) > 1 {
		sort.SliceStable(r.events, func(i, j int) bool {
			return r.events[i].t < r.events[j].t
		})
	}

	return r.events
}
