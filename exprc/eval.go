package exprc

import "math"

type funcRunner struct {
	stack     []float64
	constants []float64
	insts     []instructon
}

func (r *funcRunner) Run(x float64) float64 {
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

		default:
			panic("unexpected op")
		}
	}

	return r.stack[0]
}

func (r *funcRunner) push(v float64) {
	r.stack = append(r.stack, v)
}

func (r *funcRunner) pop3() (float64, float64, float64) {
	c := r.stack[len(r.stack)-1]
	b := r.stack[len(r.stack)-2]
	a := r.stack[len(r.stack)-3]
	r.stack = r.stack[:len(r.stack)-3]
	return a, b, c
}

func (r *funcRunner) pop2() (float64, float64) {
	b := r.stack[len(r.stack)-1]
	a := r.stack[len(r.stack)-2]
	r.stack = r.stack[:len(r.stack)-2]
	return a, b
}

func (r *funcRunner) pop() float64 {
	v := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	return v
}
