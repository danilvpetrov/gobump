package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/danilvpetrov/gobump/transformers/gofile"
	"github.com/danilvpetrov/gobump/transformers/gomodfile"
)

// walkDir walks a given directory and runs the relevant file transformers that
// apply the new go module path.
func walkDir(dir, newPath string) error {
	return filepath.WalkDir(
		dir,
		func(path string, d fs.DirEntry, err error) error {
			// Do not walk into "testdata" directories.
			if d.IsDir() && d.Name() == "testdata" {
				return filepath.SkipDir
			}

			if d.IsDir() {
				return nil
			}

			// Ignore files starting with "." or "_".
			if strings.HasPrefix(d.Name(), ".") || strings.HasPrefix(d.Name(), "_") {
				return nil
			}

			if d.Type()&fs.ModeSymlink != 0 {
				// Ignore symlinks pointing to directories.
				if t, err := os.Stat(path); err == nil && t.IsDir() {
					return nil
				}
			}

			switch {
			default:
				return nil
			case filepath.Ext(path) == ".go":
				return runTransformers(path, gofile.UpdateImports(newPath))
			case filepath.Base(path) == "go.mod":
				return runTransformers(path, gomodfile.UpdateModulePath(newPath))
			}
		},
	)
}
