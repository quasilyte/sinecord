package stage

import (
	"math"

	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/synthdb"
	"github.com/sinshu/go-meltysynth/meltysynth"
)

type SampleSet struct {
	PerSecond int
	Left      []float32
	Right     []float32
}

type musicPlayer struct {
	ctx         *Context
	instruments []*instrument
	delays      []float64
	settings    *meltysynth.SynthesizerSettings
	length      int32
	left        []float32
	right       []float32
}

func newMusicPlayer(ctx *Context, instruments []*instrument) *musicPlayer {
	p := &musicPlayer{
		ctx:         ctx,
		instruments: instruments,
		delays:      make([]float64, len(instruments)),
	}
	p.settings = meltysynth.NewSynthesizerSettings(44100)
	p.settings.EnableReverbAndChorus = false
	p.length = 20 * p.settings.SampleRate
	p.left = make([]float32, p.length)
	p.right = make([]float32, p.length)
	return p
}

func (p *musicPlayer) walkNotes(events []noteActivation, f func(i, num int)) {
	i := 0
	for i < len(events) {
		num := 1
		t := events[i].t
		for j := i + 1; j < len(events); j++ {
			if events[j].t != t {
				break
			}
			num++
		}
		f(i, num)
		i += num
	}
}

func (p *musicPlayer) createPCM(prog SynthProgram, progress *float64) *SampleSet {
	synthesizer, err := meltysynth.NewSynthesizer(assets.SoundFontTimGM6mb, p.settings)
	if err != nil {
		panic(err)
	}

	for channel, inst := range p.instruments {
		synthesizer.ProcessMidiMessage(int32(channel), 0xC0, inst.patchNumber, 0)
		synthesizer.ProcessMidiMessage(int32(channel), 0xB0, 0x07, inst.mappedVolume)
	}

	synthesizer.MasterVolume = 0.75

	samplesPerSecond := float64(p.settings.SampleRate)
	t := 0.0
	blockOffset := 0

	events := p.ctx.runner.RunProgram(prog)
	processedEvents := 0
	p.walkNotes(events, func(i, num int) {
		if processedEvents != 0 {
			*progress = float64(processedEvents) / float64(len(events))
		}
		processedEvents += num

		eventTime := events[i].t

		elapsed := eventTime - t
		t = eventTime
		blockSize := int(samplesPerSecond * elapsed)
		synthesizer.Render(p.left[blockOffset:blockOffset+blockSize], p.right[blockOffset:blockOffset+blockSize])
		blockOffset += blockSize

		for j := 0; j < num; j++ {
			e := events[i+j]
			inst := p.instruments[e.id]
			channel := int32(e.id)
			synthesizer.NoteOffAllChannel(channel, false)
			y := math.Abs(inst.compiledFx(e.t))
			if y > 3 || y < -3 {
				continue
			}
			note := int32(math.Round(y*float64(synthdb.Ocvate4EndCode-synthdb.Octave1StartCode+1)/3)) + synthdb.Octave1StartCode
			velocity := int32(40)
			synthesizer.NoteOn(channel, note, int32(velocity))
		}
	})

	synthesizer.Render(p.left[blockOffset:], p.right[blockOffset:])

	return &SampleSet{
		PerSecond: int(p.settings.SampleRate),
		Left:      p.left,
		Right:     p.right,
	}
}
