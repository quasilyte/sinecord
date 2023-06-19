package stage

import (
	"fmt"
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/sinecord/exprc"
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/synthdb"
)

type Synthesizer struct {
	scene *ge.Scene

	ctx *Context

	changed bool

	sf *synthdb.SoundFont

	player *musicPlayer

	recompileDelay float64

	instruments []*instrument

	EventRedrawPlotRequest gsignal.Event[int]
}

func NewSynthesizer(ctx *Context, sf *synthdb.SoundFont) *Synthesizer {
	instruments := make([]*instrument, 4)
	for i := range instruments {
		instruments[i] = &instrument{}
	}
	return &Synthesizer{
		ctx:         ctx,
		changed:     true,
		sf:          sf,
		instruments: instruments,
		player:      newMusicPlayer(ctx, instruments),
	}
}

func (s *Synthesizer) Init(scene *ge.Scene) {
	s.scene = scene
}

func (s *Synthesizer) IsDisposed() bool { return false }

func (s *Synthesizer) Update(delta float64) {
	s.recompileDelay = gmath.ClampMin(s.recompileDelay-delta, 0)
	if s.recompileDelay == 0 {
		s.recompileDelay = s.scene.Rand().FloatRange(0.5, 0.8)
		if i := s.needsPlotRedraw(); i != -1 {
			inst := s.instruments[i]
			if inst.fx == "" {
				inst.compiledFx = nil
				s.EventRedrawPlotRequest.Emit(i)
				return
			}
			fn, err := exprc.Compile(inst.fx)
			if err != nil {
				fmt.Printf("exprc: %v\n", err)
				return
			}
			s.changed = true
			inst.compiledFx = fn
			s.EventRedrawPlotRequest.Emit(i)
		}
	}
}

func (s *Synthesizer) HasChanges() bool {
	return s.changed
}

func (s *Synthesizer) ExportTrack() gamedata.Track {
	var t gamedata.Track
	for _, inst := range s.instruments {
		t.Instruments = append(t.Instruments, gamedata.InstrumentSettings{
			Function:       inst.fx,
			Period:         inst.period,
			Volume:         inst.volume,
			InstrumentName: synthdb.TimGM6mb.Instruments[inst.instrumentIndex].Name,
			Enabled:        inst.enabled,
		})
	}
	return t
}

func (s *Synthesizer) CreatePCM(progress *float64) (*SampleSet, SynthProgram) {
	if !s.changed {
		return nil, SynthProgram{}
	}
	s.changed = false
	prog := s.CreateProgram()
	return s.player.createPCM(prog, progress), prog
}

func (s *Synthesizer) CreateProgram() SynthProgram {
	prog := SynthProgram{
		Length:      20,
		Instruments: make([]SynthProgramInstrument, 0, 4),
	}
	for id, inst := range s.instruments {
		if !inst.enabled || inst.compiledFx == nil {
			continue
		}
		index := len(prog.Instruments)
		prog.Instruments = append(prog.Instruments, SynthProgramInstrument{
			ID:     id,
			Index:  index,
			Func:   inst.compiledFx,
			Period: inst.period,
			Kind:   inst.kind,
		})
	}
	return prog
}

func (s *Synthesizer) SetInstrumentEnabled(id int, enabled bool) {
	s.changed = true
	s.instruments[id].enabled = enabled
}

func (s *Synthesizer) SetInstrumentVolume(id int, volume float64) {
	s.changed = true
	inst := s.instruments[id]
	inst.volume = volume
	inst.mappedVolume = int32(math.Round(127.0 * volume))
}

func (s *Synthesizer) SetInstrumentPatch(id int, index int) {
	s.changed = true
	inst := s.instruments[id]
	inst.patchNumber = int32(s.sf.Instruments[index].PatchNumber)
	inst.instrumentIndex = index
	inst.kind = s.sf.Instruments[index].Kind
}

func (s *Synthesizer) SetInstrumentPeriod(id int, period float64) {
	s.changed = true
	s.instruments[id].period = period
}

func (s *Synthesizer) SetInstrumentFunction(id int, fx string) {
	s.changed = true
	s.recompileDelay = 0.75
	s.instruments[id].SetFx(fx)
}

func (s *Synthesizer) GetInstrumentFunction(id int) func(float64) float64 {
	return s.instruments[id].compiledFx
}

func (s *Synthesizer) needsPlotRedraw() int {
	for i, inst := range s.instruments {
		oldFx := inst.oldFx
		inst.oldFx = inst.fx
		if inst.fx != oldFx {
			return i
		}
	}
	return -1
}
