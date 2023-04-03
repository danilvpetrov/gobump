package pathx_test

import (
	"testing"

	. "github.com/danilvpetrov/gobump/transformers/internal/pathx"
)

func TestIsPathMajor(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "valid major version suffix", in: "v2", want: true},
		{name: "valid major version suffix (significant version number)", in: "v12345", want: true},
		{name: "invalid major version suffix", in: "v1", want: false},
		{name: "invalid major version suffix (first char correct followed by a non-digit)", in: "vv", want: false},
		{name: "invalid major version suffix (short string)", in: "v", want: false},
		{name: "invalid major version suffix (long string looks like major versoin)", in: "v12345notvalid", want: false},
		{name: "invalid major version suffix (only non-digit characters)", in: "<invalid-major-version>", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPathMajor(tt.in); got != tt.want {
				t.Errorf("IsPathMajor() = %v, want %v", got, tt.want)
			}
		})
	}
}
