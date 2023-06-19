package scenes

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/gtask"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/stage"
	"github.com/quasilyte/sinecord/styles"
	"github.com/quasilyte/sinecord/synthdb"
)

type StageController struct {
	state *session.State
	scene *ge.Scene

	config stage.Config
	track  gamedata.Track

	canvas *stage.Canvas
	synth  *stage.Synthesizer
	board  *stage.Board

	prog    stage.SynthProgram
	samples *stage.SampleSet
	player  *audio.Player

	waveUpdateDelay float64
	samplesBuf      []float64

	instrumentIcons []*ebiten.Image

	canvasWidget  *widget.Graphic
	canvasImage   *ebiten.Image
	canvasImageBg *ebiten.Image

	currentMode stageMode
	statusLabel *widget.Text
}

type stageMode int

const (
	stageReady stageMode = iota
	stagePlaying
	stageEncoding
)

func NewStageController(state *session.State, config stage.Config) *StageController {
	return &StageController{
		state:  state,
		config: config,
	}
}

func (c *StageController) Init(scene *ge.Scene) {
	c.scene = scene

	d := scene.Dict()

	c.instrumentIcons = make([]*ebiten.Image, c.config.MaxInstruments)
	for i := range c.instrumentIcons {
		c.instrumentIcons[i] = ebiten.NewImage(26, 26)
	}

	ctx := stage.NewContext(c.config)
	ctx.PlotScale = 46
	ctx.PlotOffset = gmath.Vec{
		X: 4,
		Y: 46 * 3,
	}

	c.canvasImageBg = scene.LoadImage(assets.ImagePlotBackground).Data
	c.canvasImage = ebiten.NewImage(c.canvasImageBg.Bounds().Dx(), c.canvasImageBg.Bounds().Dy())
	c.canvas = stage.NewCanvas(ctx, scene, c.canvasImage)

	smallFont := scene.Context().Loader.LoadFont(assets.FontArcadeSmall).Face

	c.synth = stage.NewSynthesizer(ctx, synthdb.TimGM6mb)
	scene.AddObject(c.synth)

	c.board = stage.NewBoard(stage.BoardConfig{
		Canvas:  c.canvas,
		Targets: c.config.Targets,
	})
	c.board.Init(scene)
	c.board.EventNote.Connect(c, func(instrumentID int) {
		// clr := styles.PlotColorByID[instrumentID]
		// c.canvas.WaveColor.R = (float32(clr.R) / 255.0)
		// c.canvas.WaveColor.G = (float32(clr.G) / 255.0)
		// c.canvas.WaveColor.B = (float32(clr.B) / 255.0)
		// c.canvas.WaveColor.R *= (float32(clr.R) / 255.0) + 1
		// c.canvas.WaveColor.G *= (float32(clr.G) / 255.0) + 1
		// c.canvas.WaveColor.B *= (float32(clr.B) / 255.0) + 1
	})

	c.synth.EventRedrawPlotRequest.Connect(nil, func(id int) {
		f := c.synth.GetInstrumentFunction(id)
		if f == nil {
			c.canvas.ClearPlot(id)
			return
		}
		c.canvas.RedrawPlot(id, f)
	})

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	outerGrid := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				// StretchHorizontal:  true,
				// StretchVertical:    true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch(nil, nil),
			widget.GridLayoutOpts.Spacing(8, 8))))

	root.AddChild(outerGrid)

	instrumentsGrid := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(6),
			widget.GridLayoutOpts.Stretch([]bool{false, false, false, true, false, false, false}, nil),
			widget.GridLayoutOpts.Spacing(4, 8),
		)))
	outerGrid.AddChild(instrumentsGrid)

	patchNames := make([]string, len(synthdb.TimGM6mb.Instruments))
	for i, inst := range synthdb.TimGM6mb.Instruments {
		patchNames[i] = inst.Name
	}

	volumeLevels := []float64{
		0.4,
		0.6,
		0.8,
		0.9,
		1.0,
	}
	periods := []float64{
		0.10,
		0.15,
		0.20,
		0.25,
		0.30,
		0.35,

		0.40,
		0.50,
		0.60,
		0.70,
		0.80,
		0.90,

		1.00,
		1.25,
		1.50,
		1.75,
		2.00,
	}
	periodLabels := make([]string, len(periods))
	for i := range periods {
		periodLabels[i] = fmt.Sprintf("%.2f", periods[i])
	}

	for i := 0; i < c.config.MaxInstruments; i++ {
		instrumentID := i

		var loadedInstrument *gamedata.InstrumentSettings
		if !c.track.IsEmpty() && len(c.track.Instruments) > instrumentID {
			if c.track.Instruments[instrumentID].Function != "" {
				loadedInstrument = &c.track.Instruments[instrumentID]
			}
		}

		colorPanel := eui.NewPanel(c.state.UIResources, 0, 0)
		c.canvas.DrawInstrumentIcon(c.instrumentIcons[instrumentID], synthdb.BassInstrument, styles.PlotColorByID[instrumentID])
		colorRect := widget.NewGraphic(
			widget.GraphicOpts.Image(c.instrumentIcons[instrumentID]),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
			),
		)
		colorPanel.AddChild(colorRect)
		instrumentsGrid.AddChild(colorPanel)

		textTnput := eui.NewTextInput(c.state.UIResources,
			widget.TextInputOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(880, 0),
				widget.WidgetOpts.ToolTip(
					widget.NewToolTip(
						widget.ToolTipOpts.Content(eui.NewTooltip(c.state.UIResources, "f(x)")),
						widget.ToolTipOpts.Delay(time.Second),
					),
				),
			),
			widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
				c.synth.SetInstrumentFunction(instrumentID, strings.ToLower(args.InputText))
			}),
			widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
				good := true
				if len(newInputText) > 60 {
					good = false
				}
				if good {
					for _, ch := range newInputText {
						if !unicode.IsPrint(ch) || ch >= utf8.RuneSelf {
							good = false
							break
						}
					}
				}
				return good, nil
			}),
		)
		instrumentsGrid.AddChild(textTnput)
		if loadedInstrument != nil {
			textTnput.InputText = loadedInstrument.Function
		}

		stepPeriodLevel := 4
		if loadedInstrument != nil {
			stepPeriodLevel = xslices.Index(periods, loadedInstrument.Period)
		}
		c.synth.SetInstrumentPeriod(instrumentID, periods[stepPeriodLevel])
		instrumentsGrid.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
			Resources:  c.state.UIResources,
			Input:      c.state.Input,
			ValueNames: periodLabels,
			Value:      &stepPeriodLevel,
			Tooltip:    eui.NewTooltip(c.state.UIResources, "step period in seconds"),
			OnPressed: func() {
				c.synth.SetInstrumentPeriod(instrumentID, periods[stepPeriodLevel])
			},
		}))

		patchIndex := 0
		if loadedInstrument != nil {
			patchIndex = xslices.IndexWhere(synthdb.TimGM6mb.Instruments, func(inst *synthdb.Instrument) bool {
				return inst.Name == loadedInstrument.InstrumentName
			})
		}
		c.selectInstrument(instrumentID, patchIndex)
		instrumentsGrid.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
			Resources:  c.state.UIResources,
			Input:      c.state.Input,
			ValueNames: patchNames,
			Value:      &patchIndex,
			MinWidth:   320,
			Tooltip:    eui.NewTooltip(c.state.UIResources, "instrument style"),
			OnPressed: func() {
				c.selectInstrument(instrumentID, patchIndex)
			},
		}))

		volumeLevel := 4
		if loadedInstrument != nil {
			volumeLevel = xslices.Index(volumeLevels, loadedInstrument.Volume)
		}
		c.synth.SetInstrumentVolume(instrumentID, volumeLevels[volumeLevel])
		instrumentsGrid.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
			Resources:  c.state.UIResources,
			Input:      c.state.Input,
			ValueNames: []string{"20%", "40%", "60%", "80%", "100%"},
			MinWidth:   120,
			Value:      &volumeLevel,
			Tooltip:    eui.NewTooltip(c.state.UIResources, "instrument volume level"),
			OnPressed: func() {
				c.synth.SetInstrumentVolume(instrumentID, volumeLevels[volumeLevel])
			},
		}))

		instrumentEnabled := (instrumentID == 0 && loadedInstrument == nil) ||
			(loadedInstrument != nil && loadedInstrument.Enabled)
		c.synth.SetInstrumentEnabled(instrumentID, instrumentEnabled)
		instrumentsGrid.AddChild(eui.NewBoolSelectButton(eui.BoolSelectButtonConfig{
			Resources:  c.state.UIResources,
			ValueNames: []string{"off", "on"},
			MinWidth:   120,
			Value:      &instrumentEnabled,
			Tooltip:    eui.NewTooltip(c.state.UIResources, "enable or disable this instrument"),
			OnPressed: func() {
				c.synth.SetInstrumentEnabled(instrumentID, instrumentEnabled)
			},
		}))
	}

	statusLabelContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.Insets{Top: 8}),
		)),
	)
	statusLabel := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.Text("status: ready", smallFont, styles.NormalTextColor),
	)
	statusLabelContainer.AddChild(statusLabel)
	outerGrid.AddChild(statusLabelContainer)
	c.statusLabel = statusLabel

	{
		playerGridContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		)

		panel := eui.NewPanelWithPadding(c.state.UIResources, 0, 0, widget.Insets{
			Left:   24,
			Right:  24,
			Top:    24,
			Bottom: 24,
		})

		playerGrid := widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
				}),
			),
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(8)),
				widget.GridLayoutOpts.Spacing(4, 8),
			)))

		buttonsGrid := widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(256, 0),
			),
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, nil),
				widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(8)),
				widget.GridLayoutOpts.Spacing(4, 8),
			)))

		buttonsGrid.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.stage.run"), c.onPlayPressed))

		stopButton := eui.NewButton(c.state.UIResources, "stop", func() {
			if c.currentMode != stagePlaying {
				return
			}
			c.setMode(stageReady)
			c.board.ClearProgram()
			c.player.Pause()
		})
		buttonsGrid.AddChild(stopButton)

		saveButton := eui.NewButton(c.state.UIResources, "save", func() {
			back := NewStageController(c.state, c.config)
			t := c.synth.ExportTrack()
			back.track = t
			c.changeScene(NewSaverController(c.state, t, back))
		})
		buttonsGrid.AddChild(saveButton)

		loadButton := eui.NewButton(c.state.UIResources, "load", func() {
			back := NewStageController(c.state, c.config)
			loader := NewLoaderController(c.state, back)
			loader.EventLoaded.Connect(nil, func(track gamedata.Track) {
				back.track = track
			})
			c.changeScene(loader)
		})
		buttonsGrid.AddChild(loadButton)

		exitButton := eui.NewButton(c.state.UIResources, "exit", func() {})
		buttonsGrid.AddChild(exitButton)

		c.canvasWidget = widget.NewGraphic(
			widget.GraphicOpts.Image(c.canvasImage),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
			),
		)
		panel.AddChild(c.canvasWidget)

		playerGrid.AddChild(buttonsGrid)
		playerGrid.AddChild(panel)

		playerGridContainer.AddChild(playerGrid)

		outerGrid.AddChild(playerGridContainer)
	}

	initUI(scene, root)

	scene.AddGraphics(c)
}

