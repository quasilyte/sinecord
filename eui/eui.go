package eui

import (
	"image/color"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/styles"
	"golang.org/x/image/font"
)

type Resources struct {
	button       *buttonResource
	selectButton *buttonResource
	textInput    *textInputResource
	panel        *panelResource
	tooltip      *tooltipResources
}

type tooltipResources struct {
	Background *image.NineSlice
	Padding    widget.Insets
	FontFace   font.Face
	TextColor  color.Color
}

type buttonResource struct {
	Image         *widget.ButtonImage
	Padding       widget.Insets
	TextColors    *widget.ButtonTextColor
	AltTextColors *widget.ButtonTextColor
	FontFace      font.Face
}

type textInputResource struct {
	Image      *widget.TextInputImage
	Padding    widget.Insets
	TextColors *widget.TextInputColor
	FontFace   font.Face
}

type panelResource struct {
	Image   *image.NineSlice
	Padding widget.Insets
}

func PrepareResources(loader *resource.Loader) *Resources {
	result := &Resources{}

	smallFont := loader.LoadFont(assets.FontArcadeSmall).Face
	mediumFont := loader.LoadFont(assets.FontArcadeNormal).Face
	monoMediumFont := loader.LoadFont(assets.FontMonospaceNormal).Face

	{
		disabled := nineSliceImage(loader.LoadImage(assets.ImageUIButtonDisabled).Data, 12, 0)
		idle := nineSliceImage(loader.LoadImage(assets.ImageUIButtonIdle).Data, 12, 0)
		hover := nineSliceImage(loader.LoadImage(assets.ImageUIButtonHover).Data, 12, 0)
		pressed := nineSliceImage(loader.LoadImage(assets.ImageUIButtonPressed).Data, 12, 0)
		selectIdle := nineSliceImage(loader.LoadImage(assets.ImageUISelectButtonIdle).Data, 12, 0)
		selectHover := nineSliceImage(loader.LoadImage(assets.ImageUISelectButtonHover).Data, 12, 0)
		selectPressed := nineSliceImage(loader.LoadImage(assets.ImageUISelectButtonPressed).Data, 12, 0)
		buttonPadding := widget.Insets{
			Left:  30,
			Right: 30,
		}
		buttonColors := &widget.ButtonTextColor{
			Idle:     styles.NormalTextColor,
			Disabled: styles.DisabledTextColor,
		}
		result.button = &buttonResource{
			Image: &widget.ButtonImage{
				Idle:     idle,
				Hover:    hover,
				Pressed:  pressed,
				Disabled: disabled,
			},
			Padding:    buttonPadding,
			TextColors: buttonColors,
			AltTextColors: &widget.ButtonTextColor{
				Idle:     styles.CompletedLevelTextColor,
				Disabled: styles.DisabledTextColor,
			},
			FontFace: mediumFont,
		}
		result.selectButton = &buttonResource{
			Image: &widget.ButtonImage{
				Idle:    selectIdle,
				Hover:   selectHover,
				Pressed: selectPressed,
			},
			Padding:    buttonPadding,
			TextColors: buttonColors,
			FontFace:   monoMediumFont,
		}
	}

	{
		idle := loader.LoadImage(assets.ImageUITextInputIdle).Data
		result.textInput = &textInputResource{
			Image: &widget.TextInputImage{
				Idle: nineSliceImage(idle, 12, 0),
			},
			Padding: widget.Insets{
				Left:   14,
				Right:  14,
				Top:    20,
				Bottom: 20,
			},
			FontFace: monoMediumFont,
			TextColors: &widget.TextInputColor{
				Idle:  styles.NormalTextColor,
				Caret: styles.CaretColor,
			},
		}
	}

	{
		idle := loader.LoadImage(assets.ImageUIPanelIdle).Data
		result.panel = &panelResource{
			Image: nineSliceImage(idle, 10, 10),
			Padding: widget.Insets{
				Left:   16,
				Right:  16,
				Top:    10,
				Bottom: 10,
			},
		}
	}

	{
		result.tooltip = &tooltipResources{
			Background: nineSliceImage(loader.LoadImage(assets.ImageUITooltip).Data, 10, 10),
			Padding: widget.Insets{
				Left:   16,
				Right:  16,
				Top:    10,
				Bottom: 10,
			},
			FontFace:  smallFont,
			TextColor: styles.NormalTextColor,
		}
	}

	return result
}

func nineSliceImage(i *ebiten.Image, centerWidth, centerHeight int) *image.NineSlice {
	w, h := i.Size()
	return image.NewNineSlice(i,
		[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
		[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight})
}

func NewRowLayoutContainerWithMinWidth(minWidth, spacing int, rowscale []bool) *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
			StretchVertical:   true,
		})),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(minWidth, 0)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, rowscale),
			widget.GridLayoutOpts.Spacing(spacing, spacing),
		)),
	)
}

type ImageButton struct {
	GraphicWidget *widget.Graphic
	Widget        widget.PreferredSizeLocateableWidget
}

