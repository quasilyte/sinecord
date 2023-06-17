package scenes

import (
	"fmt"
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

	// {
	// 	s := scene.NewSprite(assets.ImagePlotBackground)
	// 	s.Centered = false
	// 	s.Shader = scene.NewShader(assets.ShaderSine)
	// 	scene.AddGraphics(s)
	// }

	{
		for i, p := range assets.SoundFontTimGM6mb.Presets {
			fmt.Printf("[i=%d patch=%d] %s:\n", i, p.PatchNumber, p.Name)
			for _, region := range p.Regions {
				fmt.Printf(" > %s [%d, %d] (%d, %d)\n", region.Instrument.Name,
					region.GetKeyRangeStart(), region.GetKeyRangeEnd(),
					region.GetVelocityRangeStart(), region.GetVelocityRangeEnd())
			}
			// fmt.Printf("[%d] %s\n", p.PatchNumber, p.Name)
		}
		// for i, inst := range assets.SoundFontTimGM6mb.Instruments {
		// 	fmt.Println(i, inst.Name)
		// 	for _, r := range inst.Regions {
		// 		if r.GetVelocityRangeStart() != 0 {
		// 			panic("start")
		// 		}
		// 		if r.GetVelocityRangeEnd() != 127 {
		// 			panic("end")
		// 		}
		// 	}
		// }
		// inst0 := assets.SoundFontTimGM6mb.Instruments[11]
		// for _, r := range inst0.Regions {
		// 	fmt.Printf("  > keyrange (%d, %d) velocity (%d, %d)\n",
		// 		r.GetKeyRangeStart(), r.GetKeyRangeEnd(),
		// 		r.GetVelocityRangeStart(), r.GetVelocityRangeEnd())
		// }

		synthdb.TimGM6mb.Load(assets.SoundFontTimGM6mb)
		fmt.Println(synthdb.TimGM6mb.Instruments[0])
	}
}

func (c *MainMenuController) Update(delta float64) {}
