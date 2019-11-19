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

// Rewrite rewrites bytes containing Go source code to replace pkg/errors
// wrapping with the new "%w" fmt.Errorf wrapping.
func Rewrite(filename string, oldSource []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, oldSource, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}

	rewritten := astrewrite.Walk(file, visitor)

	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, rewritten)
	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w", err)
	}
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
	// Rotate err to the end of a new args list
	newArgs := make([]ast.Expr, len(ce.Args)-1)
	copy(newArgs, ce.Args[1:])
	newArgs = append(newArgs, ce.Args[0])

	// If the format string is a literal, we can rewrite it:
	//     "......" -> "......: %w"
	// Otherwise, we replace it with a binary op to add the wrap code:
	//     SomeNonLiteral -> SomeNonLiteral + ": %w"
	fmtStr, ok := newArgs[0].(*ast.BasicLit)
	if ok {
		// Strip trailing `"` and append wrap code and new trailing `"`
		fmtStr.Value = fmtStr.Value[:len(fmtStr.Value)-1] + `: %w"`
	} else {
		binOp := &ast.BinaryExpr{
			X:  newArgs[0],
			Op: token.ADD,
			Y:  &ast.BasicLit{Kind: token.STRING, Value: `": %w"`},
		}
		newArgs[0] = binOp
	}

	return newErrorfExpr(newArgs)
}

func newErrorfExpr(args []ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "fmt"},
			Sel: &ast.Ident{Name: "Errorf"},
		},
		Args: args,
	}
}
