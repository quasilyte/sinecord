package synthdb

import (
	"fmt"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

type SoundFont struct {
	Name string

	Instruments []*Instrument

	Data *meltysynth.SoundFont
}

type Instrument struct {
	Kind InstrumentKind

	Name string

	Index       int
	PatchNumber int
}

type InstrumentKind int

const (
	BassInstrument InstrumentKind = iota
	DrumInstrument
	KeyboardInstrument
	BrassInstrument
	StringInstrument
	OtherInstrument
)

var TimGM6mb = &SoundFont{
	Name: "TimGM6mb",

	Instruments: []*Instrument{
		{Kind: BassInstrument, Name: "Synth Bass 1", Index: 88},
		{Kind: BassInstrument, Name: "Synth Bass 2", Index: 87},
		{Kind: BassInstrument, Name: "Slap Bass 1", Index: 90},
		{Kind: BassInstrument, Name: "Slap Bass 2", Index: 89},
		{Kind: BassInstrument, Name: "Bass & Lead", Index: 49},
		{Kind: BassInstrument, Name: "Distortion Guitar", Index: 96},
		{Kind: BassInstrument, Name: "5th Saw Wave", Index: 50},

		{Kind: KeyboardInstrument, Name: "Piano 1", Index: 126},
		{Kind: KeyboardInstrument, Name: "Piano 2", Index: 125},
		{Kind: KeyboardInstrument, Name: "Piano 3", Index: 124},
		{Kind: KeyboardInstrument, Name: "Honky Tonk", Index: 123},
		{Kind: KeyboardInstrument, Name: "Clavinet", Index: 119},
		{Kind: KeyboardInstrument, Name: "Harpsichord", Index: 120},

		{Kind: BrassInstrument, Name: "Synth Brass 1", Index: 70},
		{Kind: BrassInstrument, Name: "Synth Brass 2", Index: 69},
		{Kind: BrassInstrument, Name: "Brass", Index: 71},
		{Kind: BrassInstrument, Name: "Bassoon", Index: 65},
		{Kind: BrassInstrument, Name: "French Horns", Index: 68},
		{Kind: BrassInstrument, Name: "Tuba", Index: 73},

		{Kind: StringInstrument, Name: "Charang", Index: 52},
		{Kind: StringInstrument, Name: "Banjo", Index: 31},
		{Kind: StringInstrument, Name: "Koto", Index: 29},
		{Kind: StringInstrument, Name: "Synth Strings", Index: 79},
		{Kind: StringInstrument, Name: "Dulcimer", Index: 111},
		{Kind: StringInstrument, Name: "Guitar Harmonics", Index: 95},

		{Kind: DrumInstrument, Name: "Timpani", Index: 80},
		{Kind: DrumInstrument, Name: "Synth Drum", Index: 18},
		{Kind: DrumInstrument, Name: "Taiko Drum", Index: 20},
		{Kind: DrumInstrument, Name: "Steel Drum", Index: 22},
		{Kind: DrumInstrument, Name: "Tom Drum", Index: 19},

		{Kind: OtherInstrument, Name: "Tinker Bell", Index: 24},
		{Kind: OtherInstrument, Name: "Voice Oohs", Index: 78},
		{Kind: OtherInstrument, Name: "Choir Aahs", Index: 131},
	},
}

func (sf *SoundFont) Load(data *meltysynth.SoundFont) {
	sf.Data = data

	presets := map[int]bool{}

	for _, inst := range sf.Instruments {
		if presets[inst.Index] {
			panic("found duplicated preset")
		}
		presets[inst.Index] = true

		instInfo := data.Presets[inst.Index]
		inst.PatchNumber = int(instInfo.PatchNumber)

		for _, region := range instInfo.Regions {
			if region.GetVelocityRangeStart() != 0 {
				panic("unexpected velocity range start")
			}
			if region.GetVelocityRangeEnd() != 127 {
				panic("unexpected velocity range end")
			}
			if region.GetKeyRangeStart() != 0 {
				panic("unexpected key range start")
			}
			if region.GetKeyRangeEnd() != 127 {
				panic("unexpected key range end")
			}
		}

		fmt.Printf("loaded %q instrument (patch name = %s)\n", inst.Name, instInfo.Name)
	}
}
