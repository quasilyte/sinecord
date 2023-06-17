package stage

import (
	"encoding/binary"
	"math"

	"bytes"

	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/synthdb"
	"github.com/sinshu/go-meltysynth/meltysynth"
)

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

func (p *musicPlayer) createPCM(prog SynthProgram) []byte {
	synthesizer, err := meltysynth.NewSynthesizer(assets.SoundFontTimGM6mb, p.settings)
	if err != nil {
		panic(err)
	}

	for channel, inst := range p.instruments {
		synthesizer.ProcessMidiMessage(int32(channel), 0xC0, int32(inst.instrumentIndex), 0)
	}

	synthesizer.MasterVolume = 0.75

	samplesPerSecond := float64(p.settings.SampleRate)
	t := 0.0
	blockOffset := 0

	events := p.ctx.runner.RunProgram(prog)
	p.walkNotes(events, func(i, num int) {
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

	return generatePCM(p.left, p.right)
}

func generatePCM(left, right []float32) []byte {
	length := len(left)

	a := float32(32768.0)

	data := make([]int16, 2*length)

	for i := 0; i < length; i++ {
		data[2*i] = int16(a * left[i])
		data[2*i+1] = int16(a * right[i])
	}

	var buf bytes.Buffer
	buf.Grow(len(data) * 2)

	binary.Write(&buf, binary.LittleEndian, data)
	return buf.Bytes()

	// length := len(left)
	// var max float64

	// for i := 0; i < length; i++ {
	// 	absLeft := math.Abs(float64(left[i]))
	// 	absRight := math.Abs(float64(right[i]))
	// 	if max < absLeft {
	// 		max = absLeft
	// 	}
	// 	if max < absRight {
	// 		max = absRight
	// 	}
	// }

	// a := 32768 * float32(0.99/max)

	// data := make([]int16, 2*length)

	// for i := 0; i < length; i++ {
	// 	data[2*i] = int16(a * left[i])
	// 	data[2*i+1] = int16(a * right[i])
	// }

	// var buf bytes.Buffer
	// buf.Grow(len(data) * 2)

	// binary.Write(&buf, binary.LittleEndian, data)
	// return buf.Bytes()
}
