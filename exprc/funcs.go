package exprc

import (
	"math"

	"github.com/quasilyte/gmath"
)

type builtinFunction struct {
	numArgs int
	op      operation
}

var builtinFuncMap = map[string]builtinFunction{
	"abs":         {numArgs: 1, op: opAbsFunc},
	"sin":         {numArgs: 1, op: opSinFunc},
	"cos":         {numArgs: 1, op: opCosFunc},
	"step":        {numArgs: 2, op: opStepFunc},
	"smoothstep":  {numArgs: 3, op: opSmootstepFunc},
	"min":         {numArgs: 2, op: opMinFunc},
	"max":         {numArgs: 2, op: opMaxFunc},
	"clamp":       {numArgs: 3, op: opClampFunc},
	"pow":         {numArgs: 2, op: opPowFunc},
	"tan":         {numArgs: 1, op: opTanFunc},
	"tanh":        {numArgs: 1, op: opTanhFunc},
	"atan":        {numArgs: 1, op: opAtanFunc},
	"atan2":       {numArgs: 2, op: opAtan2Func},
	"asin":        {numArgs: 1, op: opAsinFunc},
	"acos":        {numArgs: 1, op: opAcosFunc},
	"log":         {numArgs: 1, op: opLogFunc},
	"log2":        {numArgs: 1, op: opLog2Func},
	"sqrt":        {numArgs: 1, op: opSqrtFunc},
	"inversesqrt": {numArgs: 1, op: opInversesqrtFunc},
	"sign":        {numArgs: 1, op: opSignFunc},
	"floor":       {numArgs: 1, op: opFloorFunc},
	"ceil":        {numArgs: 1, op: opCeilFunc},
	"fract":       {numArgs: 1, op: opFractFunc},
	"mod":         {numArgs: 2, op: opModFunc},
	"until":       {numArgs: 2, op: opUntilFunc},
	"after":       {numArgs: 2, op: opAfterFunc},
}

func step(edge, x float64) float64 {
	if x < edge {
		return 0
	}
	return 1
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func smoothstep(edge0, edge1, x float64) float64 {
	t := gmath.Clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)
	return t * t * (3.0 - 2.0*t)
}

func fract(x float64) float64 {
	return x - math.Floor(x)
}

func mod(x, y float64) float64 {
	return x - y*math.Floor(x/y)
}

func until(x, v, threshold float64) float64 {
	if x+gmath.Epsilon <= threshold {
		return v
	}
	return -10
}

func after(x, v, threshold float64) float64 {
	if x+gmath.Epsilon >= threshold {
		return v
	}
	return -10
}

func sign(x float64) float64 {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func inversesqrt(x float64) float64 {
	const magic64 = 0x5FE6EB50C7B537A9
	if x < 0 {
		return math.NaN()
	}
	n2, th := x*0.5, float64(1.5)
	b := math.Float64bits(x)
	b = magic64 - (b >> 1)
	f := math.Float64frombits(b)
	f *= th - (n2 * f * f)
	return f
}
