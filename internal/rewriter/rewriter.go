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
	oldAST, err := parser.ParseFile(fset, filename, oldSource, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}

	newAST := astrewrite.Walk(oldAST, visitor)

	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, newAST)
	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w", err)
	}
	return buf.Bytes(), nil
}

func visitor(n ast.Node) (ast.Node, bool) {
	switch v := n.(type) {
	case *ast.CallExpr:
		return handleCallExpr(v)
	case *ast.GenDecl:
		return handleImportDecl(v)
	default:
		return n, true
	}
}

func handleCallExpr(ce *ast.CallExpr) (ast.Node, bool) {
	name := getCallExprLiteral(ce)
	switch name {
	case "errors.Wrap":
		return rewriteWrap(ce), true
	case "errors.Wrapf":
		return rewriteWrap(ce), true
	case "errors.Errorf":
		return newErrorfExpr(ce.Args), true
	default:
		return ce, true
	}
}

func handleImportDecl(gd *ast.GenDecl) (ast.Node, bool) {
	// Ignore GenDecl's that aren't imports.
	if gd.Tok != token.IMPORT {
		return gd, true
	}
	// Push "errors" to the front of specs so formatting will sort it with
	// core libraries and discard pkg/errors.
	newSpecs := []ast.Spec{
		&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"errors"`}},
	}
	for _, s := range gd.Specs {
		im, ok := s.(*ast.ImportSpec)
		if !ok {
			continue
		}
		if im.Path.Value == `"github.com/pkg/errors"` {
			continue
		}
		newSpecs = append(newSpecs, s)
	}
	gd.Specs = newSpecs
	return gd, true
}

func rewriteWrap(ce *ast.CallExpr) *ast.CallExpr {
	// Rotate err to the end of a new args list
	newArgs := make([]ast.Expr, len(ce.Args)-1)
	copy(newArgs, ce.Args[1:])
	newArgs = append(newArgs, ce.Args[0])

	// If the format string is a fmt.Sprintf call, we can unwrap it.
	c, name := getCallExpr(newArgs[0])
	if c != nil && name == "fmt.Sprintf" {
		newArgs = append(c.Args, newArgs[1:]...)
	}

	// If the format string is a literal, we can rewrite it:
	//     "......" -> "......: %w"
	//     `......` -> `......: %w`
	// Otherwise, we replace it with a binary op to add the wrap code:
	//     SomeNonLiteral -> SomeNonLiteral + ": %w"
	fmtStr, ok := newArgs[0].(*ast.BasicLit)
	if ok {
		// Strip closing quote character (" or `) and append it back after wrap code.
		fmtStr.Value = fmtStr.Value[:len(fmtStr.Value)-1] + `: %w` + fmtStr.Value[len(fmtStr.Value)-1:]
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

func getCallExpr(n ast.Node) (*ast.CallExpr, string) {
	c, ok := n.(*ast.CallExpr)
	if !ok {
		return nil, ""
	}
	name := getCallExprLiteral(c)
	if name == "" {
		return nil, ""
	}
	return c, name
}

func getCallExprLiteral(c *ast.CallExpr) string {
	s, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	i, ok := s.X.(*ast.Ident)
	if !ok {
		return ""
	}

	return i.Name + "." + s.Sel.Name
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
