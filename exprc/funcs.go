package exprc

import (
	"math"

	"github.com/quasilyte/gmath"
)

type BuiltinFunction struct {
	Args []string
	Doc  string
	op   operation
}

var BuiltinFuncMap = map[string]BuiltinFunction{
	"abs":         {Args: []string{"x"}, op: opAbsFunc, Doc: "Get the absolute value of the parameter"},
	"sin":         {Args: []string{"x"}, op: opSinFunc, Doc: "Compute the sine of the parameter"},
	"cos":         {Args: []string{"x"}, op: opCosFunc, Doc: "Compute the cosine of the parameter"},
	"step":        {Args: []string{"edge", "x"}, op: opStepFunc, Doc: "Generate a step function by comparing two values"},
	"smoothstep":  {Args: []string{"edge1", "edge2", "x"}, op: opSmootstepFunc, Doc: "Perform Hermite interpolation between two values"},
	"min":         {Args: []string{"value1", "value2"}, op: opMinFunc, Doc: "Find the lesser of two values"},
	"max":         {Args: []string{"value1", "value2"}, op: opMaxFunc, Doc: "Find the greater of two values"},
	"clamp":       {Args: []string{"x", "min_value", "max_value"}, op: opClampFunc, Doc: "Constrain a value to lie between two further values"},
	"pow":         {Args: []string{"x", "exponent"}, op: opPowFunc, Doc: "Calculates x^exponent"},
	"tan":         {Args: []string{"x"}, op: opTanFunc, Doc: "Compute the tangent of the parameter"},
	"tanh":        {Args: []string{"x"}, op: opTanhFunc, Doc: "Compute the hyperbolic tangent of the parameter"},
	"atan":        {Args: []string{"x"}, op: opAtanFunc, Doc: "Compute the arc-tangent of y-over-x"},
	"asin":        {Args: []string{"x"}, op: opAsinFunc, Doc: "Compute the arc-sine of the parameter"},
	"acos":        {Args: []string{"x"}, op: opAcosFunc, Doc: "Compute the arc-cosine of the parameter"},
	"log":         {Args: []string{"x"}, op: opLogFunc, Doc: "Get the natural logarithm of the parameter"},
	"log2":        {Args: []string{"x"}, op: opLog2Func, Doc: "Get the base 2 logarithm of the parameter"},
	"sqrt":        {Args: []string{"x"}, op: opSqrtFunc, Doc: "Compute the square root of the parameter"},
	"inversesqrt": {Args: []string{"x"}, op: opInversesqrtFunc, Doc: "Compute the inverse of the square root of the parameter"},
	"sign":        {Args: []string{"x"}, op: opSignFunc, Doc: "Returns -1 if x<0, 0 if x=0, and +1 if x>0"},
	"floor":       {Args: []string{"x"}, op: opFloorFunc, Doc: "Get the nearest integer less than or equal to the parameter"},
	"ceil":        {Args: []string{"x"}, op: opCeilFunc, Doc: "Get the nearest integer that is greater than or equal to the parameter"},
	"fract":       {Args: []string{"x"}, op: opFractFunc, Doc: "Get the fractional part of x"},
	"mod":         {Args: []string{"x", "divisor"}, op: opModFunc, Doc: "Compute value of one parameter modulo another, like x%divisor"},
	"gamma":       {Args: []string{"x"}, op: opGammaFunc, Doc: "Compute the Gamma function of x"},
	"until":       {Args: []string{"x", "threshold"}, op: opUntilFunc, Doc: "Returns x if x<=threshold"},
	"after":       {Args: []string{"x", "threshold"}, op: opAfterFunc, Doc: "Returns x if x>=threshold"},
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
