package scenes

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/controls"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/styles"
)

type MissionsController struct {
	state *session.State

	scene *ge.Scene
}

func NewMissionsController(state *session.State) *MissionsController {
	return &MissionsController{
		state: state,
	}
}

func (c *MissionsController) Init(scene *ge.Scene) {
	c.scene = scene

	bigFont := scene.Context().Loader.LoadFont(assets.FontArcadeBig).Face
	normalFont := scene.Context().Loader.LoadFont(assets.FontArcadeNormal).Face

	d := scene.Dict()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(640, 10, nil)
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

	missionLabels := []string{"", "I", "II", "III", "IV"}
	for actNumber, levels := range c.state.LevelsByAct {
		if actNumber == 0 {
			continue
		}
		buttonsGrid.AddChild(eui.NewCenteredLabel(fmt.Sprintf("%s %d ", d.Get("menu.play.act"), actNumber), normalFont))
		for i := range levels {
			l := levels[i]
			b := eui.NewButtonWithConfig(c.state.UIResources, eui.ButtonConfig{
				TextAltColor: c.state.Persistent.GetLevelCompletionStatus(l) >= session.LevelCompleted,
				Text:         missionLabels[l.MissionNumber],
				OnClick: func() {
					scene.Context().ChangeScene(NewMissionViewController(c.state, l))
				},
			})
			buttonsGrid.AddChild(b)
		}
		for i := len(levels) + 1; i < len(missionLabels); i++ {
			b := eui.NewButton(c.state.UIResources, "    ", func() {})
			buttonsGrid.AddChild(b)
			b.GetWidget().Disabled = true
		}
	}

	panel.AddChild(buttonsGrid)
	rowContainer.AddChild(panel)

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		c.back()
	}))

	initUI(scene, root)
}

func (c *MissionsController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.back()
	}
}

func (c *MissionsController) back() {
	c.scene.Context().ChangeScene(NewPlayController(c.state))
}
