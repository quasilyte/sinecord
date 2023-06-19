package scenes

import (
	"fmt"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/session"
	"github.com/quasilyte/sinecord/styles"
)

type SaverController struct {
	state *session.State

	track gamedata.Track

	existingTracks []gamedata.Track

	canSave           bool
	selectedSlot      int
	slotSelector      *ge.Rect
	selectedSlotLabel *widget.Text
	trackNameInput    *widget.TextInput
	saveButton        *widget.Button

	back ge.SceneController
}

func NewSaverController(state *session.State, track gamedata.Track, back ge.SceneController) *SaverController {
	return &SaverController{
		state: state,
		track: track,
		back:  back,
	}
}

func (c *SaverController) Init(scene *ge.Scene) {
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

	title := d.Get("menu.save_track")
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
		slotsGrid.AddChild(b)
	}

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	c.selectedSlotLabel = eui.NewCenteredLabel("--------------------", smallFont)
	rowContainer.AddChild(c.selectedSlotLabel)

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	nameInputAnchor := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
	nameContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Spacing(16, 0))))

	nameInputLabel := eui.NewCenteredLabel("enter track name", smallFont)
	c.trackNameInput = eui.NewTextInput(c.state.UIResources,
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(380, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.TextInputOpts.ChangedHandler(func(_ *widget.TextInputChangedEventArgs) {
			c.updateStatusLabel()
		}),
		widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
			good := true
			if len(newInputText) > 20 {
				good = false
			}
			if good {
				for _, ch := range newInputText {
					if !unicode.IsPrint(ch) || ch >= utf8.RuneSelf {
						good = false
						break
					}
				}
			}
			return good, nil
		}),
	)
	nameContainer.AddChild(nameInputLabel)
	nameContainer.AddChild(c.trackNameInput)
	nameInputAnchor.AddChild(nameContainer)
	rowContainer.AddChild(nameInputAnchor)

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, true}, nil),
			widget.GridLayoutOpts.Spacing(32, 0))))

	c.saveButton = eui.NewButton(c.state.UIResources, "save", func() {
		c.track.Slot = c.selectedSlot
		if c.trackNameInput.InputText == "" {
			c.track.Name = c.existingTracks[c.selectedSlot-1].Name
		} else {
			c.track.Name = c.trackNameInput.InputText
		}
		c.track.Date = time.Now()
		scene.Context().SaveGameData(c.track.FileName(), c.track)
		scene.Context().ChangeScene(c.back)
	})
	c.saveButton.GetWidget().Disabled = true

	backButton := eui.NewButton(c.state.UIResources, d.Get("menu.back"), func() {
		scene.Context().ChangeScene(c.back)
	})
	buttonsContainer.AddChild(c.saveButton)
	buttonsContainer.AddChild(backButton)

	rowContainer.AddChild(buttonsContainer)

	initUI(scene, root)

	c.updateStatusLabel()
}

func (c *SaverController) updateStatusLabel() {
	if c.selectedSlot == -1 {
		c.selectedSlotLabel.Label = "no save slot selected!"
		c.canSave = false
		return
	}

	selectedTrack := &c.existingTracks[c.selectedSlot-1]

	if c.trackNameInput.InputText == "" && selectedTrack.IsEmpty() {
		c.selectedSlotLabel.Label = "the track name is empty!"
		c.canSave = false
		return
	}

	var text string
	if selectedTrack.IsEmpty() {
		text = fmt.Sprintf("save to slot %d as %q", c.selectedSlot, c.trackNameInput.InputText)
	} else {
		if c.trackNameInput.InputText == "" {
			text = fmt.Sprintf("overwrite %q", selectedTrack.Name)
		} else {
			text = fmt.Sprintf("overwrite %q as %q", selectedTrack.Name, c.trackNameInput.InputText)
		}
	}
	c.selectedSlotLabel.Label = text
	c.canSave = true
}

func (c *SaverController) Update(delta float64) {
	c.saveButton.GetWidget().Disabled = !c.canSave
}
