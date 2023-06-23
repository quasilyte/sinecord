package scenes

import (
	"os"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/styles"
	"github.com/quasilyte/sinecord/synthdb"
)

type MainMenuController struct {
	state *session.State
}

func NewMainMenuController(state *session.State) *MainMenuController {
	return &MainMenuController{
		state: state,
	}
}

func (c *MainMenuController) Init(scene *ge.Scene) {
	bigFont := scene.Context().Loader.LoadFont(assets.FontArcadeBig).Face
	smallFont := scene.Context().Loader.LoadFont(assets.FontArcadeSmall).Face

	d := scene.Dict()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(400, 10, nil)
	root.AddChild(rowContainer)

	rowContainer.AddChild(eui.NewCenteredLabel("sinecord", bigFont))
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.play"), func() {
		scene.Context().ChangeScene(NewPlayController(c.state))
	}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.profile"), func() {}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.manual"), func() {}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.settings"), func() {}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.credits"), func() {}))

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.exit"), func() {
		os.Exit(0)
	}))

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))
	rowContainer.AddChild(eui.NewCenteredLabel("Ebitengine Game Jam 2023 edition", smallFont))

	initUI(scene, root)

	synthdb.TimGM6mb.Load(assets.SoundFontTimGM6mb)

	assets.ReadLevelsData()
}

func (c *MainMenuController) Update(delta float64) {}
