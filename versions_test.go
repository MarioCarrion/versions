package versions_test

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/versions"
)

func Test_SameVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected bool
	}{
		{
			"All present",
			[]string{"v0.3.1", "v0.3.1", "v0.3.1"},
			true,
		},
		{
			"No values",
			[]string{},
			true,
		},
		{
			"Index Zero blank",
			[]string{"", "v0.0.0-20160331181800-b5bfa59ec0ad", "v0.0.0-20160331181800-b5bfa59ec0ad"},
			true,
		},
		{
			"Only one value #1",
			[]string{"", "v0.21.0", ""},
			true,
		},
		{
			"Only one value #2",
			[]string{"", "", "v1.4.1"},
			true,
		},
		{
			"One different",
			[]string{"", "v1.0.0", "v5.0.69"},
			false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			actual := versions.SameVersion(test.input)

			if actual != test.expected {
				t.Fatalf("expected %t, actual %t", test.expected, actual)
			}
		})
	}
}

func Test_Versions(t *testing.T) {
	tests := []struct {
		name     string
		input    versions.Versions
		expected versions.Versions
	}{
		{
			"OK",
			versions.Versions{
				{Name: "bbb"},
				{Name: "aaa"},
			},
			versions.Versions{
				{Name: "aaa"},
				{Name: "bbb"},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			sort.Sort(test.input)
			if !cmp.Equal(test.input, test.expected) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(test.input, test.expected))
			}
		})
	}
}
