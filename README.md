# versions

[![Go Report Card](https://goreportcard.com/badge/github.com/MarioCarrion/versions)](https://goreportcard.com/report/github.com/MarioCarrion/versions)
[![Circle CI](https://circleci.com/gh/MarioCarrion/versions.svg?style=svg)](https://circleci.com/gh/MarioCarrion/versions)

Dependencies report generator for Go projects.

## Installing

`versions` requires Go 1.14 or greater, install it using:

```
go install github.com/MarioCarrion/versions/cmd/versions
```

For projects depending on `versions` you could use the [`tools.go` paradigm](https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md):

```go
// +build tools

package tools

import (
	_ "github.com/MarioCarrion/versions/cmd/versions"
)
```

## Using

After installing you can use:

```
versions <full path to 1 go.mod> <full path to 2 go.mod> <full path to N go.mod>
```

## Example

:warning: New outputs are currently in development, at the moment Flavored Markdown is the only supported one.

Using:

```
versions ~/Repositories/versions/go.mod ~/Repositories/nit/go.mod
```

The following output will be generated:

```
|                                                              |    github.com/MarioCarrion/nit    |                                    github.com/MarioCarrion/versions                                    |
|--------------------------------------------------------------|-----------------------------------|--------------------------------------------------------------------------------------------------------|
| :white_check_mark: Go                                        |                              1.14 |                                                                                                   1.14 |
| :white_check_mark: github.com/MarioCarrion/nit               |                                   | v0.6.5                                                                                                 |
| :white_check_mark: github.com/go-enry/go-license-detector/v4 |                                   | v4.0.0<br>Apache License 2.0 permissive                                                                |
| :white_check_mark: github.com/golangci/golangci-lint         | v1.23.8                           | v1.23.8                                                                                                |
| :white_check_mark: github.com/google/go-cmp                  | v0.4.0<br>BSD-3-Clause permissive | v0.4.0<br>BSD-3-Clause permissive                                                                      |
| :white_check_mark: github.com/olekukonko/tablewriter         |                                   | v0.0.4<br>MIT License permissive                                                                       |
| :white_check_mark: github.com/pkg/errors                     | v0.9.1<br>BSD-2-Clause permissive |                                                                                                        |
| :white_check_mark: github.com/senseyeio/diligent             |                                   | v0.0.0-20191014201558-431d9a760f2d github.com/MarioCarrion/diligent v0.0.0-20200617184744-03fbc970a7f7 |
| :white_check_mark: golang.org/x/mod                          |                                   | v0.2.0<br>BSD-3-Clause permissive                                                                      |
```

Which renders like this in Markdown

|                                                              |    github.com/MarioCarrion/nit    |                                    github.com/MarioCarrion/versions                                    |
|--------------------------------------------------------------|-----------------------------------|--------------------------------------------------------------------------------------------------------|
| :white_check_mark: Go                                        |                              1.14 |                                                                                                   1.14 |
| :white_check_mark: github.com/MarioCarrion/nit               |                                   | v0.6.5                                                                                                 |
| :white_check_mark: github.com/go-enry/go-license-detector/v4 |                                   | v4.0.0<br>Apache License 2.0 permissive                                                                |
| :white_check_mark: github.com/golangci/golangci-lint         | v1.23.8                           | v1.23.8                                                                                                |
| :white_check_mark: github.com/google/go-cmp                  | v0.4.0<br>BSD-3-Clause permissive | v0.4.0<br>BSD-3-Clause permissive                                                                      |
| :white_check_mark: github.com/olekukonko/tablewriter         |                                   | v0.0.4<br>MIT License permissive                                                                       |
| :white_check_mark: github.com/pkg/errors                     | v0.9.1<br>BSD-2-Clause permissive |                                                                                                        |
| :white_check_mark: github.com/senseyeio/diligent             |                                   | v0.0.0-20191014201558-431d9a760f2d github.com/MarioCarrion/diligent v0.0.0-20200617184744-03fbc970a7f7 |
| :white_check_mark: golang.org/x/mod                          |                                   | v0.2.0<br>BSD-3-Clause permissive                                                                      |

## Features

* [X] Packages: license support.
* [ ] Packages: update availables support.
    * [ ] Merge Requests creation for Gitlab.
    * [ ] Pull Requests creation for Github.
* [ ] Packages: efferent and afferent metrics support.
* [ ] Output: Graphviz .
* [ ] Output: JSON.

## Development requirements

Go >= 1.14

## Project dependencies

* For determining the LICENSE used by the project:
  * [senseyeio/diligent](https://github.com/senseyeio/diligent): Get the licenses associated with your software dependencies
  * [go-enry/go-license-detector](https://github.com/go-enry/go-license-detector): Reliable project licenses detector.
