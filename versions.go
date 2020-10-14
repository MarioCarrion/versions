package versions

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-enry/go-license-detector/v4/licensedb"
	"github.com/go-enry/go-license-detector/v4/licensedb/api"
	"github.com/go-enry/go-license-detector/v4/licensedb/filer"
	"github.com/senseyeio/diligent"
	"golang.org/x/mod/modfile"
)

type (
	// GoVersion defines the module version of Go used by the Module.
	GoVersion string

	// ModuleName defines the name of the module.
	ModuleName string

	// PackageName defines the name of the package.
	PackageName string

	//-

	// License represents the LICENSE used by a Package.
	License struct {
		Identifier string
		Name       string
		ShortName  string
		Type       diligent.Type
		Category   diligent.Category
	}

	// Package represents an imported Go packaged in a Module.
	Package struct {
		Name            PackageName
		Version         string
		IsIndirect      bool
		ReplacedPath    string
		ReplacedVersion string
		License         License
	}

	//-

	// GoVersions handles Go versions used by different Modules.
	GoVersions struct {
		values map[ModuleName]GoVersion
		order  []ModuleName
		last   GoVersion
		same   bool
	}

	// Packages handles Packages used by different modules.
	Packages struct {
		packages     map[PackageName]map[ModuleName]Package
		lastVersions map[PackageName]Package
		sameVersions map[PackageName]bool
		names        []PackageName
	}

	//-

	// ModuleGoVersion represents a module and its Go version.
	ModuleGoVersion struct {
		Name      ModuleName
		GoVersion GoVersion
	}

	//-

	// Module represents the contents of a go.mod file.
	Module struct {
		ModuleGoVersion
		DependencyRequirements map[PackageName]Package
	}

	// Versions contains the parsed go.mod files.
	Versions struct {
		Modules    map[ModuleName]Module
		GoVersions GoVersions
		Packages   Packages
	}
)

// New returns the parsed versions used by all the mod files.
func New(files []string) (Versions, error) {
	parsed, err := newModFiles(files)
	if err != nil {
		return Versions{}, err
	}

	result := Versions{
		Modules: make(map[ModuleName]Module),
	}

	licenses := make(map[string]License)

	for _, modfile := range parsed {
		module := newModule(modfile)

		result.Modules[module.Name] = module
		result.GoVersions.Set(module.Name, module.GoVersion)

		for k, pkg := range module.DependencyRequirements {
			license, ok := licenses[pkg.Path()]
			if !ok {
				license = newLicense(pkg.Path())
				licenses[pkg.Path()] = license
			}

			pkg.License = license

			module.DependencyRequirements[k] = pkg

			result.Packages.Set(module.Name, pkg)
		}
	}

	return result, nil
}

func goModCache() string {
	if gomodcache := os.Getenv("GOMODCACHE"); gomodcache != "" {
		return gomodcache
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	return filepath.Join(gopath, "pkg", "mod")
}

func newLicense(path string) License {
	filer, err := filer.FromDirectory(path)
	if err != nil {
		return License{}
	}

	licenses, err := licensedb.Detect(filer)
	if err != nil {
		return License{}
	}

	var (
		match api.Match
		name  string
	)

	for k, v := range licenses {
		if name == "" || v.Confidence > match.Confidence {
			match = v
			name = k
		}
	}

	if name == "" {
		return License{}
	}

	license, err := diligent.GetLicenseFromIdentifier(name)
	if err != nil {
		return License{}
	}

	return License{
		Identifier: license.Identifier,
		Name:       license.Name,
		ShortName:  license.ShortName,
		Type:       license.Type,
		Category:   license.Category,
	}
}

func newModFiles(files []string) ([]*modfile.File, error) {
	parsed := make([]*modfile.File, len(files))

	for i, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		f, err := modfile.Parse(file, data, nil)
		if err != nil {
			return nil, err
		}

		parsed[i] = f
	}

	return parsed, nil
}

func newModule(modfile *modfile.File) Module {
	module := Module{
		ModuleGoVersion: ModuleGoVersion{
			Name:      ModuleName(modfile.Module.Mod.Path),
			GoVersion: GoVersion(modfile.Go.Version),
		},
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

// IsSame returns true when all Modules use the same Go Version.
func (g *GoVersions) IsSame() bool {
	return g.same
}

// Set sets the Go Version being used by the defined Module.
func (g *GoVersions) Set(name ModuleName, version GoVersion) {
	if g.last == "" {
		g.last = version
		g.same = true
	} else if g.same {
		g.same = g.last == version
	}

	if g.values == nil {
		g.values = make(map[ModuleName]GoVersion)
	}

	if _, ok := g.values[name]; !ok {
		// to keep the Set order
		g.order = append(g.order, name)
	}

	g.values[name] = version
}

// Values returns a copy of currently set modules and their Go versions.
func (g *GoVersions) Values() []ModuleGoVersion {
	result := make([]ModuleGoVersion, len(g.order))

	for i, name := range g.order {
		version := g.values[name]
		result[i].Name = name
		result[i].GoVersion = version
	}

	return result
}

// Path returns the full filesystem path pointing to the package
func (p Package) Path() string {
	version := func(v string) string {
		if v == "" {
			return ""
		}

		return fmt.Sprintf("@%s", v)
	}

	var newVersion string
	if p.ReplacedPath != "" {
		newVersion = fmt.Sprintf("%s%s", p.ReplacedPath, version(p.ReplacedVersion))
	} else {
		newVersion = fmt.Sprintf("%s%s", p.Name, version(p.Version))
	}

	return filepath.Join(goModCache(), newVersion)
}

// IsSame returns true when all Modules use the same Package Version.
func (p *Packages) IsSame(value PackageName) bool {
	if p.sameVersions == nil {
		return false
	}

	return p.sameVersions[value]
}

// Names returns a slice of all package names used in total.
func (p *Packages) Names() []PackageName {
	if p.names == nil {
		return nil
	}

	res := make([]PackageName, len(p.packages))
	copy(res, p.names)

	return res
}

// Set sets the Package being used by the defined Module.
func (p *Packages) Set(name ModuleName, pkg Package) {
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
		p.names = append(p.names, pkg.Name)
	}

	mods[name] = pkg

	p.packages[pkg.Name] = mods

	if p.sameVersions[pkg.Name] && pkg != p.lastVersions[pkg.Name] {
		p.sameVersions[pkg.Name] = false
	}
}

// Values returns a copy of currently set modules and their packages by package.
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
