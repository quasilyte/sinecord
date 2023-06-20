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
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/styles"
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
	plots           []*ebiten.Image
	periods         []*ebiten.Image
	plotsHidden     []bool

	ctx *Context

	Running bool
}

func NewCanvas(ctx *Context, scene *ge.Scene, img *ebiten.Image) *Canvas {
	canvas := &Canvas{
		ctx:             ctx,
		scene:           scene,
		canvasImage:     img,
		scratch:         ebiten.NewImage(img.Size()),
		waves:           ebiten.NewImage(img.Size()),
		scratchVertices: make([]ebiten.Vertex, 6000),
		scratchIndices:  make([]uint16, 0, 8000),
		objects:         make([]ge.SceneGraphics, 0, 32),
		plots:           make([]*ebiten.Image, ctx.config.MaxInstruments),
		periods:         make([]*ebiten.Image, ctx.config.MaxInstruments),
		plotsHidden:     make([]bool, ctx.config.MaxInstruments),
	}
	for i := range canvas.plots {
		canvas.plots[i] = ebiten.NewImage(img.Size())
	}
	for i := range canvas.periods {
		canvas.periods[i] = ebiten.NewImage(img.Size())
	}
	return canvas
}

func (c *Canvas) RedrawPlot(id int, f func(x float64) float64, points []gmath.Vec) {
	img := c.plots[id]
	img.Clear()

	height := float64(img.Bounds().Dy())
	var clr ge.ColorScale
	clr.SetColor(styles.PlotColorByID[id])

	dx := 1.0 / 30.0
	smallDx := 1.0 / 60.0
	tinyDx := 1.0 / 180.0
	x := 0.0
	for x < 20 {
		y := f(x)
		scaled := c.ctx.Scaler.ScaleXY(x, y)
		if scaled.Y >= 0 && scaled.Y <= height {
			var p vector.Path
			for x < 20 {
				y := f(x)
				scaled := c.ctx.Scaler.ScaleXY(x, y)
				if scaled.Y < 0 || scaled.Y > height {
					break
				}
				p.LineTo(float32(scaled.X), float32(scaled.Y))
				if math.Abs(f(x+dx)-y) > 0.2 {
					x += tinyDx
				} else {
					x += dx
				}
			}
			c.DrawPath(img, p, 2, clr)
		}
		x += smallDx
	}

	img = c.periods[id]
	img.Clear()
	periodColor := styles.PlotColorByID[id]
	periodColor.A /= 2
	for _, p := range points {
		scaled := c.ctx.Scaler.ScalePos(p)
		{
			x1 := float32(scaled.X)
			y1 := float32(scaled.Y) - 6
			x2 := x1
			y2 := float32(scaled.Y) + 6
			vector.StrokeLine(img, x1, y1, x2, y2, 2, periodColor, true)
		}
		{
			x1 := float32(scaled.X) - 6
			y1 := float32(scaled.Y)
			x2 := float32(scaled.X) + 6
			y2 := y1
			vector.StrokeLine(img, x1, y1, x2, y2, 2, periodColor, true)
		}
	}
}

func (c *Canvas) ClearPlot(id int) {
	c.plots[id].Clear()
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
		y := ((4.5 * data[sampleIndex] * 3.0) * 46.0) + (46 * 3)
		p.LineTo(float32(x), float32(y))
	}
	var cs ge.ColorScale
	cs.SetColor(styles.SoundWaveColor)
	cs.A = 0.8
	c.DrawPath(c.waves, p, 2, cs)
}

func (c *Canvas) SetPlotHidden(id int, hidden bool) {
	c.plotsHidden[id] = hidden
}

func (c *Canvas) DrawInstrumentIcon(dst *ebiten.Image, kind gamedata.InstrumentKind, clr color.RGBA) {
	shape := gamedata.InstrumentShape(kind)
	width := dst.Bounds().Dx()
	var colorScale ge.ColorScale
	colorScale.SetColor(clr)
	c.drawFilledShape(dst, shape, float32(width/2), float32(width/2), float32(width/2), 0, colorScale)
}

func (c *Canvas) translate(p pos32, x, y float32) (float32, float32) {
	return x + p.x, y + p.y
}

func (c *Canvas) createShapePath(shape gamedata.Shape, x, y, r float32, angle gmath.Rad) vector.Path {
	var p vector.Path
	switch shape {
	case gamedata.ShapeCircle:
		p.Arc(x, y, r, 0, 2*math.Pi, vector.Clockwise)
	case gamedata.ShapeSquare:
		p.MoveTo(c.translate(rotate(-r, -r, angle), x, y))
		p.LineTo(c.translate(rotate(-r, +r, angle), x, y))
		p.LineTo(c.translate(rotate(+r, +r, angle), x, y))
		p.LineTo(c.translate(rotate(+r, -r, angle), x, y))
		p.Close()
	case gamedata.ShapeTriangle:
		p.MoveTo(c.translate(rotate(-r, -r, angle), x, y))
		p.LineTo(c.translate(rotate(0, r, angle), x, y))
		p.LineTo(c.translate(rotate(r, -r, angle), x, y))
		p.Close()
	case gamedata.ShapeHexagon:
		r2 := r / 2
		p.MoveTo(c.translate(rotate(-r2, -r, angle), x, y))
		p.LineTo(c.translate(rotate(-r, 0, angle), x, y))
		p.LineTo(c.translate(rotate(-r2, r, angle), x, y))
		p.LineTo(c.translate(rotate(r2, r, angle), x, y))
		p.LineTo(c.translate(rotate(r, 0, angle), x, y))
		p.LineTo(c.translate(rotate(r2, -r, angle), x, y))
		p.Close()
	case gamedata.ShapeStar:
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
	case gamedata.ShapeCross:
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

func (c *Canvas) drawFilledShape(dst *ebiten.Image, shape gamedata.Shape, x, y, r float32, angle gmath.Rad, clr ge.ColorScale) {
	p := c.createShapePath(shape, x, y, r, angle)
	c.DrawPathFilled(dst, p, 1, clr)
}

func (c *Canvas) drawShape(dst *ebiten.Image, shape gamedata.Shape, x, y, r float32, angle gmath.Rad, clr ge.ColorScale) {
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

	var drawPlotOptions ebiten.DrawImageOptions
	if c.Running {
		drawPlotOptions.ColorM.Scale(1, 1, 1, 0.2)
	}
	for i, p := range c.plots {
		if c.plotsHidden[i] {
			continue
		}
		c.canvasImage.DrawImage(p, &drawPlotOptions)
	}
	for i, p := range c.periods {
		if c.plotsHidden[i] {
			continue
		}
		c.canvasImage.DrawImage(p, &drawPlotOptions)
	}

	// if !c.Running {
	// 	width := bg.Bounds().Dx()
	// 	height := bg.Bounds().Dy()
	// 	for _, shader := range c.fnShaders {
	// 		c.scratch.Clear()
	// 		c.scratch.DrawImage(c.canvasImage, &drawOptions)
	// 		var options ebiten.DrawRectShaderOptions
	// 		if shader == nil {
	// 			continue
	// 		}
	// 		options.Images[0] = c.scratch
	// 		options.CompositeMode = ebiten.CompositeModeCopy
	// 		c.canvasImage.DrawRectShader(width, height, shader, &options)
	// 	}
	// }

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