func (c *StageController) selectInstrument(instrumentID int, patchIndex int) {
	c.synth.SetInstrumentPatch(instrumentID, patchIndex)
	kind := synthdb.TimGM6mb.Instruments[patchIndex].Kind
	c.instrumentIcons[instrumentID].Clear()
	c.canvas.DrawInstrumentIcon(c.instrumentIcons[instrumentID], kind, styles.PlotColorByID[instrumentID])
}

func (c *StageController) IsDisposed() bool { return false }

func (c *StageController) Draw(*ebiten.Image) {
	c.canvas.Draw()
}

func (c *StageController) Update(delta float64) {
	c.canvas.Running = c.currentMode == stagePlaying

	if c.currentMode == stagePlaying {
		c.waveUpdateDelay = gmath.ClampMin(c.waveUpdateDelay-delta, 0)
		if c.waveUpdateDelay == 0 {
			c.waveUpdateDelay = 0.1
			c.canvas.RenderWave(c.waveSamples())
		}

		if c.board.ProgramTick(delta) {
			c.board.ClearProgram()
			c.setMode(stageReady)
		}
	}

	c.canvas.Update(delta)
}

func (c *StageController) waveSamples() []float64 {
	c.samplesBuf = c.samplesBuf[:0]

	sampleRate := float64(c.samples.PerSecond)
	numSamples := len(c.samples.Left)

	currentSecond := c.player.Current().Seconds()
	currentSample := int(math.Round(currentSecond * sampleRate))
	samplesPerHalf := c.samples.PerSecond / 60

	fromSample := currentSample - samplesPerHalf
	toSample := currentSample + samplesPerHalf

	if fromSample < 0 || toSample >= numSamples {
		return nil
	}

	for i := fromSample; i < toSample; i++ {
		v := float64(c.samples.Left[i] + c.samples.Right[i])
		c.samplesBuf = append(c.samplesBuf, v)
	}

	return c.samplesBuf
}

