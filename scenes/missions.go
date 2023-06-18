package scenes

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/stage"
	"github.com/quasilyte/sinecord/styles"
	"github.com/quasilyte/sinecord/synthdb"
)

type MissionsController struct {
	state *session.State
}

func NewMissionsController(state *session.State) *MissionsController {
	return &MissionsController{
		state: state,
	}
}

func (c *MissionsController) Init(scene *ge.Scene) {
	bigFont := scene.Context().Loader.LoadFont(assets.FontArcadeBig).Face
	normalFont := scene.Context().Loader.LoadFont(assets.FontArcadeNormal).Face

	d := scene.Dict()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(400, 10, nil)
	root.AddChild(rowContainer)

	title := d.Get("menu.play.missions")
	rowContainer.AddChild(eui.NewCenteredLabel(title, bigFont))
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	panel := eui.NewPanel(c.state.UIResources, 0, 0)

	buttonsGrid := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(5),
			widget.GridLayoutOpts.Stretch([]bool{false, true, true, true, true}, nil),
			widget.GridLayoutOpts.Spacing(8, 8))))

	labels := []string{"I", "II", "III", "IV"}
	buttonsGrid.AddChild(eui.NewCenteredLabel(fmt.Sprintf("%s 1 ", d.Get("menu.play.act")), normalFont))
	for i := 0; i < 4; i++ {
		buttonsGrid.AddChild(eui.NewButton(c.state.UIResources, labels[i], func() {
			scene.Context().ChangeScene(NewStageController(c.state, stage.Config{
				MaxInstruments: 4,
				Targets: []stage.Target{
					{
						Instrument: synthdb.BassInstrument,
						Pos:        gmath.Vec{X: 1, Y: 1},
						Size:       stage.SmallTarget,
					},
					{
						Instrument: synthdb.BassInstrument,
						Pos:        gmath.Vec{X: 2, Y: 0.5},
						Size:       stage.SmallTarget,
					},
					{
						Instrument: synthdb.BassInstrument,
						Pos:        gmath.Vec{X: 4, Y: 1},
						Size:       stage.SmallTarget,
					},
				},
			}))
		}))
	}
	buttonsGrid.AddChild(eui.NewCenteredLabel(fmt.Sprintf("%s 2 ", d.Get("menu.play.act")), normalFont))
	for i := 0; i < 4; i++ {
		buttonsGrid.AddChild(eui.NewButton(c.state.UIResources, labels[i], func() {
		}))
	}

	panel.AddChild(buttonsGrid)
	rowContainer.AddChild(panel)

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		scene.Context().ChangeScene(NewPlayController(c.state))
	}))

	initUI(scene, root)
}

func (c *MissionsController) Update(delta float64) {}
