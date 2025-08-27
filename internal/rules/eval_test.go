package rules

import "testing"

func TestCallMatcher_UnknownReturnsFalse(t *testing.T) {
	prog, err := Compile("Unknown(`value`)")
	if err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}
	ctx := Context{Name: "n", Provider: "p"}
	if prog.Match(ctx) {
		t.Error("unknown matcher should evaluate to false")
	}
}

func TestCallMatcher_UnknownMixedCaseReturnsFalse(t *testing.T) {
	// Mixed case name ensures callMatcher lowercases before switch
	prog, err := Compile("NoSuchMatcher(`x`)")
	if err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}
	if prog.Match(Context{}) {
		t.Error("mixed-case unknown matcher should evaluate to false")
	}
}

func TestMatchRegexp_EmptyPatternTrue(t *testing.T) {
	if !matchRegexp("", "anything") {
		t.Error("empty pattern should match true")
	}
}

func TestEval_UnexpectedOpsReturnFalse(t *testing.T) {
	ctx := Context{}
	if eval(BinaryExpr{Op: ILLEGAL}, ctx) {
		t.Error("unexpected binary op should be false")
	}
	if eval(UnaryExpr{Op: ILLEGAL}, ctx) {
		t.Error("unexpected unary op should be false")
	}
}