func (c *StageController) onPlayPressed() {
	switch c.currentMode {
	case stageReady, stagePlaying:
		// OK
	default:
		return
	}

	// Don't spawn a task if we don't need to.
	if !c.synth.HasChanges() {
		c.runPlayer()
		return
	}

	c.setMode(stageEncoding)

	encodeTask := gtask.StartTask(func(ctx *gtask.TaskContext) {
		ctx.Progress.Total = 1.0
		samples, prog := c.synth.CreatePCM(&ctx.Progress.Current)
		if samples != nil {
			if c.player != nil {
				c.player.Play()
				c.player.Close()
			}
			pcm := generatePCM(samples.Left, samples.Right)
			c.player = c.scene.Audio().GetContext().NewPlayerFromBytes(pcm)
			c.samples = samples
			c.prog = prog
		}
	})
	encodeTask.EventProgress.Connect(nil, func(p gtask.TaskProgress) {
		c.statusLabel.Label = fmt.Sprintf("status: encoding (%d%%)", int(100*p.Current))
	})
	encodeTask.EventCompleted.Connect(nil, func(gsignal.Void) {
		c.runPlayer()
	})
	c.scene.AddObject(encodeTask)
}

func (c *StageController) runPlayer() {
	c.player.Rewind()
	c.player.Play()
	c.board.StartProgram(c.prog)
	c.canvas.Reset()
	c.setMode(stagePlaying)
}

func (c *StageController) changeScene(newScene ge.SceneController) {
	if c.player != nil {
		c.player.Pause()
	}
	c.scene.Context().ChangeScene(newScene)
}

func (c *StageController) setMode(m stageMode) {
	if c.currentMode == m {
		return
	}
	c.currentMode = m
	var modeText string
	switch m {
	case stageReady:
		modeText = "ready"
	case stagePlaying:
		modeText = "playing"
	case stageEncoding:
		modeText = "encoding"
	default:
		modeText = "unknown"
	}
	c.statusLabel.Label = "status: " + modeText
}
