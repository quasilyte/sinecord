package scenes

import (
	"fmt"
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

type StageConfig struct {
	MaxInstruments int
}

type StageController struct {
	state *session.State

	config StageConfig

	running bool

	canvas *stage.Canvas
	synth  *stage.Synthesizer
	board  *stage.Board

	canvasUpdater *canvasUpdater

	player *audio.Player

	canvasWidget  *widget.Graphic
	canvasImage   *ebiten.Image
	canvasImageBg *ebiten.Image
}

func NewStageController(state *session.State, config StageConfig) *StageController {
	return &StageController{
		state:  state,
		config: config,
	}
}

func (c *StageController) Init(scene *ge.Scene) {
	d := scene.Dict()

	ctx := stage.NewContext()

	c.canvas = stage.NewCanvas()

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

	c.synth.EventRecompileShaderRequest.Connect(nil, func(id int) {
		fx := c.synth.GetInstrumentFunction(id)
		if fx == "" {
			c.canvasUpdater.fnShaders[id] = nil
			return
		}
		shader, err := stage.CompilePlotShader(id, fx)
		if err != nil {
			fmt.Println(err)
		}
		c.canvasUpdater.fnShaders[id] = shader
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
				widget.WidgetOpts.MinSize(640, 0),
				widget.WidgetOpts.ToolTip(
					widget.NewToolTip(
						widget.ToolTipOpts.Content(eui.NewTooltip(c.state.UIResources, "f(x)")),
						widget.ToolTipOpts.Delay(time.Second),
					),
				),
			),
			widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
				c.synth.SetInstrumentFunction(instrumentID, args.InputText)
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

		stepPeriodLevel := 1
		c.synth.SetInstrumentPeriod(instrumentID, 0.1*float64(stepPeriodLevel)+0.1)
		instrumentsGrid.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
			Resources:  c.state.UIResources,
			Input:      c.state.Input,
			ValueNames: []string{"0.1", "0.2", "0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0"},
			Value:      &stepPeriodLevel,
			Tooltip:    eui.NewTooltip(c.state.UIResources, "step period in seconds"),
			OnPressed: func() {
				c.synth.SetInstrumentPeriod(instrumentID, 0.1*float64(stepPeriodLevel)+0.1)
			},
		}))

		patchIndex := 0
		c.synth.SetInstrumentPatch(instrumentID, patchIndex)
		instrumentsGrid.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
			Resources:  c.state.UIResources,
			Input:      c.state.Input,
			ValueNames: patchNames,
			Value:      &patchIndex,
			Tooltip:    eui.NewTooltip(c.state.UIResources, "instrument style"),
			OnPressed: func() {
				c.synth.SetInstrumentPatch(instrumentID, patchIndex)
			},
		}))

		instrumentsGrid.AddChild(eui.NewButton(c.state.UIResources, "100%", nil))

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

	outerGrid.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))
	outerGrid.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.stage.run"), func() {
		c.running = true
		pcm := c.synth.CreatePCM()
		if pcm != nil {
			if c.player != nil {
				c.player.Play()
				c.player.Close()
			}
			c.player = audio.NewPlayerFromBytes(scene.Audio().GetContext(), pcm)
		}
		c.player.Rewind()
		c.player.Play()
		c.board.StartProgram(c.synth.CreateProgram())
	}))

	{
		width := 1536
		height := 320
		panel := eui.NewPanel(c.state.UIResources, width, height)

		c.canvasImageBg = scene.LoadImage(assets.ImagePlotBackground).Data

		c.canvasImage = ebiten.NewImage(c.canvasImageBg.Bounds().Dx(), c.canvasImageBg.Bounds().Dy())

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

		outerGrid.AddChild(panel)
	}

	c.canvasUpdater = &canvasUpdater{
		scene:       scene,
		canvasImage: c.canvasImage,
		stageCanvas: c.canvas,
		scratch:     ebiten.NewImage(c.canvasImage.Size()),
		fnShaders:   make([]*ebiten.Shader, c.config.MaxInstruments),
	}
	scene.AddGraphics(c.canvasUpdater)

	initUI(scene, root)
}

func (c *StageController) Update(delta float64) {
	c.canvasUpdater.DrawShaders = !c.running
	if c.running {
		if c.board.ProgramTick(delta) {
			c.board.ClearProgram()
			fmt.Println("finished")
			c.running = false
		}
	}
	c.canvasUpdater.Update(delta)
}

type canvasUpdater struct {
	scene       *ge.Scene
	canvasImage *ebiten.Image
	stageCanvas *stage.Canvas
	scratch     *ebiten.Image

	time float64

	fnShaders []*ebiten.Shader

	DrawShaders bool
}

func (c *canvasUpdater) IsDisposed() bool { return false }

func (c *canvasUpdater) Update(delta float64) {
	c.time += delta
}

func (c *canvasUpdater) Draw(*ebiten.Image) {
	c.canvasImage.Clear()

	plotBackground := c.scene.LoadImage(assets.ImagePlotBackground).Data

	var drawOptions ebiten.DrawImageOptions
	c.canvasImage.DrawImage(plotBackground, &drawOptions)

	if c.DrawShaders {
		width := plotBackground.Bounds().Dx()
		height := plotBackground.Bounds().Dy()
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

	c.stageCanvas.Draw(c.canvasImage)

	// {
	// 	c.scratch.Clear()
	// 	c.scratch.DrawImage(c.canvas, &drawOptions)

	// 	var options ebiten.DrawRectShaderOptions
	// 	options.Uniforms = map[string]any{
	// 		"Seed": c.scene.Rand().FloatRange(0, 9999999999999),
	// 		"Tick": c.time,
	// 	}
	// 	options.CompositeMode = ebiten.CompositeModeCopy
	// 	options.Images[0] = c.scratch
	// 	videoDistortionShader := c.scene.Context().Loader.LoadShader(assets.ShaderVideoDistortion).Data
	// 	c.canvas.DrawRectShader(width, height, videoDistortionShader, &options)
	// }
}
