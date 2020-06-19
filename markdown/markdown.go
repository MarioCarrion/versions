// Package markdown allows versions to be rendered as basic Flavored Markdown,
// supported by GitLab and GitHub.
package markdown

import (
	"strings"

	"github.com/MarioCarrion/versions"
)

type (
	// ModulesSorting is the enum for sorting modules options.
	ModulesSorting uint

	// PackagesSorting is the enum for sorting packages options.
	PackagesSorting uint

	//-

	// Markdown renders versions as basic flavored Markdown.
	Markdown struct {
		versions            versions.Versions
		modulesSortBy       ModulesSorting
		packagesSortBy      PackagesSorting
		packagesShowLicense bool
	}

	// Option is configuration option for this renderer.
	Option func(*Markdown)
)

const (
	// PackagesSortingAsFound indicates packages are not sorted and
	// instead they are rendered in the order the were found.
	PackagesSortingAsFound PackagesSorting = iota

	// PackagesSortingAlphabeticallySupported indicates packages are sorted
	// alphabetically in two groups by their name. The first group includes all
	// modules having the same version. The second group includes all modules
	// having different versions.
	PackagesSortingAlphabeticallySupported

	// PackagesSortingAlphabetically indicates packages are sorted alphabetically
	// by their name.
	PackagesSortingAlphabetically
)

const (
	// ModulesSortingAsInput indicates modules are rendered in the order they
	// were parsed.
	ModulesSortingAsInput ModulesSorting = iota

	// ModulesSortingAlphabetically indicates modules are rendered alphabetically
	// by their name.
	ModulesSortingAlphabetically
)

// NewMarkdown instantiates a new template for rendering in Markdown.
func NewMarkdown(v versions.Versions, opts ...Option) Markdown {
	md := Markdown{
		versions: v,
	}

	for _, opt := range opts {
		opt(&md)
	}

	return md
}

// WithModulesSorting allows specifyig the sorting option for modules.
func WithModulesSorting(opt ModulesSorting) Option {
	return func(m *Markdown) {
		m.modulesSortBy = opt
	}
}

// WithPackagesLicense allows display the package License when present.
func WithPackagesLicense(opt bool) Option {
	return func(m *Markdown) {
		m.packagesShowLicense = opt
	}
}

// WithPackagesSorting allows specifying the sorting option for packages.
func WithPackagesSorting(opt PackagesSorting) Option {
	return func(m *Markdown) {
		m.packagesSortBy = opt
	}
}

// String returns versions in Markdown format.
func (m Markdown) String() string {
	mods := make([]versions.Module, len(m.versions.Modules))
	index := 0

	for _, mod := range m.versions.Modules {
		mods[index] = mod
		index++
	}

	header := newHeader(m.modulesSortBy, m.versions.GoVersions.IsSame(), mods)
	pkgs := newPackages(m.versions, header.modules, m.packagesSortBy, m.packagesShowLicense)

	var str strings.Builder

	str.WriteString(header.String())
	str.WriteString(pkgs.String())

	return str.String()
}
