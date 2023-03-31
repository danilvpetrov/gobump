package gofile_test

import (
	"bytes"
	"testing"

	. "github.com/danilvpetrov/gobump/transformers/gofile"
)

func TestUpdateImports(t *testing.T) {
	tests := []struct {
		name          string
		gofile        string
		newImportPath string
		wantOk        bool
		wantErr       bool
		wantOut       string
	}{
		{
			name: "updates .go file import path to a newer version",
			gofile: `package main

import (
	"example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			wantOk:        true,
			newImportPath: "example.org/foo/bar/v2",
			wantOut: `package main

import (
	"example.org/foo/bar/v2"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
		},
		{
			name: "updates .go file import path to an older version",
			gofile: `package main

import (
	"example.org/foo/bar/v2"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			wantOk:        true,
			newImportPath: "example.org/foo/bar",
			wantOut: `package main

import (
	"example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
		},
		{
			name: "updates .go file named import path to a newer version",
			gofile: `package main

import (
	foobar "example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			wantOk:        true,
			newImportPath: "example.org/foo/bar/v2",
			wantOut: `package main

import (
	foobar "example.org/foo/bar/v2"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
		},
		{
			name: "updates .go file named import path to an older version",
			gofile: `package main

import (
	foobar "example.org/foo/bar/v2"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			wantOk:        true,
			newImportPath: "example.org/foo/bar",
			wantOut: `package main

import (
	foobar "example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
		},
		{
			name: "updates .go file dot import path to a newer version",
			gofile: `package main

import (
	. "example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			wantOk:        true,
			newImportPath: "example.org/foo/bar/v2",
			wantOut: `package main

import (
	. "example.org/foo/bar/v2"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
		},
		{
			name: "updates .go file dot import path to an older version",
			gofile: `package main

import (
	. "example.org/foo/bar/v2"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			wantOk:        true,
			newImportPath: "example.org/foo/bar",
			wantOut: `package main

import (
	. "example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
		},
		{
			name: "updates .go file underscore import path to a newer version",
			gofile: `package main

import (
	_ "example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			wantOk:        true,
			newImportPath: "example.org/foo/bar/v2",
			wantOut: `package main

import (
	_ "example.org/foo/bar/v2"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
		},
		{
			name: "updates .go file underscore import path to an older version",
			gofile: `package main

import (
	_ "example.org/foo/bar/v2"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			wantOk:        true,
			newImportPath: "example.org/foo/bar",
			wantOut: `package main

import (
	_ "example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
		},
		{
			name:    "returns an error if .go file is not valid",
			gofile:  "<invalid-go-file>",
			wantErr: true,
		},
		{
			name: "returns an error if a new import path is invalid",
			gofile: `package main

import (
	foobar "example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			newImportPath: "example.org/foo/bar/v1",
			wantErr:       true,
		},
		{
			name: "returns ok as false if import path is not matching",
			gofile: `package main

import (
	foobar "example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			newImportPath: "example.org/bar/foo/v2",
			wantOk:        false,
		},
		{
			name: "returns ok as false if import path is the same",
			gofile: `package main

import (
	foobar "example.org/foo/bar"
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`,
			newImportPath: "example.org/bar/foo",
			wantOk:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := bytes.NewBufferString(tt.gofile), &bytes.Buffer{}

			ok, err := UpdateImports(tt.newImportPath)(r, w)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateModulePath() error = %v, wantErr %v", err, tt.wantErr)
			}

			if ok != tt.wantOk {
				t.Fatalf("UpdateModulePath() ok = %v, wantOk %v", ok, tt.wantOk)
			}

			if out := w.String(); out != tt.wantOut {
				t.Fatalf("UpdateModulePath() out = %s, wantOut %s", out, tt.wantOut)
			}
		})
	}
}
