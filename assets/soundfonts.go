package assets

import (
	"bytes"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/sinshu/go-meltysynth/meltysynth"
)

var (
	SoundFontTimGM6mb *meltysynth.SoundFont
)

func decodeSoundfonts(ctx *ge.Context) {
	decode := func(id resource.RawID) *meltysynth.SoundFont {
		data := ctx.Loader.LoadRaw(id).Data
		soundFont, err := meltysynth.NewSoundFont(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		return soundFont
	}

	SoundFontTimGM6mb = decode(RawSoundFontTimGM6mb)
}
