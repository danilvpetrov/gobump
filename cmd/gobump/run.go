package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/danilvpetrov/gobump"
	"github.com/danilvpetrov/gobump/transformers/gofile"
	"github.com/danilvpetrov/gobump/transformers/gomodfile"
	"github.com/danilvpetrov/gobump/transformers/protofile"
	"golang.org/x/mod/module"
)

func run() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %s", err)
	}

	var noGoGet bool
	flag.BoolVar(&noGoGet, "n", false, "don't run 'go get' for the new module path")
	flag.Usage = usage
	flag.Parse()

	newPath := flag.Arg(0)

	if err := checkPath(wd, newPath); err != nil {
		return err
	}

	if err := gobump.WalkDir(
		os.DirFS(wd),
		func(path string) error {
			switch {
			default:
				return nil
			case filepath.Base(path) == "go.mod":
				return runTransformers(path, gomodfile.UpdateModulePath(newPath))
			case filepath.Ext(path) == ".go":
				return runTransformers(path, gofile.UpdateImports(newPath))
			case filepath.Ext(path) == ".proto":
				return runTransformers(path, protofile.UpdateModulePath(newPath))
			}
		},
	); err != nil {
		return err
	}

	if !noGoGet {
		ok, err := shouldRunGoGet(wd, newPath)
		if err != nil {
			return err
		}

		if ok {
			if err := runGoGet(newPath); err != nil {
				return err
			}

			if err := runGoModTidy(); err != nil {
				return err
			}
		}
	}

	fmt.Println("done")

	return nil
}

func checkPath(wd, path string) error {
	if err := module.CheckPath(path); err != nil {
		return fmt.Errorf("invalid module path %q: %w", path, err)
	}

	mp, dr, err := gobump.ParseModules(wd)
	if err != nil {
		return err
	}

	pn, po := modulePrefix(mp), modulePrefix(path)
	if pn == po {
		return nil
	}

	for _, m := range dr {
		if pd := modulePrefix(m); pd == po {
			return nil
		}
	}

	return fmt.Errorf(
		"module path '%s' does not match module '%s' or any of its direct dependencies",
		path,
		mp,
	)
}

func runGoGet(module string) error {
	args := []string{"go", "get", fmt.Sprintf("%s@latest", module)}

	fmt.Printf("running '%s'...", strings.Join(args, " "))

	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		return fmt.Errorf("error running '%s': %v", strings.Join(args, " "), err)
	}

	os.Stdout.Write(out)
	return nil
}

func runGoModTidy() error {
	fmt.Println("running 'go mod tidy'...")

	out, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		return fmt.Errorf("error running 'go mod tidy': %v", err)
	}

	os.Stdout.Write(out)
	return nil
}

func shouldRunGoGet(wd, path string) (bool, error) {
	mp, _, err := gobump.ParseModules(wd)
	if err != nil {
		return false, err
	}

	if pn, po := modulePrefix(mp), modulePrefix(path); pn == po {
		return false, nil
	}

	// The module path is assumed to be a direct dependency of the module.
	return true, nil
}

func modulePrefix(path string) string {
	pfx, _, ok := module.SplitPathVersion(path)
	if !ok {
		return ""
	}

	return pfx
}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		`
This tool allows managing the major version in the Go module paths. The module
path can be the path of the module itself or one of the module's direct dependencies.

usage: gobump [flags] <new go module path>

`,
	)
	flag.PrintDefaults()
}
