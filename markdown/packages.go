package markdown

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MarioCarrion/versions"
)

type (
	packageSet struct {
		same        bool
		showLicense bool
		Name        versions.PackageName
		packages    []versions.Package
	}

	packageSets []packageSet

	packages struct {
		same      packageSets
		different packageSets
	}
)

func newPackages(vs versions.Versions, modules []module, sorting PackagesSorting, showLicense bool) packages {
	var res packages

	for _, name := range vs.Packages.Names() {
		set := packageSet{
			Name:        name,
			showLicense: showLicense,
			same:        vs.Packages.IsSame(name),
			packages:    make([]versions.Package, len(modules)),
		}

		for i, mod := range modules {
			vMod := vs.Modules[mod.Name]
			pkg, ok := vMod.DependencyRequirements[name]

			if ok {
				set.packages[i] = pkg
			}
		}

		dest := &res.same
		if sorting == PackagesSortingAlphabeticallySupported && !set.same {
			dest = &res.different
		}

		*dest = append(*dest, set)
	}

	if sorting != PackagesSortingAsFound {
		sort.Sort(&res.same)
		sort.Sort(&res.different)
	}

	return res
}

func (p packageSet) Values() []string {
	res := make([]string, len(p.packages)+1)

	var str string

	if p.same { // XXX: Make emoji option configurable
		// str = "✓ "
		// str = "✅"
		str = ":white_check_mark: "
	}

	res[0] = fmt.Sprintf("%s%s", str, p.Name)

	for i, v := range p.packages {
		var b strings.Builder

		b.WriteString(v.Version)

		if v.ReplacedPath != "" {
			b.WriteString(" ")
			b.WriteString(v.ReplacedPath)
		}

		if v.ReplacedVersion != "" {
			b.WriteString(" ")
			b.WriteString(v.ReplacedVersion)
		}

		if p.showLicense && v.License.Identifier != "" {
			b.WriteString("<br>")
			b.WriteString(string(v.License.Category))
			b.WriteString(" ")
			b.WriteString(v.License.Name)
		}

		res[i+1] = b.String()
	}

	return res
}

func (p packageSets) Len() int {
	return len([]packageSet(p))
}

func (p packageSets) Less(i, j int) bool {
	arr := []packageSet(p)
	return arr[i].Name < arr[j].Name
}

func (p packageSets) Swap(i, j int) {
	arr := []packageSet(p)
	arr[i], arr[j] = arr[j], arr[i]
}

//-

func (p packages) Values() [][]string {
	var res [][]string

	if len(p.same) > 0 {
		setsSame := []packageSet(p.same)

		for _, val := range setsSame {
			res = append(res, val.Values())
		}
	}

	if len(p.different) > 0 {
		setsSame := []packageSet(p.different)

		for _, val := range setsSame {
			res = append(res, val.Values())
		}
	}

	return res
}
