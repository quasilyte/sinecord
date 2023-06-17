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
		{Name: "Synth Drum", Index: 18},
		{Name: "Taiko Drum", Index: 20},
		{Name: "Steel Drum", Index: 22},
		{Name: "Tom Drum", Index: 19},
		{Name: "Distortion Guitar", Index: 96},
	},
}

func (sf *SoundFont) Load(data *meltysynth.SoundFont) {
	sf.Data = data

	for _, inst := range sf.Instruments {
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
