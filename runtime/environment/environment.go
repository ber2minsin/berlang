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
	variables map[string]Variable
}

// Create a variable type that stores map[name]variable which has the const/let type and the value
type Variable struct {
    value values.RtVal
    // isConst bool
    varType string
}

func NewVariable(value values.RtVal, varType string) Variable {
    return Variable{value: value, varType: varType}
}

func NewEnvironment(parent *Environment) Environment {
	return Environment{parent: parent, variables: make(map[string]Variable)}
}

func (env *Environment) Resolve(ident *ast.Identifier) (values.RtVal, error) {
	fmt.Printf("Trying to resolve variable %v\n", ident.Name)
	fmt.Printf("Current env variables %v\n", env.variables)

	if env.variables == nil {
		fmt.Println("There is no variables in this environment, looking to the parent")
		return env.parent.Resolve(ident)
	}
	if variable, found := env.variables[ident.Name]; found {
		return variable.value, nil
	}
	if env.parent != nil {
		return env.parent.Resolve(ident)
	}
	return nil, fmt.Errorf("identifier '%s' not found", ident.Name)
}

func (env *Environment) DeclareVar(decl *ast.VarDecl, r EvalInterface) (values.RtVal, error) {
	if decl.Value == nil {
        println("Requested to declare a None variable")
		val := &values.NoneVal{}
		env.variables[decl.Name] = NewVariable(val, decl.ValType)
		return val, nil
	}

	val, err := r.Evaluate(*decl.Value)
	if err != nil {
		return nil, err
	}
	env.variables[decl.Name] = NewVariable(val, decl.ValType)

	return val, nil
}

func (env *Environment) AssignVar(assign *ast.VarAssign, r EvalInterface) (values.RtVal, error) {
    val, err := r.Evaluate(*assign.Value)
    if err != nil {
        return nil, err
    }
    if _, found := env.variables[assign.Name]; !found {
        return nil, fmt.Errorf("variable '%s' not found", assign.Name)
    }

    // Check if the variable is a constant
    if env.variables[assign.Name].varType == "const" {
        return nil, fmt.Errorf("variable '%s' is a constant and cannot be reassigned", assign.Name)
    }

    env.variables[assign.Name] = NewVariable(val, "let")

    return val, nil
}

// TODO Keeping in mind to have RAII in berlang, we need to pass the variables into the child context (as value),
// and drop it at the end of the execution of the current context
