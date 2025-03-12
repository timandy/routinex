# RoutineX Compiler

[![Build Status](https://github.com/timandy/routinex/actions/workflows/build.yml/badge.svg)](https://github.com/timandy/routinex/actions)
[![Codecov](https://codecov.io/gh/timandy/routinex/branch/main/graph/badge.svg)](https://app.codecov.io/gh/timandy/routinex)
[![Go Report Card](https://goreportcard.com/badge/github.com/timandy/routinex)](https://goreportcard.com/report/github.com/timandy/routinex)
[![Documentation](https://pkg.go.dev/badge/github.com/timandy/routinex.svg)](https://pkg.go.dev/github.com/timandy/routinex)
[![Release](https://img.shields.io/github/release/timandy/routinex.svg)](https://github.com/timandy/routinex/releases)
[![License](https://img.shields.io/github/license/timandy/routinex.svg)](https://github.com/timandy/routinex/blob/main/LICENSE)

> [中文版](README_zh.md)

**`RoutineX`** is a compilation tool designed to enhance the goroutine capabilities of the `runtime` package in `Golang` standard library.

As a compilation tool for the `routine` static mode, it supports advanced features such as native goroutine-local storage (`TLS`).

This enables `routine` in static mode to achieve higher performance and a safer memory access mechanism.

## :house:Introduction

In the normal mode, `routine` stores the `labelMap` and goroutine context data in the `g.labels` field based on the structure layout.

This requires additional validation and judgment during reading and inheritance, which impacts performance and may not be compatible with future go versions.

To further improve performance and compatibility, the `RoutineX` compiler uses source code enhancement techniques to export some functions at the lower level,
allowing `routine` to efficiently and safely access the data of the goroutine structure `g`.

## :hammer_and_wrench:Usage Instructions

### Installation

```shell
go install -a github.com/timandy/routinex@latest
```

Before calling, you need to set the environment variable to append `$(go env GOPATH)/bin` to `PATH`, so that `routinex` can be run directly in the console.

- windows

```powershell
#!/bin/pwsh

$env:Path="$env:Path;$(go env GOPATH)\bin"
```

- linux

```shell
#!/bin/bash

export PATH="$PATH:$(go env GOPATH)/bin"
```

### Enable Debugging

Using debugging parameters can output logs from `routinex`.

- `--debug` or `-d` outputs detailed logs.
- `--verbose` or `-v` outputs brief logs.

- windows

```powershell
#!/bin/pwsh

# Set environment variable
$env:Path="$env:Path;$(go env GOPATH)\bin"
# Add parameters
go build -a -toolexec='routinex -v' -o main.exe .
```

- linux

```shell
#!/bin/bash

# Set environment variable
export PATH="$PATH:$(go env GOPATH)/bin"
# Add parameters
go build -a -toolexec='routinex -v' -o main.exe .
```

### Multiple Toolchains

If you have already used another toolchain before using `routinex`, take `abc` as an example.

Since it cannot be guaranteed that `abc` has chain transfer functionality, `routinex` should be placed before the `abc` toolchain.

After `routinex` is executed, it will automatically call the `abc` tool.

- windows

```powershell
#!/bin/pwsh

# Set environment variable
$env:Path="$env:Path;$(go env GOPATH)\bin"
# Add parameters
go build -a -toolexec='routinex -v abc' -o main.exe .
```

- linux

```shell
#!/bin/bash

# Set environment variable
export PATH="$PATH:$(go env GOPATH)/bin"
# Add parameters
go build -a -toolexec='routinex -v abc' -o main.exe .
```

## :bulb:Implementation Principle

[See](docs/PRINCIPLE.md)

## :rocket:Performance Improvement

[See](docs/PERFORMANCE.md)

## :scroll:*License*

`RoutineX`is released under the [Apache License 2.0](LICENSE).

```
Copyright 2021-2025 TimAndy

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
