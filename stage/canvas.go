package stage

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/assets"
)

type Canvas struct {
	time float64

	scene *ge.Scene

	sprites []*ge.Sprite

	canvasImage *ebiten.Image
	scratch     *ebiten.Image

	fnShaders []*ebiten.Shader

	Running bool
}

func NewCanvas(ctx *Context, scene *ge.Scene, img *ebiten.Image) *Canvas {
	return &Canvas{
		scene:       scene,
		canvasImage: img,
		scratch:     ebiten.NewImage(img.Size()),
		sprites:     make([]*ge.Sprite, 0, 16),
		fnShaders:   make([]*ebiten.Shader, ctx.config.MaxInstruments),
	}
}

func (c *Canvas) SetShader(id int, shader *ebiten.Shader) {
	c.fnShaders[id] = shader
}

func (c *Canvas) AddSprite(s *ge.Sprite) {
	c.sprites = append(c.sprites, s)
}

func (c *Canvas) Update(delta float64) {
	c.time += delta
}

func (c *Canvas) IsDisposed() bool { return false }

func (c *Canvas) Draw() {
	c.canvasImage.Clear()

	plotBackground := c.scene.LoadImage(assets.ImagePlotBackground).Data

	var drawOptions ebiten.DrawImageOptions
	c.canvasImage.DrawImage(plotBackground, &drawOptions)

	if !c.Running {
		width := plotBackground.Bounds().Dx()
		height := plotBackground.Bounds().Dy()
		for _, shader := range c.fnShaders {
			c.scratch.Clear()
			c.scratch.DrawImage(c.canvasImage, &drawOptions)

			var options ebiten.DrawRectShaderOptions
			if shader == nil {
				continue
			}
			options.Images[0] = c.scratch
			options.CompositeMode = ebiten.CompositeModeCopy
			c.canvasImage.DrawRectShader(width, height, shader, &options)
		}
	}

	c.drawSprites()
}

func (c *Canvas) drawSprites() {
	liveSprites := c.sprites[:0]
	for _, s := range c.sprites {
		if s.IsDisposed() {
			continue
		}
		s.Draw(c.canvasImage)
		liveSprites = append(liveSprites, s)
	}
	c.sprites = liveSprites
}
