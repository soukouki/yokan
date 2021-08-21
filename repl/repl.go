package repl

import (
	"bufio"
	"fmt"
	"io"
	
	"yokan/lexer"
	"yokan/parser"
	"yokan/object"
	"yokan/evaluator"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evalated := evaluator.Eval(program, env)
		if evalated != nil {
			ret := evalated
			if ret.Type() == object.SHOULD_NOT_VIEWABLE_OBJ { continue }
			io.WriteString(out, ret.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}