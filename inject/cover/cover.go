package cover

import (
	"github.com/timandy/routinex/inject/api"
	"github.com/timandy/routinex/tools/flag"
	"github.com/timandy/routinex/tools/opt"
	"github.com/timandy/routinex/tools/stringutil"
)

var flags = []string{"-o", "-outfilelist"}

type CoverCmd struct {
	injectors []api.Injector
	options   *CoverOptions
}

// NewCoverCmd 构造函数
func NewCoverCmd(injectors []api.Injector) *CoverCmd {
	return &CoverCmd{injectors: injectors}
}

// Resolve 解析参数
func (c *CoverCmd) Resolve(args []string, app *opt.AppOptions) {
	options := c.resolve(args)
	if options != nil {
		options.Debug = app.Debug
		options.Verbose = app.Verbose
		options.Args = args
	}
	c.options = options
}

// resolve 解析参数-递归
func (c *CoverCmd) resolve(args []string) *CoverOptions {
	if len(args) == 0 {
		return nil
	}
	options := &CoverOptions{}
	flagSet := flag.ParseStruct(options, args[0], args[1:])
	options.ReadConfig(flagSet.Args())
	if options.IsValid(flagSet.Name()) {
		return options
	}
	remainArgs := flagSet.Args()
	return c.resolve(remainArgs)
}

// IsValid 是否有效
func (c *CoverCmd) IsValid() bool {
	return c.options != nil
}

// Execute 执行注入器, 返回修改后的参数
func (c *CoverCmd) Execute() []string {
	options := c.options
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
func (c *CoverCmd) indexPath(args []string) int {
	for _, f := range flags {
		i := stringutil.LastIndexOf(args, f)
		if i != -1 {
			return i + 2 //skip two item is file path
		}
	}
	return -1
}
