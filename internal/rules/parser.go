// Package rules implements the matcher DSL AST and parser.
package rules

import "fmt"

// ErrSyntax indicates an invalid matcher rule syntax.
var ErrSyntax = fmt.Errorf("invalid rule syntax")

type parser struct {
	toks []Token
	pos  int
	len  int
}

func newParser(toks []Token) *parser { return &parser{toks: toks, len: len(toks)} }

func (p *parser) peek() Token {
	if p.pos >= p.len {
		return Token{Type: EOF}
	}
	return p.toks[p.pos]
}

func (p *parser) next() Token {
	t := p.peek()
	p.pos++
	return t
}

func (p *parser) expect(tt TokenType) (Token, error) {
	t := p.next()
	if t.Type != tt {
		return t, ErrSyntax
	}
	return t, nil
}

// lint:ignore ireturn Parser intentionally returns the Expr interface for its AST nodes.
func parse(input string) (Expr, error) {
	toks, err := LexAll(input)
	if err != nil {
		return nil, err
	}
	p := newParser(toks)
	expr, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if p.peek().Type != EOF {
		return nil, ErrSyntax
	}
	return expr, nil
}

func (p *parser) parseOr() (Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == OR {
		p.next()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = BinaryExpr{Op: OR, Left: left, Right: right}
	}
	return left, nil
}

func (p *parser) parseAnd() (Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == AND {
		p.next()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = BinaryExpr{Op: AND, Left: left, Right: right}
	}
	return left, nil
}

func (p *parser) parseUnary() (Expr, error) {
	if p.peek().Type == NOT {
		p.next()
		e, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return UnaryExpr{Op: NOT, Expr: e}, nil
	}
	return p.parsePrimary()
}

func (p *parser) parsePrimary() (Expr, error) {
	t := p.peek()
	switch t.Type {
	case IDENT:
		name := p.next().Lexeme
		if _, err := p.expect(LPAREN); err != nil {
			return nil, err
		}
		argTok, err := p.expect(STRING)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(RPAREN); err != nil {
			return nil, err
		}
		return CallExpr{Name: name, Arg: argTok.Lexeme}, nil
	case LPAREN:
		p.next()
		e, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(RPAREN); err != nil {
			return nil, err
		}
		return e, nil
	default:
		return nil, ErrSyntax
	}
}

// Compile compiles a rule string to an executable Program.
func Compile(rule string) (*Program, error) {
	if rule == "" {
		return &Program{expr: nil}, nil
	}
	e, err := parse(rule)
	if err != nil {
		return nil, err
	}
	return &Program{expr: e}, nil
}
