package gobump

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// ParseModules parses go modules in the go.mod of the Go module. This function
// returns the module path and a list of direct dependency module paths declared
// using "require" directive.
//
// For more info on the go.mod file structure refer to this resource:
// https://go.dev/doc/modules/gomod-ref.
func ParseModules(moduleDir string) (string, []string, error) {
	p := filepath.Join(moduleDir, "go.mod")
	bb, err := os.ReadFile(p)
	if err != nil {
		return "", nil, fmt.Errorf(
			"failed to open go.mod file in %q directory: %s",
			moduleDir,
			err,
		)
	}

	mf, err := modfile.Parse(p, bb, nil)
	if err != nil {
		return "", nil, fmt.Errorf(
			"invalid go.mod file in %q directory: %s",
			moduleDir,
			err,
		)
	}

	var directRequires []string
	for _, req := range mf.Require {
		if !req.Indirect {
			directRequires = append(directRequires, req.Mod.Path)
		}
	}

	return mf.Module.Mod.Path, directRequires, nil
}
