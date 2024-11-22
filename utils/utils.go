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
	TOKEN_PLUS     TokenType = "PLUS"
	TOKEN_MINUS    TokenType = "MINUS"
	TOKEN_MULT     TokenType = "MULTIPLY"
	TOKEN_DIV      TokenType = "DIVIDE"
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
	')': TOKEN_RPAREN,
	'(': TOKEN_LPAREN,
    '+': TOKEN_PLUS,
    '-': TOKEN_MINUS,
    '*': TOKEN_MULT,
    '/': TOKEN_DIV,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type TokenQueue struct {
	lock   sync.Mutex
	tokens []Token
}

func NewTokenQueue() *TokenQueue {
	return &TokenQueue{}
}

func (ts *TokenQueue) Push(t Token) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	ts.tokens = append(ts.tokens, t)
}

func (ts *TokenQueue) Pop() (Token, error) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	l := len(ts.tokens)
	if l == 0 {
		return Token{}, empty_stack
	}

	res := ts.tokens[0]
	ts.tokens = ts.tokens[1:l]

	return res, nil
}

func (ts *TokenQueue) Len() int {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	return len(ts.tokens)
}

// WARN Use this for testing or debugging
func (ts *TokenQueue) Tokens() []Token {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	copiedTokens := make([]Token, len(ts.tokens))
	copy(copiedTokens, ts.tokens)
	return copiedTokens
}
