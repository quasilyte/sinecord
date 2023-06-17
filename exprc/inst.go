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
	opMinFunc
	opMaxFunc
	opClampFunc
	opPowFunc
	opTanFunc
	opAtanFunc
	opAtan2Func
	opAsinFunc
	opAcosFunc
	opLogFunc
	opLog2Func
	opSqrtFunc
	opInversesqrtFunc
	opSignFunc
	opFloorFunc
	opCeilFunc
	opFractFunc
	opModFunc

	opAdd
	opMul
	opSub
	opDiv
)
