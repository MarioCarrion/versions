package markdown

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/senseyeio/diligent"

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
			sorting     PackagesSorting
			packages    []inputPkg
			modules     inputModules
			showLicense bool
		}
	)

	newInput := func(sorting PackagesSorting, showLicense bool) input {
		return input{
			sorting:     sorting,
			showLicense: showLicense,
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
								Name:         "pkg1",
								Version:      "v1",
								ReplacedPath: "fixtures/license/valid",
							},
							"abc": {
								Name:    "abc",
								Version: "v1",
								License: versions.License{
									Identifier: "Test",
									Name:       "LicenseName",
									Category:   diligent.Permissive,
								},
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
								Name:         "pkg1",
								Version:      "v1",
								ReplacedPath: "fixtures/license/valid",
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
						Name:         "pkg1",
						Version:      "v1",
						ReplacedPath: "fixtures/license/valid",
					},
				},
				{
					"Module1",
					versions.Package{
						Name:    "abc",
						Version: "v1",
						License: versions.License{
							Identifier: "Test",
							Name:       "LicenseName",
							Category:   diligent.Permissive,
						},
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
						Name:         "pkg1",
						Version:      "v1",
						ReplacedPath: "fixtures/license/valid",
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
		expected [][]string
	}{
		{
			"OK: PackagesSortingAsFound",
			newInput(PackagesSortingAsFound, false),
			[][]string{
				{":white_check_mark: pkg1", "v1 fixtures/license/valid", "v1 fixtures/license/valid"},
				{":white_check_mark: abc", "v1", ""},
				{"diff", "v2", "v1"},
				{"adiff", "v2", "v1"},
			},
		},
		{
			"OK: PackagesSortingAsFound with License",
			newInput(PackagesSortingAsFound, true),
			[][]string{
				{":white_check_mark: pkg1", "v1 fixtures/license/valid", "v1 fixtures/license/valid"},
				{":white_check_mark: abc", "v1<br>permissive LicenseName", ""},
				{"diff", "v2", "v1"},
				{"adiff", "v2", "v1"},
			},
		},
		{
			"OK: PackagesSortingAlphabeticallySupported",
			newInput(PackagesSortingAlphabeticallySupported, false),
			[][]string{
				{":white_check_mark: abc", "v1", ""},
				{":white_check_mark: pkg1", "v1 fixtures/license/valid", "v1 fixtures/license/valid"},
				{"adiff", "v2", "v1"},
				{"diff", "v2", "v1"},
			},
		},
		{
			"OK: PackagesSortingAlphabetically",
			newInput(PackagesSortingAlphabetically, false),
			[][]string{
				{":white_check_mark: abc", "v1", ""},
				{"adiff", "v2", "v1"},
				{"diff", "v2", "v1"},
				{":white_check_mark: pkg1", "v1 fixtures/license/valid", "v1 fixtures/license/valid"},
			},
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

			pkgs := newPackages(versions, test.input.modules.values, test.input.sorting, test.input.showLicense)
			if got := pkgs.Values(); !cmp.Equal(got, test.expected) {
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
