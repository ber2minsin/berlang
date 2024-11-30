package lexer

import (
	"berlang/utils"
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Lexer struct {
	reader *bufio.Reader
	ch     byte
	line   int
	column int
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(r),
		line:   1,
		column: 0,
	}
}

func (l *Lexer) readChar() error {
	ch, err := l.reader.ReadByte()
	if err == io.EOF {
		l.ch = 0
		return err
	} else if err != nil {
		return err
	}

	l.ch = ch
	l.column++
	if ch == '\n' {
		l.line++
		l.column = 0
	}
	return nil
}

func (l *Lexer) unreadChar() error {
	l.column--
	fmt.Println("Unreading character")
	// TODO Handle where we are unreading the first token so
	// column is len(line) and line--
	return l.reader.UnreadByte()
}

func (l *Lexer) Lex() (*utils.TokenQueue, error) {
	var tokens = utils.NewTokenQueue()

	// Prime the first character
	err := l.readChar()
	if err != nil && err != io.EOF {
		return utils.NewTokenQueue(), err
	}

	for l.ch != 0 {
		l.skipWhitespace()

		tok, err := l.nextToken()
		if err != io.EOF && err != nil {
			return utils.NewTokenQueue(), err
		}

		tokens.Push(tok)

		if tok.Type == utils.TOKEN_EOF {
			break
		}
	}

	return tokens, nil
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if err := l.readChar(); err != nil {
			return
		}
	}
}

func (l *Lexer) nextToken() (utils.Token, error) {
	var tok utils.Token
	tok.Line = l.line
	tok.Column = l.column

	if l.ch == 0 {
		tok.Type = utils.TOKEN_EOF
		return tok, io.EOF
	}

    // TODO add k that would mean look k characters ahead to see if anything more matches
    // This needs to be done, at least k=1 because we need to, for example discriminate > and >=

	if tokenType, ok := utils.SingleCharTokens[l.ch]; ok {
		tok.Type = tokenType
		tok.Literal = string(l.ch)

		l.readChar()

		return tok, nil
	}

	// Handle other cases
	switch {
	case isLetter(l.ch):
		return l.lexIdentifier()
	case isDigit(l.ch):
		return l.lexNumber()
	default:
		tok.Type = utils.TOKEN_ILLEGAL
		tok.Literal = string(l.ch)

		return tok, nil
	}
}

func (l *Lexer) lexIdentifier() (utils.Token, error) {
	var tok utils.Token
	tok.Line = l.line
	tok.Column = l.column

	var sb strings.Builder
	for isLetter(l.ch) || isDigit(l.ch) {
		sb.WriteByte(l.ch)
		if err := l.readChar(); err != nil {
			break
		}
	}

	tok.Literal = sb.String()
	tok.Type = lookupIdentifier(tok.Literal)

	return tok, nil
}

func (l *Lexer) lexNumber() (utils.Token, error) {
	var tok utils.Token
	tok.Line = l.line
	tok.Column = l.column

	var sb strings.Builder
	for isDigit(l.ch) {
		sb.WriteByte(l.ch)
		if err := l.readChar(); err != nil {
			break
		}
	}

	tok.Literal = sb.String()
	tok.Type = utils.TOKEN_NUMBER

	return tok, nil
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func lookupIdentifier(ident string) utils.TokenType {
	if tok, ok := utils.Keywords[ident]; ok {
		return tok
	}
	return utils.TOKEN_IDENT
}
