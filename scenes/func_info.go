package scenes

import (
	"fmt"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/controls"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/exprc"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/stage"
	"github.com/quasilyte/sinecord/styles"
)

type FuncInfoController struct {
	state *session.State

	scene *ge.Scene

	funcIndex int
	funcs     []exprcFunc

	canvas *stage.Canvas

	titleWidget    *widget.Text
	infoTextWidget *widget.Text

	canvasImage   *ebiten.Image
	canvasImageBg *ebiten.Image
	canvasWidget  *widget.Graphic

	funcLabel *ge.Label

	backController ge.SceneController
}

func NewFuncInfoController(state *session.State, funcIndex int, back ge.SceneController) *FuncInfoController {
	return &FuncInfoController{
		state:          state,
		backController: back,
		funcIndex:      funcIndex,
	}
}

func (c *FuncInfoController) Init(scene *ge.Scene) {
	c.scene = scene

	bigFont := scene.Context().Loader.LoadFont(assets.FontArcadeBig).Face
	monoFont := scene.Context().Loader.LoadFont(assets.FontMonospaceNormal).Face

	d := scene.Dict()

	funcs := sortedFuncList()
	c.funcs = funcs

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(640, 10, nil)
	root.AddChild(rowContainer)

	switcherGrid := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{false, true, false}, nil),
			widget.GridLayoutOpts.Spacing(8, 0))))

	var funcIndexSlider gmath.Slider
	funcIndexSlider.SetBounds(0, len(c.funcs)-1)
	funcIndexSlider.TrySetValue(c.funcIndex)

	prevButton := eui.NewButton(c.state.UIResources, "<", func() {
		funcIndexSlider.Dec()
		c.funcIndex = funcIndexSlider.Value()
		c.updateInfo()
	})

	nextButton := eui.NewButton(c.state.UIResources, ">", func() {
		funcIndexSlider.Inc()
		c.funcIndex = funcIndexSlider.Value()
		c.updateInfo()
	})

	c.titleWidget = eui.NewCenteredLabel(funcs[c.funcIndex].Name, bigFont)

	switcherGrid.AddChild(prevButton)
	switcherGrid.AddChild(c.titleWidget)
	switcherGrid.AddChild(nextButton)

	rowContainer.AddChild(switcherGrid)
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	panel := eui.NewPanelWithPadding(c.state.UIResources, 1024, 160, widget.NewInsetsSimple(24))
	rowContainer.AddChild(panel)

	c.infoTextWidget = eui.NewLabel("", monoFont, widget.TextOpts.MaxWidth(960))
	panel.AddChild(c.infoTextWidget)

	c.canvasImageBg = scene.LoadImage(assets.ImagePlotBackground).Data
	c.canvasImage = ebiten.NewImage(c.canvasImageBg.Bounds().Dx(), c.canvasImageBg.Bounds().Dy())
	canvasPanel := eui.NewPanelWithPadding(c.state.UIResources, 0, 0, widget.Insets{
		Left:   24,
		Right:  24,
		Top:    24,
		Bottom: 24,
	})
	c.canvasWidget = widget.NewGraphic(
		widget.GraphicOpts.Image(c.canvasImage),
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)
	canvasPanel.AddChild(c.canvasWidget)

	rowContainer.AddChild(canvasPanel)

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		c.back()
	}))

	ctx := stage.NewContext(stage.Config{
		MaxInstruments: 1,
	})
	ctx.Scaler = c.state.PlotScaler
	c.canvas = stage.NewCanvas(ctx, scene, c.canvasImage)

	c.funcLabel = scene.NewLabel(assets.FontMonospaceNormal)
	c.funcLabel.Width = float64(c.canvasImage.Bounds().Dx())
	c.funcLabel.Height = 64
	c.funcLabel.AlignHorizontal = ge.AlignHorizontalRight
	c.funcLabel.AlignVertical = ge.AlignVerticalCenter
	c.funcLabel.Pos.Offset = gmath.Vec{X: 400, Y: 700}
	c.funcLabel.ColorScale.SetColor(styles.NormalTextColor)
	scene.AddGraphicsAbove(c.funcLabel, 1)

	c.updateInfo()

	initUI(scene, root)
}

func (c *FuncInfoController) funcExample() string {
	fn := c.funcs[c.funcIndex]
	var args []string
	addHalf := false

	complexArg := false
	switch fn.Name {
	case "abs", "after", "until", "mod", "sign":
		complexArg = true
	}

	for _, a := range fn.Args {
		if a == "x" {
			if complexArg {
				args = append(args, "sin(x)")
			} else {
				args = append(args, "x")
			}
			continue
		}
		addHalf = true
		switch a {
		case "min_value":
			args = append(args, "1")
		case "max_value":
			args = append(args, "2")
		case "edge1", "edge":
			args = append(args, "2")
		case "edge2", "threshold":
			args = append(args, "10")
		case "value1":
			args = append(args, "x")
		case "value2", "x2":
			args = append(args, "1.5")
		case "exponent":
			args = append(args, "2")
		case "divisor":
			args = append(args, "3")
		default:
			fmt.Println(a)
			args = append(args, "1")
		}
	}
	s := fn.Name + "(" + strings.Join(args, ", ") + ")"
	if addHalf {
		s += " + 0.5"
	}
	return s
}

func (c *FuncInfoController) updateInfo() {
	snippet := c.funcExample()

	fn := c.funcs[c.funcIndex]
	c.titleWidget.Label = fn.Name
	c.infoTextWidget.Label = c.funcInfoText()

	c.canvas.ClearPlot(0)

	compiled, err := exprc.Compile(snippet)
	if err != nil {
		panic(err)
	}

	c.canvas.RedrawPlot(0, compiled, nil)
	c.canvas.Draw()

	c.funcLabel.Text = "y = " + snippet
}

func (c *FuncInfoController) funcInfoText() string {
	fn := c.funcs[c.funcIndex]
	lines := []string{
		"Function " + fn.Name + "(" + strings.Join(fn.Args, ", ") + ")",
		"",
		"Description: " + fn.Doc + ".",
	}
	return strings.Join(lines, "\n")
}

func (c *FuncInfoController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.back()
	}
}

func (c *FuncInfoController) back() {
	c.scene.Context().ChangeScene(NewFuncListController(c.state, c.backController))
}
