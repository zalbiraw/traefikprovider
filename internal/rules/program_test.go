package rules

import "testing"

func TestProgram_Match_NilReceiver(t *testing.T) {
	var p *Program
	if !p.Match(Context{}) {
		t.Error("nil Program should match (return true)")
	}
}
