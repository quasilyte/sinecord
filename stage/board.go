package stage

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/sinecord/styles"
)

type Board struct {
	scene  *ge.Scene
	canvas *Canvas

	finished bool
	length   float64
	t        float64
	prog     SynthProgram
	events   []noteActivation
	runner   programRunner

	effects []*circleWaveNode

	signals []*signalNode

	EventNote gsignal.Event[int]
}

type BoardConfig struct {
	Canvas *Canvas
}

func NewBoard(scene *ge.Scene, config BoardConfig) *Board {
	return &Board{
		canvas:  config.Canvas,
		length:  20,
		scene:   scene,
		signals: make([]*signalNode, 0, 4),
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

	plotScale := b.canvas.ctx.PlotScale
	plotOffset := b.canvas.ctx.PlotOffset

	if b.t+delta >= b.length {
		delta = b.length - b.t
		b.finished = true
	}

	{
		liveEffects := b.effects[:0]
		for _, effect := range b.effects {
			if effect.IsDisposed() {
				continue
			}
			liveEffects = append(liveEffects, effect)
			effect.Update(delta)
		}
		b.effects = liveEffects
	}

	for len(b.events) != 0 {
		e := b.events[0]
		if e.t > b.t {
			break
		}
		b.events = b.events[1:]
		y := b.prog.Instruments[e.index].Func(e.t)
		pos := gmath.Vec{
			X: e.t * plotScale,
			Y: -(y * plotScale),
		}
		pos = pos.Add(plotOffset)
		effect := newCircleWaveNode(b.canvas, pos, styles.PlotColorByID[e.id], b.prog.Instruments[e.index].Period)
		b.effects = append(b.effects, effect)
		b.canvas.AddGraphics(effect)
		b.EventNote.Emit(e.id)
	}

	x := b.t
	for i, sig := range b.signals {
		y := b.prog.Instruments[i].Func(x)
		sig.sprite.Visible = y >= -3 && y <= 3
		if sig.sprite.Visible {
			sig.pos.X = (x * plotScale)
			sig.pos.Y = -(y * plotScale)
			sig.pos = sig.pos.Add(plotOffset)
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

	for _, effect := range b.effects {
		effect.Dispose()
	}
	b.effects = b.effects[:0]

	b.finished = false
	b.t = 0
}
