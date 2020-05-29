package versions

import (
	"golang.org/x/mod/modfile"
)

type (
	GoVersion string

	ModuleName string

	PackageName string

	//-

	GoVersions struct {
		values map[ModuleName]GoVersion
		last   GoVersion
		same   bool
	}

	ModuleV2 struct {
		Name                   ModuleName
		GoVersion              GoVersion
		DependencyRequirements map[PackageName]Package
	}

	Packages struct {
		packages     map[PackageName]map[ModuleName]Package
		lastVersions map[PackageName]Package
		sameVersions map[PackageName]bool
	}

	VersionV2 struct {
		Modules map[ModuleName]string //RepositoryVersion
	}

	//-

	Result struct {
		Modules    map[ModuleName]ModuleV2
		GoVersions GoVersions
		Packages   Packages
	}
)

func NewResult(files []string) (Result, error) { // XXX: -> NewModules?
	parsed, err := NewModFiles(files)
	if err != nil {
		return Result{}, err
	}

	result := Result{
		Modules: make(map[ModuleName]ModuleV2),
	}

	for _, modfile := range parsed {
		module := NewModule(modfile)

		result.Modules[module.Name] = module
		result.GoVersions.Add(module.Name, module.GoVersion)

		for _, pkg := range module.DependencyRequirements {
			result.Packages.Add(module.Name, pkg)
		}
	}

	return result, nil
}

//-

func NewModule(modfile *modfile.File) ModuleV2 {
	module := ModuleV2{
		Name:      ModuleName(modfile.Module.Mod.Path),
		GoVersion: GoVersion(modfile.Go.Version),
	}

	dependencies := make(map[PackageName]Package)

	for _, require := range modfile.Require {
		dependencies[PackageName(require.Mod.Path)] = Package{
			Name:       PackageName(require.Mod.Path),
			Version:    require.Mod.Version,
			IsIndirect: require.Indirect,
		}
	}

	for _, replace := range modfile.Replace {
		pkg, ok := dependencies[PackageName(replace.Old.Path)]
		if !ok {
			continue
		}

		pkg.ReplacedPath = replace.New.Path
		pkg.ReplacedVersion = replace.New.Version

		dependencies[PackageName(replace.Old.Path)] = pkg
	}

	module.DependencyRequirements = dependencies

	return module
}

func (g *GoVersions) Add(name ModuleName, version GoVersion) {
	if g.last == "" {
		g.last = version
		g.same = true
	} else if g.same {
		g.same = g.last == version
	}

	if g.values == nil {
		g.values = make(map[ModuleName]GoVersion)
	}

	g.values[name] = version
}

func (g *GoVersions) IsSame() bool {
	return g.same
}

func (g *GoVersions) Values() map[ModuleName]GoVersion {
	result := make(map[ModuleName]GoVersion)

	for k, v := range g.values {
		result[k] = v
	}

	return result
}

func (p *Packages) Add(name ModuleName, pkg Package) {
	if p.packages == nil {
		p.packages = make(map[PackageName]map[ModuleName]Package)
		p.lastVersions = make(map[PackageName]Package)
		p.sameVersions = make(map[PackageName]bool)
	}

	mods, ok := p.packages[pkg.Name]
	if !ok {
		mods = make(map[ModuleName]Package)
		p.lastVersions[pkg.Name] = pkg
		p.sameVersions[pkg.Name] = true
	}

	mods[name] = pkg

	p.packages[pkg.Name] = mods

	if p.sameVersions[pkg.Name] && pkg != p.lastVersions[pkg.Name] {
		p.sameVersions[pkg.Name] = false
	}
}

func (p *Packages) IsSame(value PackageName) bool {
	if p.sameVersions == nil {
		return false
	}

	return p.sameVersions[value]
}

func (p *Packages) Values(value PackageName) map[ModuleName]Package {
	result := make(map[ModuleName]Package)

	if p.packages == nil {
		return result
	}

	for k, v := range p.packages[value] {
		result[k] = v
	}

	return result
}
