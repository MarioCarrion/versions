package markdown

import (
	"fmt"
	"sort"

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

func (h header) GoVersions() []string {
	res := make([]string, len(h.modules)+1)

	var str string

	if h.same {
		str = ":white_check_mark: "
	}

	res[0] = fmt.Sprintf("%sGo", str)

	for i, mod := range h.modules {
		res[i+1] = string(mod.GoVersion)
	}

	return res
}

func (h header) Names() []string {
	res := make([]string, len(h.modules)+1)

	for i, mod := range h.modules {
		res[i+1] = string(mod.Name)
	}

	return res
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
