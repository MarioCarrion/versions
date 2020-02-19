package versions_test

import (
	"reflect"
	"testing"

	"github.com/MarioCarrion/versions"
)

func Test_NewModFiles(t *testing.T) {
	type expected struct {
		withErr     bool
		moduleNames []string
	}

	tests := []struct {
		name     string
		input    []string
		expected expected
	}{
		{
			"Valid",
			[]string{"./fixtures/valid.mod"},
			expected{
				moduleNames: []string{"fixture.com/valid"},
			},
		},
		{
			"Sorted",
			[]string{
				"./fixtures/valid.mod",
				"./fixtures/avalid.mod",
			},
			expected{
				moduleNames: []string{
					"fixture.com/avalid",
					"fixture.com/valid",
				},
			},
		},
		{
			"Invalid: not found",
			[]string{
				"./does/not/exist.mod",
			},
			expected{
				withErr: true,
			},
		},
		{
			"Invalid: parse error",
			[]string{
				"./fixtures/invalid.mod",
			},
			expected{
				withErr: true,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			modfiles, err := versions.NewModFiles(test.input)

			if test.expected.withErr == (err == nil) {
				t.Fatalf("expected error: %t, got %s", test.expected.withErr, err)
			}

			if err != nil {
				return
			}

			actual := make([]string, len(modfiles))
			for i, modfile := range modfiles {
				actual[i] = modfile.Module.Mod.Path
			}

			if !reflect.DeepEqual(actual, test.expected.moduleNames) {
				t.Fatalf("expected: %+v, actual: %+v\n", test.expected.moduleNames, actual)
			}
		})
	}
}
