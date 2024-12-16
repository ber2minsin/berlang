package environment

import (
	"berlang/frontend/ast"
	"berlang/runtime/values"
	"fmt"
)
type EvalInterface interface {
	Evaluate(stmt ast.Stmt) (values.RtVal, error)
}
type Environment struct {
	parent    *Environment
	variables map[string]values.RtVal
}

func NewEnvironment(parent *Environment) Environment {
	return Environment{parent: parent, variables: make(map[string]values.RtVal)}
}

func (env *Environment) Resolve(ident *ast.Identifier) (values.RtVal, error) {
    fmt.Printf("Trying to resolve variable %v\n", ident.Name)
    fmt.Printf("Current env variables %v\n", env.variables)

	if env.variables == nil {
		fmt.Println("There is no variables in this environment, looking to the parent")
		return env.parent.Resolve(ident)
	}
	if value, found := env.variables[ident.Name]; found {
		return value, nil
	}
	if env.parent != nil {
		return env.parent.Resolve(ident)
	}
	return nil, fmt.Errorf("identifier '%s' not found", ident.Name)
}

func (env *Environment) DeclareVar(decl *ast.VarDecl, r EvalInterface) (values.RtVal, error) {
	val, err := r.Evaluate(*decl.Value)
	if err != nil {
		return nil, err
	}
	env.variables[decl.Name] = val

    fmt.Printf("Current env variables %v\n", env.variables)
	return val, nil
}

// TODO Keeping in mind to have RAII in berlang, we need to pass the variables into the child context (as value),
// and drop it at the end of the execution of the current context
