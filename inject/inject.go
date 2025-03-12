package inject

import (
	"github.com/timandy/routinex/inject/api"
	"github.com/timandy/routinex/inject/compile"
	compileInjector "github.com/timandy/routinex/inject/compile/injector"
	"github.com/timandy/routinex/inject/cover"
	coverInjector "github.com/timandy/routinex/inject/cover/injector"
	"github.com/timandy/routinex/tools/opt"
)

var (
	coverInjectors   = []api.Injector{coverInjector.NewRoutineXInjector()}
	compileInjectors = []api.Injector{compileInjector.NewRuntimeInjector(), compileInjector.NewRoutineXInjector()}
	cmds             = []api.Cmd{cover.NewCoverCmd(coverInjectors), compile.NewCompileCmd(compileInjectors)}
)

func Execute(args []string, app *opt.AppOptions) []string {
	for _, cmd := range cmds {
		cmd.Resolve(args, app)
		if !cmd.IsValid() {
			continue
		}
		return cmd.Execute()
	}
	return args
}
