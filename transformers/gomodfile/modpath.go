package gomodfile

import (
	"fmt"
	"io"

	"github.com/danilvpetrov/gobump/transformers"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// UpdateModulePath replaces the module path in a go.mod file.
func UpdateModulePath(
	modulePath string,
) transformers.Transformer {
	return func(in io.Reader, out io.Writer) (ok bool, err error) {
		bb, err := io.ReadAll(in)
		if err != nil {
			return false, err
		}

		mf, err := modfile.Parse("", bb, nil)
		if err != nil {
			return false, err
		}

		pfx, _, ok := module.SplitPathVersion(modulePath)
		if !ok {
			return false, fmt.Errorf("target path %s is invalid", modulePath)
		}

		if op, _, _ := module.SplitPathVersion(mf.Module.Mod.Path); op != pfx {
			return false, nil
		}

		if err := mf.AddModuleStmt(modulePath); err != nil {
			return false, err
		}

		bb, err = mf.Format()
		if err != nil {
			return false, err
		}

		if _, err := out.Write(bb); err != nil {
			return false, err
		}

		return true, nil
	}
}
