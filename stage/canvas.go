package stage

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/synthdb"
)

type Canvas struct {
	time float64

	scene *ge.Scene

	objects []ge.SceneGraphics

	canvasImage *ebiten.Image
	scratch     *ebiten.Image
	waves       *ebiten.Image

	scratchVertices []ebiten.Vertex
	scratchIndices  []uint16
	fnShaders       []*ebiten.Shader

	ctx *Context

	Running bool
}

func NewCanvas(ctx *Context, scene *ge.Scene, img *ebiten.Image) *Canvas {
	return &Canvas{
		ctx:             ctx,
		scene:           scene,
		canvasImage:     img,
		scratch:         ebiten.NewImage(img.Size()),
		waves:           ebiten.NewImage(img.Size()),
		scratchVertices: make([]ebiten.Vertex, 6000),
		scratchIndices:  make([]uint16, 0, 8000),
		objects:         make([]ge.SceneGraphics, 0, 32),
		fnShaders:       make([]*ebiten.Shader, ctx.config.MaxInstruments),
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
	c.DrawPath(c.waves, p, 2, ge.ColorScale{R: 0.616, G: 0.843, B: 0.576, A: 1})
}

func (c *Canvas) DrawInstrumentIcon(dst *ebiten.Image, kind synthdb.InstrumentKind, clr color.RGBA) {
	shape := instrumentWaveShape(kind)
	width := dst.Bounds().Dx()
	var colorScale ge.ColorScale
	colorScale.SetColor(clr)
	c.drawFilledShape(dst, shape, float32(width/2), float32(width/2), float32(width/2), 0, colorScale)
}

func (c *Canvas) translate(p pos32, x, y float32) (float32, float32) {
	return x + p.x, y + p.y
}

func (c *Canvas) createShapePath(shape waveShape, x, y, r float32, angle gmath.Rad) vector.Path {
	var p vector.Path
	switch shape {
	case waveCircle:
		p.Arc(x, y, r, 0, 2*math.Pi, vector.Clockwise)
	case waveSquare:
		p.MoveTo(c.translate(rotate(-r, -r, angle), x, y))
		p.LineTo(c.translate(rotate(-r, +r, angle), x, y))
		p.LineTo(c.translate(rotate(+r, +r, angle), x, y))
		p.LineTo(c.translate(rotate(+r, -r, angle), x, y))
		p.Close()
	case waveTriangle:
		p.MoveTo(c.translate(rotate(-r, -r, angle), x, y))
		p.LineTo(c.translate(rotate(0, r, angle), x, y))
		p.LineTo(c.translate(rotate(r, -r, angle), x, y))
		p.Close()
	case waveOctagon:
		r2 := r / 2
		p.MoveTo(c.translate(rotate(-r, -r2, angle), x, y))
		p.LineTo(c.translate(rotate(-r, r2, angle), x, y))
		p.LineTo(c.translate(rotate(-r2, r, angle), x, y))
		p.LineTo(c.translate(rotate(r2, r, angle), x, y))
		p.LineTo(c.translate(rotate(r, r2, angle), x, y))
		p.LineTo(c.translate(rotate(r, -r2, angle), x, y))
		p.LineTo(c.translate(rotate(r2, -r, angle), x, y))
		p.LineTo(c.translate(rotate(-r2, -r, angle), x, y))
		p.Close()
	case waveStar:
		r3 := r / 3
		p.MoveTo(c.translate(rotate(-r3, -r3, angle), x, y))
		p.LineTo(c.translate(rotate(-r, 0, angle), x, y))
		p.LineTo(c.translate(rotate(-r3, r3, angle), x, y))
		p.LineTo(c.translate(rotate(0, r, angle), x, y))
		p.LineTo(c.translate(rotate(r3, r3, angle), x, y))
		p.LineTo(c.translate(rotate(r, 0, angle), x, y))
		p.LineTo(c.translate(rotate(r3, -r3, angle), x, y))
		p.LineTo(c.translate(rotate(0, -r, angle), x, y))
		p.Close()
	case waveCross:
		r2 := r / 2
		p.MoveTo(c.translate(rotate(-r2, -r2, angle), x, y))
		p.LineTo(c.translate(rotate(-r, -r2, angle), x, y))
		p.LineTo(c.translate(rotate(-r, r2, angle), x, y))
		p.LineTo(c.translate(rotate(-r2, r2, angle), x, y))
		p.LineTo(c.translate(rotate(-r2, r, angle), x, y))
		p.LineTo(c.translate(rotate(r2, r, angle), x, y))
		p.LineTo(c.translate(rotate(r2, r2, angle), x, y))
		p.LineTo(c.translate(rotate(r, r2, angle), x, y))
		p.LineTo(c.translate(rotate(r, -r2, angle), x, y))
		p.LineTo(c.translate(rotate(r2, -r2, angle), x, y))
		p.LineTo(c.translate(rotate(r2, -r, angle), x, y))
		p.LineTo(c.translate(rotate(-r2, -r, angle), x, y))
		p.Close()
	}
	return p
}

func (c *Canvas) drawFilledShape(dst *ebiten.Image, shape waveShape, x, y, r float32, angle gmath.Rad, clr ge.ColorScale) {
	p := c.createShapePath(shape, x, y, r, angle)
	c.DrawPathFilled(dst, p, 1, clr)
}

func (c *Canvas) drawShape(dst *ebiten.Image, shape waveShape, x, y, r float32, angle gmath.Rad, clr ge.ColorScale) {
	p := c.createShapePath(shape, x, y, r, angle)
	c.DrawPath(dst, p, 1, clr)
}

func (c *Canvas) DrawPath(dst *ebiten.Image, p vector.Path, width float32, clr ge.ColorScale) {
	var strokeOptions vector.StrokeOptions
	strokeOptions.Width = width
	c.scratchVertices, c.scratchIndices = p.AppendVerticesAndIndicesForStroke(c.scratchVertices[:0], c.scratchIndices[:0], &strokeOptions)
	vs := c.scratchVertices
	is := c.scratchIndices
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = clr.R
		vs[i].ColorG = clr.G
		vs[i].ColorB = clr.B
		vs[i].ColorA = clr.A
	}
	op := ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	}
	dst.DrawTriangles(vs, is, whiteSubImage, &op)
}

func (c *Canvas) DrawPathFilled(dst *ebiten.Image, p vector.Path, width float32, clr ge.ColorScale) {
	c.scratchVertices, c.scratchIndices = p.AppendVerticesAndIndicesForFilling(c.scratchVertices[:0], c.scratchIndices[:0])
	vs := c.scratchVertices
	is := c.scratchIndices
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = clr.R
		vs[i].ColorG = clr.G
		vs[i].ColorB = clr.B
		vs[i].ColorA = clr.A
	}
	op := ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	}
	dst.DrawTriangles(vs, is, whiteSubImage, &op)
}

func (c *Canvas) AddGraphics(o ge.SceneGraphics) {
	c.objects = append(c.objects, o)
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

	c.drawObjects()
}

func (c *Canvas) drawObjects() {
	liveObjects := c.objects[:0]
	for _, o := range c.objects {
		if o.IsDisposed() {
			continue
		}
		o.Draw(c.canvasImage)
		liveObjects = append(liveObjects, o)
	}
	c.objects = liveObjects
}

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}
