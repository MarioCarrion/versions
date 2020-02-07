# version

[![Go Report Card](https://goreportcard.com/badge/github.com/MarioCarrion/versions)](https://goreportcard.com/report/github.com/MarioCarrion/versions)
[![Circle CI](https://circleci.com/gh/MarioCarrion/versions.svg?style=svg)](https://circleci.com/gh/MarioCarrion/versions)

`versions` is a Go program for generating a report in Markdown of used packages from different repositories.

## Installing

* Using `go` (< 11.1): `go install github.com/MarioCarrion/version` **or** download a precompiled release.
* Using `go` (>= 1.11): `GO111MODULE=on go get github.com/MarioCarrion/versions`,
* Using it as a dependency with the [`tools.go` paradigm](https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md):

```go
// +build tools

package tools

import (
	_ "github.com/MarioCarrion/versions"
)
```

## Using

After installing you can use:

```
versions <full path to go.mod> <full path to go.mod>
```

## Example

Using

```
versions ~/Repositories/versions/go.mod ~/Repositories/nit/go.mod
```

The following output will be generated:

```
| | github.com/MarioCarrion/nit  | github.com/MarioCarrion/versions  |
|--- | :---:  | :---:  |
| Go version | | 1.13  | 1.13  |
| :white_check_mark: github.com/golangci/golangci-lint |  v1.23.3  |  v1.23.3  |  |
| github.com/google/go-cmp |  v0.2.0  |    |  |
| github.com/pkg/errors |  v0.8.1  |    |  |
| :white_check_mark: golang.org/x/mod |    |  v0.2.0  |  |
```

Which renders like this in Markdown

| | github.com/MarioCarrion/nit  | github.com/MarioCarrion/versions  |
|--- | :---:  | :---:  |
| Go version | | 1.13  | 1.13  |
| :white_check_mark: github.com/golangci/golangci-lint |  v1.23.3  |  v1.23.3  |  |
| github.com/google/go-cmp |  v0.2.0  |    |  |
| github.com/pkg/errors |  v0.8.1  |    |  |
| :white_check_mark: golang.org/x/mod |    |  v0.2.0  |  |

## Development requirements

Go >= 1.13.6
