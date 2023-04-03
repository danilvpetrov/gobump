package pathx

import (
	"fmt"
	"path"
	"strings"

	"golang.org/x/mod/module"
)

// UpdateImportPath updates Golang package import path to the new module. If no
// updates were performed to the import path, ok is returned as false and the
// new import path is returned as empty string.
//
// It returns an error if the given module path is invalid.
func UpdateImportPath(newModule, importPath string) (_ string, ok bool, _ error) {
	if newModule == "" || importPath == "" || newModule == importPath {
		return "", false, nil
	}

	pfx, vp, ok := module.SplitPathVersion(newModule)
	if !ok {
		return "", false, fmt.Errorf("module path %s is invalid", newModule)
	}

	if !strings.HasPrefix(importPath, pfx) {
		// The targeted import path is not related to the new model.
		// Return the import path unchanged.
		return "", false, nil
	}

	if strings.HasPrefix(newModule, "gopkg.in/") {
		// gopkg.in paths are treated differently
		return newModule, true, nil
	}

	ee := pathElementsAfterPrefix(pfx, importPath)
	if len(ee) > 0 {
		if IsPathMajor(ee[0]) {
			ee[0] = vp
		} else {
			for _, el := range ee[1:] {
				// If the position of the major version in the original import
				// path does not match to the path, bail out.
				if IsPathMajor(el) {
					return "", false, nil
				}
			}

			ee = append([]string{vp}, ee...)
		}
	} else {
		ee = []string{vp}
	}

	ee = append([]string{pfx}, ee...)

	return path.Join(ee...), true, nil
}

func pathElementsAfterPrefix(prefix, path string) []string {
	pfxEls, pathEls := strings.Split(prefix, "/"),
		strings.Split(path, "/")

	return pathEls[len(pfxEls):]
}
