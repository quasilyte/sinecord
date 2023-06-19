package styles

import (
	"image/color"

	"github.com/quasilyte/ge"
)

var (
	TransparentColor = color.RGBA{}

	TargetColor = ge.RGB(0x9dd793)

	NormalTextColor   = ge.RGB(0x9dd793)
	DisabledTextColor = ge.RGB(0x5a7a91)
	CaretColor        = SeparatorColor

	SeparatorColor = ge.RGB(0x79badc)

	SoundWaveColor = SeparatorColor

	Plot1Color = ge.RGB(0x5987f2)
	Plot2Color = ge.RGB(0xad92d6)
	Plot3Color = ge.RGB(0xd7c781)
	Plot4Color = ge.RGB(0xeb1683)
	Plot5Color = ge.RGB(0x0ac36f)
)

var PlotColorByID = [...]color.RGBA{
	Plot1Color,
	Plot2Color,
	Plot3Color,
	Plot4Color,
	Plot5Color,
}
