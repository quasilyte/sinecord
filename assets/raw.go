package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerRawResources(ctx *ge.Context) {
	rawResources := map[resource.RawID]resource.RawInfo{
		RawSoundFontTimGM6mb: {Path: "raw/TimGM6mb.sf2"},
		RawEnDictBase:        {Path: "raw/lang/en.txt"},
		RawLevelTilesetJSON:  {Path: "raw/leveledit.tsj"},
		RawAct1Mission1JSON:  {Path: "raw/level/act1_level1.tmj"},
	}

	for id, res := range rawResources {
		ctx.Loader.RawRegistry.Set(id, res)
		ctx.Loader.LoadRaw(id)
	}
}

const (
	RawNone resource.RawID = iota

	RawEnDictBase

	RawSoundFontTimGM6mb

	RawLevelTilesetJSON
	RawAct1Mission1JSON
)
