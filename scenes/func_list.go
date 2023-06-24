package scenes

import (
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/controls"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/styles"
)

type FuncListController struct {
	state *session.State

	scene *ge.Scene

	backController ge.SceneController
}

func NewFuncListController(state *session.State, back ge.SceneController) *FuncListController {
	return &FuncListController{
		state:          state,
		backController: back,
	}
}

func (c *FuncListController) Init(scene *ge.Scene) {
	c.scene = scene

	bigFont := scene.Context().Loader.LoadFont(assets.FontArcadeBig).Face
	monoFont := scene.Context().Loader.LoadFont(assets.FontMonospaceNormal).Face

	d := scene.Dict()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(640, 10, nil)
	root.AddChild(rowContainer)

	title := d.Get("menu.manual.functions")
	rowContainer.AddChild(eui.NewCenteredLabel(title, bigFont))
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	panel := eui.NewPanelWithPadding(c.state.UIResources, 0, 0, widget.NewInsetsSimple(24))
	rowContainer.AddChild(panel)

	funcs := sortedFuncList()

	funcRows := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch(nil, nil),
			widget.GridLayoutOpts.Spacing(0, 8))))
	panel.AddChild(funcRows)

	for i := range funcs {
		fn := funcs[i]
		funcIndex := i

		pairGrid := widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal: true,
			})),
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Stretch([]bool{false, true}, nil),
				widget.GridLayoutOpts.Spacing(8, 0))))

		funcLabel := fn.Name + "(" + strings.Join(fn.Args, ", ") + ")"
		pairGrid.AddChild(eui.NewTextButton(c.state.UIResources, funcLabel, func() {
			scene.Context().ChangeScene(NewFuncInfoController(c.state, funcIndex, c.backController))
		}))
		pairGrid.AddChild(eui.NewLabel(fn.Doc, monoFont))

		funcRows.AddChild(pairGrid)
	}

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		c.back()
	}))

	initUI(scene, root)
}

func (c *FuncListController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.back()
	}
}

func (c *FuncListController) back() {
	c.scene.Context().ChangeScene(NewManualController(c.state, c.backController))
}
