package versions

import (
	"sort"
)

type (
	GoMods struct {
		Names             []string
		SameGoVersions    bool
		Modules           []Module
		AllVersions       Versions
		SameVersions      Versions
		DifferentVersions Versions
	}

	Version struct {
		Name   string
		IsSame bool
	}

	Versions []Version

	//-

	Module struct {
		GoVersion string
		Packages  map[string]Package
	}

	// Package is the module being used by the repository Go module.
	Package struct {
		Name            PackageName
		Version         string
		IsIndirect      bool
		ReplacedPath    string
		ReplacedVersion string
	}
)

// NewGoMods returns categorized Go Modules depending on versions being used.
func NewGoMods(files []string) (GoMods, error) { //nolint: funlen
	parsed, err := NewModFiles(files)
	if err != nil {
		return GoMods{}, err
	}

	gomods := GoMods{
		Names:   make([]string, len(parsed)),
		Modules: make([]Module, len(parsed)),
	}

	//-

	allpackages := make(map[string]interface{})

	for i, file := range parsed {
		gomods.Names[i] = file.Module.Mod.Path

		gomods.Modules[i].GoVersion = file.Go.Version

		modules := make(map[string]Package)

		for _, req := range file.Require {
			allpackages[req.Mod.Path] = nil

			pkg := Package{
				Version:    req.Mod.Version,
				IsIndirect: req.Indirect,
			}

			modules[req.Mod.Path] = pkg
		}

		for _, rep := range file.Replace {
			orig, ok := modules[rep.Old.Path]
			if !ok {
				continue
			}

			orig.ReplacedPath = rep.New.Path
			orig.ReplacedVersion = rep.New.Version

			modules[rep.Old.Path] = orig
		}

		gomods.Modules[i].Packages = modules
	}

	gomods.AllVersions = make([]Version, len(allpackages))

	var (
		index int
		pkgs  []string = make([]string, len(gomods.Modules))
	)

	for pkg := range allpackages {
		for i := range gomods.Names {
			v, ok := gomods.Modules[i].Packages[pkg]
			if ok {
				pkgs[i] = v.Version
				if v.ReplacedVersion != "" {
					pkgs[i] = v.ReplacedVersion
				}
			} else {
				pkgs[i] = ""
			}
		}

		gomods.AllVersions[index].Name = pkg

		version := Version{
			Name: pkg,
		}

		if SameVersion(pkgs) {
			version.IsSame = true
			gomods.SameVersions = append(gomods.SameVersions, version)
			gomods.AllVersions[index].IsSame = true
		} else {
			gomods.DifferentVersions = append(gomods.DifferentVersions, version)
		}

		index++
	}

	sort.Sort(gomods.AllVersions)
	sort.Sort(gomods.SameVersions)
	sort.Sort(gomods.DifferentVersions)

	versions := make([]string, len(gomods.Modules))
	for i, module := range gomods.Modules {
		versions[i] = module.GoVersion
	}

	gomods.SameGoVersions = SameVersion(versions)

	return gomods, nil
}

// SameVersion indicates whether all values in `versions` are the same,  empty
// strings are ignored.
func SameVersion(versions []string) bool {
	var last string

	for _, current := range versions {
		if current == "" {
			continue
		}

		if last != "" && current != last {
			return false
		}

		if current != "" {
			last = current
		}
	}

	return true
}

// Len is the number of elements in the collection.
func (v Versions) Len() int {
	return len(v)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (v Versions) Less(i, j int) bool {
	return v[i].Name < v[j].Name
}

// Swap swaps the elements with indexes i and j.
func (v Versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
