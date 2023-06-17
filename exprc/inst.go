package exprc

type function struct {
	constants []float64
}

type instructon struct {
	op  operation
	arg uint8
}

type operation uint8

const (
	opUnknown operation = iota

	// $arg - const index
	opFloatConst

	opArg

	opNeg

	opAbsFunc
	opSinFunc
	opCosFunc
	opStepFunc
	opSmootstepFunc

	opAdd
	opMul
	opSub
	opDiv
)
