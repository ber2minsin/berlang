package terminal

import (
	"berlang/frontend/lexer"
	"berlang/frontend/parser"
	"berlang/runtime/interpreter"
	"fmt"
	"strings"
	"sync"
)

type Terminal struct {
    runtime interpreter.Runtime
    mu      sync.Mutex
    history []string
}


func NewTerminal() *Terminal {
    return &Terminal{
        runtime: interpreter.NewRuntime(),
        history: make([]string, 0),
    }
}

type CommandResult struct {
    Command string
    Output  string
    Error   string
}

func (t *Terminal) ExecuteCommand(command string) CommandResult {
    t.mu.Lock()
    defer t.mu.Unlock()

    command = strings.TrimSpace(command)
    if command == "" {
        return CommandResult{}
    }

    t.history = append(t.history, command)

    lexer := lexer.NewLexer(strings.NewReader(command))
    ts, err := lexer.Lex()
    if err != nil {
        return CommandResult{
            Command: command,
            Error: "Lexing error: " + err.Error(),
        }
    }

    parser := parser.NewParser(ts)
    result, err := parser.Parse()
    if err != nil {
        return CommandResult{
            Command: command,
            Error: "Parsing error: " + err.Error(),
        }
    }

    rtresult, err := t.runtime.Evaluate(result)
    if err != nil {
        return CommandResult{
            Command: command,
            Error: "Runtime error: " + err.Error(),
        }
    }

    return CommandResult{
        Command: command,
        Output: fmt.Sprintf("%+v", rtresult),
    }
}
