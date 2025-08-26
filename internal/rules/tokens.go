package rules

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL
	IDENT     // matcher name
	STRING    // backtick string literal
	LPAREN    // (
	RPAREN    // )
	AND       // &&
	OR        // ||
	NOT       // !
	COMMA     // , (not used now but reserved)
)

type Token struct {
	Type  TokenType
	Lexeme string
}
