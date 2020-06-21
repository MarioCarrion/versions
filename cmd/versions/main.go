// +build go1.14

package main

import (
	"fmt"
	"os"

	"github.com/MarioCarrion/versions"
	"github.com/MarioCarrion/versions/markdown"
)

func main() {
	params := os.Args[1:]
	if len(params) == 0 {
		fmt.Println("path to go.mod files required")
		os.Exit(1)
	}

	gomods, err := versions.New(params)
	if err != nil {
		fmt.Printf("error parsing files %s\n", err)
		os.Exit(1)
	}

	md := markdown.NewMarkdown(gomods,
		markdown.WithModulesSorting(markdown.ModulesSortingAlphabetically),
		markdown.WithPackagesSorting(markdown.PackagesSortingAlphabeticallySupported),
		markdown.WithPackagesLicense(true))

	fmt.Println(md.String())
}
