package learngoast

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"os"
	"testing"
)

func TestAst(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "foo.go", nil, parser.ParseComments)
	require.NoError(t, err)

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
	require.NoError(t, printer.Fprint(os.Stdout, fset, newAst))
}
