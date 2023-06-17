package exprc

import "github.com/quasilyte/gmath"

type builtinFunction struct {
	numArgs int
	op      operation
}

var builtinFuncMap = map[string]builtinFunction{
	"abs":        {numArgs: 1, op: opAbsFunc},
	"sin":        {numArgs: 1, op: opSinFunc},
	"cos":        {numArgs: 1, op: opCosFunc},
	"step":       {numArgs: 2, op: opStepFunc},
	"smoothstep": {numArgs: 3, op: opSmootstepFunc},
}

func step(edge, x float64) float64 {
	if x < edge {
		return 0
	}
	return 1
}

func smoothstep(edge0, edge1, x float64) float64 {
	t := gmath.Clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)
	return t * t * (3.0 - 2.0*t)
}
