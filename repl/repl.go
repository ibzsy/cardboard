package repl

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/ibzsy/cardboard/eval"
	"github.com/ibzsy/cardboard/lexer"
	"github.com/ibzsy/cardboard/object"
	"github.com/ibzsy/cardboard/parser"
)

func StartREPL(buffer *bytes.Buffer, env *object.Environment) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()

	input := strings.TrimSpace(scanner.Text())

	if input == "" {
		return
	} else if input == ":q" {
		os.Exit(0)
	}

	lex := lexer.CreateLexer(input)
	parser := parser.CreateParser(lex)
	program := parser.ParseCardBoard()

	if checkParserErrors(parser, buffer) {
		return
	}

	evaluatedProgram := eval.Eval(program, env)
	buffer.WriteString(evaluatedProgram.Inspect())
}

func checkParserErrors(p *parser.Parser, buffer *bytes.Buffer) bool {
	errs := p.GetErrors()
	if len(errs) > 0 {
		for _, err := range errs {
			buffer.WriteString(err)
		}
		return true
	}
	return false
}
