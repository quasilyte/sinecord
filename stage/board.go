package stage

import (
	"image/color"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type SynthProgram struct {
	Instruments []SynthProgramInstrument
}

type SynthProgramInstrument struct {
	Color color.RGBA
	Func  func(x float64) float64
}

type Board struct {
	scene  *ge.Scene
	canvas *Canvas

	plotScale  float64
	plotOffset gmath.Vec

	length float64
	t      float64
	prog   SynthProgram

	signals []*signalNode
}

type BoardConfig struct {
	Canvas     *Canvas
	PlotScale  float64
	PlotOffset gmath.Vec
}

func NewBoard(scene *ge.Scene, config BoardConfig) *Board {
	return &Board{
		canvas:     config.Canvas,
		plotOffset: config.PlotOffset,
		plotScale:  config.PlotScale,
		length:     20,
		scene:      scene,
		signals:    make([]*signalNode, 0, 4),
	}
}

func (b *Board) StartProgram(prog SynthProgram) {
	b.reset()
	b.initProgram(prog)
}

func (b *Board) ClearProgram() {
	b.reset()
}

func (b *Board) ProgramTick(delta float64) bool {
	finished := false
	if b.t+delta >= b.length {
		delta = b.length - b.t
		finished = true
	}

	x := b.t
	for i, sig := range b.signals {
		y := b.prog.Instruments[i].Func(x)
		sig.sprite.Visible = y >= -3 && y <= 3
		if sig.sprite.Visible {
			sig.pos.X = (x * b.plotScale)
			sig.pos.Y = -(y * b.plotScale)
			sig.pos = sig.pos.Add(b.plotOffset)
		}
	}
	b.t += delta

	return finished
}

func (b *Board) initProgram(prog SynthProgram) {
	b.prog = prog

	for _, inst := range prog.Instruments {
		sig := newSignalNode(b.canvas, inst.Color)
		b.scene.AddObject(sig)
		b.signals = append(b.signals, sig)
	}
}

func (b *Board) reset() {
	for _, sig := range b.signals {
		sig.Dispose()
	}
	b.signals = b.signals[:0]

	b.t = 0
}
