package stage

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/assets"
)

type Canvas struct {
	time float64

	scene *ge.Scene

	sprites []*ge.Sprite

	canvasImage *ebiten.Image
	scratch     *ebiten.Image
	waves       *ebiten.Image

	fnShaders []*ebiten.Shader

	Running bool
}

func NewCanvas(ctx *Context, scene *ge.Scene, img *ebiten.Image) *Canvas {
	return &Canvas{
		scene:       scene,
		canvasImage: img,
		scratch:     ebiten.NewImage(img.Size()),
		waves:       ebiten.NewImage(img.Size()),
		sprites:     make([]*ge.Sprite, 0, 16),
		fnShaders:   make([]*ebiten.Shader, ctx.config.MaxInstruments),
	}
}

func (c *Canvas) SetShader(id int, shader *ebiten.Shader) {
	c.fnShaders[id] = shader
}

func (c *Canvas) Reset() {
	c.waves.Clear()
}

func (c *Canvas) RenderWave(data []float64) {
	c.waves.Clear()
	if data == nil {
		return
	}

	offsetX := 4
	width := c.canvasImage.Bounds().Dx()
	widthAvailable := width - offsetX
	samplesPerPixel := len(data) / widthAvailable
	var p vector.Path
	for x := offsetX; x < width; x++ {
		sampleIndex := x * samplesPerPixel
		if sampleIndex > len(data) {
			break
		}
		y := ((5 * data[sampleIndex] * 3.0) * 46.0) + (46 * 3)
		p.LineTo(float32(x), float32(y))
	}
	var strokeOptions vector.StrokeOptions
	strokeOptions.Width = 2
	vs, is := p.AppendVerticesAndIndicesForStroke(nil, nil, &strokeOptions)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 0.616
		vs[i].ColorG = 0.843
		vs[i].ColorB = 0.576
		vs[i].ColorA = 1
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true}
	c.waves.DrawTriangles(vs, is, whiteSubImage, op)
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

	var bg *ebiten.Image
	if c.Running {
		bg = c.scene.LoadImage(assets.ImagePlayBackground).Data
	} else {
		bg = c.scene.LoadImage(assets.ImagePlotBackground).Data
	}

	var drawOptions ebiten.DrawImageOptions
	c.canvasImage.DrawImage(bg, &drawOptions)

	if !c.Running {
		width := bg.Bounds().Dx()
		height := bg.Bounds().Dy()
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

	if c.Running {
		c.canvasImage.DrawImage(c.waves, &drawOptions)
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

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}
