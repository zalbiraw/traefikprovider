package rules

// TokenType represents the lexical token type.
type TokenType int

const (
	// EOF marks end of input.
	EOF TokenType = iota
	// ILLEGAL marks an invalid or unexpected token.
	ILLEGAL
	// IDENT is a matcher name identifier.
	IDENT // matcher name
	// STRING is a backtick string literal.
	STRING // backtick string literal
	// LPAREN is the left parenthesis '('.
	LPAREN // (
	// RPAREN is the right parenthesis ')'.
	RPAREN // )
	// AND is the logical and operator '&&'.
	AND // &&
	// OR is the logical or operator '||'.
	OR // ||
	// NOT is the logical not operator '!'.
	NOT // !
	// COMMA is the comma ',' (reserved).
	COMMA // , (not used now but reserved)
)

// Token represents a lexical token with its type and literal lexeme.
type Token struct {
	Type   TokenType
	Lexeme string
}
