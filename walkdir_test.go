package gobump_test

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	. "github.com/danilvpetrov/gobump"
)

func TestWalkDir(t *testing.T) {
	tests := []struct {
		name    string
		fs      fs.FS
		f       func(file string) error
		wantErr bool
	}{
		{
			name: "should locate a correct file",
			fs:   os.DirFS("internal/testdata/walkdir/dira"),
			f: func(file string) error {
				if file != "cmd/main.go" {
					t.Errorf("WalkDir() file = %v, want %v", file, "cmd/main.go")
				}

				return nil
			},
			wantErr: false,
		},
		{
			name: "should ignore idiomatic directories",
			fs:   os.DirFS("internal/testdata/walkdir/dirb"),
			f: func(file string) error {
				t.Error("WalkDir() should not have called f()")
				return nil
			},
			wantErr: false,
		},
		{
			name: "should ignore idiomatic files",
			fs:   os.DirFS("internal/testdata/walkdir/dirc"),
			f: func(file string) error {
				t.Error("WalkDir() should not have called f()")
				return nil
			},
			wantErr: false,
		},
		{
			name: "should return error if f() returns error",
			fs:   os.DirFS("internal/testdata/walkdir/dira"),
			f: func(file string) error {
				return errors.New("<error>")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WalkDir(tt.fs, tt.f); (err != nil) != tt.wantErr {
				t.Errorf("WalkDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
