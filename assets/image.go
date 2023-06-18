package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerImageResources(ctx *ge.Context) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageUIButtonDisabled:      {Path: "image/ebitenui/button-disabled.png"},
		ImageUIButtonIdle:          {Path: "image/ebitenui/button-idle.png"},
		ImageUIButtonHover:         {Path: "image/ebitenui/button-hover.png"},
		ImageUIButtonPressed:       {Path: "image/ebitenui/button-pressed.png"},
		ImageUISelectButtonIdle:    {Path: "image/ebitenui/select-button-idle.png"},
		ImageUISelectButtonHover:   {Path: "image/ebitenui/select-button-hover.png"},
		ImageUISelectButtonPressed: {Path: "image/ebitenui/select-button-pressed.png"},
		ImageUITextInputIdle:       {Path: "image/ebitenui/text-input-idle.png"},
		ImageUIPanelIdle:           {Path: "image/ebitenui/panel-idle.png"},
		ImageUITooltip:             {Path: "image/ebitenui/tooltip.png"},

		ImagePlotBackground: {Path: "image/plot_background.png"},
		ImagePlayBackground: {Path: "image/play_background.png"},
		ImageSignal:         {Path: "image/signal.png"},
	}

	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
		ctx.Loader.LoadImage(id)
	}
}

const (
	ImageNone resource.ImageID = iota

	ImageUIButtonDisabled
	ImageUIButtonIdle
	ImageUIButtonHover
	ImageUIButtonPressed
	ImageUISelectButtonIdle
	ImageUISelectButtonHover
	ImageUISelectButtonPressed
	ImageUITextInputIdle
	ImageUIPanelIdle
	ImageUITooltip

	ImagePlotBackground
	ImagePlayBackground
	ImageSignal
)
