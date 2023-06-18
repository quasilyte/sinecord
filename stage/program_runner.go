package stage

import "github.com/quasilyte/sinecord/synthdb"

type SynthProgram struct {
	Length      float64
	Instruments []SynthProgramInstrument
}

type SynthProgramInstrument struct {
	ID     int // The channel is identical
	Index  int
	Func   func(x float64) float64
	Period float64
	Kind   synthdb.InstrumentKind
}

type programRunner struct {
	delays []float64
	events []noteActivation
}

type noteActivation struct {
	index int
	id    int
	t     float64
}

func (r *programRunner) RunProgram(prog SynthProgram) []noteActivation {
	r.delays = r.delays[:0]
	r.events = r.events[:0]

	for _, inst := range prog.Instruments {
		r.delays = append(r.delays, inst.Period)
	}

	const stepsPerSecond = 180
	dt := 1.0 / stepsPerSecond
	t := 0.0

	for t < prog.Length {
		t += dt

		for i, inst := range prog.Instruments {
			delay := r.delays[i]
			if delay > dt {
				r.delays[i] -= dt
				continue
			}
			extraTime := dt - delay
			realTime := t + extraTime
			r.delays[i] = inst.Period - extraTime
			r.events = append(r.events, noteActivation{
				index: i,
				id:    inst.ID,
				t:     realTime,
			})
		}
	}

	return r.events
}
