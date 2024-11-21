package main

import (
	"berlang/lexer"
	"berlang/parser"
	"bufio"
	"fmt"
	"os"
	"strings"
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

		result, err := parser.Parse(0)
		if err != nil {
			fmt.Printf("Parsing error: %v", err)
		}

		fmt.Printf("Parsed Result: %+v\n", result)

	}

}
