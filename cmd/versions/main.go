package main

import (
	"fmt"
	"os"

	"github.com/MarioCarrion/versions"
)

func main() {
	params := os.Args[1:]
	if len(params) == 0 {
		fmt.Println("path to go.mod files required")
		os.Exit(1)
	}

	gomods, err := versions.NewGoMods(params)
	if err != nil {
		fmt.Printf("error parsing files %s\n", err)
		os.Exit(1)
	}

	//-

	versions.PrintMarkdown(gomods)
}
