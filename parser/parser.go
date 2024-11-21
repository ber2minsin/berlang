package parser

import (
	"berlang/utils"
)

type Parser struct {
	tokenStack *utils.TokenStack
	curToken   utils.Token
}

func NewParser(ts *utils.TokenStack) *Parser {
	token, err := ts.Pop()
	if err != nil {
	}
	return &Parser{tokenStack: ts, curToken: token}
}

func (p *Parser) Parse() {


}
