// +build tools

package tools

// go install github.com/golangci/golangci-lint/cmd/golangci-lint github.com/MarioCarrion/nit/cmd/nit

import (
	_ "github.com/MarioCarrion/nit/cmd/nit"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