func NewImageButton(res *Resources, img *ebiten.Image, config ButtonConfig) ImageButton {
	result := ImageButton{}
	result.GraphicWidget = widget.NewGraphic(
		widget.GraphicOpts.Image(img),
	)

	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)
	b := NewButtonWithConfig(res, config)
	container.AddChild(b)
	container.AddChild(result.GraphicWidget)
	result.Widget = container
	return result
}

type ButtonConfig struct {
	Text          string
	TextAlignLeft bool
	TextAltColor  bool
	OnClick       func()
	TooltipLabel  string
	LayoutData    any
	MinWidth      int
	Font          font.Face
}

func NewButtonWithConfig(res *Resources, config ButtonConfig) *widget.Button {
	ff := config.Font
	if ff == nil {
		ff = res.button.FontFace
	}
	options := []widget.ButtonOpt{
		widget.ButtonOpts.Image(res.button.Image),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if config.OnClick != nil {
				config.OnClick()
			}
		}),
	}
	colors := res.button.TextColors
	if config.TextAltColor {
		colors = res.button.AltTextColors
	}
	if config.TextAlignLeft {
		options = append(options,
			widget.ButtonOpts.TextSimpleLeft(config.Text, ff, colors, res.button.Padding))
	} else {
		options = append(options,
			widget.ButtonOpts.Text(config.Text, ff, colors),
			widget.ButtonOpts.TextPadding(res.button.Padding))
	}
	if config.LayoutData != nil {
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(config.LayoutData)))
	}
	if config.MinWidth != 0 {
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(config.MinWidth, 0)))
	}
	if config.TooltipLabel != "" {
		tooltipContent := NewTooltip(res, config.TooltipLabel)
		tt := widget.NewToolTip(
			widget.ToolTipOpts.Content(tooltipContent),
			widget.ToolTipOpts.Delay(time.Second),
		)
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.ToolTip(tt)))
	}
	return widget.NewButton(options...)
}

func NewButton(res *Resources, text string, onclick func()) *widget.Button {
	return NewButtonWithConfig(res, ButtonConfig{
		Text:    text,
		OnClick: onclick,
	})
}

type BoolSelectButtonConfig struct {
	Resources *Resources

	Value      *bool
	Label      string
	ValueNames []string

	MinWidth int

	Tooltip *widget.Container

	OnPressed func()
	OnHover   func()
}

func NewBoolSelectButton(config BoolSelectButtonConfig) *widget.Button {
	value := config.Value
	key := config.Label
	valueNames := config.ValueNames

	var slider gmath.Slider
	slider.SetBounds(0, 1)
	if *value {
		slider.TrySetValue(1)
	}
	makeLabel := func() string {
		if key == "" {
			return valueNames[slider.Value()]
		}
		return key + ": " + valueNames[slider.Value()]
	}

	var buttonOptions []widget.ButtonOpt
	if config.Tooltip != nil {
		tt := widget.NewToolTip(
			widget.ToolTipOpts.Content(config.Tooltip),
			widget.ToolTipOpts.Delay(time.Second),
		)
		buttonOptions = append(buttonOptions, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.ToolTip(tt)))
	}
	if config.MinWidth != 0 {
		buttonOptions = append(buttonOptions, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(config.MinWidth, 0)))
	}
	button := newSelectButton(config.Resources, makeLabel(), buttonOptions...)

	button.ClickedEvent.AddHandler(func(args interface{}) {
		slider.Inc()
		*value = slider.Value() == 1

		button.Text().Label = makeLabel()
		if config.OnPressed != nil {
			config.OnPressed()
		}
	})

	if config.OnHover != nil {
		button.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			config.OnHover()
		})
	}

	return button
}

type SelectButtonConfig struct {
	Resources *Resources
	Input     *input.Handler

	Value          *int
	Label          string
	ValueNames     []string
	DisabledValues []int

	MinWidth int

	Tooltip *widget.Container

	OnPressed func()
	OnHover   func()
}

func NewSelectButton(config SelectButtonConfig) *widget.Button {
	maxValue := len(config.ValueNames) - 1
	value := config.Value
	key := config.Label
	valueNames := config.ValueNames

	var slider gmath.Slider
	slider.SetBounds(0, maxValue)
	slider.TrySetValue(*value)
	makeLabel := func() string {
		if key == "" {
			return valueNames[slider.Value()]
		}
		return key + ": " + valueNames[slider.Value()]
	}

	var buttonOptions []widget.ButtonOpt
	if config.Tooltip != nil {
		tt := widget.NewToolTip(
			widget.ToolTipOpts.Content(config.Tooltip),
			widget.ToolTipOpts.Delay(time.Second),
		)
		buttonOptions = append(buttonOptions, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.ToolTip(tt)))
	}
	if config.MinWidth != 0 {
		buttonOptions = append(buttonOptions, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(config.MinWidth, 0)))
	}
	button := newSelectButton(config.Resources, makeLabel(), buttonOptions...)

	button.ClickedEvent.AddHandler(func(args interface{}) {
		increase := false
		{
			cursorPos := config.Input.CursorPos()
			buttonRect := button.GetWidget().Rect
			buttonWidth := buttonRect.Dx()
			if cursorPos.X >= float64(buttonRect.Min.X)+float64(buttonWidth)*0.5 {
				increase = true
			}
		}

		for {
			if increase {
				slider.Inc()
			} else {
				slider.Dec()
			}
			*value = slider.Value()
			if !xslices.Contains(config.DisabledValues, *value) {
				break
			}
		}

		button.Text().Label = makeLabel()
		if config.OnPressed != nil {
			config.OnPressed()
		}
	})

	if config.OnHover != nil {
		button.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			config.OnHover()
		})
	}

	return button
}

