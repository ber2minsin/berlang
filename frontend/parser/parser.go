package parser

import (
	"berlang/frontend/ast"
	"berlang/utils"
	"errors"
	"fmt"
)

var rules map[utils.TokenType]ParseRule
var UnexpectedTokenError = errors.New("Unexpected token")

type ParseFunc func(p *Parser, left ast.Expr) (ast.Expr, error)

type ParseRule struct {
	LBP int8
	NUD ParseFunc
	LED ParseFunc
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

func (p *Parser) Parse() (ast.Stmt, error) {

	program := ast.NewProgram()
	switch p.currentToken().Type {

	case utils.TOKEN_LET:
		// return p.parseVariableDeclaration()
	case utils.TOKEN_FUNCTION:
	default:
        parsedStmt, err := p.parseExpr(0)
        if err != nil {
            return nil, err
        }
		program.Body = append(program.Body, parsedStmt)
        return program, nil
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

	right, err := p.parseExpr(0)
	if err != nil {
		return nil, UnexpectedTokenError
	}
	varDeclaration["value"] = right
	return varDeclaration, nil
}

func (p *Parser) parseExpr(precedence int8) (ast.Expr, error) {
	currentToken := p.currentToken()
	currentTokenRule := rules[currentToken.Type]

	lhs, err := currentTokenRule.NUD(p, nil)
	if err != nil {
		fmt.Printf("No NUD was found for the token")
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
			NUD: func(p *Parser, _ ast.Expr) (ast.Expr, error) {
				numericLiteral := ast.NewNumericLiteral(p.curToken.Literal)
				return numericLiteral, nil
			},
		},
		utils.TOKEN_PLUS: {
			LBP: 10,
			// NUD: func(p *Parser, _ ast.Expr) (ast.Expr, error) {
			// right, err := p.parseStatement(100)
			// if err != nil {
			// 	return nil, err
			// }
			// TODO add UnaryExpression
			// return nil, err
			// },
			LED: func(p *Parser, left ast.Expr) (ast.Expr, error) {

				right, err := p.parseExpr(10)

				if err != nil {
					return nil, err
				}
				newBinExpr := ast.NewBinaryExpr(left, right, "+")
				return newBinExpr, nil
			},
		},
		utils.TOKEN_MINUS: {
			LBP: 10,
			// NUD: func(p *Parser, _ ast.Expr) (ast.Expr, error) {
			// 	right, err := p.parseStatement(100)
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// 	return map[string]ast.Expr{"op": "-", "right": right}, nil
			// },
			LED: func(p *Parser, left ast.Expr) (ast.Expr, error) {
				right, err := p.parseExpr(10)

				if err != nil {
					return nil, err
				}
				newBinExpr := ast.NewBinaryExpr(left, right, "-")
				return newBinExpr, nil
			},
		},
		utils.TOKEN_MULT: {
			LBP: 20,
			NUD: nil,
			LED: func(p *Parser, left ast.Expr) (ast.Expr, error) {
				right, err := p.parseExpr(20)

				if err != nil {
					return nil, err
				}
				newBinExpr := ast.NewBinaryExpr(left, right, "*")
				return newBinExpr, nil
			},
		},
		utils.TOKEN_DIV: {
			LBP: 20,
			NUD: nil,
			LED: func(p *Parser, left ast.Expr) (ast.Expr, error) {
				right, err := p.parseExpr(20)

				if err != nil {
					return nil, err
				}
				newBinExpr := ast.NewBinaryExpr(left, right, "/")
				return newBinExpr, nil
			},
		},
		utils.TOKEN_IDENT: {
			LBP: 0,
			NUD: func(p *Parser, _ ast.Expr) (ast.Expr, error) {
				value := p.currentToken().Literal
				return ast.NewIdentifier(value), nil
			},
		},
		utils.TOKEN_LPAREN: {
			LBP: 0,
			NUD: func(p *Parser, _ ast.Expr) (ast.Expr, error) {
				p.nextToken()
				expr, err := p.parseExpr(0)
				if err != nil {
					return nil, err
				}

				if p.currentToken().Type != utils.TOKEN_RPAREN {
					return nil, utils.NewParseError(")", p.currentToken().Literal, float64(p.curToken.Line), float64(p.currentToken().Column))
				}

				return expr, nil
			},
		},
		utils.TOKEN_RPAREN: {
			LBP: 0,
			NUD: nil,
			LED: nil,
		},
	}
}
