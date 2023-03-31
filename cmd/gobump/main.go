package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
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

	if err := walkDir(wd, newPath); err != nil {
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
	fmt.Fprintf(os.Stderr, "usage: gobump [flags] <new go module path>\n")
	flag.PrintDefaults()
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
