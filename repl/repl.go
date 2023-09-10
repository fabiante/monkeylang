package repl

import (
	"bufio"
	"fmt"
	"github.com/fabiante/monkeylang/lexer"
	"github.com/fabiante/monkeylang/token"
	"io"
)

const Prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(Prompt)
		scan := scanner.Scan()
		if !scan {
			return
		}

		line := scanner.Text()

		lex := lexer.NewLexer(line)

		for t := lex.NextToken(); t.Type != token.EOF; t = lex.NextToken() {
			_, _ = fmt.Fprintf(out, "%+v\n", t)
		}
	}
}
