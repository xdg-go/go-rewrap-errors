package rewriter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/fatih/astrewrite"
)

func Rewrite(filename string, oldSource []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, oldSource, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}

	rewritten := astrewrite.Walk(file, visitor)

	buf := &bytes.Buffer{}
	format.Node(buf, fset, rewritten)
	return buf.Bytes(), nil
}

func visitor(n ast.Node) (ast.Node, bool) {
	c, ok := n.(*ast.CallExpr)
	if !ok {
		return n, true
	}

	s, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return n, true
	}

	i, ok := s.X.(*ast.Ident)
	if !ok {
		return n, true
	}

	if i.Name != "errors" {
		return n, true
	}

	switch s.Sel.Name {
	case "Wrap":
		return rewriteWrap(c), true
	case "Wrapf":
		return rewriteWrap(c), true
	default:
	}

	return n, true
}

func rewriteWrap(ce *ast.CallExpr) *ast.CallExpr {
	// fmt.Errorf("......: %w", ..., err)

	newArgs := make([]ast.Expr, len(ce.Args)-1)
	copy(newArgs, ce.Args[1:])
	newArgs = append(newArgs, ce.Args[0])

	fmtStr := newArgs[0].(*ast.BasicLit)
	fmtStr.Value = fmtStr.Value[:len(fmtStr.Value)-1] + ": %w\""

	fe := newErrorfExpr()
	fe.Args = newArgs
	return fe
}

func newErrorfExpr() *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "fmt"},
			Sel: &ast.Ident{Name: "Errorf"},
		},
	}
}
