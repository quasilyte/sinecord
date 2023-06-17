package stage

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/styles"
)

type Board struct {
	scene  *ge.Scene
	canvas *Canvas

	plotScale  float64
	plotOffset gmath.Vec

	finished bool
	length   float64
	t        float64
	prog     SynthProgram
	events   []noteActivation
	runner   programRunner

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
	if b.finished {
		panic("running a finished program")
	}

	if b.t+delta >= b.length {
		delta = b.length - b.t
		b.finished = true
	}

	for len(b.events) != 0 {
		e := b.events[0]
		if e.t > b.t {
			break
		}
		b.events = b.events[1:]
		y := b.prog.Instruments[e.index].Func(e.t)
		pos := gmath.Vec{
			X: e.t * b.plotScale,
			Y: -(y * b.plotScale),
		}
		pos = pos.Add(b.plotOffset)
		effect := newEffectNode(b.canvas, pos, assets.ImageCircleExplosion, styles.PlotColorByID[e.id])
		b.scene.AddObject(effect)
		effect.anim.SetAnimationSpan(b.prog.Instruments[e.index].Period)
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

	return b.finished
}

func (b *Board) initProgram(prog SynthProgram) {
	b.prog = prog
	b.events = b.runner.RunProgram(prog)

	for _, inst := range prog.Instruments {
		sig := newSignalNode(b.canvas, styles.PlotColorByID[inst.ID])
		b.scene.AddObject(sig)
		b.signals = append(b.signals, sig)
	}
}

func (b *Board) reset() {
	for _, sig := range b.signals {
		sig.Dispose()
	}
	b.signals = b.signals[:0]

	b.finished = false
	b.t = 0
}
