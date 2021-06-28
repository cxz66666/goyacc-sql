package parser

import (
	"io"
	"goyacc-sql/lexer"
	"goyacc-sql/types"
)

// Parse returns parsed Spanner DDL statements.
func Parse(r io.Reader,channel chan<- types.DStatements) error {
	impl := lexer.NewLexerImpl(r, &keywordTokenizer{})
	l := newLexerWrapper(impl,channel)
	yyParse(l)
	if l.err != nil {
		return  l.err
	}
	return nil
}
