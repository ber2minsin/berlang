package parser

import (
	"berlang/utils"
	"errors"
	"fmt"
)

var rules map[utils.TokenType]ParseRule
var UnexpectedTokenError = errors.New("Unexpected token")

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
	tokenStack *utils.TokenQueue
	curToken   utils.Token
}

func NewParser(ts *utils.TokenQueue) *Parser {
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

func (p *Parser) Parse() (interface{}, error) {

    switch p.currentToken().Type {

    case utils.TOKEN_LET:
        return p.parseVariableDeclaration()
	case utils.TOKEN_FUNCTION:
	case utils.TOKEN_IDENT:
	default:
		panic("unexpected utils.TokenType")
	}
    return nil, UnexpectedTokenError

}

func (p *Parser) expectToken(expectedType utils.TokenType) error {
    if err := p.nextToken(); err != nil {
        return err
    }

    if p.currentToken().Type != expectedType {
        return utils.NewParseError(
            string(expectedType),
            string(p.currentToken().Type),
            float64(p.currentToken().Line),
            float64(p.currentToken().Column),
        )
    }

    return nil
}

func (p *Parser) parseVariableDeclaration() (interface{}, error) {

    var varDeclaration map[string]interface{} = make(map[string]interface{})

    if err := p.expectToken(utils.TOKEN_IDENT); err != nil {
        return nil, err
    }
    varDeclaration["name"] = string(p.currentToken().Literal)

    if err := p.expectToken(utils.TOKEN_COLON); err != nil {
        return nil, err
    }

    if err := p.expectToken(utils.TOKEN_TYPE); err != nil {
        return nil, err
    }

	if _, ok := utils.Keywords[p.currentToken().Literal]; ok {
        varDeclaration["type"] = p.currentToken().Literal
	}

    if err := p.expectToken(utils.TOKEN_ASSIGN); err != nil {
        return nil, err
    }

    p.nextToken()

    right, err := p.parseStatement(0)
    if err != nil {
        return nil, UnexpectedTokenError
    }
    varDeclaration["value"] = right
    return varDeclaration, nil
}

func (p *Parser) parseStatement(precedence int8) (interface{}, error) {
	currentToken := p.currentToken()
	currentTokenRule := rules[currentToken.Type]

	lhs, err := currentTokenRule.NUD(p, nil)
	if err != nil {
		return nil, err
	}

	if err := p.nextToken(); err != nil ||
		p.currentToken().Type == utils.TOKEN_EOF ||
		p.currentToken().Type == utils.TOKEN_SEMI {
		return lhs, nil
	}

	for p.tokenStack.Len() > 0 {
		currentToken = p.currentToken()
		currentTokenRule = rules[currentToken.Type]

		if currentTokenRule.LBP <= precedence {
			break
		}

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
				right, err := p.parseStatement(100) // High precedence for unary operations
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"op": "+", "right": right}, nil
			},
			LED: func(p *Parser, left interface{}) (interface{}, error) {
				right, err := p.parseStatement(10) // Same precedence
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"op": "+", "left": left, "right": right}, nil
			},
		},
		utils.TOKEN_MINUS: {
			LBP: 10,
			NUD: func(p *Parser, _ interface{}) (interface{}, error) {
				right, err := p.parseStatement(100) // High precedence for unary minus
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"op": "-", "right": right}, nil
			},
			LED: func(p *Parser, left interface{}) (interface{}, error) {
				right, err := p.parseStatement(10) // Same precedence
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
				right, err := p.parseStatement(20) // Same precedence
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
				right, err := p.parseStatement(20) // Same precedence
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
				p.nextToken()
				expr, err := p.parseStatement(0)
				if err != nil {
					return nil, err
				}

				if p.currentToken().Type != utils.TOKEN_RPAREN {
					return nil, fmt.Errorf("expected ')', got %v", p.currentToken().Type)
				}

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
