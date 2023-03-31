package gomodfile_test

import (
	"bytes"
	"testing"

	. "github.com/danilvpetrov/gobump/transformers/gomodfile"
)

func TestUpdateModulePath(t *testing.T) {
	tests := []struct {
		name       string
		modulePath string
		modfile    string
		wantOk     bool
		wantErr    bool
		wantOut    string
	}{
		{
			name:       "should update module path to a newer version",
			modulePath: "example.com/foo/bar/v2",
			modfile: `module example.com/foo/bar

go 1.20
`,
			wantOk: true,
			wantOut: `module example.com/foo/bar/v2

go 1.20
`,
		},
		{
			name:       "should update module path to an older version",
			modulePath: "example.com/foo/bar",
			modfile: `module example.com/foo/bar/v2

go 1.20
`,
			wantOk: true,
			wantOut: `module example.com/foo/bar

go 1.20
`,
		},
		{
			name:       "should not update non-matchin module path",
			modulePath: "example.com/bar/foo",
			modfile: `module example.com/foo/bar

go 1.20
`,
			wantOk: false,
		},
		{
			name:       "should not update partially matching module path",
			modulePath: "example.com/foo/bar/v2",
			modfile: `module example.com/foo/bar/foobar

go 1.20
`,
			wantOk: false,
		},
		{
			name:    "should return an error if a go.mod file is invalid",
			modfile: "<invalid-gomod-file>",
			wantErr: true,
		},
		{
			name:       "should return an error if the module path is invalid",
			modulePath: "example.com/bar/foo/v1",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := bytes.NewBufferString(tt.modfile), &bytes.Buffer{}

			ok, err := UpdateModulePath(tt.modulePath)(r, w)
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
