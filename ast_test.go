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
	"testing"
)

func TestRenameFunction(t *testing.T) {
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

	var buf bytes.Buffer
	require.NoError(t, printer.Fprint(&buf, fset, newAst))

	goldenMaster, err := ioutil.ReadFile("foo.go.golden")
	require.NoError(t, err)

	assert.Equal(t, string(goldenMaster), buf.String())
}
