package synthdb

import (
	"fmt"

	"github.com/quasilyte/sinecord/gamedata"
	"github.com/sinshu/go-meltysynth/meltysynth"
)

type SoundFont struct {
	Name string

	Instruments []*Instrument

	Data *meltysynth.SoundFont
}

type Instrument struct {
	Kind gamedata.InstrumentKind

	Name string

	Index       int
	PatchNumber int
}

var TimGM6mb = &SoundFont{
	Name: "TimGM6mb",

	Instruments: []*Instrument{
		{Kind: gamedata.BassInstrument, Name: "Synth Bass 1", Index: 88},
		{Kind: gamedata.BassInstrument, Name: "Synth Bass 2", Index: 87},
		{Kind: gamedata.BassInstrument, Name: "Slap Bass 1", Index: 90},
		{Kind: gamedata.BassInstrument, Name: "Slap Bass 2", Index: 89},
		{Kind: gamedata.BassInstrument, Name: "Bass & Lead", Index: 49},
		{Kind: gamedata.BassInstrument, Name: "Distortion Guitar", Index: 96},
		{Kind: gamedata.BassInstrument, Name: "5th Saw Wave", Index: 50},
		{Kind: gamedata.BassInstrument, Name: "Saw Wave", Index: 55},
		{Kind: gamedata.BassInstrument, Name: "Square Wave", Index: 56},

		{Kind: gamedata.KeyboardInstrument, Name: "Piano 1", Index: 126},
		{Kind: gamedata.KeyboardInstrument, Name: "Piano 2", Index: 125},
		{Kind: gamedata.KeyboardInstrument, Name: "Piano 3", Index: 124},
		{Kind: gamedata.KeyboardInstrument, Name: "Piano 4", Index: 8},
		{Kind: gamedata.KeyboardInstrument, Name: "Honky Tonk", Index: 123},
		{Kind: gamedata.KeyboardInstrument, Name: "Clavinet", Index: 119},
		{Kind: gamedata.KeyboardInstrument, Name: "Harpsichord", Index: 120},
		{Kind: gamedata.KeyboardInstrument, Name: "Polysynth", Index: 46},

		{Kind: gamedata.BrassInstrument, Name: "Synth Brass 1", Index: 70},
		{Kind: gamedata.BrassInstrument, Name: "Synth Brass 2", Index: 69},
		{Kind: gamedata.BrassInstrument, Name: "Brass", Index: 71},
		{Kind: gamedata.BrassInstrument, Name: "Bassoon", Index: 65},
		{Kind: gamedata.BrassInstrument, Name: "French Horns", Index: 68},
		{Kind: gamedata.BrassInstrument, Name: "Tuba", Index: 73},
		{Kind: gamedata.BrassInstrument, Name: "Bagpipe", Index: 27},
		{Kind: gamedata.BrassInstrument, Name: "Mute Trumpet", Index: 72},

		{Kind: gamedata.StringInstrument, Name: "Charang", Index: 52},
		{Kind: gamedata.StringInstrument, Name: "Banjo", Index: 31},
		{Kind: gamedata.StringInstrument, Name: "Koto", Index: 29},
		{Kind: gamedata.StringInstrument, Name: "Synth Strings", Index: 79},
		{Kind: gamedata.StringInstrument, Name: "Dulcimer", Index: 111},
		{Kind: gamedata.StringInstrument, Name: "Guitar Harmonics", Index: 95},
		{Kind: gamedata.StringInstrument, Name: "Electronic", Index: 5},
		{Kind: gamedata.StringInstrument, Name: "Clean Guitar", Index: 99},
		{Kind: gamedata.StringInstrument, Name: "Nylon Guitar", Index: 102},

		{Kind: gamedata.DrumInstrument, Name: "Timpani", Index: 80},
		{Kind: gamedata.DrumInstrument, Name: "Synth Drum", Index: 18},
		{Kind: gamedata.DrumInstrument, Name: "Taiko Drum", Index: 20},
		{Kind: gamedata.DrumInstrument, Name: "Steel Drum", Index: 22},
		{Kind: gamedata.DrumInstrument, Name: "Tom Drum", Index: 19},

		{Kind: gamedata.OtherInstrument, Name: "Tinker Bell", Index: 24},
		{Kind: gamedata.OtherInstrument, Name: "Voice Oohs", Index: 78},
		{Kind: gamedata.OtherInstrument, Name: "Choir Aahs", Index: 131},
		{Kind: gamedata.OtherInstrument, Name: "Soundtrack", Index: 39},
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
