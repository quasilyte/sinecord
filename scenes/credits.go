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

type CreditsController struct {
	state *session.State

	scene *ge.Scene
}

func NewCreditsController(state *session.State) *CreditsController {
	return &CreditsController{
		state: state,
	}
}

func (c *CreditsController) Init(scene *ge.Scene) {
	c.scene = scene

	bigFont := scene.Context().Loader.LoadFont(assets.FontArcadeBig).Face
	monoFont := scene.Context().Loader.LoadFont(assets.FontMonospaceNormal).Face

	d := scene.Dict()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(400, 10, nil)
	root.AddChild(rowContainer)

	title := d.Get("menu.credits")
	rowContainer.AddChild(eui.NewCenteredLabel(title, bigFont))
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	panel := eui.NewPanelWithPadding(c.state.UIResources, 0, 0, widget.NewInsetsSimple(24))
	rowContainer.AddChild(panel)

	creditsText := `Sinecord - a game by quasilyte.

Created during Ebitengine Game Jam 2023.

Special thanks:

* Hajime Hoshi and Ebitengine community
* Mark Caprenter: ebitenui GUI library
* Nobuaki Tanaka: meltysynth SoundFont synthesizer library

Assets:

* TimGM6mb SoundFont by Tim Brechbill
* Arcade font by Yuji Adachi
* White Rabbit font by Matthew Welch

Thank you for playing my game.`
	panel.AddChild(eui.NewLabel(creditsText, monoFont, widget.TextOpts.MaxWidth(1024)))

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		c.back()
	}))

	initUI(scene, root)
}

func (c *CreditsController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.back()
	}
}

func (c *CreditsController) back() {
	c.scene.Context().ChangeScene(NewMainMenuController(c.state))
}
