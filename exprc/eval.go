package exprc

import (
	"math"

	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type FuncRunner struct {
	stack     []float64
	constants []float64
	insts     []instructon
	funcsUsed []string
}

func (r *FuncRunner) UsesFunc(name string) bool {
	return xslices.Contains(r.funcsUsed, name)
}

func (r *FuncRunner) Run(x float64) float64 {
	r.stack = r.stack[:0]

	for _, inst := range r.insts {
		switch inst.op {
		case opFloatConst:
			r.push(r.constants[inst.arg])
		case opNeg:
			r.push(-r.pop())
		case opArg:
			r.push(x)

		case opAdd:
			a, b := r.pop2()
			r.push(a + b)
		case opMul:
			a, b := r.pop2()
			r.push(a * b)
		case opSub:
			a, b := r.pop2()
			r.push(a - b)
		case opDiv:
			a, b := r.pop2()
			r.push(a / b)

		case opAbsFunc:
			r.push(math.Abs(r.pop()))
		case opSinFunc:
			r.push(math.Sin(r.pop()))
		case opCosFunc:
			r.push(math.Cos(r.pop()))
		case opStepFunc:
			r.push(step(r.pop2()))
		case opSmootstepFunc:
			r.push(smoothstep(r.pop3()))
		case opMinFunc:
			r.push(min(r.pop2()))
		case opMaxFunc:
			r.push(max(r.pop2()))
		case opClampFunc:
			r.push(gmath.Clamp(r.pop3()))
		case opPowFunc:
			r.push(math.Pow(r.pop2()))
		case opTanFunc:
			r.push(math.Tan(r.pop()))
		case opTanhFunc:
			r.push(math.Tanh(r.pop()))
		case opAtanFunc:
			r.push(math.Atan(r.pop()))
		case opAsinFunc:
			r.push(math.Asin(r.pop()))
		case opAcosFunc:
			r.push(math.Acos(r.pop()))
		case opLogFunc:
			r.push(math.Log(r.pop()))
		case opLog2Func:
			r.push(math.Log2(r.pop()))
		case opSqrtFunc:
			r.push(math.Sqrt(r.pop()))
		case opInversesqrtFunc:
			r.push(inversesqrt(r.pop()))
		case opSignFunc:
			r.push(sign(r.pop()))
		case opFloorFunc:
			r.push(math.Floor(r.pop()))
		case opCeilFunc:
			r.push(math.Ceil(r.pop()))
		case opFractFunc:
			r.push(fract(r.pop()))
		case opModFunc:
			r.push(mod(r.pop2()))
		case opGammaFunc:
			r.push(math.Gamma(r.pop()))
		case opUntilFunc:
			v, threshold := r.pop2()
			r.push(until(x, v, threshold))
		case opAfterFunc:
			v, threshold := r.pop2()
			r.push(after(x, v, threshold))

		default:
			panic("unexpected op")
		}
	}

	return r.stack[0]
}

func (r *FuncRunner) push(v float64) {
	r.stack = append(r.stack, v)
}

func (r *FuncRunner) pop3() (float64, float64, float64) {
	c := r.stack[len(r.stack)-1]
	b := r.stack[len(r.stack)-2]
	a := r.stack[len(r.stack)-3]
	r.stack = r.stack[:len(r.stack)-3]
	return a, b, c
}

func (r *FuncRunner) pop2() (float64, float64) {
	b := r.stack[len(r.stack)-1]
	a := r.stack[len(r.stack)-2]
	r.stack = r.stack[:len(r.stack)-2]
	return a, b
}

func (r *FuncRunner) pop() float64 {
	v := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	return v
}
