package main

import (
	"os"

	"github.com/timandy/routinex/inject"
	"github.com/timandy/routinex/tools/exec"
	"github.com/timandy/routinex/tools/flag"
	"github.com/timandy/routinex/tools/log"
	"github.com/timandy/routinex/tools/opt"
	"github.com/timandy/routinex/tools/slices"
)

func main() {
	// parse app options
	args := slices.Filter(os.Args, func(s string) bool { return s != "" })
	appOpt := &opt.AppOptions{}
	flagSet := flag.ParseStruct(appOpt, args[0], args[1:])
	// print entry args
	if appOpt.Debug {
		log.PrintArgs("entry", args)
	}
	// print usage
	if appOpt.Help {
		flagSet.SortFlags = false
		flag.PrintUsage(flagSet)
		return
	}
	// print before args
	remainArgs := flagSet.Args()
	if appOpt.Debug {
		log.PrintArgs("before", remainArgs)
	}
	// return when no remain args
	if len(remainArgs) == 0 {
		if appOpt.Debug {
			log.Info("no remain args and exit")
		}
		return
	}
	// exists remained args, run the cmd finally
	defer func() {
		// print after args
		if appOpt.Debug {
			log.PrintArgs("after", remainArgs)
		}
		exec.RunCmd(remainArgs)
	}()
	// exec inject
	remainArgs = inject.Execute(remainArgs, appOpt)
}
