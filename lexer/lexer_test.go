package lexer

import (
	"os"
	"testing"
)

func TestLex(t *testing.T) {
	file, err := os.Open("../berlang/000-variable.bl")
	if err != nil {
		t.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a new lexer with the file
	l := NewLexer(file)

	// Lex the file content
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Error during lexing: %v", err)
	}

	// Print the tokens
	for _, token := range tokens {
		t.Logf("%+v", token)
	}
}
