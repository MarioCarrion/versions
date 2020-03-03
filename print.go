package versions

import (
	"fmt"
	"strings"
)

// PrintMarkdown outputs the parsed Go modules as Markdown.
func PrintMarkdown(gomods GoMods) {
	var (
		header0 strings.Builder
		header1 strings.Builder
		header2 strings.Builder
	)

	header0.WriteString("|")
	header1.WriteString("|---")

	header2.WriteString("|")

	if gomods.SameGoVersions {
		header2.WriteString(" :white_check_mark: ")
	}

	header2.WriteString("Go ")

	for i, name := range gomods.Names {
		header0.WriteString(" | ")
		header0.WriteString(name)
		header0.WriteString(" ")

		header1.WriteString(" | :---: ")

		header2.WriteString("| ")
		header2.WriteString(gomods.Modules[i].GoVersion)
	}

	header0.WriteString(" |\n")
	header1.WriteString(" |\n")
	header2.WriteString(" |\n")

	fmt.Printf("%s", header0.String())
	fmt.Printf("%s", header1.String())
	fmt.Printf("%s", header2.String())

	//-

	print := func(versions Versions) {
		for _, version := range versions {
			var line strings.Builder

			line.WriteString("| ")

			if version.IsSame {
				line.WriteString(":white_check_mark: ")
			}

			line.WriteString(version.Name)
			line.WriteString(" | ")

			for i := 0; i < len(gomods.Names); i++ {
				pkg := gomods.Modules[i].Packages[version.Name]

				line.WriteString(pkg.Version)

				if pkg.IsIndirect {
					line.WriteString(" :question:")
				}

				if pkg.ReplacedVersion != "" {
					line.WriteString(" :exclamation:")
				}

				line.WriteString(" | ")
			}

			fmt.Println(line.String())
		}
	}

	// printVersions(gomods.AllVersions)

	print(gomods.SameVersions)
	print(gomods.DifferentVersions)
}
