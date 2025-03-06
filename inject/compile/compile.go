package compile

import (
	"github.com/timandy/routinex/inject/api"
	"github.com/timandy/routinex/tools/flag"
	"github.com/timandy/routinex/tools/log"
	"github.com/timandy/routinex/tools/opt"
	"github.com/timandy/routinex/tools/stringutil"
)

var flags = []string{"-asmhdr", "-pack"}

type CompileCmd struct {
	injectors []api.Injector
	options   *CompileOptions
}

// NewCompileCmd 构造函数
func NewCompileCmd(injectors []api.Injector) *CompileCmd {
	return &CompileCmd{injectors: injectors}
}

// Resolve 解析参数
func (c *CompileCmd) Resolve(args []string, app *opt.AppOptions) {
	options := c.resolve(args)
	if options != nil {
		options.Debug = app.Debug
		options.Verbose = app.Verbose
		options.Args = args
	}
	c.options = options
}

// resolve 解析参数-递归
func (c *CompileCmd) resolve(args []string) *CompileOptions {
	if len(args) == 0 {
		return nil
	}
	options := &CompileOptions{}
	flagSet := flag.ParseStruct(options, args[0], args[1:])
	if options.IsValid(flagSet.Name()) {
		return options
	}
	remainArgs := flagSet.Args()
	return c.resolve(remainArgs)
}

// IsValid 是否有效
func (c *CompileCmd) IsValid() bool {
	return c.options != nil
}

// Execute 执行注入器, 返回修改后的参数
func (c *CompileCmd) Execute() []string {
	options := c.options
	if options.Debug {
		log.PrintArg("workdir", options.GetWorkDir())
	}
	pathIdx := c.indexPath(options.Args)
	if pathIdx == -1 {
		return options.Args
	}
	for _, injector := range c.injectors {
		api.ExecInjector(injector, options, pathIdx)
	}
	return options.Args
}

// indexPath 返回路径所在的索引
func (c *CompileCmd) indexPath(args []string) int {
	for _, f := range flags {
		i := stringutil.LastIndexOf(args, f)
		if i != -1 {
			return i + 1 //next item is file path
		}
	}
	return -1
}
