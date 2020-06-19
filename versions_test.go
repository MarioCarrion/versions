package versions_test

import (
	"fmt"
	"go/build"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/versions"
)

func Test_GoVersions(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			name    versions.ModuleName
			version versions.GoVersion
		}

		expected struct {
			isSame bool
			values []versions.ModuleGoVersion
		}
	)

	tests := []struct {
		name     string
		input    []input
		expected expected
	}{
		{
			"OK",
			[]input{
				{
					versions.ModuleName("Name"),
					versions.GoVersion("Version"),
				},
			},
			expected{
				true,
				[]versions.ModuleGoVersion{
					{
						Name:      versions.ModuleName("Name"),
						GoVersion: versions.GoVersion("Version"),
					},
				},
			},
		},
		{
			"NotSame",
			[]input{
				{
					versions.ModuleName("Name1"),
					versions.GoVersion("Version1"),
				},
				{
					versions.ModuleName("Name2"),
					versions.GoVersion("Version2"),
				},
				{
					versions.ModuleName("Name3"),
					versions.GoVersion("Version"),
				},
			},
			expected{
				false,
				[]versions.ModuleGoVersion{
					{
						Name:      versions.ModuleName("Name1"),
						GoVersion: versions.GoVersion("Version1"),
					},
					{
						Name:      versions.ModuleName("Name2"),
						GoVersion: versions.GoVersion("Version2"),
					},
					{
						Name:      versions.ModuleName("Name3"),
						GoVersion: versions.GoVersion("Version"),
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			goversions := versions.GoVersions{}

			for _, input := range test.input {
				goversions.Set(input.name, input.version)
			}

			if got := goversions.IsSame(); got != test.expected.isSame {
				t.Fatalf("expected %T, got %T", test.expected, got)
			}

			if values := goversions.Values(); !cmp.Equal(values, test.expected.values) {
				t.Fatalf("expected values do not match: %s", cmp.Diff(values, test.expected.values))
			}
		})
	}
}

func Test_New(t *testing.T) {
	t.Parallel()

	type (
		expected struct {
			modules    map[versions.ModuleName]versions.Module
			goVersions []versions.ModuleGoVersion
			packages   map[versions.PackageName]map[versions.ModuleName]versions.Package
			err        bool
		}
	)

	tests := []struct {
		name     string
		input    []string
		expected expected
	}{
		{
			"Invalid mod",
			[]string{"fixtures/does_not_exist.mod"},
			expected{
				err:        true,
				goVersions: []versions.ModuleGoVersion{},
			},
		},
		{
			"OK",
			[]string{
				"fixtures/new_module_indirect.mod",
				"fixtures/new_module_replace.mod",
				"fixtures/new_module_simple.mod",
			},
			expected{
				modules: map[versions.ModuleName]versions.Module{
					"fixture.com/new_module_indirect": {
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_indirect",
							GoVersion: "1.14",
						},
						DependencyRequirements: map[versions.PackageName]versions.Package{
							"github.com/MarioCarrion/indirect": {
								Name:       "github.com/MarioCarrion/indirect",
								Version:    "v0.0.1",
								IsIndirect: true,
							},
						},
					},
					"fixture.com/new_module_replace": {
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_replace",
							GoVersion: "1.14",
						},
						DependencyRequirements: map[versions.PackageName]versions.Package{
							"github.com/MarioCarrion/nit": {
								Name:            "github.com/MarioCarrion/nit",
								Version:         "v1.23.3",
								ReplacedPath:    "replaced/MarioCarrion/nit",
								ReplacedVersion: "v9.0.0",
							},
						},
					},
					"fixture.com/new_module_simple": {
						ModuleGoVersion: versions.ModuleGoVersion{
							Name:      "fixture.com/new_module_simple",
							GoVersion: "1.13",
						},
						DependencyRequirements: map[versions.PackageName]versions.Package{
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
				goVersions: []versions.ModuleGoVersion{
					{
						Name:      "fixture.com/new_module_indirect",
						GoVersion: "1.14",
					},
					{
						Name:      "fixture.com/new_module_replace",
						GoVersion: "1.14",
					},
					{
						Name:      "fixture.com/new_module_simple",
						GoVersion: "1.13",
					},
				},
				packages: map[versions.PackageName]map[versions.ModuleName]versions.Package{
					"github.com/MarioCarrion/indirect": {
						"fixture.com/new_module_indirect": {
							Name:       "github.com/MarioCarrion/indirect",
							Version:    "v0.0.1",
							IsIndirect: true,
						},
					},
					"github.com/MarioCarrion/nit": {
						"fixture.com/new_module_replace": {
							Name:            "github.com/MarioCarrion/nit",
							Version:         "v1.23.3",
							ReplacedPath:    "replaced/MarioCarrion/nit",
							ReplacedVersion: "v9.0.0",
						},
						"fixture.com/new_module_simple": {
							Name:    "github.com/MarioCarrion/nit",
							Version: "v1.23.1",
						},
					},
					"github.com/MarioCarrion/swagger-lint": {
						"fixture.com/new_module_simple": {
							Name:    "github.com/MarioCarrion/swagger-lint",
							Version: "v1.0.0",
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

			got, err := versions.New(test.input)
			if (err != nil) != test.expected.err {
				t.Fatalf("expected error %t, got %t", test.expected.err, err != nil)
			}

			if !cmp.Equal(got.Modules, test.expected.modules) {
				t.Fatalf("expected modules do not match: %s", cmp.Diff(got, test.expected.modules))
			}

			if goversions := got.GoVersions.Values(); !cmp.Equal(goversions, test.expected.goVersions) {
				t.Fatalf("expected goversions do not match: %s", cmp.Diff(goversions, test.expected.goVersions))
			}

			for pkg, expectedPkg := range test.expected.packages {
				if packages := got.Packages.Values(pkg); !cmp.Equal(packages, expectedPkg) {
					t.Fatalf("expected goversions do not match: %s", cmp.Diff(packages, expectedPkg))
				}
			}
		})
	}
}

func Test_Package(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    versions.Package
		setup    func() func()
		expected string
	}{
		{
			"Yes GOPATH + Yes Version + Yes ReplacedPath + Yes ReplacedVersion",
			versions.Package{
				Name:            "Name",
				Version:         "Version",
				ReplacedPath:    "NewPath",
				ReplacedVersion: "NewVersion",
			},
			func() func() {
				old := os.Getenv("GOPATH")
				os.Setenv("GOPATH", "/blah")
				return func() {
					os.Setenv("GOPATH", old)
				}
			},
			"/blah/pkg/mod/NewPath@NewVersion",
		},
		{
			"Yes GOPATH + Yes Version + Yes ReplacedPath + No ReplacedVersion",
			versions.Package{
				Name:         "Name",
				Version:      "Version",
				ReplacedPath: "NewPath",
			},
			func() func() {
				old := os.Getenv("GOPATH")
				os.Setenv("GOPATH", "/blah")
				return func() {
					os.Setenv("GOPATH", old)
				}
			},
			"/blah/pkg/mod/NewPath",
		},
		{
			"Yes GOPATH + Yes Version + No ReplacedPath + Yes ReplacedVersion",
			versions.Package{
				Name:            "Name",
				Version:         "Version",
				ReplacedVersion: "NewVersion",
			},
			func() func() {
				old := os.Getenv("GOPATH")
				os.Setenv("GOPATH", "/blah")
				return func() {
					os.Setenv("GOPATH", old)
				}
			},
			"/blah/pkg/mod/Name@Version",
		},
		{
			"Yes GOPATH + No Version + Yes ReplacedPath + Yes ReplacedVersion",
			versions.Package{
				Name:            "Name",
				ReplacedPath:    "NewPath",
				ReplacedVersion: "NewVersion",
			},
			func() func() {
				old := os.Getenv("GOPATH")
				os.Setenv("GOPATH", "/blah")
				return func() {
					os.Setenv("GOPATH", old)
				}
			},
			"/blah/pkg/mod/NewPath@NewVersion",
		},
		{
			"Yes GOPATH + Yes Version + No ReplacedPath + Yes ReplacedVersion",
			versions.Package{
				Name:            "Name",
				Version:         "Version",
				ReplacedVersion: "NewVersion",
			},
			func() func() {
				old := os.Getenv("GOPATH")
				os.Setenv("GOPATH", "/blah")
				return func() {
					os.Setenv("GOPATH", old)
				}
			},
			"/blah/pkg/mod/Name@Version",
		},
		{
			"No GOPATH",
			versions.Package{
				Name:    "Name",
				Version: "Version",
			},
			func() func() {
				return func() {}
			},
			fmt.Sprintf("%s/pkg/mod/Name@Version", build.Default.GOPATH),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			teardown := test.setup()
			defer teardown()
			if got := test.input.Path(); got != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, got)
			}
		})
	}
}

func Test_Packages(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			name versions.ModuleName
			pkg  versions.Package
		}

		expected struct {
			isSame map[versions.PackageName]bool
			values map[versions.PackageName]map[versions.ModuleName]versions.Package
			names  []versions.PackageName
		}
	)

	tests := []struct {
		name     string
		input    []input
		expected expected
	}{
		{
			"All the same",
			[]input{
				{
					"Module1",
					versions.Package{
						Name:            "pkg1",
						Version:         "v1",
						IsIndirect:      true,
						ReplacedPath:    "pkg2",
						ReplacedVersion: "v2",
					},
				},
				{
					"Module2",
					versions.Package{
						Name:            "pkg1",
						Version:         "v1",
						IsIndirect:      true,
						ReplacedPath:    "pkg2",
						ReplacedVersion: "v2",
					},
				},
			},
			expected{
				map[versions.PackageName]bool{
					"pkg1": true,
				},
				map[versions.PackageName]map[versions.ModuleName]versions.Package{
					"pkg1": {
						"Module1": versions.Package{
							Name:            "pkg1",
							Version:         "v1",
							IsIndirect:      true,
							ReplacedPath:    "pkg2",
							ReplacedVersion: "v2",
						},
						"Module2": versions.Package{
							Name:            "pkg1",
							Version:         "v1",
							IsIndirect:      true,
							ReplacedPath:    "pkg2",
							ReplacedVersion: "v2",
						},
					},
				},
				[]versions.PackageName{
					"pkg1",
				},
			},
		},
		{
			"Different versions",
			[]input{
				{
					"Module1",
					versions.Package{
						Name:    "pkg1",
						Version: "v1",
					},
				},
				{
					"Module2",
					versions.Package{
						Name:    "pkg1",
						Version: "v2",
					},
				},
			},
			expected{
				map[versions.PackageName]bool{
					"pkg1": false,
				},
				map[versions.PackageName]map[versions.ModuleName]versions.Package{
					"pkg1": {
						"Module1": versions.Package{
							Name:    "pkg1",
							Version: "v1",
						},
						"Module2": versions.Package{
							Name:    "pkg1",
							Version: "v2",
						},
					},
				},
				[]versions.PackageName{
					"pkg1",
				},
			},
		},
		{
			"Different modules",
			[]input{
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
						Name:    "pkgA",
						Version: "v1",
					},
				},
				{
					"Module2",
					versions.Package{
						Name:    "pkg2",
						Version: "v2",
					},
				},
			},
			expected{
				map[versions.PackageName]bool{
					"pkg1": true,
					"pkgA": true,
					"pkg2": true,
				},
				map[versions.PackageName]map[versions.ModuleName]versions.Package{
					"pkg1": {
						"Module1": versions.Package{
							Name:    "pkg1",
							Version: "v1",
						},
					},
					"pkgA": {
						"Module1": versions.Package{
							Name:    "pkgA",
							Version: "v1",
						},
					},
					"pkg2": {
						"Module2": versions.Package{
							Name:    "pkg2",
							Version: "v2",
						},
					},
				},
				[]versions.PackageName{
					"pkg1",
					"pkgA",
					"pkg2",
				},
			},
		},
		{
			"Uninitialized maps",
			nil,
			expected{
				map[versions.PackageName]bool{
					"pkg1": false,
				},
				map[versions.PackageName]map[versions.ModuleName]versions.Package{
					"pkg1": make(map[versions.ModuleName]versions.Package),
				},
				nil,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			packages := versions.Packages{}

			for _, input := range test.input {
				packages.Set(input.name, input.pkg)
			}

			for expectedIsSamePkgName, expectedIsSameValue := range test.expected.isSame {
				if got := packages.IsSame(expectedIsSamePkgName); got != expectedIsSameValue {
					t.Fatalf("expected %s %T, got %T", expectedIsSamePkgName, expectedIsSamePkgName, got)
				}
			}

			for expectedValuesPkgName, expectedValues := range test.expected.values {
				if values := packages.Values(expectedValuesPkgName); !cmp.Equal(values, expectedValues) {
					t.Fatalf("expected values for %s do not match: %s", expectedValuesPkgName, cmp.Diff(values, expectedValues))
				}
			}

			if names := packages.Names(); !cmp.Equal(names, test.expected.names) {
				t.Fatalf("expected names do not match: %s", cmp.Diff(names, test.expected.names))
			}
		})
	}
}
