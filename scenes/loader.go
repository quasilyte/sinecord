package scenes

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/styles"
)

type LoaderController struct {
	state *session.State

	existingTracks []gamedata.Track

	canLoad           bool
	selectedSlot      int
	slotSelector      *ge.Rect
	selectedSlotLabel *widget.Text
	loadButton        *widget.Button

	back ge.SceneController

	EventLoaded gsignal.Event[gamedata.Track]
}

func NewLoaderController(state *session.State, back ge.SceneController) *LoaderController {
	return &LoaderController{
		state: state,
		back:  back,
	}
}

func (c *LoaderController) Init(scene *ge.Scene) {
	c.selectedSlot = -1

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

	title := d.Get("menu.load_track")
	rowContainer.AddChild(eui.NewCenteredLabel(title, bigFont))
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	slotsGrid := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch(nil, nil),
			widget.GridLayoutOpts.Spacing(16, 24))))
	slotsGridAnchor := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
	slotsGridAnchor.AddChild(slotsGrid)
	rowContainer.AddChild(slotsGridAnchor)

	c.existingTracks = gamedata.DiscoverTracks(scene.Context())

	for i := range c.existingTracks {
		t := c.existingTracks[i]
		var label string
		if t.IsEmpty() {
			label = fmt.Sprintf("%d. empty slot", t.Slot)
		} else {
			label = fmt.Sprintf("%d. [%s] %s", t.Slot, formatDateISO8601(t.Date, true), t.Name)
		}
		var b *widget.Button
		b = eui.NewButtonWithConfig(c.state.UIResources, eui.ButtonConfig{
			Text:          label,
			TextAlignLeft: true,
			OnClick: func() {
				rect := b.GetWidget().Rect
				if c.slotSelector == nil {
					width := float64(rect.Dx()) + 12
					height := float64(rect.Dy()) + 12
					c.slotSelector = ge.NewRect(scene.Context(), width, height)
					c.slotSelector.Centered = false
					c.slotSelector.FillColorScale.SetRGBA(0, 0, 0, 0)
					c.slotSelector.OutlineWidth = 2
					c.slotSelector.OutlineColorScale.SetColor(styles.NormalTextColor)
					scene.AddGraphicsAbove(c.slotSelector, 1)
				}
				c.slotSelector.Pos.Offset = gmath.Vec{
					X: float64(rect.Min.X) - 6,
					Y: float64(rect.Min.Y) - 6,
				}
				c.selectedSlot = t.Slot
				c.updateStatusLabel()
			},
			MinWidth: 720,
			Font:     monoFont,
		})
		b.GetWidget().Disabled = t.IsEmpty()
		slotsGrid.AddChild(b)
	}

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	c.selectedSlotLabel = eui.NewCenteredLabel("--------------------", smallFont)
	rowContainer.AddChild(c.selectedSlotLabel)

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, true}, nil),
			widget.GridLayoutOpts.Spacing(32, 0))))

	c.loadButton = eui.NewButton(c.state.UIResources, "load", func() {
		track := c.existingTracks[c.selectedSlot-1]
		c.EventLoaded.Emit(track)
		scene.Context().ChangeScene(c.back)
	})
	c.loadButton.GetWidget().Disabled = true

	backButton := eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		scene.Context().ChangeScene(c.back)
	})
	buttonsContainer.AddChild(c.loadButton)
	buttonsContainer.AddChild(backButton)

	rowContainer.AddChild(buttonsContainer)

	initUI(scene, root)

	c.updateStatusLabel()
}

func (c *LoaderController) updateStatusLabel() {
	if c.selectedSlot == -1 {
		c.selectedSlotLabel.Label = "no slot selected!"
		c.canLoad = false
		return
	}

	c.selectedSlotLabel.Label = fmt.Sprintf("load %q", c.existingTracks[c.selectedSlot-1].Name)
	c.canLoad = true
}

func (c *LoaderController) Update(delta float64) {
	c.loadButton.GetWidget().Disabled = !c.canLoad
}
