package utils

import (
	"errors"
	"sync"
)

var empty_stack = errors.New("Tried popping on an empty stack")

type TokenType string

const (
	TOKEN_LET      TokenType = "LET"
	TOKEN_FUNCTION TokenType = "FUNCTION"
	TOKEN_IDENT    TokenType = "IDENT"
	TOKEN_COLON    TokenType = "COLON"
	TOKEN_TYPE     TokenType = "TYPE"
	TOKEN_ASSIGN   TokenType = "ASSIGN"
	TOKEN_NUMBER   TokenType = "NUMBER"
	TOKEN_SEMI     TokenType = "SEMI"
	TOKEN_LBRACE   TokenType = "LBRACE"
	TOKEN_RBRACE   TokenType = "RBRACE"
	TOKEN_RPAREN   TokenType = "TOKEN_RPAREN"
	TOKEN_LPAREN   TokenType = "TOKEN_LPAREN"
	TOKEN_EOF      TokenType = "EOF"
	TOKEN_ILLEGAL  TokenType = "ILLEGAL"
	TOKEN_TRUE     TokenType = "TRUE"
	TOKEN_FALSE    TokenType = "FALSE"
)

var Keywords = map[string]TokenType{
	"let":    TOKEN_LET,
	"def":    TOKEN_FUNCTION,
	"int":    TOKEN_TYPE,
	"string": TOKEN_TYPE,
	"bool":   TOKEN_TYPE,
	"true":   TOKEN_TRUE,
	"false":  TOKEN_FALSE,
}

var SingleCharTokens = map[byte]TokenType{
	':': TOKEN_COLON,
	'=': TOKEN_ASSIGN,
	';': TOKEN_SEMI,
	'{': TOKEN_LBRACE,
	'}': TOKEN_RBRACE,
	'(': TOKEN_RPAREN,
	')': TOKEN_LPAREN,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type TokenStack struct {
	lock   sync.Mutex
	tokens []Token
}

func NewTokenStack() *TokenStack {
	return &TokenStack{}
}

func (ts *TokenStack) Push(t Token) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	ts.tokens = append(ts.tokens, t)
}

func (ts *TokenStack) Pop() (Token, error) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	l := len(ts.tokens)
	if l == 0 {
		return Token{}, empty_stack
	}

	res := ts.tokens[l-1]
	ts.tokens = ts.tokens[:l-1]

	return res, nil
}

/// WARN Use this for testing or debugging
func (ts *TokenStack) Tokens() []Token {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	copiedTokens := make([]Token, len(ts.tokens))
	copy(copiedTokens, ts.tokens)
	return copiedTokens
}
