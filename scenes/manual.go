package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/controls"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/styles"
)

type ManualController struct {
	state *session.State

	scene *ge.Scene

	backController ge.SceneController
}

func NewManualController(state *session.State, back ge.SceneController) *ManualController {
	return &ManualController{
		state:          state,
		backController: back,
	}
}

func (c *ManualController) Init(scene *ge.Scene) {
	c.scene = scene

	bigFont := scene.Context().Loader.LoadFont(assets.FontArcadeBig).Face

	d := scene.Dict()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(640, 10, nil)
	root.AddChild(rowContainer)

	title := d.Get("menu.manual")
	rowContainer.AddChild(eui.NewCenteredLabel(title, bigFont))
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.manual.howto"), func() {
		scene.Context().ChangeScene(NewHowtoController(c.state, c.backController))
	}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.manual.functions"), func() {
		scene.Context().ChangeScene(NewFuncListController(c.state, c.backController))
	}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.manual.key_bindings"), func() {
		scene.Context().ChangeScene(NewKeyBindingsController(c.state, c.backController))

	}))

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		c.back()
	}))

	initUI(scene, root)
}

func (c *ManualController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.back()
	}
}

func (c *ManualController) back() {
	c.scene.Context().ChangeScene(c.backController)
}
