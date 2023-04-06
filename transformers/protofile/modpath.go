package protofile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/danilvpetrov/gobump/transformers"
	"github.com/danilvpetrov/gobump/transformers/internal/pathx"
)

// UpdateModulePath updates a corresponding a Go module path in *.proto  file.
func UpdateModulePath(
	modulePath string,
) transformers.Transformer {
	return func(in io.Reader, out io.Writer) (ok bool, err error) {
		var (
			buf        bytes.Buffer
			isModified bool
		)
		s := bufio.NewScanner(in)

		for s.Scan() {
			l := s.Text()

			switch {
			case strings.HasPrefix(l, "import"):
				nl, ok, err := updateImport(l, modulePath)
				if err != nil {
					return false, err
				}
				if ok {
					isModified = true
					l = nl
				}
			case strings.HasPrefix(l, "option go_package"):
				nl, ok, err := updateGoPkgOption(l, modulePath)
				if err != nil {
					return false, err
				}
				if ok {
					isModified = true
					l = nl
				}
			}

			// Write the resultant line to a temp buffer.
			fmt.Fprintln(&buf, l)
		}

		if isModified {
			if _, err := buf.WriteTo(out); err != nil {
				return false, err
			}
		}

		return isModified, nil
	}
}

// updateImport updates and returns an import statement with the new Golang
// module path. If the import statement is not updated, ok is returned as false.
func updateImport(importStmt, newModulePath string) (_ string, ok bool, _ error) {
	ss := strings.Split(importStmt, `"`)
	if len(ss) < 3 {
		return "", false, nil
	}

	np, ok, err := pathx.UpdateImportPath(newModulePath, ss[1])
	if err != nil {
		return "", false, err
	}

	if !ok {
		return "", false, nil
	}

	ss[1] = np

	return strings.Join(ss, `"`), true, nil
}

// updateGoPkgOption updates a module path in 'go_package' option. If the import
// option is not updated, ok is returned as false.
//
// See [this link](https://protobuf.dev/reference/go/go-generated/#package) for
// reference.
func updateGoPkgOption(option, newModulePath string) (_ string, ok bool, _ error) {
	ss := strings.Split(option, `"`)
	if len(ss) < 3 {
		return "", false, nil
	}

	path := ss[1]
	var alias string
	// Try to detect if there is any package alias specified.
	if pp := strings.Split(ss[1], `;`); len(pp) == 2 {
		path, alias = pp[0], pp[1]
	}

	np, ok, err := pathx.UpdateImportPath(newModulePath, path)
	if err != nil {
		return "", false, err
	}

	if !ok {
		return "", false, nil
	}

	if alias != "" {
		ss[1] = np + ";" + alias
	} else {
		ss[1] = np
	}

	return strings.Join(ss, `"`), true, nil
}
