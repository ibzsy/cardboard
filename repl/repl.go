package repl

import (
	"bufio"
	"github.com/ibzsy/cardboard/eval"
	"github.com/ibzsy/cardboard/lexer"
	"github.com/ibzsy/cardboard/object"
	"github.com/ibzsy/cardboard/parser"
	"fmt"
	"os"
	"strings"
)

func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.CreateEnvironment()

	fmt.Println("Cardboard v1.0! type :q to quit REPL.")

	for {
		fmt.Print(">>> ")
		scanner.Scan()

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			return
		} else if input == ":q" {
			fmt.Println("Ending REPL.")
			os.Exit(0)
		}

		lex := lexer.CreateLexer(input)
		parser := parser.CreateParser(lex)
		program := parser.ParseCardBoard()

		if checkParserErrors(parser) {
			continue
		}

		evaluatedProgram := eval.Eval(program, env)
		fmt.Println(evaluatedProgram.Inspect())
	}
}

func checkParserErrors(p *parser.Parser) bool {
	errs := p.GetErrors()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return true
	}
	return false
}
