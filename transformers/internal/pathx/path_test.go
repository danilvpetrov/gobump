package pathx_test

import (
	"testing"

	. "github.com/danilvpetrov/gobump/transformers/internal/pathx"
)

func TestUpdateImportPath(t *testing.T) {
	tests := []struct {
		name       string
		newModule  string
		importPath string
		want       string
		wantOK     bool
		wantErr    bool
	}{
		{
			name:       "v1 -> v2 (no subdirs)",
			newModule:  "example.org/foo/bar/v2",
			importPath: "example.org/foo/bar",
			want:       "example.org/foo/bar/v2",
			wantOK:     true,
		},
		{
			name:       "v1 -> v2 (with subdirs)",
			newModule:  "example.org/foo/bar/v2",
			importPath: "example.org/foo/bar/sub/dir",
			want:       "example.org/foo/bar/v2/sub/dir",
			wantOK:     true,
		},
		{
			name:       "v2 -> v3 (no subdirs)",
			newModule:  "example.org/foo/bar/v3",
			importPath: "example.org/foo/bar/v2",
			want:       "example.org/foo/bar/v3",
			wantOK:     true,
		},
		{
			name:       "v2 -> v3 (with subdirs)",
			newModule:  "example.org/foo/bar/v3",
			importPath: "example.org/foo/bar/v2/sub/dir",
			want:       "example.org/foo/bar/v3/sub/dir",
			wantOK:     true,
		},
		{
			name:       "v2 -> v1 (no subdirs)",
			newModule:  "example.org/foo/bar",
			importPath: "example.org/foo/bar/v2",
			want:       "example.org/foo/bar",
			wantOK:     true,
		},
		{
			name:       "v2 -> v1 (with subdirs)",
			newModule:  "example.org/foo/bar",
			importPath: "example.org/foo/bar/v2/sub/dir",
			want:       "example.org/foo/bar/sub/dir",
			wantOK:     true,
		},
		{
			name:       "no match",
			newModule:  "example.org/foo/bar/v2",
			importPath: "example.org/bar/foo",
			wantOK:     false,
		},
		{
			name:       "new path equal to old",
			newModule:  "example.org/foo/bar/v2",
			importPath: "example.org/foo/bar/v2",
			wantOK:     false,
		},
		{
			name:       "new module is empty",
			newModule:  "",
			importPath: "example.org/foo/bar",
			wantOK:     false,
		},
		{
			name:       "import path is empty",
			newModule:  "example.org/foo/bar/v2",
			importPath: "",
			wantOK:     false,
		},
		{
			name:       "partially matching but major versions in different positions",
			newModule:  "example.org/foo/bar/v2",
			importPath: "example.org/foo/bar/v/v2v/v3",
			wantOK:     false,
		},
		{
			name:       "invalid new module",
			newModule:  "example.org/foo/bar/v1",
			importPath: "example.org/foo/bar",
			wantErr:    true,
		},
		{
			name:       "gopkg.in/ module path v1 -> v2",
			newModule:  "gopkg.in/yaml.v2",
			importPath: "gopkg.in/yaml.v1",
			want:       "gopkg.in/yaml.v2",
			wantOK:     true,
		},
		{
			name:       "gopkg.in/ module path v2 -> v1",
			newModule:  "gopkg.in/yaml.v1",
			importPath: "gopkg.in/yaml.v2",
			want:       "gopkg.in/yaml.v1",
			wantOK:     true,
		},
		{
			name:       "gopkg.in/ no match",
			newModule:  "gopkg.in/yaml.v2",
			importPath: "example.org/foo/bar",
			wantOK:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := UpdateImportPath(tt.newModule, tt.importPath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateImportPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if ok != tt.wantOK {
				t.Fatalf("UpdateImportPath() ok = %v, wantOK %v", ok, tt.wantOK)
			}
			if got != tt.want {
				t.Fatalf("UpdateImportPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
