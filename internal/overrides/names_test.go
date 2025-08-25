package overrides

import (
	"testing"
)

func TestStripProvider(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"noat", "noat"},
		{"svc@file", "svc"},
		{"ns/svc@kubernetes@file", "ns/svc@kubernetes"},
	}
	for _, c := range cases {
		if got := stripProvider(c.in); got != c.out {
			t.Fatalf("stripProvider(%q)=%q want %q", c.in, got, c.out)
		}
	}
}
