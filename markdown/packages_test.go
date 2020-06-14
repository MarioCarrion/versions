package markdown

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/versions"
)

func Test_newPackages(t *testing.T) {
	t.Parallel()

	type (
		inputPkg struct {
			name versions.ModuleName
			pkg  versions.Package
		}

		inputModules struct {
			values       []module
			dependencies map[versions.ModuleName]versions.Module
		}

		input struct {
			sorting  PackagesSorting
			packages []inputPkg
			modules  inputModules
		}
	)

	newInput := func(sorting PackagesSorting) input {
		return input{
			sorting: sorting,
			modules: inputModules{
				values: []module{
					{
						Module: versions.Module{
							ModuleGoVersion: versions.ModuleGoVersion{Name: "Module1"},
						},
					},
					{
						Module: versions.Module{
							ModuleGoVersion: versions.ModuleGoVersion{Name: "Module2"}},
					},
				},
				dependencies: map[versions.ModuleName]versions.Module{
					"Module1": {
						DependencyRequirements: map[versions.PackageName]versions.Package{
							"pkg1": {
								Name:    "pkg1",
								Version: "v1",
							},
							"abc": {
								Name:    "abc",
								Version: "v1",
							},
							"diff": {
								Name:    "diff",
								Version: "v2",
							},
							"adiff": {
								Name:    "adiff",
								Version: "v2",
							},
						},
					},
					"Module2": {
						DependencyRequirements: map[versions.PackageName]versions.Package{
							"pkg1": {
								Name:    "pkg1",
								Version: "v1",
							},
							"diff": {
								Name:    "diff",
								Version: "v1",
							},
							"adiff": {
								Name:    "adiff",
								Version: "v1",
							},
						},
					},
				},
			},
			packages: []inputPkg{
				{
					"Module1",
					versions.Package{
						Name:    "pkg1",
						Version: "v1",
					},
				},
				{
					"Module1",
					versions.Package{
						Name:    "abc",
						Version: "v1",
					},
				},
				{
					"Module1",
					versions.Package{
						Name:    "diff",
						Version: "v2",
					},
				},
				{
					"Module1",
					versions.Package{
						Name:    "adiff",
						Version: "v2",
					},
				},
				{
					"Module2",
					versions.Package{
						Name:    "pkg1",
						Version: "v1",
					},
				},
				{
					"Module2",
					versions.Package{
						Name:    "diff",
						Version: "v1",
					},
				},
				{
					"Module2",
					versions.Package{
						Name:    "adiff",
						Version: "v1",
					},
				},
			},
		}
	}

	tests := []struct {
		name     string
		input    input
		expected string
	}{
		{
			"OK: PackagesSortingAsFound",
			newInput(PackagesSortingAsFound),
			`| :white_check_mark: pkg1 | v1 | v1 |
| :white_check_mark: abc | v1 |  |
| diff | v2 | v1 |
| adiff | v2 | v1 |
`,
		},
		{
			"OK: PackagesSortingAlphabeticallySupported",
			newInput(PackagesSortingAlphabeticallySupported),
			`| :white_check_mark: abc | v1 |  |
| :white_check_mark: pkg1 | v1 | v1 |
| adiff | v2 | v1 |
| diff | v2 | v1 |
`,
		},
		{
			"OK: PackagesSortingAlphabetically",
			newInput(PackagesSortingAlphabetically),
			`| :white_check_mark: abc | v1 |  |
| adiff | v2 | v1 |
| diff | v2 | v1 |
| :white_check_mark: pkg1 | v1 | v1 |
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			packages := versions.Packages{}
			for _, p := range test.input.packages {
				packages.Set(p.name, p.pkg)
			}

			versions := versions.Versions{
				Packages: packages,
				Modules:  test.input.modules.dependencies,
			}

			pkgs := newPackages(versions, test.input.modules.values, test.input.sorting)
			if got := pkgs.String(); !cmp.Equal(got, test.expected) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(got, test.expected))
			}
		})
	}
}

func Test_packageSets_Sort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    packageSets
		expected packageSets
	}{
		{
			"OK",
			[]packageSet{
				{
					Name: "bbbb",
				},
				{
					Name: "aaaa",
				},
			},
			[]packageSet{
				{
					Name: "aaaa",
				},
				{
					Name: "bbbb",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			sort.Sort(test.input)

			if !cmp.Equal(test.input, test.expected, cmp.AllowUnexported(packageSet{})) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(test.input, test.expected, cmp.AllowUnexported(packageSet{})))
			}
		})
	}
}

func Test_packageSets_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    packageSets
		expected string
	}{
		{
			"OK",
			[]packageSet{
				{
					same: true,
					Name: "name",
					packages: []versions.Package{
						{
							Name:            "name",
							Version:         "v1",
							IsIndirect:      true,
							ReplacedPath:    "rpath",
							ReplacedVersion: "rversion",
						},
						{
							Name:            "name",
							Version:         "v2",
							IsIndirect:      true,
							ReplacedPath:    "rpath",
							ReplacedVersion: "rversion",
						},
					},
				},
				{
					same: false,
					Name: "another",
					packages: []versions.Package{
						{
							Name:       "another",
							Version:    "v1",
							IsIndirect: true,
						},
						{
							Name:            "another",
							Version:         "v2",
							IsIndirect:      true,
							ReplacedPath:    "rpath",
							ReplacedVersion: "rversion",
						},
					},
				},
			},
			`| :white_check_mark: name | v1/rversion | v2/rversion |
| another | v1 | v2/rversion |
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.input.String(); !cmp.Equal(got, test.expected) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(got, test.expected))
			}
		})
	}
}
