package lexer

import (
    "bufio"
    "io"
    "strings"
)

type TokenType string

const (
    TOKEN_LET     TokenType = "LET"
    TOKEN_IDENT   TokenType = "IDENT"
    TOKEN_COLON   TokenType = "COLON"
    TOKEN_TYPE    TokenType = "TYPE"
    TOKEN_ASSIGN  TokenType = "ASSIGN"
    TOKEN_NUMBER  TokenType = "NUMBER"
    TOKEN_SEMI    TokenType = "SEMI"
    TOKEN_LBRACE  TokenType = "LBRACE"
    TOKEN_RBRACE  TokenType = "RBRACE"
    TOKEN_EOF     TokenType = "EOF"
    TOKEN_ILLEGAL TokenType = "ILLEGAL"
)

var keywords = map[string]TokenType{
    "let":    TOKEN_LET,
    "int":    TOKEN_TYPE,
    "string": TOKEN_TYPE,
    "bool":   TOKEN_TYPE,
}

var singleCharTokens = map[byte]TokenType{
    ':': TOKEN_COLON,
    '=': TOKEN_ASSIGN,
    ';': TOKEN_SEMI,
    '{': TOKEN_LBRACE,
    '}': TOKEN_RBRACE,
}

type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}

type Lexer struct {
    reader     *bufio.Reader
    ch         byte
    line       int
    column     int
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
    return l.reader.UnreadByte()
}

func (l *Lexer) Lex() ([]Token, error) {
    var tokens []Token

    // Prime the first character
    err := l.readChar()
    if err != nil && err != io.EOF {
        return nil, err
    }

    for {
        l.skipWhitespace()

        tok, err := l.nextToken()
        if err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }

        tokens = append(tokens, tok)

        if tok.Type == TOKEN_EOF {
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

func (l *Lexer) nextToken() (Token, error) {
    var tok Token
    tok.Line = l.line
    tok.Column = l.column

    // Check single-char tokens
    if tokenType, ok := singleCharTokens[l.ch]; ok {
        tok.Type = tokenType
        tok.Literal = string(l.ch)
        if err := l.readChar(); err != nil {
            if err == io.EOF {
                tok.Type = TOKEN_EOF
            } else {
                return tok, err
            }
        }
        return tok, nil
    }

    // Handle other cases
    switch {
    case l.ch == 0:
        tok.Type = TOKEN_EOF
        return tok, io.EOF
    case isLetter(l.ch):
        return l.lexIdentifier()
    case isDigit(l.ch):
        return l.lexNumber()
    default:
        tok.Type = TOKEN_ILLEGAL
        tok.Literal = string(l.ch)
        if err := l.readChar(); err != nil {
            if err == io.EOF {
                tok.Type = TOKEN_EOF
            } else {
                return tok, err
            }
        }
        return tok, nil
    }
}

func (l *Lexer) lexIdentifier() (Token, error) {
    var tok Token
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

    // Unread the last character that broke the loop
    if err := l.unreadChar(); err != nil {
        return tok, err
    }

    return tok, nil
}

func (l *Lexer) lexNumber() (Token, error) {
    var tok Token
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
    tok.Type = TOKEN_NUMBER

    // Unread the last character that broke the loop
    if err := l.unreadChar(); err != nil {
        return tok, err
    }

    return tok, nil
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func lookupIdentifier(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return TOKEN_IDENT
}
