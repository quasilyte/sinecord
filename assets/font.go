package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerFontResources(ctx *ge.Context) {
	fontResources := map[resource.FontID]resource.FontInfo{
		FontArcadeSmall:  {Path: "font/arcade.otf", Size: 14},
		FontArcadeNormal: {Path: "font/arcade.otf", Size: 20},
		FontArcadeBig:    {Path: "font/arcade.otf", Size: 30},
	}

	for id, res := range fontResources {
		ctx.Loader.FontRegistry.Set(id, res)
		ctx.Loader.LoadFont(id)
	}
}

const (
	FontNone resource.FontID = iota

	FontArcadeSmall
	FontArcadeNormal
	FontArcadeBig
)
