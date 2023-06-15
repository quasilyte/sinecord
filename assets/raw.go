package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerRawResources(ctx *ge.Context) {
	rawResources := map[resource.RawID]resource.RawInfo{
		RawSoundfontBassSynth: {Path: "raw/bass_synth_classic.sf2"},
		RawEnDictBase:         {Path: "raw/lang/en.txt"},
	}

	for id, res := range rawResources {
		ctx.Loader.RawRegistry.Set(id, res)
		ctx.Loader.LoadRaw(id)
	}
}

const (
	RawNone resource.RawID = iota

	RawEnDictBase

	RawSoundfontBassSynth
)
