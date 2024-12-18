package interpreter


import (
	"berlang/frontend/ast"
	"berlang/runtime/environment"
	"berlang/runtime/values"
	"fmt"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

type Runtime struct {
	// TODO maybe we can handle this differently
	// this is not good for multithreading
	CurEnv environment.Environment
}

func NewRuntime() Runtime {
	return Runtime{environment.NewEnvironment(nil)}
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
	rhs, err := r.Evaluate(be.Right)
	if err != nil {
		return nil, err
	}

	lhs, err := r.Evaluate(be.Left)
	if err != nil {
		return nil, err
	}

	if lhs.GetType() == values.NumberValue && rhs.GetType() == values.NumberValue {
		result, err := r.evaluateNumeric(lhs, rhs, be.Operator)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("Binary expression did not evaluate to a NumericLiteral")

}

func (r *Runtime) evaluateNumeric(lhs, rhs values.RtVal, operator string) (values.RtVal, error) {
	if lhs.GetType() == values.NumberValue && rhs.GetType() == values.NumberValue {
		return r.evalNumericBinaryExpr(lhs.(*values.NumVal), rhs.(*values.NumVal), operator)
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
	case ast.IdentifierType:
        fmt.Printf("Trying to resolve %v+\n", stmt)
		return r.CurEnv.Resolve(stmt.(*ast.Identifier))
    case ast.VarDeclType:
        fmt.Println("Declaring variable")
        return r.CurEnv.DeclareVar((stmt.(*ast.VarDecl)), r)
	default:
		return nil, fmt.Errorf("Unrecognized expression %+v", stmt.GetKind())
	}
}