func newSelectButton(res *Resources, text string, opts ...widget.ButtonOpt) *widget.Button {
	options := []widget.ButtonOpt{
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(res.selectButton.Image),
		widget.ButtonOpts.Text(text, res.selectButton.FontFace, res.selectButton.TextColors),
		widget.ButtonOpts.TextPadding(res.selectButton.Padding),
	}
	options = append(options, opts...)
	return widget.NewButton(options...)
}

func NewColoredLabel(text string, ff font.Face, clr color.RGBA, options ...widget.TextOpt) *widget.Text {
	opts := []widget.TextOpt{
		widget.TextOpts.Text(text, ff, clr),
	}
	if len(options) != 0 {
		opts = append(opts, options...)
	}
	return widget.NewText(opts...)
}

func NewLabel(text string, ff font.Face, options ...widget.TextOpt) *widget.Text {
	return NewColoredLabel(text, ff, styles.NormalTextColor, options...)
}

func NewCenteredLabelWithMaxWidth(text string, ff font.Face, width float64) *widget.Text {
	options := []widget.TextOpt{
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.Text(text, ff, styles.NormalTextColor),
	}
	if width != -1 {
		options = append(options, widget.TextOpts.MaxWidth(width))
	}
	return widget.NewText(options...)
}

func NewCenteredLabel(text string, ff font.Face) *widget.Text {
	return NewCenteredLabelWithMaxWidth(text, ff, -1)
}

func NewSeparator(ld interface{}, clr color.RGBA) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.ImageNineSlice(image.NewNineSliceColor(clr)),
	))

	return c
}

func NewTooltip(res *Resources, text string) *widget.Container {
	tt := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.tooltip.Background),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.tooltip.Padding),
			widget.RowLayoutOpts.Spacing(2),
		)))
	tt.AddChild(widget.NewText(
		widget.TextOpts.MaxWidth(800),
		widget.TextOpts.Text(text, res.tooltip.FontFace, res.tooltip.TextColor),
	))
	return tt
}

type FunctionInputConfig struct {
	MinWidth      int
	TooltipLabel  string
	MaxTextLength int
	OnChange      func(s string)
}

func NewFunctionInput(res *Resources, config FunctionInputConfig) *widget.TextInput {
	return NewTextInput(res,
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(config.MinWidth, 0),
			widget.WidgetOpts.ToolTip(
				widget.NewToolTip(
					widget.ToolTipOpts.Content(NewTooltip(res, config.TooltipLabel)),
					widget.ToolTipOpts.Delay(time.Second),
				),
			),
		),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			if config.OnChange != nil {
				config.OnChange(args.InputText)
			}
		}),
		widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
			good := true
			if len(newInputText) > config.MaxTextLength {
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
}

func NewTextInput(res *Resources, opts ...widget.TextInputOpt) *widget.TextInput {
	options := []widget.TextInputOpt{
		widget.TextInputOpts.Image(res.textInput.Image),
		widget.TextInputOpts.Color(res.textInput.TextColors),
		widget.TextInputOpts.Padding(res.textInput.Padding),
		widget.TextInputOpts.Face(res.textInput.FontFace),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(res.textInput.FontFace, 2),
		),
		widget.TextInputOpts.AllowDuplicateSubmit(true),
	}
	options = append(options, opts...)
	t := widget.NewTextInput(options...)
	return t
}

func NewPanelWithPadding(res *Resources, minWidth, minHeight int, padding widget.Insets) *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.Image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(padding),
		)),
		// widget.ContainerOpts.Layout(widget.NewRowLayout(
		// 	widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		// 	widget.RowLayoutOpts.Spacing(4),
		// 	widget.RowLayoutOpts.Padding(res.panel.Padding),
		// )),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(minWidth, minHeight),
		),
	)
}

func NewPanel(res *Resources, minWidth, minHeight int) *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.Image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(res.panel.Padding),
		)),
		// widget.ContainerOpts.Layout(widget.NewRowLayout(
		// 	widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		// 	widget.RowLayoutOpts.Spacing(4),
		// 	widget.RowLayoutOpts.Padding(res.panel.Padding),
		// )),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(minWidth, minHeight),
		),
	)
}
