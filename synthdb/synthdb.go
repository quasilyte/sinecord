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
	Name string

	Index       int
	PatchNumber int
}

var TimGM6mb = &SoundFont{
	Name: "TimGM6mb",

	Instruments: []*Instrument{
		{Name: "Synth Bass 1", Index: 88},
		{Name: "Synth Bass 2", Index: 87},
		{Name: "Slap Bass 1", Index: 90},
		{Name: "Slap Bass 2", Index: 89},

		{Name: "Synth Brass 1", Index: 70},
		{Name: "Synth Brass 2", Index: 69},
		{Name: "Brass", Index: 71},

		{Name: "Piano 1", Index: 126},
		{Name: "Piano 2", Index: 125},
		{Name: "Piano 3", Index: 124},
		{Name: "Honky Tonk", Index: 123},

		{Name: "Tinker Bell", Index: 24},

		{Name: "Clavinet", Index: 119},
		{Name: "Bass & Lead", Index: 49},
		{Name: "5th Saw Wave", Index: 50},
		{Name: "Charang", Index: 52},
		{Name: "Distortion Guitar", Index: 96},

		{Name: "Bassoon", Index: 65},
		{Name: "French Horns", Index: 68},
		{Name: "Tuba", Index: 73},

		{Name: "Banjo", Index: 31},
		{Name: "Koto", Index: 29},
		{Name: "Synth Strings", Index: 79},

		{Name: "Dulcimer", Index: 111},

		{Name: "Guitar Harmonics", Index: 95},
		{Name: "Harpsichord", Index: 120},

		{Name: "Voice Oohs", Index: 78},
		{Name: "Choir Aahs", Index: 131},

		{Name: "Timpani", Index: 80},
		{Name: "Synth Drum", Index: 18},
		{Name: "Taiko Drum", Index: 20},
		{Name: "Steel Drum", Index: 22},
		{Name: "Tom Drum", Index: 19},
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
