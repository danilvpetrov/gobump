package protofile_test

import (
	"bytes"
	"testing"

	. "github.com/danilvpetrov/gobump/transformers/protofile"
)

func TestUpdateModulePath(t *testing.T) {
	tests := []struct {
		name       string
		modulePath string
		protofile  string
		wantOut    string
		wantOk     bool
		wantErr    bool
	}{
		{
			name:       "should update module import path to a newer version",
			modulePath: "example.org/foo/bar/v2",
			protofile: `syntax = "proto3";
package foobar;

import "example.org/foo/bar/blah/some.proto";

message FooBar{
	string foo = 1;
}
`,
			wantOut: `syntax = "proto3";
package foobar;

import "example.org/foo/bar/v2/blah/some.proto";

message FooBar{
	string foo = 1;
}
`,
			wantOk: true,
		},
		{
			name:       "should update module import path to an older version",
			modulePath: "example.org/foo/bar",
			protofile: `syntax = "proto3";
package foobar;

import "example.org/foo/bar/v2/blah/some.proto";

message FooBar{
	string foo = 1;
}
`,
			wantOut: `syntax = "proto3";
package foobar;

import "example.org/foo/bar/blah/some.proto";

message FooBar{
	string foo = 1;
}
`,
			wantOk: true,
		},
		{

			name:       "should return an error if a new import module path is invalid",
			modulePath: "example.org/foo/bar/v1",
			protofile: `syntax = "proto3";
package foobar;

import "example.org/bar/foo/blah/some.proto";

message FooBar{
string foo = 1;
}
`,
			wantErr: true,
		},
		{

			name:       "should not update a partially matching import path",
			modulePath: "example.org/foo/bar/v3",
			protofile: `syntax = "proto3";
package foobar;

import "example.org/bar/foo/blah/v2/some.proto";

message FooBar{
string foo = 1;
}
`,
			wantOk: false,
		},
		{

			name:       "should not update a non-matching import path",
			modulePath: "example.org/foo/bar/v2",
			protofile: `syntax = "proto3";
package foobar;

import "example.org/bar/foo/blah/some.proto";

message FooBar{
string foo = 1;
}
`,
			wantOk: false,
		},
		{
			name:       "should update go package option to a new version",
			modulePath: "example.org/foo/bar/v2",
			protofile: `syntax = "proto3";
package foobar;

option go_package = "example.org/foo/bar";

message FooBar{
	string foo = 1;
}
`,
			wantOut: `syntax = "proto3";
package foobar;

option go_package = "example.org/foo/bar/v2";

message FooBar{
	string foo = 1;
}
`,
			wantOk: true,
		},
		{
			name:       "should update go package option to an older version",
			modulePath: "example.org/foo/bar",
			protofile: `syntax = "proto3";
package foobar;

option go_package = "example.org/foo/bar/v2";

message FooBar{
	string foo = 1;
}
`,
			wantOut: `syntax = "proto3";
package foobar;

option go_package = "example.org/foo/bar";

message FooBar{
	string foo = 1;
}
`,
			wantOk: true,
		},
		{
			name:       "should update go package option with alias to a new version",
			modulePath: "example.org/foo/bar/v2",
			protofile: `syntax = "proto3";
package foobar;

option go_package = "example.org/foo/bar;foobar";

message FooBar{
	string foo = 1;
}
`,
			wantOut: `syntax = "proto3";
package foobar;

option go_package = "example.org/foo/bar/v2;foobar";

message FooBar{
	string foo = 1;
}
`,
			wantOk: true,
		},
		{
			name:       "should update go package option with alias to an older version",
			modulePath: "example.org/foo/bar",
			protofile: `syntax = "proto3";
package foobar;

option go_package = "example.org/foo/bar/v2;foobar";

message FooBar{
	string foo = 1;
}
`,
			wantOut: `syntax = "proto3";
package foobar;

option go_package = "example.org/foo/bar;foobar";

message FooBar{
	string foo = 1;
}
`,
			wantOk: true,
		},
		{

			name:       "should return an error if a new module path for go package option is invalid",
			modulePath: "example.org/foo/bar/v1",
			protofile: `syntax = "proto3";
package foobar;

option go_package = "example.org/foo/bar";

message FooBar{
string foo = 1;
}
`,
			wantErr: true,
		},
		{

			name:       "should not update a non-matching import path for go package option",
			modulePath: "example.org/foo/bar/v2",
			protofile: `syntax = "proto3";
package foobar;

option go_package = "example.org/bar/foo";

message FooBar{
string foo = 1;
}
`,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := bytes.NewBufferString(tt.protofile), &bytes.Buffer{}

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
