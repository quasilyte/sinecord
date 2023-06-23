package styles

import (
	"image/color"

	"github.com/quasilyte/ge"
)

var (
	TransparentColor = color.RGBA{}

	TargetColor          = ge.RGB(0x9dd793)
	TargetMissColorRed   = ge.RGB(0xff2948)
	TargetMissColorGreen = ge.RGB(0x03ff20)
	TargetMissColorBlue  = ge.RGB(0x0326ff)
	TargetColorBonus     = ge.RGB(0xd7d7d7)

	NormalTextColor   = ge.RGB(0x9dd793)
	DisabledTextColor = ge.RGB(0x5a7a91)
	CaretColor        = ge.RGB(0xfed846)

	CompletedLevelTextColor = CaretColor

	SeparatorColor = ge.RGB(0x79badc)

	SoundWaveColor        = SeparatorColor
	VictorySoundWaveColor = NormalTextColor

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
