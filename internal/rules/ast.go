package rules

// Expr is a node in the expression AST.
type Expr interface {
    isExpr()
}

type (
    // BinaryExpr represents left OP right
    BinaryExpr struct {
        Op    TokenType
        Left  Expr
        Right Expr
    }

    // UnaryExpr represents OP Expr
    UnaryExpr struct {
        Op   TokenType
        Expr Expr
    }

    // CallExpr represents Ident(arg)
    CallExpr struct {
        Name string
        Arg  string
    }
)

func (BinaryExpr) isExpr() {}
func (UnaryExpr) isExpr()  {}
func (CallExpr) isExpr()   {}
