package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerImageResources(ctx *ge.Context) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageUIButtonIdle:     {Path: "image/ebitenui/button-idle.png"},
		ImageUIButtonHover:    {Path: "image/ebitenui/button-hover.png"},
		ImageUIButtonPressed:  {Path: "image/ebitenui/button-pressed.png"},
		ImageUIButtonDisabled: {Path: "image/ebitenui/button-disabled.png"},
	}

	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
		ctx.Loader.LoadImage(id)
	}
}

const (
	ImageNone resource.ImageID = iota

	ImageUIButtonIdle
	ImageUIButtonHover
	ImageUIButtonPressed
	ImageUIButtonDisabled
)
