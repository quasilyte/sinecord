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

	targets        []*targetNode
	effects        []*waveNode
	pendingEffects []*waveNode

	signals []*signalNode

	config BoardConfig

	EventNote gsignal.Event[int]
}

type BoardConfig struct {
	Canvas *Canvas

	MaxInstruments int

	Targets []Target
}

func NewBoard(config BoardConfig) *Board {
	return &Board{
		config:  config,
		canvas:  config.Canvas,
		length:  20,
		signals: make([]*signalNode, 0, config.MaxInstruments),
	}
}

func (b *Board) Init(scene *ge.Scene) {
	b.scene = scene
	b.deployTargets()
}

func (b *Board) StartProgram(prog SynthProgram) {
	b.reset()
	b.initProgram(prog)
}

func (b *Board) ClearProgram() {
	b.reset()
	b.deployTargets()
}

func (b *Board) ProgramTick(delta float64) bool {
	if b.finished {
		panic("running a finished program")
	}

	if b.t+delta >= b.length {
		delta = b.length - b.t
		b.finished = true
	}

	{
		liveTargets := b.targets[:0]
		for _, target := range b.targets {
			if target.IsDisposed() {
				continue
			}
			liveTargets = append(liveTargets, target)
			target.Update(delta)
		}
		b.targets = liveTargets
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
		pos := b.canvas.scaleXY(e.t, y)
		inst := b.prog.Instruments[e.index]
		shape := instrumentWaveShape(inst.Kind)
		effect := newWaveNode(b.canvas, shape, pos, styles.PlotColorByID[e.id], inst.Period*0.95)
		b.addWaveEffect(effect)
		b.EventNote.Emit(e.id)

		effect.EventFinished.Connect(nil, func(r float64) {
			for _, t := range b.targets {
				if t.instrument != inst.Kind {
					continue
				}
				if t.pos.DistanceTo(pos) < 0.5*(float64(t.r)+r) {
					t.Dispose()
					d := 2 * (float64(t.r) / b.canvas.ctx.PlotScale)
					shape := instrumentWaveShape(t.instrument)
					offset := gmath.Vec{X: 2, Y: 2}
					b.addWaveEffect(newWaveNode(b.canvas, shape, t.pos.Sub(offset), styles.TargetColor, d))
					b.addWaveEffect(newWaveNode(b.canvas, shape, t.pos, styles.TargetColor, d))
					b.addWaveEffect(newWaveNode(b.canvas, shape, t.pos.Add(offset), styles.TargetColor, d))
				}
			}
		})
	}

	x := b.t
	for i, sig := range b.signals {
		y := b.prog.Instruments[i].Func(x)
		sig.sprite.Visible = y >= -3 && y <= 3
		if sig.sprite.Visible {
			sig.pos = b.canvas.scaleXY(x, y)
		}
	}
	b.t += delta

	if len(b.pendingEffects) != 0 {
		b.effects = append(b.effects, b.pendingEffects...)
		b.pendingEffects = b.pendingEffects[:0]
	}

	return b.finished
}

func (b *Board) addWaveEffect(effect *waveNode) {
	b.pendingEffects = append(b.pendingEffects, effect)
	b.canvas.AddGraphics(effect)
}

func (b *Board) initProgram(prog SynthProgram) {
	b.prog = prog
	b.events = b.runner.RunProgram(prog)

	for _, inst := range prog.Instruments {
		sig := newSignalNode(b.canvas, styles.PlotColorByID[inst.ID])
		b.scene.AddObject(sig)
		b.signals = append(b.signals, sig)
	}

	b.deployTargets()
}

func (b *Board) deployTargets() {
	for _, t := range b.config.Targets {
		n := newTargetNode(b, t)
		b.canvas.AddGraphics(n)
		b.targets = append(b.targets, n)
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

	for _, target := range b.targets {
		target.Dispose()
	}
	b.targets = b.targets[:0]

	b.finished = false
	b.t = 0
}
