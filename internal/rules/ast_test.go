package rules

import "testing"

// These tests invoke the no-op isExpr methods to ensure coverage on the AST node types.
func TestAST_isExpr_Binary(t *testing.T) {
	var e BinaryExpr
	e.isExpr()
}

func TestAST_isExpr_Unary(t *testing.T) {
	var e UnaryExpr
	e.isExpr()
}

func TestAST_isExpr_Call(t *testing.T) {
	var e CallExpr
	e.isExpr()
}
