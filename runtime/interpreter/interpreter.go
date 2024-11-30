package interpreter

import (
	"berlang/frontend/ast"
	"berlang/runtime/values"
	"fmt"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

type Runtime struct {
	Variables map[string]interface{} // TODO too broad, perhaps use an interface
}

func (r *Runtime) evalProgramType(p *ast.Program) (values.RtVal, error) {

	var lastEvaluated values.RtVal

	for _, stmt := range p.Body {
		var err error
		spew.Printf("Looking at stmt %+v", stmt)
		lastEvaluated, err = r.Evaluate(stmt)
		if err != nil {

			return nil, err
		}
	}
	return lastEvaluated, nil
}
func (r *Runtime) evalNumericVal(nl *ast.NumericLiteral) (values.RtVal, error) {

	casted, err := strconv.ParseFloat(nl.Value, 64)
	if err != nil {
		return nil, err
	}
	return &values.NumVal{Value: casted, Type: values.NumberValue}, nil
}

func (r *Runtime) evalNumericBinaryExpr(lhs *values.NumVal, rhs *values.NumVal, op string) (values.RtVal, error) {
    switch op {
    case "+":
        return &values.NumVal{Value: lhs.Value + rhs.Value, Type: values.NumberValue}, nil
    case "-":
        return &values.NumVal{Value: lhs.Value - rhs.Value, Type: values.NumberValue}, nil
    case "*":
        return &values.NumVal{Value: lhs.Value * rhs.Value, Type: values.NumberValue}, nil
    case "/":
        if rhs.Value == 0 {
            return nil, fmt.Errorf("division by zero")
        }
        return &values.NumVal{Value: lhs.Value / rhs.Value, Type: values.NumberValue}, nil
    default:
        return nil, fmt.Errorf("unsupported operator: %s", op)
    }
}

func (r *Runtime) evalBinaryExpr(be *ast.BinaryExpr) (values.RtVal, error) {
    fmt.Printf("Evaluating BinaryExpr: %+v\n", be)

    rhs, err := r.Evaluate(be.Right)
    if err != nil {
        fmt.Printf("Error evaluating RHS: %v\n", err)
        return nil, err
    }
    fmt.Printf("RHS evaluated to: %+v\n", rhs)

    lhs, err := r.Evaluate(be.Left)
    if err != nil {
        fmt.Printf("Error evaluating LHS: %v\n", err)
        return nil, err
    }
    fmt.Printf("LHS evaluated to: %+v\n", lhs)

    if lhs.GetType() == values.NumberValue && rhs.GetType() == values.NumberValue {
        lhsNumVal := lhs.(*values.NumVal)
        rhsNumVal := rhs.(*values.NumVal)

        result, err := r.evalNumericBinaryExpr(lhsNumVal, rhsNumVal, be.Operator)
        if err != nil {
            fmt.Printf("Error in numeric binary expression: %v\n", err)
            return nil, err
        }
        fmt.Printf("BinaryExpr result: %+v\n", result)
        return result, nil
    }

    return nil, fmt.Errorf("unsupported binary expression types: %T and %T", lhs, rhs)
}

func (r *Runtime) Evaluate(stmt ast.Stmt) (values.RtVal, error) {
    spew.Printf("Requested to evaluate %+v\n", stmt)

	switch stmt.GetKind() {
	case ast.BinaryExprType:
		return r.evalBinaryExpr(stmt.(*ast.BinaryExpr))
	case ast.ProgramType:
		return r.evalProgramType(stmt.(*ast.Program))
	case ast.NumericLiteralType:
		return r.evalNumericVal(stmt.(*ast.NumericLiteral))
	default:

        return nil, fmt.Errorf("Unrecognized expression %+v", stmt)
	}
}

func NewRuntime() Runtime {
	return Runtime{Variables: make(map[string]interface{})}
}
