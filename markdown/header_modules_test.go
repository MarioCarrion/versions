package markdown

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/versions"
)

func Test_module(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    modules
		expected modules
	}{
		{
			"OK",
			modules{
				{
					Module: versions.Module{
						ModuleGoVersion: versions.ModuleGoVersion{Name: "one"},
					},
				},
				{
					Module: versions.Module{
						ModuleGoVersion: versions.ModuleGoVersion{Name: "two"}},
				},
				{
					Module: versions.Module{
						ModuleGoVersion: versions.ModuleGoVersion{Name: "abc"},
					},
				},
			},
			modules{
				{
					Module: versions.Module{
						ModuleGoVersion: versions.ModuleGoVersion{Name: "abc"},
					},
				},
				{
					Module: versions.Module{
						ModuleGoVersion: versions.ModuleGoVersion{Name: "one"},
					},
				},
				{
					Module: versions.Module{
						ModuleGoVersion: versions.ModuleGoVersion{Name: "two"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			sort.Sort(test.input)

			if !cmp.Equal(test.input, test.expected) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(test.input, test.expected))
			}
		})
	}
}

//- header

func Test_header_String(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			sorting ModulesSorting
			same    bool
			values  []versions.Module
		}
	)

	tests := []struct {
		name     string
		input    input
		expected string
	}{
		{
			"OK: ModulesSortingAsInput and Same",
			input{
				same: true,
				values: []versions.Module{
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_replace",
							GoVersion: "1.14",
						},
					},
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_indirect",
							GoVersion: "1.14",
						},
					},
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_simple",
							GoVersion: "1.14",
						},
					},
				},
			},
			`| | fixture.com/new_module_replace | fixture.com/new_module_indirect | fixture.com/new_module_simple |
| :white_check_mark: Go | 1.14 | 1.14 | 1.14 |
`,
		},
		{
			"OK: ModulesSortingAlphabetically and Different",
			input{
				sorting: ModulesSortingAlphabetically,
				values: []versions.Module{
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_replace",
							GoVersion: "1.14",
						},
					},
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_indirect",
							GoVersion: "1.15",
						},
					},
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_simple",
							GoVersion: "1.13",
						},
					},
				},
			},
			`| | fixture.com/new_module_indirect | fixture.com/new_module_replace | fixture.com/new_module_simple |
| Go | 1.15 | 1.14 | 1.13 |
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := newHeader(test.input.sorting, test.input.same, test.input.values).String()
			if got != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, got)
			}
		})
	}
}

func Test_newHeader(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			sorting ModulesSorting
			same    bool
			values  []versions.Module
		}
	)

	tests := []struct {
		name     string
		input    input
		expected header
	}{
		{
			"OK: ModulesSortingAsInput",
			input{
				same: true,
				values: []versions.Module{
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_replace",
							GoVersion: "1.14",
						},
					},
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_indirect",
							GoVersion: "1.14",
						},
					},
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_simple",
							GoVersion: "1.14",
						},
					},
				},
			},
			header{
				same: true,
				modules: modules{
					{
						Module: versions.Module{
							ModuleGoVersion: versions.ModuleGoVersion{
								Name:      "fixture.com/new_module_replace",
								GoVersion: "1.14",
							},
						},
					},
					{
						Module: versions.Module{
							ModuleGoVersion: versions.ModuleGoVersion{
								Name:      "fixture.com/new_module_indirect",
								GoVersion: "1.14",
							},
						},
					},
					{
						Module: versions.Module{
							ModuleGoVersion: versions.ModuleGoVersion{
								Name:      "fixture.com/new_module_simple",
								GoVersion: "1.14",
							},
						},
					},
				},
			},
		},
		{
			"OK: ModulesSortingAlphabetically",
			input{
				sorting: ModulesSortingAlphabetically,
				values: []versions.Module{
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_replace",
							GoVersion: "1.14",
						},
					},
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_indirect",
							GoVersion: "1.15",
						},
					},
					{
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_simple",
							GoVersion: "1.13",
						},
					},
				},
			},
			header{
				modules: modules{
					{
						Module: versions.Module{
							ModuleGoVersion: versions.ModuleGoVersion{
								Name:      "fixture.com/new_module_indirect",
								GoVersion: "1.15",
							},
						},
					},
					{
						Module: versions.Module{
							ModuleGoVersion: versions.ModuleGoVersion{
								Name:      "fixture.com/new_module_replace",
								GoVersion: "1.14",
							},
						},
					},
					{
						Module: versions.Module{
							ModuleGoVersion: versions.ModuleGoVersion{
								Name:      "fixture.com/new_module_simple",
								GoVersion: "1.13",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := newHeader(test.input.sorting, test.input.same, test.input.values)
			if !cmp.Equal(got, test.expected, cmp.AllowUnexported(header{})) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(got, test.expected, cmp.AllowUnexported(header{})))
			}
		})
	}
}
