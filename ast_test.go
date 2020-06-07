package learngoast

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"io/ioutil"
	"os"
	"testing"
)

func TestRenameFunction(t *testing.T) {
	fset, node := parseExampleFile(t)

	assert.Equal(t, "learngoast", node.Name.Name)

	rename := func(c *astutil.Cursor) bool {
		if identNode, ok := c.Node().(*ast.Ident); ok {
			if identNode.Name == "IsApplesauce" {
				newIdent := *identNode
				newIdent.Name = "renamed"
				c.Replace(&newIdent)
			}
		}
		return true
	}
	newAst := astutil.Apply(node, rename, nil)
	require.NotNil(t, newAst)

	var buf bytes.Buffer
	require.NoError(t, printer.Fprint(&buf, fset, newAst))

	goldenMaster, err := ioutil.ReadFile("foo.go.golden")
	require.NoError(t, err)

	assert.Equal(t, string(goldenMaster), buf.String())
}

func TestDuplicateLine(t *testing.T) {
	fset, node := parseExampleFile(t)

	duplicateLine := func(c *astutil.Cursor) bool {
		exprStmt, ok := c.Node().(*ast.ExprStmt)
		if !ok {
			return true
		}
		c.InsertAfter(exprStmt)
		return true
	}
	newAst := astutil.Apply(node, duplicateLine, nil)
	require.NotNil(t, newAst)

	//require.NoError(t, ast.Print(fset, newAst))
	require.NoError(t, printer.Fprint(os.Stdout, fset, newAst))
}

func parseExampleFile(t *testing.T) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "foo.go", nil, parser.ParseComments)
	require.NoError(t, err)
	return fset, node
}
