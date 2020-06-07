package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"os"
)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "x.go", os.Stdin, parser.ParseComments)
	perr(err)

	funcToInline, ok := node.Scope.Lookup("IsApplesauce").Decl.(*ast.FuncDecl)
	if !ok {
		panic("function to inline not found")
	}

	newAst := astutil.Apply(node, inlineFunc(funcToInline), nil)
	perr(printer.Fprint(os.Stdout, fset, newAst))
}

func inlineFunc(funcToInline *ast.FuncDecl) astutil.ApplyFunc {
	return func(c *astutil.Cursor) bool {
		if c.Node() == funcToInline {
			c.Delete()
			return false
		}
		callExpr, ok := c.Node().(*ast.CallExpr)
		if !ok {
			return true
		}
		calledFunc, ok := callExpr.Fun.(*ast.Ident)
		if !ok || calledFunc.Name != funcToInline.Name.Name {
			return true
		}
		stmtToInline := funcToInline.Body.List[0].(*ast.ReturnStmt)
		c.Replace(stmtToInline.Results[0])
		return true
	}
}

func perr(err error) {
	if err != nil {
		panic(err)
	}
}
