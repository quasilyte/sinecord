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
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/stage"
	"github.com/quasilyte/sinecord/styles"
	"github.com/quasilyte/sinecord/synthdb"
)

type StageController struct {
	state *session.State

	config stage.Config

	running bool

	canvas *stage.Canvas
	synth  *stage.Synthesizer
	board  *stage.Board

	prog    stage.SynthProgram
	samples *stage.SampleSet
	player  *audio.Player

	waveUpdateDelay float64
	samplesBuf      []float64

	canvasWidget  *widget.Graphic
	canvasImage   *ebiten.Image
	canvasImageBg *ebiten.Image
}

func NewStageController(state *session.State, config stage.Config) *StageController {
	return &StageController{
		state:  state,
		config: config,
	}
}

func (c *StageController) Init(scene *ge.Scene) {
	d := scene.Dict()

	ctx := stage.NewContext(c.config)

	c.canvasImageBg = scene.LoadImage(assets.ImagePlotBackground).Data
	c.canvasImage = ebiten.NewImage(c.canvasImageBg.Bounds().Dx(), c.canvasImageBg.Bounds().Dy())
	c.canvas = stage.NewCanvas(ctx, scene, c.canvasImage)

	smallFont := scene.Context().Loader.LoadFont(assets.FontArcadeSmall).Face

	c.synth = stage.NewSynthesizer(ctx, synthdb.TimGM6mb)
	scene.AddObject(c.synth)

	c.board = stage.NewBoard(scene, stage.BoardConfig{
		Canvas:    c.canvas,
		PlotScale: 46,
		PlotOffset: gmath.Vec{
			X: 4,
			Y: 46 * 3,
		},
	})
	c.board.EventNote.Connect(c, func(instrumentID int) {
		// clr := styles.PlotColorByID[instrumentID]
		// c.canvas.WaveColor.R = (float32(clr.R) / 255.0)
		// c.canvas.WaveColor.G = (float32(clr.G) / 255.0)
		// c.canvas.WaveColor.B = (float32(clr.B) / 255.0)
		// c.canvas.WaveColor.R *= (float32(clr.R) / 255.0) + 1
		// c.canvas.WaveColor.G *= (float32(clr.G) / 255.0) + 1
		// c.canvas.WaveColor.B *= (float32(clr.B) / 255.0) + 1
	})

	c.synth.EventRecompileShaderRequest.Connect(nil, func(id int) {
		fx := c.synth.GetInstrumentFunction(id)
		if fx == "" {
			c.canvas.SetShader(id, nil)
			return
		}
		shader, err := stage.CompilePlotShader(id, fx)
		if err != nil {
			fmt.Println(err)
		}
		c.canvas.SetShader(id, shader)
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

		0.40,
		0.50,
		0.60,
		0.70,
		0.80,

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

		colorPanel := eui.NewPanel(c.state.UIResources, 0, 0)
		colorIcon := ebiten.NewImage(20, 20)
		ebitenutil.DrawRect(colorIcon, 0, 0, 20, 20, styles.PlotColorByID[i])
		colorRect := widget.NewGraphic(
			widget.GraphicOpts.Image(colorIcon),
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
				widget.WidgetOpts.MinSize(760, 0),
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
				if len(newInputText) > 48 {
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

		stepPeriodLevel := 2
		c.synth.SetInstrumentPeriod(instrumentID, 0.1*float64(stepPeriodLevel)+0.1)
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
		c.synth.SetInstrumentPatch(instrumentID, patchIndex)
		instrumentsGrid.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
			Resources:  c.state.UIResources,
			Input:      c.state.Input,
			ValueNames: patchNames,
			Value:      &patchIndex,
			MinWidth:   320,
			Tooltip:    eui.NewTooltip(c.state.UIResources, "instrument style"),
			OnPressed: func() {
				c.synth.SetInstrumentPatch(instrumentID, patchIndex)
			},
		}))

		volumeLevel := 4
		c.synth.SetInstrumentVolume(instrumentID, volumeLevels[volumeLevel])
		instrumentsGrid.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
			Resources:  c.state.UIResources,
			Input:      c.state.Input,
			ValueNames: []string{"20%", "40%", "60%", "80%", "100%"},
			Value:      &volumeLevel,
			Tooltip:    eui.NewTooltip(c.state.UIResources, "instrument volume level"),
			OnPressed: func() {
				c.synth.SetInstrumentVolume(instrumentID, volumeLevels[volumeLevel])
			},
		}))

		instrumentEnabled := instrumentID == 0
		c.synth.SetInstrumentEnabled(instrumentID, instrumentEnabled)
		instrumentsGrid.AddChild(eui.NewBoolSelectButton(eui.BoolSelectButtonConfig{
			Resources:  c.state.UIResources,
			ValueNames: []string{"off", "on"},
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

		buttonsGrid.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.stage.run"), func() {
			c.running = true
			samples, prog := c.synth.CreatePCM()
			if samples != nil {
				if c.player != nil {
					c.player.Play()
					c.player.Close()
				}
				pcm := generatePCM(samples.Left, samples.Right)
				c.player = scene.Audio().GetContext().NewPlayerFromBytes(pcm)
				c.samples = samples
				c.prog = prog
			}
			c.player.Rewind()
			c.player.Play()
			c.board.StartProgram(c.prog)
		}))

		stopButton := eui.NewButton(c.state.UIResources, "stop", func() {
			if !c.running {
				return
			}
			c.running = false
			c.board.ClearProgram()
			c.player.Pause()
		})
		buttonsGrid.AddChild(stopButton)

		saveButton := eui.NewButton(c.state.UIResources, "save", func() {})
		saveButton.GetWidget().Disabled = true
		buttonsGrid.AddChild(saveButton)

		loadButton := eui.NewButton(c.state.UIResources, "load", func() {})
		loadButton.GetWidget().Disabled = true
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

func (c *StageController) IsDisposed() bool { return false }

func (c *StageController) Draw(*ebiten.Image) {
	c.canvas.Draw()
}

func (c *StageController) Update(delta float64) {
	c.canvas.Running = c.running

	if c.running {
		c.waveUpdateDelay = gmath.ClampMin(c.waveUpdateDelay-delta, 0)
		if c.waveUpdateDelay == 0 {
			c.waveUpdateDelay = 0.1
			c.canvas.RenderWave(c.waveSamples())
		}

		if c.board.ProgramTick(delta) {
			c.board.ClearProgram()
			c.running = false
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
