package versions

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_newModFiles(t *testing.T) {
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
					"fixture.com/valid",
					"fixture.com/avalid",
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
			modfiles, err := newModFiles(test.input)

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

			if !cmp.Equal(actual, test.expected.moduleNames) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(actual, test.expected.moduleNames))
			}
		})
	}
}

func Test_newModule(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Module
	}{
		{
			"Simple",
			"fixtures/new_module_simple.mod",
			Module{
				ModuleGoVersion: ModuleGoVersion{
					Name:      "fixture.com/new_module_simple",
					GoVersion: "1.13",
				},
				DependencyRequirements: map[PackageName]Package{
					"github.com/MarioCarrion/nit": {
						Name:    "github.com/MarioCarrion/nit",
						Version: "v1.23.1",
					},
					"github.com/MarioCarrion/swagger-lint": {
						Name:    "github.com/MarioCarrion/swagger-lint",
						Version: "v1.0.0",
					},
				},
			},
		},
		{
			"Replace",
			"fixtures/new_module_replace.mod",
			Module{
				ModuleGoVersion: ModuleGoVersion{
					Name:      "fixture.com/new_module_replace",
					GoVersion: "1.14",
				},
				DependencyRequirements: map[PackageName]Package{
					"github.com/MarioCarrion/nit": {
						Name:            "github.com/MarioCarrion/nit",
						Version:         "v1.23.3",
						ReplacedPath:    "replaced/MarioCarrion/nit",
						ReplacedVersion: "v9.0.0",
					},
				},
			},
		},
		{
			"Replace: not found",
			"fixtures/new_module_replace_2.mod",
			Module{
				ModuleGoVersion: ModuleGoVersion{
					Name:      "fixture.com/new_module_replace_2",
					GoVersion: "1.14",
				},
				DependencyRequirements: map[PackageName]Package{},
			},
		},
		{
			"Indirect",
			"fixtures/new_module_indirect.mod",
			Module{
				ModuleGoVersion: ModuleGoVersion{
					Name:      "fixture.com/new_module_indirect",
					GoVersion: "1.14",
				},
				DependencyRequirements: map[PackageName]Package{
					"github.com/MarioCarrion/indirect": {
						Name:       "github.com/MarioCarrion/indirect",
						Version:    "v0.0.1",
						IsIndirect: true,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			modfile, err := newModFiles([]string{test.input})
			if err != nil {
				t.Fatalf("parsing modfile %s", err)
			}

			if got := newModule(modfile[0]); !cmp.Equal(got, test.expected) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(got, test.expected))
			}
		})
	}
}
