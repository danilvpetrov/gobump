package gobump_test

import (
	"reflect"
	"testing"

	. "github.com/danilvpetrov/gobump"
)

func TestParseModules(t *testing.T) {
	tests := []struct {
		name               string
		moduleDir          string
		wantModulePath     string
		wantDirectRequires []string
		wantErr            bool
	}{
		{
			name:           "should parse valid go.mod correctly",
			moduleDir:      "internal/testdata/modules/dira",
			wantModulePath: "example.com/foo/bar",
			wantDirectRequires: []string{
				"example.com/foobar/foobar",
				"example.com/barfoo/barfoo",
			},
		},
		{
			name:      "should return error if go.mod is not found",
			moduleDir: "internal/testdata/modules/dirb",
			wantErr:   true,
		},
		{
			name:      "should return error if go.mod is invalid",
			moduleDir: "internal/testdata/modules/dirc",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseModules(tt.moduleDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseModules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.wantModulePath {
				t.Errorf("ParseModules() got = %v, want %v", got, tt.wantModulePath)
			}

			if !reflect.DeepEqual(got1, tt.wantDirectRequires) {
				t.Errorf("ParseModules() got1 = %v, want %v", got1, tt.wantDirectRequires)
			}
		})
	}
}
