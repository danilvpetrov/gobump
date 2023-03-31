package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// parseGoModFile parses "go.mod" file and returns the main module path and all
// direct dependency module paths declared with "require" directive.
//
// For more info refer to this resource: https://go.dev/doc/modules/gomod-ref.
func parseGoModFile(moduleDir string) ([]string, error) {
	gomodPath := filepath.Join(moduleDir, "go.mod")
	if fi, err := os.Stat(gomodPath); err != nil || fi.IsDir() {
		return nil, fmt.Errorf("directory %q is not a valid go module", moduleDir)
	}

	bb, err := os.ReadFile(gomodPath)
	if err != nil {
		return nil, err
	}

	mf, err := modfile.Parse(gomodPath, bb, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid go.mod file: %s", err)
	}

	modules := []string{mf.Module.Mod.Path}
	for _, req := range mf.Require {
		if !req.Indirect {
			modules = append(modules, req.Mod.Path)
		}
	}

	return modules, nil
}

// checkModulePath checks the module path for its validity and relation to any
// of the available module paths.
func checkModulePath(modulePath string, availableModules []string) error {
	if err := module.CheckPath(modulePath); err != nil {
		return fmt.Errorf("invalid module path %q: %w", modulePath, err)
	}

	pnew, _, _ := module.SplitPathVersion(modulePath)
	for _, am := range availableModules {
		if pold, _, _ := module.SplitPathVersion(am); pold == pnew {
			return nil
		}
	}

	return fmt.Errorf(
		"failed to find a corresponding module for %q",
		modulePath,
	)
}
