package markdown

import (
	"sort"
	"strings"

	"github.com/MarioCarrion/versions"
)

type (
	header struct {
		same    bool
		modules modules
	}

	module struct {
		versions.Module
	}

	modules []module
)

func newHeader(sorting ModulesSorting, same bool, values []versions.Module) header {
	var mods modules = make([]module, len(values))

	for i, mod := range values {
		mods[i] = module{mod}
	}

	if sorting == ModulesSortingAlphabetically {
		sort.Sort(&mods)
	}

	return header{
		same:    same,
		modules: mods,
	}
}

func (h header) String() string {
	var str strings.Builder

	// Modules header

	str.WriteString("|")

	for _, mod := range h.modules {
		str.WriteString(" | ")
		str.WriteString(string(mod.Name))
	}

	str.WriteString(" |\n")

	// Modules Go Versions

	str.WriteString("| ")

	if h.same {
		str.WriteString(":white_check_mark: ")
	}

	str.WriteString("Go")

	for _, mod := range h.modules {
		str.WriteString(" | ")
		str.WriteString(string(mod.GoVersion))
	}

	str.WriteString(" |\n")

	return str.String()
}

func (m modules) Len() int {
	return len([]module(m))
}

func (m modules) Less(i, j int) bool {
	arr := []module(m)
	return arr[i].Name < arr[j].Name
}

func (m modules) Swap(i, j int) {
	arr := []module(m)
	arr[i], arr[j] = arr[j], arr[i]
}
