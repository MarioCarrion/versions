# versions

[![Go Report Card](https://goreportcard.com/badge/github.com/MarioCarrion/versions)](https://goreportcard.com/report/github.com/MarioCarrion/versions)
[![Circle CI](https://circleci.com/gh/MarioCarrion/versions.svg?style=svg)](https://circleci.com/gh/MarioCarrion/versions)

`versions` is a Go program for generating a report in Markdown of used packages from different repositories.

## Installing

* Using `go` (< 11.1): `go install github.com/MarioCarrion/versions` **or** download a precompiled release.
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
versions <full path to 1 go.mod> <full path to 2 go.mod> <full path to N go.mod>
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
| Go version  | 1.13  | 1.13  |
| github.com/golangci/golangci-lint |  v1.23.3  |  v1.23.2  |
| :white_check_mark: github.com/google/go-cmp |  v0.2.0  |    |
| :white_check_mark: github.com/pkg/errors |  v0.8.1  |    |
| :white_check_mark: golang.org/x/mod |    |  v0.2.0  |
```

Which renders like this in Markdown

| | github.com/MarioCarrion/nit  | github.com/MarioCarrion/versions  |
|--- | :---:  | :---:  |
| Go version  | 1.13  | 1.13  |
| github.com/golangci/golangci-lint |  v1.23.3  |  v1.23.2  |
| :white_check_mark: github.com/google/go-cmp |  v0.2.0  |    |
| :white_check_mark: github.com/pkg/errors |  v0.8.1  |    |
| :white_check_mark: golang.org/x/mod |    |  v0.2.0  |

## Development requirements

Go >= 1.13.6
