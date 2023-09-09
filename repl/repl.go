package repl

import (
	"bufio"
	"fmt"
	"github.com/fabiante/monkeylang/lexer"
	"github.com/fabiante/monkeylang/token"
	"io"
	"os"
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

		for t, err := lex.NextToken(); t.Type != token.EOF; t, err = lex.NextToken() {
			if err != nil {
				_, _ = fmt.Fprintf(out, "encountered error: %s", err)
				os.Exit(1)
			}
			_, _ = fmt.Fprintf(out, "%+v\n", t)
		}
	}
}
