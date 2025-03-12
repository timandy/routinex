## Implementation Principle

### :wrench:Build Parameters of `go` Command

#### Parameter `-toolexec`

Specifies an external command or script to invoke the `Go` toolchain.

- Example

```shell
go build -a -toolexec='routinex go-agent' -o main.exe .
```

`go build` is a high-level tool that executes a series of processes such as `asm.exe`, `compile.exe` and `link.exe`.

The `-toolexec` parameter mainly intercepts the compilation process and then calls underlying tools like `compile.exe`, with parameters being the built-in library source code and user source code.

Taking the compilation of the `runtime` package as an example, the compilation process goes through the following steps:

```shell
# (1)
routinex go-agent compile.exe -o xx.a -trimpath $WORK\xx=> -p runtime ...... -asmhdr runtime.go runtime2.go
# (2)
go-agent compile.exe -o xx.a -trimpath $WORK\xx=> -p runtime ...... -asmhdr runtime.go runtime2.go
# (3)
compile.exe -o xx.a -trimpath $WORK\xx=> -p runtime ...... -asmhdr runtime.go runtime2.go
```

`go build` will call `(1)`, but step `(2)` requires `routinex` in `(1)` to initiate the `(2)` call itself.

Then, `go-agent` in `(2)` needs to initiate the `(3)` call.

#### Parameter `-a`

Forces rebuilding all dependent packages, even if they are already up to date.

- Example

```shell
go build -a
```

`compile.exe` caches compilation products by package. Once a package is compiled, the result is cached.

If the next compilation process uses the compilation product of that package, it will skip the compilation and directly use the cached product.

Since `routinex` replaces the paths of library files, to ensure that the modified files take effect, the `-a` parameter must be specified.

After building once and if `routinex` is no longer needed for the next build, the `-a` parameter must be specified again to re-cache the compilation products of the standard library.

#### Parameter `-x`

Displays detailed commands of the build process.

`go build -x` outputs every command executed during the build process but does not change the build process itself.

It is merely a debugging tool to understand which commands the Go compiler and toolchain run in the background.

- Example

```shell
go build -x
```

- Output

```text
WORK=/tmp/go-build1234567890
mkdir -p $WORK/b001/
cat >$WORK/b001/_gomod_.go << 'EOF' ...
cd /path/to/project
/usr/local/go/pkg/tool/linux_amd64/compile -o ...
/usr/local/go/pkg/tool/linux_amd64/link -o myapp ...
```

You will see steps like `mkdir`, `compile`, and `link`.

The working directory or intermediate files are not retained; these files are automatically deleted after the build.

#### Parameter `-work`

Retains the temporary working directory.

`go build -work` retains the temporary working directory generated during the build process and outputs the path, allowing you to manually enter the directory to view intermediate files.

- Example

```shell
go build -work
```

- Output

```text
WORK=/tmp/go-build1234567890
```

The temporary directory `/tmp/go-build1234567890` is retained, containing compilation intermediate files (e.g. `.o`, `.a` files).

These files are typically used in the compilation and linking process of modules. Retaining these files can be useful for further debugging or analysis.

### :zap:Execution Logic

The main logic of `routinex` is as follows:

1. Filters the packages to be modified, such as the `runtime` and `routine` packages.
2. Uses `ast` to parse the source file structure, modifies the syntax, and saves the content to temporary files.
3. Modifies command-line parameters, replacing old paths in the command line with new file paths.
4. Calls the next toolchain using the modified parameters.

### :bulb:Core Principle

The `runtime` package is extended during the compilation phase using `routinex`, making the `routine` safer and more efficient at runtime.

The modifications to the `runtime` package are as follows:

- Injects extended fields `threadLocals` and `inheritableThreadLocals` into the `runtime.g` structure.
- Resets the extended fields when the goroutine exit function `runtime.goexit0` is executed, preventing the reuse of the `g` structure from polluting new goroutines.
- Adds the function `runtime.getgp` to get the current goroutine's pointer.
- Adds the function `runtime.getg0` to get the type of the goroutine structure `g`.
