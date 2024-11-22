package parser

import (
	"berlang/utils"
	"fmt"
	"testing"
)
func TestParse(t *testing.T) {
    tokenStack := utils.NewTokenQueue()
    // Push tokens in reverse order
    tokenStack.Push(utils.Token{Type: utils.TOKEN_NUMBER, Literal: "2"})
    tokenStack.Push(utils.Token{Type: utils.TOKEN_MULT, Literal: "*"})
    tokenStack.Push(utils.Token{Type: utils.TOKEN_NUMBER, Literal: "5"})
    tokenStack.Push(utils.Token{Type: utils.TOKEN_PLUS, Literal: "+"})
    tokenStack.Push(utils.Token{Type: utils.TOKEN_NUMBER, Literal: "3"})
    tokenStack.Push(utils.Token{Type: utils.TOKEN_PLUS, Literal: "+"})
    tokenStack.Push(utils.Token{Type: utils.TOKEN_NUMBER, Literal: "3"})

    parser := NewParser(tokenStack)
    result, err := parser.Parse()
    if err != nil {
        t.Fatalf("Parsing error: %v", err)
    }

    fmt.Printf("Parsed Result: %+v\n", result)
}
