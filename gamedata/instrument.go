package gamedata

type InstrumentKind int

const (
	BassInstrument InstrumentKind = iota
	DrumInstrument
	KeyboardInstrument
	BrassInstrument
	StringInstrument
	OtherInstrument
	AnyInstrument
)

//go:generate stringer -type=Shape -trimprefix=Shape
type Shape int

const (
	ShapeCircle Shape = iota
	ShapeSquare
	ShapeTriangle
	ShapeHexagon
	ShapeStar
	ShapeCross
	ShapeBonus
)

func InstrumentShape(kind InstrumentKind) Shape {
	switch kind {
	case BassInstrument:
		return ShapeStar
	case KeyboardInstrument:
		return ShapeSquare
	case BrassInstrument:
		return ShapeTriangle
	case StringInstrument:
		return ShapeHexagon
	case DrumInstrument:
		return ShapeCircle
	case AnyInstrument:
		return ShapeBonus
	default:
		return ShapeCross
	}
}
