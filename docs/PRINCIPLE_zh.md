## 实现原理

### :wrench:`go`命令的构建参数

#### 参数`-toolexec`

指定一个外部命令或脚本调用`Go`工具链。

- 示例

```shell
go build -a -toolexec='routinex go-agent' -o main.exe .
```

`go build`是个高级的工具，其会执行`asm`、`compile`和`link`等一系列过程。

参数`-toolexec`主要拦截编译过程，再调用`compile`等底层工具，其参数为自带的`库源码 和`用户源码`。

以编译`runtime`包为例，编译过程会经历以下的步骤：

```shell
# (1)
routinex go-agent compile.exe -o xx.a -trimpath $WORK\xx=> -p runtime ...... -asmhdr runtime.go runtime2.go
# (2)
go-agent compile.exe -o xx.a -trimpath $WORK\xx=> -p runtime ...... -asmhdr runtime.go runtime2.go
# (3)
compile.exe -o xx.a -trimpath $WORK\xx=> -p runtime ...... -asmhdr runtime.go runtime2.go
```

`go build`会调用`(1)`，但是第`(2)`步需要`(1)`中的`routinex`自身逻辑发起 (2)`调用。

然后需要`(2)`中的`go-agent`发起`(3)`调用。

#### 参数`-a`

强制重新构建所有依赖包，即使它们已经是最新的。

- 示例

```shell
go build -a
```

`compile.exe`会按包缓存编译产物。当编译过一次后会将结果缓存起来。

下一个编译过程如果用到了该包的编译产物，将跳过编译直接使用上次编译产物。

由于`routinex`会替换库文件的路径，为确保修改后的文件生效，所以要指定`-a`参数。

当构建一次后并且下次构建不再需要`routinex`，需要再次指定`-a`参数，以重新缓存标准库的编译产物。

#### 参数`-x`

显示构建过程的详细命令。

`go build -x`会输出构建过程中执行的每一个命令，但不会改变构建过程本身。

它只是一个调试工具，用来了解`Go`编译器和工具链在后台运行了哪些命令。

- 示例

```shell
go build -x
```

- 输出

```text
WORK=/tmp/go-build1234567890
mkdir -p $WORK/b001/
cat >$WORK/b001/_gomod_.go << 'EOF' ...
cd /path/to/project
/usr/local/go/pkg/tool/linux_amd64/compile -o ...
/usr/local/go/pkg/tool/linux_amd64/link -o myapp ...
```

你会看到`mkdir`、`compile`和`link`等步骤。

不会保留工作目录或中间文件，这些文件会在构建结束后自动删除。

#### 参数`-work`

保留临时工作目录。

`go build -work`会保留构建过程中生成的临时工作目录，并输出路径，让你可以手动进入该目录查看中间文件。

- 示例

```shell
go build -work
```

- 输出

```text
WORK=/tmp/go-build1234567890
```

临时目录`/tmp/go-build1234567890`会被保留，里面包括编译中间文件（如`.o`、`.a`文件）。

这些文件通常用于模块的编译和链接过程，保留这些文件可以用来进一步调试或分析。

### :zap:执行逻辑

`routinex`的主要逻辑如下：

1. 过滤要修改的包，例如`runtime`和`routine`包。
2. 使用`ast`解析源文件结构，修改语法并将内容存到临时文件。
3. 修改命令行参数，使用新的文件路径替换命令行中的旧路径。
4. 使用修改后的参数，调用下一个工具链。

### :bulb:核心原理

编译阶段使用`routinex`扩展`runtime`包，从而使`routine`在运行时更安全、高效。

对`runtime`包修改内容如下：

- 为`runtime.g`结构体注入扩展字段`threadLocals`和`inheritableThreadLocals`。
- 协程退出执行函数`runtime.goexit0`时重置扩展的字段，防止`g`结构复用污染新的协程。
- 增加函数`runtime.getgp`以获取当前协程的指针。
- 增加函数`runtime.getg0`以获取协程结构`g`的类型。
