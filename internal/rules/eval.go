package rules

import (
	"regexp"
	"strings"
)

// Context describes the fields available to matchers.
type Context struct {
	Name        string
	Provider    string
	Entrypoints []string
	Service     string
}

// Program is a compiled rule expression.
type Program struct {
	expr Expr
}

// Match evaluates the program against a context.
func (p *Program) Match(ctx Context) bool {
	if p == nil || p.expr == nil { return true }
	return eval(p.expr, ctx)
}

func eval(e Expr, ctx Context) bool {
	switch n := e.(type) {
	case BinaryExpr:
		switch n.Op {
		case AND:
			return eval(n.Left, ctx) && eval(n.Right, ctx)
		case OR:
			return eval(n.Left, ctx) || eval(n.Right, ctx)
		}
		return false
	case UnaryExpr:
		if n.Op == NOT { return !eval(n.Expr, ctx) }
		return false
	case CallExpr:
		return callMatcher(n.Name, n.Arg, ctx)
	default:
		return false
	}
}

func anyEntry(entries []string, pred func(string) bool) bool {
	for _, e := range entries { if pred(e) { return true } }
	return false
}

func callMatcher(name, arg string, ctx Context) bool {
	switch strings.ToLower(name) {
	case "name":
		return ctx.Name == arg
	case "nameregexp":
		return matchRegexp(arg, ctx.Name)
	case "provider":
		return ctx.Provider == arg
	case "providerregexp":
		return matchRegexp(arg, ctx.Provider)
	case "entrypoint":
		return anyEntry(ctx.Entrypoints, func(e string) bool { return e == arg })
	case "entrypointregexp":
		return anyEntry(ctx.Entrypoints, func(e string) bool { return matchRegexp(arg, e) })
	case "service":
		return ctx.Service == arg
	case "serviceregexp":
		return matchRegexp(arg, ctx.Service)
	default:
		return false
	}
}

func matchRegexp(pattern, value string) bool {
	if pattern == "" { return true }
	re, err := regexp.Compile(pattern)
	if err != nil { return false }
	return re.MatchString(value)
}
