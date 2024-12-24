package interpreter_test

// Initialize the Lexer, parser and runtime for the test, and load in files from ./berlang/*.berl and test if they are valid

import (
	"berlang/frontend/ast"
	"berlang/frontend/lexer"
	"berlang/frontend/parser"
	"berlang/runtime/interpreter"
	"berlang/runtime/values"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestInterpreter(t *testing.T) {
	// Have seperate tests for each case, lets start with var_decl.berl
	t.Run("var_decl.berl", func(t *testing.T) {

		runtime := interpreter.NewRuntime()
		parsed := parseString("let x: int = 5", t)
		result, err := runtime.Evaluate(parsed)
		if err != nil {
			t.Fatalf("Error evaluating file: %v", err)
		}

		expected := values.NumVal{Value: 5, Type: values.NumberValue}
		if !reflect.DeepEqual(result, &expected) {
			t.Fatalf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("var_assign.berl", func(t *testing.T) {

		runtime := interpreter.NewRuntime()
		parsed := parseString("let x: int = 5\n\nx = 10", t)
		result, err := runtime.Evaluate(parsed)
		if err != nil {
			t.Fatalf("Error evaluating file: %v", err)
		}

		expected := values.NumVal{Value: 10, Type: values.NumberValue}
		if !reflect.DeepEqual(result, &expected) {
			t.Fatalf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("recursive_var_assign.berl", func(t *testing.T) {
		runtime := interpreter.NewRuntime()
		parsed := parseString("let x: int = 5\nx = x + 5", t)

		result, err := runtime.Evaluate(parsed)
		if err != nil {
			t.Fatalf("Error evaluating file: %v", err)
		}

		expected := values.NumVal{Value: 10, Type: values.NumberValue}
		if !reflect.DeepEqual(result, &expected) {
			t.Fatalf("Expected %v, got %v", expected, result)
		}
	})

}

func BenchmarkInterpreter(b *testing.B) {
    var expressions []string
    for i := 0; i < b.N; i++ {
        expressions = append(expressions, generateExpression(100, b))
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        expression := expressions[i]
        runtime := interpreter.NewRuntime()
        parsed := parseString(expression, b)
        _, err := runtime.Evaluate(parsed)
        if err != nil {
            b.Fatalf("Error evaluating expression: %v", err)
        }
    }
}

func parse(input io.Reader, tb testing.TB) ast.Stmt {
    tb.Helper()

	lexer := lexer.NewLexer(input)
	tq, err := lexer.Lex()
	if err != nil {
		tb.Fatalf("Error lexing input: %v", err)
	}

	parser := parser.NewParser(tq)
	p, err := parser.Parse()
	if err != nil {
		tb.Fatalf("Error parsing input: %v", err)
	}

	return p
}


var operators = []string{"+", "-", "*", "/"}

func generateExpression(length int, tb testing.TB) string {
    tb.Helper()

	var sb strings.Builder
	openParentheses := 0

	for i := 0; i < length; i++ {
		if rand.Intn(4) == 0 { // 25% chance
			sb.WriteString("(")
			openParentheses++
		}

        // Randomize something > 0 and < 2^31
		sb.WriteString(fmt.Sprintf("%d", rand.Intn(1<<31)))

		if openParentheses > 0 && rand.Intn(4) == 0 { // 25% chance
			sb.WriteString(")")
			openParentheses--
		}

		if i < length-1 {
			sb.WriteString(fmt.Sprintf(" %s ", operators[rand.Intn(len(operators))]))
		}
	}

	// Close any remaining open parentheses
	for openParentheses > 0 {
		sb.WriteString(")")
		openParentheses--
	}

	return sb.String()
}

func parseString(input string, tb testing.TB) ast.Stmt {
	return parse(strings.NewReader(input), tb)
}

func parseFile(filePath string, tb testing.TB) ast.Stmt {
	file, err := os.Open(filePath)
	if err != nil {
		tb.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	return parse(file, tb)
}
