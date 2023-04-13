package gobump

import (
	"io/fs"
	"strings"
)

// WalkDir walks a given implementation of fs.FS and runs f() on
// each file where updates of the import module path are possible.
//
// This function ignores files and directories as per go command's convention.
// See https://pkg.go.dev/cmd/go for more details.
func WalkDir(
	fsys fs.FS,
	f func(file string) error,
) error {
	return fs.WalkDir(
		fsys,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Ignore the root directory.
			if path == "." {
				return nil
			}

			if d.IsDir() && shouldIgnoreDir(d.Name()) {
				return fs.SkipDir
			}

			// Do not process directories.
			if d.IsDir() {
				return nil
			}

			if shouldIgnoreFile(d.Name()) {
				return nil
			}

			return f(path)
		},
	)
}

func shouldIgnoreDir(name string) bool {
	switch {
	case name == "testdata":
		return true
	case strings.HasPrefix(name, "."), strings.HasPrefix(name, "_"):
		return true
	default:
		return false
	}
}

func shouldIgnoreFile(name string) bool {
	switch {
	case strings.HasPrefix(name, "."), strings.HasPrefix(name, "_"):
		return true
	default:
		return false
	}
}
