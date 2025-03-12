# RoutineX 编译器

[![Build Status](https://github.com/timandy/routinex/actions/workflows/build.yml/badge.svg)](https://github.com/timandy/routinex/actions)
[![Codecov](https://codecov.io/gh/timandy/routinex/branch/main/graph/badge.svg)](https://app.codecov.io/gh/timandy/routinex)
[![Go Report Card](https://goreportcard.com/badge/github.com/timandy/routinex)](https://goreportcard.com/report/github.com/timandy/routinex)
[![Documentation](https://pkg.go.dev/badge/github.com/timandy/routinex.svg)](https://pkg.go.dev/github.com/timandy/routinex)
[![Release](https://img.shields.io/github/release/timandy/routinex.svg)](https://github.com/timandy/routinex/releases)
[![License](https://img.shields.io/github/license/timandy/routinex.svg)](https://github.com/timandy/routinex/blob/main/LICENSE)

> [English Version](README.md)

**`RoutineX`** 是一个编译工具，旨在增强`Golang`标准库中`runtime`包的协程能力。

作为`routine`静态模式的编译工具，它支持原生协程级存储（`TLS`）等高级特性。

这使`routine`在静态模式下，具有更高的性能、更安全的内存访问机制。

## :house:介绍

`routine` 在普通的模式下，根据结构体布局，将`labelMap`和协程上下文数据存储在`g.labels`字段。

这需要再读取和继承时，要做额外的验证和判断，对性能具有一定的影响，并且可能不兼容将来的`go`版本。

为了进一步提升性能和兼容性，`RoutineX`编译器通过源码增强技术，在底层导出部分函数，
使`routine`能高效，安全的访问协程结构`g`的数据。

## :hammer_and_wrench:使用说明

### 安装

```shell
go install -a github.com/timandy/routinex@latest
```

调用前需要设置环境变量，把`$(go env GOPATH)/bin`追加到`PATH`，以便在控制台直接运行`routinex`。

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

### 开启调试

使用调试参数可以输出`routinex`的日志。

- `--debug`或`-d`输出详细的日志。
- `--verbose`或`-v`输出简略的日志。

- windows

```powershell
#!/bin/pwsh

# 设置环境变量
$env:Path="$env:Path;$(go env GOPATH)\bin"
# 添加参数
go build -a -toolexec='routinex -v' -o main.exe .
```

- linux

```shell
#!/bin/bash

# 设置环境变量
export PATH="$PATH:$(go env GOPATH)/bin"
# 添加参数
go build -a -toolexec='routinex -v' -o main.exe .
```

### 多工具链

如果在使用`routinex`之前已经使用了其他的工具链，以`abc`为例。

因为不能保证`abc`有链式传递功能，所以要将`routinex`放在`abc`工具链前边。

`routinex`执行完后，会自动调用`abc`工具。

- windows

```powershell
#!/bin/pwsh

# 设置环境变量
$env:Path="$env:Path;$(go env GOPATH)\bin"
# 添加参数
go build -a -toolexec='routinex -v abc' -o main.exe .
```

- linux

```shell
#!/bin/bash

# 设置环境变量
export PATH="$PATH:$(go env GOPATH)/bin"
# 添加参数
go build -a -toolexec='routinex -v abc' -o main.exe .
```

## :bulb:实现原理

[参见](docs/PRINCIPLE_zh.md)

## :rocket:性能提升

[参见](docs/PERFORMANCE_zh.md)

## :scroll:*许可证*

`RoutineX`是在 [Apache License 2.0](LICENSE) 下发布的。

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
