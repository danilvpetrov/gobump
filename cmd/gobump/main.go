package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/danilvpetrov/gobump"
	"github.com/danilvpetrov/gobump/transformers/gofile"
	"github.com/danilvpetrov/gobump/transformers/gomodfile"
	"github.com/danilvpetrov/gobump/transformers/protofile"
)

func run() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %s", err)
	}

	var goTidy bool
	flag.BoolVar(&goTidy, "t", false, "run 'go mod tidy' command at the end of path update")
	flag.Usage = usage
	flag.Parse()

	modules, err := parseGoModFile(wd)
	if err != nil {
		return err
	}

	newPath := flag.Arg(0)

	if err := checkModulePath(newPath, modules); err != nil {
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

	if goTidy {
		if err := runTidy(); err != nil {
			return err
		}
	}

	fmt.Println("done")
	return nil
}

func runTidy() error {
	fmt.Println("running 'go mod tidy'...")

	out, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		return fmt.Errorf("error running 'go mod tidy': %v", err)
	}

	os.Stdout.Write(out)
	return nil
}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		`
This tool allows managing the major version in the Go module paths. The module
path can be the path of the module itself or one of the module's dependencies.

usage: gobump [flags] <new go module path>

`,
	)
	flag.PrintDefaults()
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
