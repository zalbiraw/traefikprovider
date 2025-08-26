package rules

import "testing"

func TestLexAll_Errors(t *testing.T) {
	cases := []string{
		"&",           // lone &
		"|",           // lone |
		"`unterminated", // unterminated backtick string
		"$",           // invalid starting char
	}
	for _, in := range cases {
		if _, err := LexAll(in); err == nil {
			t.Errorf("LexAll(%q) expected error, got nil", in)
		}
	}
}

func TestLexAll_IdentAndString(t *testing.T) {
	// happy path to ensure coverage for IDENT and STRING scanning
	toks, err := LexAll("Name(`abc_123`) && Provider(`file`)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(toks) == 0 {
		t.Fatal("expected some tokens")
	}
}

func TestLexAll_CommaAndOperators(t *testing.T) {
	// ensure COMMA and operators are recognized
	input := "! ( ) , && ||"
	toks, err := LexAll(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// expected sequence (ignoring EOF at end): NOT, LPAREN, RPAREN, COMMA, AND, OR
	want := []TokenType{NOT, LPAREN, RPAREN, COMMA, AND, OR}
	for i, tt := range want {
		if i >= len(toks)-1 { // -1 to ignore EOF
			t.Fatalf("missing token at %d, toks=%v", i, toks)
		}
		if toks[i].Type != tt {
			t.Fatalf("token %d type=%v want=%v", i, toks[i].Type, tt)
		}
	}
}
