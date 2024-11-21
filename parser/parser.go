package parser

import (
	"berlang/utils"
	"fmt"
)

var rules map[utils.TokenType]ParseRule

type ParseFunc func(p *Parser, left interface{}) (interface{}, error)

type ParseRule struct {
	LBP int8
	NUD ParseFunc
	LED ParseFunc
}

type ParseTree struct {
	NodeType string
	Value    string
	Children *[]ParseTree
}

type Parser struct {
	tokenStack *utils.TokenStack
	curToken   utils.Token
}

func NewParser(ts *utils.TokenStack) *Parser {
	token, err := ts.Pop()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize parser: %v", err))
	}
	return &Parser{tokenStack: ts, curToken: token}
}

func (p *Parser) currentToken() utils.Token {
	return p.curToken
}

func (p *Parser) nextToken() error {
	token, err := p.tokenStack.Pop()
	if err != nil {
		return err
	}
	p.curToken = token
	return nil
}

func (p *Parser) Parse(precedence int8) (interface{}, error) {
	currentToken := p.currentToken()
	currentTokenRule := rules[currentToken.Type]

	lhs, err := currentTokenRule.NUD(p, nil)
	if err != nil {
		return nil, err
	}

	// Advance to the next token after NUD parsing
	if err := p.nextToken(); err != nil {
		// If no more tokens, return the current lhs
		return lhs, nil
	}

	// Continue parsing while there are tokens and precedence allows
	for p.tokenStack.Len() > 0 {
		currentToken = p.currentToken()
		currentTokenRule = rules[currentToken.Type]

		if currentTokenRule.LBP <= precedence {
			break
		}

		// Move to the next token before LED parsing
		if err := p.nextToken(); err != nil {
			break
		}

		lhs, err = currentTokenRule.LED(p, lhs)
		if err != nil {
			return nil, err
		}
	}

	return lhs, nil
}

func init() {
	rules = map[utils.TokenType]ParseRule{
		utils.TOKEN_NUMBER: {
			LBP: 0,
			NUD: func(p *Parser, _ interface{}) (interface{}, error) {
				value := p.currentToken().Literal
				return value, nil
			},
		},
		utils.TOKEN_PLUS: {
			LBP: 10,
			NUD: func(p *Parser, _ interface{}) (interface{}, error) {
				right, err := p.Parse(100) // High precedence for unary operations
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"op": "+", "right": right}, nil
			},
			LED: func(p *Parser, left interface{}) (interface{}, error) {
				right, err := p.Parse(10) // Same precedence
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"op": "+", "left": left, "right": right}, nil
			},
		},
		utils.TOKEN_MINUS: {
			LBP: 10,
			NUD: func(p *Parser, _ interface{}) (interface{}, error) {
				right, err := p.Parse(100) // High precedence for unary minus
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"op": "-", "right": right}, nil
			},
			LED: func(p *Parser, left interface{}) (interface{}, error) {
				right, err := p.Parse(10) // Same precedence
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"op": "-", "left": left, "right": right}, nil
			},
		},
		utils.TOKEN_MULT: {
			LBP: 20,
			NUD: nil, // Multiplication not valid as prefix
			LED: func(p *Parser, left interface{}) (interface{}, error) {
				right, err := p.Parse(20) // Same precedence
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"op": "*", "left": left, "right": right}, nil
			},
		},
		utils.TOKEN_DIV: {
			LBP: 20,
			NUD: nil, // Division not valid as prefix
			LED: func(p *Parser, left interface{}) (interface{}, error) {
				right, err := p.Parse(20) // Same precedence
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"op": "/", "left": left, "right": right}, nil
			},
		},
		utils.TOKEN_IDENT: {
			LBP: 0,
			NUD: func(p *Parser, _ interface{}) (interface{}, error) {
				value := p.currentToken().Literal
				return value, nil
			},
		},
		utils.TOKEN_LPAREN: {
			LBP: 0,
			NUD: func(p *Parser, _ interface{}) (interface{}, error) {
				// Consume the left parenthesis and parse the expression inside
				p.nextToken()        // Consume '('
				expr, err := p.Parse(0) // Parse with low precedence to support expressions inside the parentheses
				if err != nil {
					return nil, err
				}

				// Expect a closing parenthesis
				if p.currentToken().Type != utils.TOKEN_RPAREN {
					return nil, fmt.Errorf("expected ')', got %v", p.currentToken().Type)
				}
				p.nextToken() // Consume ')'

				return expr, nil
			},
		},
		utils.TOKEN_RPAREN: {
			LBP: 0,
			NUD: nil, // Parenthesis closing is not handled directly by NUD
			LED: nil, // Not applicable for LED
		},
	}
}
