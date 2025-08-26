package rules

import "testing"

func TestCompileAndMatch_BasicMatchers(t *testing.T) {
	ctx := Context{
		Name:        "web-router",
		Provider:    "file",
		Entrypoints: []string{"web", "websecure"},
		Service:     "api-service",
	}

	cases := []struct {
		rule string
		exp  bool
		name string
	}{
		{"", true, "empty rule matches all"},
		{"Name(`web-router`)", true, "name exact true"},
		{"Name(`api-router`)", false, "name exact false"},
		{"NameRegexp(`web-.*`)", true, "name regexp true"},
		{"NameRegexp(`^api-`)", false, "name regexp false"},
		{"Provider(`file`)", true, "provider exact true"},
		{"ProviderRegexp(`fi.*`)", true, "provider regexp true"},
		{"Entrypoint(`web`)", true, "entrypoint exact true"},
		{"Entrypoint(`admin`)", false, "entrypoint exact false"},
		{"EntrypointRegexp(`web.*`)", true, "entrypoint regexp true"},
		{"Service(`api-service`)", true, "service exact true"},
		{"ServiceRegexp(`api-.*`)", true, "service regexp true"},
	}

	for _, tc := range cases {
		prog, err := Compile(tc.rule)
		if err != nil {
			t.Fatalf("Compile(%q) error: %v", tc.rule, err)
		}
		if got := prog.Match(ctx); got != tc.exp {
			t.Errorf("%s: rule=%q match got=%v want=%v", tc.name, tc.rule, got, tc.exp)
		}
	}
}

func TestCompileAndMatch_LogicalOps(t *testing.T) {
	ctx := Context{Name: "web-router", Provider: "file", Entrypoints: []string{"web"}, Service: "svc"}

	cases := []struct {
		rule string
		exp  bool
		name string
	}{
		{"Name(`web-router`) && Provider(`file`)", true, "and true"},
		{"Name(`web-router`) && Provider(`consul`)", false, "and false"},
		{"Name(`web-router`) || Provider(`consul`)", true, "or true"},
		{"!(Provider(`consul`))", true, "not true"},
		{"!(Name(`web-router`))", false, "not false"},
		{"Name(`x`) || (Provider(`file`) && Entrypoint(`web`))", true, "precedence with parens"},
	}

	for _, tc := range cases {
		prog, err := Compile(tc.rule)
		if err != nil {
			t.Fatalf("Compile(%q) error: %v", tc.rule, err)
		}
		if got := prog.Match(ctx); got != tc.exp {
			t.Errorf("%s: rule=%q match got=%v want=%v", tc.name, tc.rule, got, tc.exp)
		}
	}
}

func TestCompile_InvalidSyntax(t *testing.T) {
	bad := []string{
		"Name(`a`",                    // missing )
		"Unknown(`a`)",                // unknown ident is allowed but will never match (handled at eval), still valid syntax
		"Name(`a`) &&& Provider(`b`)", // bad operator
		"Name(`a`) Provider(`b`)",     // missing operator
	}

	for _, rule := range bad {
		if rule == "Unknown(`a`)" {
			// This should compile, syntax-wise.
			if _, err := Compile(rule); err != nil {
				t.Errorf("expected Unknown() to be syntactically valid, got error: %v", err)
			}
			continue
		}
		if _, err := Compile(rule); err == nil {
			t.Errorf("Compile(%q) expected error, got nil", rule)
		}
	}
}

func TestMatch_InvalidRegexPattern(t *testing.T) {
	prog, err := Compile("NameRegexp(`(`)")
	if err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}
	ctx := Context{Name: "any"}
	if prog.Match(ctx) {
		t.Error("invalid regex should result in no match")
	}
}
