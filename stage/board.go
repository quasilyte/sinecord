package stage

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/styles"
)

type Board struct {
	scene  *ge.Scene
	canvas *Canvas

	ctx *Context

	finished bool
	victory  bool
	penalty  bool
	length   float64
	t        float64
	prog     SynthProgram
	events   []noteActivation
	runner   programRunner

	optionalHits   int
	targetsLeft    int
	targets        []*targetNode
	effects        []*waveNode
	pendingEffects []*waveNode

	signals []*signalNode

	config BoardConfig

	EventNote    gsignal.Event[int]
	EventVictory gsignal.Event[bool]
}

type BoardConfig struct {
	Canvas *Canvas

	Level *gamedata.LevelData

	MaxInstruments int

	Targets []gamedata.Target
}

func NewBoard(ctx *Context, config BoardConfig) *Board {
	return &Board{
		ctx:     ctx,
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

func (b *Board) onVictory() {
	b.victory = true
	b.EventVictory.Emit(b.isBonusAchieved())
}

func (b *Board) isBonusAchieved() bool {
	objectives := b.config.Level.Bonus

	if objectives.AvoidOptional && b.optionalHits > 0 {
		return false
	}

	if objectives.AllTargets && len(b.targets) != 0 {
		return false
	}

	numInstrumentsUsed := len(b.prog.Instruments)
	if numInstrumentsUsed > objectives.MaxInstruments {
		return false
	}

	for _, fn := range objectives.ForbiddenFuncs {
		for _, inst := range b.prog.Instruments {
			if inst.Func.UsesFunc(fn) {
				return false
			}
		}
	}

	return true
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
	if !b.victory && b.targetsLeft == 0 && b.config.Level != nil {
		b.onVictory()
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
		y := b.prog.Instruments[e.index].Func.Run(e.t)
		pos := b.ctx.Scaler.ScaleXY(e.t, y)
		inst := b.prog.Instruments[e.index]
		shape := gamedata.InstrumentShape(inst.Kind)
		effect := newWaveNode(b.canvas, shape, pos, styles.PlotColorByID[e.id], inst.Period*0.95)
		b.addWaveEffect(effect)
		b.EventNote.Emit(e.id)

		effect.EventFinished.Connect(nil, func(r float64) {
			for _, t := range b.targets {
				canHit := t.outline || t.instrument == inst.Kind || t.instrument == gamedata.AnyInstrument
				if !canHit {
					continue
				}
				if t.pos.DistanceTo(pos) > float64(t.r*0.7)+(r*0.7) {
					continue
				}
				if !t.OnDamage(inst.ID) {
					continue
				}
				clr1 := styles.TargetColor
				clr2 := styles.TargetColor
				clr3 := styles.TargetColor
				isBonus := t.instrument == gamedata.AnyInstrument
				if t.outline && t.instrument != inst.Kind {
					b.penalty = true
					clr1 = styles.TargetMissColorRed
					clr2 = styles.TargetMissColorGreen
					clr3 = styles.TargetMissColorBlue
				}
				if isBonus {
					clr1 = styles.TargetColorBonus
					clr2 = styles.TargetColorBonus
					clr3 = styles.TargetColorBonus
					b.optionalHits++
				} else {
					b.targetsLeft--
				}
				d := 2 * (float64(t.r) / b.ctx.Scaler.Factor)
				shape := gamedata.InstrumentShape(t.instrument)
				offset := gmath.Vec{X: 2, Y: 2}
				b.addWaveEffect(newWaveNode(b.canvas, shape, t.pos.Sub(offset), clr1, d))
				b.addWaveEffect(newWaveNode(b.canvas, shape, t.pos, clr2, d))
				b.addWaveEffect(newWaveNode(b.canvas, shape, t.pos.Add(offset), clr3, d))
			}
		})
	}

	x := b.t
	for i, sig := range b.signals {
		y := b.prog.Instruments[i].Func.Run(x)
		sig.sprite.Visible = y >= -3 && y <= 3
		if sig.sprite.Visible {
			sig.pos = b.ctx.Scaler.ScaleXY(x, y)
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
		if t.Instrument != gamedata.AnyInstrument {
			b.targetsLeft++
		}
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
	b.victory = false
	b.penalty = false
	b.t = 0
	b.optionalHits = 0
	b.targetsLeft = 0
}
