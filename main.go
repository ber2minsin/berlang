package main

import (
	"berlang/frontend/lexer"
	"berlang/frontend/parser"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Interactive Berlang Shell - Press Ctrl+C to exit")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		lexer := lexer.NewLexer(strings.NewReader(line))

		ts, err := lexer.Lex()
		if err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
        fmt.Printf("Lexer returned: %+v\n", ts)

		parser := parser.NewParser(ts)

		result, err := parser.Parse()
		if err != nil {
			fmt.Printf("Parsing error: %v", err)
		}

		spew.Printf("Parsed Result: %+v\n", result)

	}

}
