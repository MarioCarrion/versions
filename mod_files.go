package versions

import (
	"io/ioutil"
	"sort"

	"golang.org/x/mod/modfile"
)

type (
	parse struct {
		parsed []*modfile.File
	}
)

// NewModFiles returns a parsed and sorted (by module name) slice a modfiles.
func NewModFiles(files []string) ([]*modfile.File, error) {
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

	p := parse{
		parsed: parsed,
	}

	sort.Sort(&p)

	return p.parsed, nil
}

// Len is the number of elements in the collection.
func (p *parse) Len() int {
	return len(p.parsed)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (p *parse) Less(i, j int) bool {
	return p.parsed[i].Module.Mod.Path < p.parsed[j].Module.Mod.Path
}

// Swap swaps the elements with indexes i and j.
func (p *parse) Swap(i, j int) {
	p.parsed[i], p.parsed[j] = p.parsed[j], p.parsed[i]
}
