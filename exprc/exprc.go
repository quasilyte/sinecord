package exprc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"strconv"
)

func Compile(src string) (func(float64) float64, error) {
	var c compiler
	c.src = src
	runner, err := c.CompileRoot()
	if err != nil {
		return nil, err
	}
	f := func(x float64) float64 {
		return runner.Run(x)
	}
	return f, nil
}

type compiler struct {
	src string

	insts         []instructon
	constants     []float64
	constantsPool map[float64]uint8
}

func (c *compiler) CompileRoot() (runner *funcRunner, err error) {
	defer func() {
		rv := recover()
		if rv != nil {
			if recoveredError, ok := rv.(error); ok {
				err = recoveredError
			} else {
				panic(rv)
			}
		}
	}()

	astExpr, err := parser.ParseExpr(c.src)
	if err != nil {
		return nil, err
	}
	c.compileExpr(astExpr)

	runner = &funcRunner{
		stack:     make([]float64, 0, 4),
		constants: c.constants,
		insts:     c.insts,
	}
	return runner, nil
}

func (c *compiler) throwf(format string, args ...any) {
	panic(fmt.Errorf(format, args...))
}

func (c *compiler) internConst(v float64) uint8 {
	if c.constantsPool == nil {
		c.constantsPool = map[float64]uint8{}
	}
	if id, ok := c.constantsPool[v]; ok {
		return uint8(id)
	}
	id := uint8(len(c.constants))
	c.constants = append(c.constants, v)
	c.constantsPool[v] = id
	return id
}

func (c *compiler) emit0(op operation) {
	c.insts = append(c.insts, instructon{
		op: op,
	})
}

func (c *compiler) emit1(op operation, arg uint8) {
	c.insts = append(c.insts, instructon{
		op:  op,
		arg: arg,
	})
}

func (c *compiler) compileExpr(e ast.Expr) {
	switch e := e.(type) {
	case *ast.ParenExpr:
		c.compileExpr(e.X)
	case *ast.BasicLit:
		c.compileBasicLit(e)
	case *ast.UnaryExpr:
		c.compileUnaryExpr(e)
	case *ast.BinaryExpr:
		c.compileBinaryExpr(e)
	case *ast.Ident:
		c.compileIdent(e)
	case *ast.CallExpr:
		c.compileCallExpr(e)
	default:
		c.throwf("unexpected or malformed expression")
	}
}

func (c *compiler) compileCallExpr(e *ast.CallExpr) {
	fn, ok := e.Fun.(*ast.Ident)
	if !ok {
		c.throwf("expected a function name, found something else")
	}

	funcInfo, ok := builtinFuncMap[fn.Name]
	if !ok {
		c.throwf("unknown function %q", fn.Name)
	}
	if len(e.Args) != funcInfo.numArgs {
		c.throwf("%q expects %d arguments, found %d", fn.Name, funcInfo.numArgs, len(e.Args))
	}

	for _, arg := range e.Args {
		c.compileExpr(arg)
	}
	c.emit0(funcInfo.op)
}

func (c *compiler) compileIdent(e *ast.Ident) {
	switch e.Name {
	case "x":
		c.emit0(opArg)
	case "pi":
		c.emit1(opFloatConst, c.internConst(math.Pi))
	case "phi":
		c.emit1(opFloatConst, c.internConst(math.Phi))
	case "e":
		c.emit1(opFloatConst, c.internConst(math.E))
	default:
		c.throwf("unknown variable %q", e.Name)
	}
}

func (c *compiler) compileBasicLit(e *ast.BasicLit) {
	switch e.Kind {
	case token.INT:
		v, err := strconv.ParseInt(e.Value, 0, 64)
		if err != nil {
			panic(err)
		}
		c.emit1(opFloatConst, c.internConst(float64(v)))
	case token.FLOAT:
		v, err := strconv.ParseFloat(e.Value, 64)
		if err != nil {
			panic(err)
		}

		c.emit1(opFloatConst, c.internConst(v))
	default:
		c.throwf("unexpected literal: %v", e.Value)
	}
}

func (c *compiler) compileBinaryExpr(e *ast.BinaryExpr) {
	c.compileExpr(e.X)
	c.compileExpr(e.Y)

	switch e.Op {
	case token.ADD:
		c.emit0(opAdd)
	case token.SUB:
		c.emit0(opSub)
	case token.MUL:
		c.emit0(opMul)
	case token.QUO:
		c.emit0(opDiv)
	default:
		c.throwf("unexpected binary operator: %s", e.Op)
	}
}

func (c *compiler) compileUnaryExpr(e *ast.UnaryExpr) {
	c.compileExpr(e.X)

	switch e.Op {
	case token.SUB:
		c.emit0(opNeg)
	default:
		c.throwf("unexpected unary operator: %s", e.Op)
	}
}
