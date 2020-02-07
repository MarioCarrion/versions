package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"golang.org/x/mod/modfile"
)

type (
	Table struct {
		Modules  []string
		Versions []string
		Packages map[string][]Module // key is "Package Name"
	}

	Module struct {
		Version         string
		IsIndirect      bool
		ReplacedPath    string
		ReplacedVersion string
	}
)

func main() {
	params := os.Args[1:]
	if len(params) == 0 {
		fmt.Println("path to go.mod files required")
		os.Exit(1)
	}

	parsed, err := parse(params)
	if err != nil {
		fmt.Printf("error parsing files %s\n", err)
		os.Exit(1)
	}

	table := newTable(parsed)

	//-

	printMarkdown(table)
}

//-

func newTable(parsed []*modfile.File) Table {
	table := Table{}
	table.Modules = make([]string, len(parsed))
	table.Versions = make([]string, len(parsed))
	table.Packages = make(map[string][]Module)

	//-

	for i, file := range parsed {
		table.Modules[i] = file.Module.Mod.Path
		table.Versions[i] = file.Go.Version

		for _, req := range file.Require {
			pkg := req.Mod.Path

			var modules []Module
			if modules = table.Packages[pkg]; modules == nil {
				modules = make([]Module, len(parsed))
			}

			modules[i].Version = req.Mod.Version
			modules[i].IsIndirect = req.Indirect

			table.Packages[pkg] = modules
		}

		for _, rep := range file.Replace {
			old, ok := table.Packages[rep.Old.Path]
			if !ok {
				continue
			}

			old[i].ReplacedPath = rep.New.Path
			old[i].ReplacedVersion = rep.New.Version
			table.Packages[rep.Old.Path] = old
		}
	}

	return table
}

func parse(files []string) ([]*modfile.File, error) {
	sort.Strings(files)

	parse := func(file string) (*modfile.File, error) {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		f, err := modfile.Parse(file, data, nil)
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	var err error

	parsed := make([]*modfile.File, len(files))

	for i, f := range files {
		parsed[i], err = parse(f)
		if err != nil {
			return nil, err
		}
	}

	return parsed, nil
}

func printMarkdown(table Table) {
	line0 := "|"
	line1 := "|---"
	line2 := "| Go version "

	for i, name := range table.Modules {
		line0 = fmt.Sprintf("%s | %s ", line0, name)
		line1 = fmt.Sprintf("%s | :---: ", line1)
		line2 = fmt.Sprintf("%s | %s ", line2, table.Versions[i])
	}

	fmt.Printf("%s |\n", line0)
	fmt.Printf("%s |\n", line1)
	fmt.Printf("%s |\n", line2)

	sortedpkgs := make([]string, len(table.Packages))

	var index int

	for k := range table.Packages {
		sortedpkgs[index] = k
		index++
	}

	sort.Strings(sortedpkgs)

	//-

	for _, pkg := range sortedpkgs {
		v := table.Packages[pkg]

		var line, lastversion string

		same := true

		for _, p := range v {
			var version string
			if p.ReplacedVersion != "" {
				version = p.ReplacedVersion
			} else {
				version = p.Version
			}

			if lastversion == "" {
				lastversion = version
			}

			if same && lastversion != version {
				same = false
			}

			line = fmt.Sprintf("%s %s ", line, version)
			if p.IsIndirect {
				line = fmt.Sprintf("%s :question: ", line)
			}

			if p.ReplacedVersion != "" {
				line = fmt.Sprintf("%s :exclamation: ", line)
			}

			line = fmt.Sprintf("%s | ", line)
		}

		var prefix string
		if same {
			prefix = fmt.Sprintf(":white_check_mark: %s", pkg)
		} else {
			prefix = pkg
		}

		line = fmt.Sprintf("| %s | %s", prefix, line)

		fmt.Println(line)
	}
}
