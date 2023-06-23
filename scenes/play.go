package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/controls"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/stage"
	"github.com/quasilyte/sinecord/styles"
)

type PlayController struct {
	state *session.State

	scene *ge.Scene
}

func NewPlayController(state *session.State) *PlayController {
	return &PlayController{
		state: state,
	}
}

func (c *PlayController) Init(scene *ge.Scene) {
	c.scene = scene

	bigFont := scene.Context().Loader.LoadFont(assets.FontArcadeBig).Face

	d := scene.Dict()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(400, 10, nil)
	root.AddChild(rowContainer)

	title := d.Get("menu.play")
	rowContainer.AddChild(eui.NewCenteredLabel(title, bigFont))
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	rowContainer.AddChild(eui.NewButtonWithConfig(c.state.UIResources, eui.ButtonConfig{
		Text:         d.Get("menu.play.missions"),
		TooltipLabel: d.Get("menu.play.missions.desc"),
		OnClick: func() {
			scene.Context().ChangeScene(NewMissionsController(c.state))
		},
	}))

	rowContainer.AddChild(eui.NewButtonWithConfig(c.state.UIResources, eui.ButtonConfig{
		Text:         d.Get("menu.play.sandbox"),
		TooltipLabel: d.Get("menu.play.sandbox.desc"),
		OnClick: func() {
			scene.Context().ChangeScene(NewStageController(c.state, stage.Config{
				MaxInstruments: 5,
				Mode:           gamedata.SandboxMode,
			}))
		},
	}))

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		c.back()
	}))

	initUI(scene, root)
}

func (c *PlayController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.back()
	}
}

func (c *PlayController) back() {
	c.scene.Context().ChangeScene(NewMainMenuController(c.state))
}
