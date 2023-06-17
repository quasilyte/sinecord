package exprc

import (
	"math"
	"testing"
)

func TestSimple(t *testing.T) {
	type functionRun struct {
		arg    float64
		result float64
	}

	tests := []struct {
		src  string
		runs []functionRun
	}{
		// Basic literals.
		{
			src: "1",
			runs: []functionRun{
				{arg: 0, result: 1},
				{arg: -1, result: 1},
			},
		},
		{
			src: "-1.5",
			runs: []functionRun{
				{arg: 0, result: -1.5},
				{arg: -1, result: -1.5},
			},
		},

		// Simple variable references.
		{
			src: "x",
			runs: []functionRun{
				{arg: 1, result: 1},
				{arg: -1.6, result: -1.6},
			},
		},

		// Binary expressions.
		{
			src: "x+x",
			runs: []functionRun{
				{arg: 1, result: 2},
				{arg: -1.5, result: -3},
			},
		},
		{
			src: "x-(x+2)",
			runs: []functionRun{
				{arg: 1, result: -2},
				{arg: -1.5, result: -2},
			},
		},
		{
			src: "2-x",
			runs: []functionRun{
				{arg: 2, result: 0},
				{arg: 3.5, result: -1.5},
			},
		},
		{
			src: "x*x",
			runs: []functionRun{
				{arg: 1, result: 1},
				{arg: -4, result: 16},
			},
		},
		{
			src: "x*2.5",
			runs: []functionRun{
				{arg: 1, result: 2.5},
				{arg: -4, result: (-4 * 2.5)},
			},
		},
		{
			src: "x/3",
			runs: []functionRun{
				{arg: 1, result: 1.0 / 3.0},
				{arg: 4, result: 4.0 / 3.0},
			},
		},

		// Functions.
		{
			src: "sin(x)",
			runs: []functionRun{
				{arg: 2, result: math.Sin(2)},
				{arg: -1.5, result: math.Sin(-1.5)},
			},
		},
		{
			src: "cos(x)",
			runs: []functionRun{
				{arg: 2, result: math.Cos(2)},
				{arg: -1.5, result: math.Cos(-1.5)},
			},
		},
	}

	for _, test := range tests {
		f, err := Compile(test.src)
		if err != nil {
			t.Fatalf("%q: %v", test.src, err)
		}
		for _, r := range test.runs {
			result := f(r.arg)
			if result != r.result {
				t.Fatalf("%q:\nf(%v)\nwant: %v\nhave: %v", test.src, r.arg, r.result, result)
			}
		}
	}
}
