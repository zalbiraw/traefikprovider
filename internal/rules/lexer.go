package rules

import (
	"unicode"
)

type lexer struct {
	src []rune
	pos int
	len int
}

func newLexer(input string) *lexer {
	r := []rune(input)
	return &lexer{src: r, len: len(r)}
}

func (l *lexer) next() rune {
	if l.pos >= l.len {
		return 0
	}
	ch := l.src[l.pos]
	l.pos++
	return ch
}

func (l *lexer) peek() rune {
	if l.pos >= l.len {
		return 0
	}
	return l.src[l.pos]
}

func (l *lexer) skipSpace() {
	for unicode.IsSpace(l.peek()) {
		l.next()
	}
}

func (l *lexer) scanIdent() string {
	start := l.pos - 1
	for {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			l.next()
			continue
		}
		break
	}
	return string(l.src[start:l.pos])
}

func (l *lexer) scanBacktickString() (string, bool) {
	// already consumed opening backtick
	start := l.pos
	for {
		ch := l.next()
		if ch == 0 {
			return "", false
		}
		if ch == '`' {
			return string(l.src[start : l.pos-1]), true
		}
	}
}

// LexAll tokenizes the input into a slice of tokens.
//
//nolint:gocyclo // The small state machine is clearer un-factored.
func LexAll(input string) ([]Token, error) {
	l := newLexer(input)
	toks := []Token{}
	for {
		l.skipSpace()
		ch := l.next()
		switch ch {
		case 0:
			return append(toks, Token{Type: EOF}), nil
		case '(':
			toks = append(toks, Token{Type: LPAREN, Lexeme: "("})
		case ')':
			toks = append(toks, Token{Type: RPAREN, Lexeme: ")"})
		case '!':
			toks = append(toks, Token{Type: NOT, Lexeme: "!"})
		case '&':
			if l.peek() == '&' {
				l.next()
				toks = append(toks, Token{Type: AND, Lexeme: "&&"})
				break
			}
			return nil, ErrSyntax
		case '|':
			if l.peek() == '|' {
				l.next()
				toks = append(toks, Token{Type: OR, Lexeme: "||"})
				break
			}
			return nil, ErrSyntax
		case ',':
			toks = append(toks, Token{Type: COMMA, Lexeme: ","})
		case '`':
			if s, ok := l.scanBacktickString(); ok {
				toks = append(toks, Token{Type: STRING, Lexeme: s})
				break
			}
			return nil, ErrSyntax
		default:
			if unicode.IsLetter(ch) {
				ident := l.scanIdent()
				toks = append(toks, Token{Type: IDENT, Lexeme: ident})
				break
			}
			return nil, ErrSyntax
		}
	}
}
