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
	instruments []*instrument
	delays      []float64
	settings    *meltysynth.SynthesizerSettings
	length      int32
	left        []float32
	right       []float32
}

func newMusicPlayer(instruments []*instrument) *musicPlayer {
	p := &musicPlayer{
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

func (p *musicPlayer) createPCM() []byte {
	synthesizer, err := meltysynth.NewSynthesizer(assets.SoundFontTimGM6mb, p.settings)
	if err != nil {
		panic(err)
	}

	synthesizer.MasterVolume = 0.75

	const stepsPerSecond = 180
	dt := 1.0 / stepsPerSecond
	t := 0.0
	duration := 20.0
	blockOffset := 0

	for i := range p.delays {
		p.delays[i] = p.instruments[i].period
	}

	samplesPerDt := int(p.settings.SampleRate) / int(stepsPerSecond)

	for t < duration {
		t += dt

		for i, inst := range p.instruments {
			channel := int32(i)
			if !inst.enabled || inst.compiledFx == nil {
				continue
			}
			delay := p.delays[i]
			if delay > dt {
				p.delays[i] -= dt
				continue
			}
			// synthesizer.ProcessMidiMessage(0, 0xB0, 0x27, 0x3FFF)
			y := math.Abs(inst.compiledFx(t))
			// 48 notes
			// From 60 to 107
			// value range is 0 to 3
			synthesizer.NoteOffAllChannel(channel, false)
			if y > 3 || y < -3 {
				continue
			}
			note := int32(math.Round(y*float64(synthdb.Ocvate4EndCode-synthdb.Octave1StartCode+1)/3)) + synthdb.Octave1StartCode

			// volume := int32(0)
			// if y <= 3 && y >= -3 {
			// 	volume = int32(math.Round(y * (127.0 / 3)))
			// }
			// synthesizer.ProcessMidiMessage(0, 0xB0, 0x07, volume)
			synthesizer.ProcessMidiMessage(channel, 0xC0, int32(inst.instrumentIndex), 0) // guitar
			p.delays[i] = inst.period - (dt - delay)
			velocity := int32(40)
			synthesizer.NoteOn(channel, note, int32(velocity))
		}

		synthesizer.Render(p.left[blockOffset:blockOffset+samplesPerDt], p.right[blockOffset:blockOffset+samplesPerDt])
		blockOffset += samplesPerDt
	}

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
