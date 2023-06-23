package scenes

import (
	"fmt"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/stage"
	"github.com/quasilyte/sinecord/styles"
)

type MissionViewController struct {
	state *session.State

	levelData *gamedata.LevelData
}

func NewMissionViewController(state *session.State, levelData *gamedata.LevelData) *MissionViewController {
	return &MissionViewController{
		state:     state,
		levelData: levelData,
	}
}

func (c *MissionViewController) Init(scene *ge.Scene) {
	bigFont := scene.Context().Loader.LoadFont(assets.FontArcadeBig).Face
	smallFont := scene.Context().Loader.LoadFont(assets.FontArcadeSmall).Face

	monoFont := scene.Context().Loader.LoadFont(assets.FontMonospaceNormal).Face

	d := scene.Dict()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(1280, 10, nil)
	root.AddChild(rowContainer)

	title := fmt.Sprintf("%s %d - %s %d", d.Get("menu.play.act"), c.levelData.ActNumber, d.Get("menu.mission.title"), c.levelData.MissionNumber)
	rowContainer.AddChild(eui.NewCenteredLabel(title, bigFont))
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	infoPanel := eui.NewPanelWithPadding(c.state.UIResources, 0, 0, widget.NewInsetsSimple(24))
	rowContainer.AddChild(infoPanel)

	infoTextWidget := eui.NewLabel(c.levelData.Description, monoFont, widget.TextOpts.MaxWidth(1200))
	infoPanelRows := c.panelRowsContainer()
	infoPanelRows.AddChild(eui.NewCenteredLabel("Mission overview", smallFont))
	infoPanelRows.AddChild(infoTextWidget)
	infoPanel.AddChild(infoPanelRows)

	bonusPanel := eui.NewPanelWithPadding(c.state.UIResources, 0, 0, widget.NewInsetsSimple(24))
	bonusPanelRows := c.panelRowsContainer()
	rowContainer.AddChild(bonusPanel)
	bonusPanel.AddChild(bonusPanelRows)

	bonusPanelRows.AddChild(eui.NewCenteredLabel("Bonus conditions", smallFont))
	bonusPanelRows.AddChild(eui.NewLabel(c.formatBonusRequirements(), monoFont))

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	buttonsGrid := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true}, nil),
			widget.GridLayoutOpts.Spacing(32, 8))))

	rowContainer.AddChild(buttonsGrid)

	stageConfig := stage.Config{
		Data:           c.levelData,
		MaxInstruments: c.levelData.MaxInstruments,
		Targets:        c.levelData.Targets,
		Mode:           gamedata.MissionMode,
	}

	buttonsGrid.AddChild(eui.NewButtonWithConfig(c.state.UIResources, eui.ButtonConfig{
		Text: d.Get("menu.mission.start"),
		OnClick: func() {
			scene.Context().ChangeScene(NewStageController(c.state, stageConfig))
		},
	}))

	showSolution := eui.NewButtonWithConfig(c.state.UIResources, eui.ButtonConfig{
		Text:         d.Get("menu.mission.show_solution"),
		TooltipLabel: "see one of the intended solutions (satisfy bonus conditions to unlock it)",
		OnClick: func() {
			stageConfig.Track = c.levelData.Solution
			scene.Context().ChangeScene(NewStageController(c.state, stageConfig))
		},
	})
	// showSolution.GetWidget().Disabled = c.state.Persistent.GetLevelCompletionStatus(c.levelData) < session.LevelCompletedWithBonus
	buttonsGrid.AddChild(showSolution)

	buttonsGrid.AddChild(eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		scene.Context().ChangeScene(NewMissionsController(c.state))
	}))

	initUI(scene, root)
}

func (c *MissionViewController) formatBonusRequirements() string {
	var rules []string

	objectives := c.levelData.Bonus

	pluralSuffix := ""
	if objectives.MaxInstruments != 1 {
		pluralSuffix = "s"
	}
	rules = append(rules, fmt.Sprintf("* Use no more than %d instrument%s", objectives.MaxInstruments, pluralSuffix))

	if len(objectives.ForbiddenFuncs) != 0 {
		rules = append(rules, fmt.Sprintf("* Don't use any of these functions: %s", strings.Join(objectives.ForbiddenFuncs, ", ")))
	}

	if objectives.AvoidPenalty {
		rules = append(rules, "* Avoid misplays")
	}

	if objectives.AllTargets {
		rules = append(rules, "* Mark all optional targets too")
	}

	return strings.Join(rules, "\n")
}

func (c *MissionViewController) panelRowsContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
			StretchVertical:   true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, nil),
			widget.GridLayoutOpts.Spacing(32, 32),
		)),
	)
}

func (c *MissionViewController) Update(delta float64) {}
