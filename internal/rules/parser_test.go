package rules

import "testing"

func TestParse_ParenthesesAndEOF(t *testing.T) {
	// valid with parentheses
	if _, err := parse("(Name(`a`) && Provider(`b`))"); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	// invalid: extra token after a valid expr
	if _, err := parse("Name(`a`) )"); err == nil {
		t.Fatalf("expected syntax error for trailing token, got nil")
	}
}

func TestParse_UnaryOperator(t *testing.T) {
	if _, err := parse("!Name(`x`)"); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
}

func TestParse_Errors_PrimaryAndUnary(t *testing.T) {
    // Bad primary: starts with RPAREN
    if _, err := parse(")"); err == nil {
        t.Fatalf("expected syntax error for stray ')', got nil")
    }
    // Missing closing RPAREN after group
    if _, err := parse("(Name(`x`)"); err == nil {
        t.Fatalf("expected syntax error for missing ')', got nil")
    }
    // Dangling NOT without primary
    if _, err := parse("!"); err == nil {
        t.Fatalf("expected syntax error for dangling '!', got nil")
    }
}

func TestParse_Errors_AndOperatorDangling(t *testing.T) {
    // Left side ok, dangling '&&' without right operand
    if _, err := parse("Name(`a`) &&"); err == nil {
        t.Fatalf("expected syntax error for dangling &&, got nil")
    }
    // Also for OR
    if _, err := parse("Name(`a`) ||"); err == nil {
        t.Fatalf("expected syntax error for dangling ||, got nil")
    }
}
