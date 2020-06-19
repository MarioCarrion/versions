package markdown

import (
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

func (p packageSet) String() string {
	var str strings.Builder

	str.WriteString("| ")

	if p.same {
		str.WriteString(":white_check_mark: ")
	}

	str.WriteString(string(p.Name))

	for _, v := range p.packages {
		str.WriteString(" | ")
		str.WriteString(v.Version)

		if v.ReplacedPath != "" {
			str.WriteString(" - ")
			str.WriteString(v.ReplacedPath)
		}

		if v.ReplacedVersion != "" {
			str.WriteString("/")
			str.WriteString(v.ReplacedVersion)
		}

		if p.showLicense && v.License.Identifier != "" {
			str.WriteString("<br>")
			str.WriteString(v.License.Name)
			str.WriteString(" - ")
			str.WriteString(string(v.License.Category))
		}
	}

	str.WriteString(" |\n")

	return str.String()
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

func (p packageSets) String() string {
	var str strings.Builder

	sets := []packageSet(p)

	for _, s := range sets {
		str.WriteString(s.String())
	}

	return str.String()
}

func (p packages) String() string {
	var str strings.Builder

	if p.same != nil {
		str.WriteString(p.same.String())
	}

	if p.different != nil {
		str.WriteString(p.different.String())
	}

	return str.String()
}
