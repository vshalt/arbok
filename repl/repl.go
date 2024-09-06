package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/vshalt/arbok/evaluator"
	"github.com/vshalt/arbok/lexer"
	"github.com/vshalt/arbok/object"
	"github.com/vshalt/arbok/parser"
)

const PROMPT = `Hello, welcome to arbok!
>> `

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	fmt.Printf(PROMPT)
	for {
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == "exit" {
			break
		}
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n>> ")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "parser ran into errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
