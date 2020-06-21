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

func Test_header(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			sorting ModulesSorting
			same    bool
			values  []versions.Module
		}

		expected struct {
			goVersions []string
			names      []string
		}
	)

	tests := []struct {
		name     string
		input    input
		expected expected
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
			expected{
				goVersions: []string{":white_check_mark: Go", "1.14", "1.14", "1.14"},
				names:      []string{"", "fixture.com/new_module_replace", "fixture.com/new_module_indirect", "fixture.com/new_module_simple"},
			},
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
			expected{
				goVersions: []string{"Go", "1.15", "1.14", "1.13"},
				names:      []string{"", "fixture.com/new_module_indirect", "fixture.com/new_module_replace", "fixture.com/new_module_simple"},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			header := newHeader(test.input.sorting, test.input.same, test.input.values)

			goVersions := header.GoVersions()
			if !cmp.Equal(goVersions, test.expected.goVersions) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(goVersions, test.expected.goVersions))
			}

			names := header.Names()
			if !cmp.Equal(names, test.expected.names) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(names, test.expected.names))
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
