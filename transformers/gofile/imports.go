package gofile

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"strconv"

	"github.com/danilvpetrov/gobump/transformers"
	"github.com/danilvpetrov/gobump/transformers/internal/pathx"
	"golang.org/x/tools/go/ast/astutil"
)

// UpdateImports replaces the import path in a .go file with a new import path
// if the latter is applicable.
func UpdateImports(
	newImportPath string,
) transformers.Transformer {
	return func(in io.Reader, out io.Writer) (ok bool, err error) {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", in, parser.ParseComments)
		if err != nil {
			return false, err
		}

		var rewrote bool
		for _, i := range f.Imports {
			p := importPath(i)
			np, ok, err := pathx.UpdateImportPath(newImportPath, p)
			if err != nil {
				return false, err
			}
			if ok {
				rewrote = astutil.RewriteImport(fset, f, p, np)
			}
		}
		if !rewrote {
			return false, nil
		}

		if err := format.Node(out, fset, f); err != nil {
			return false, err
		}

		return true, nil
	}
}

func importPath(i *ast.ImportSpec) string {
	p, err := strconv.Unquote(i.Path.Value)
	if err != nil {
		return ""
	}

	return p
}
