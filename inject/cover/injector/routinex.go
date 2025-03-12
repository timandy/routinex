package injector

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/timandy/routinex/inject/api"
	"github.com/timandy/routinex/inject/cover"
	"github.com/timandy/routinex/tools/exec"
	"github.com/timandy/routinex/tools/log"
	"github.com/timandy/routinex/tools/os"
	"github.com/timandy/routinex/tools/slices"
)

type RoutineXInjector struct {
}

func NewRoutineXInjector() api.Injector {
	return &RoutineXInjector{}
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) PreHandlePackage(options api.CmdOptions, result *api.InjectResult) bool {
	pkg := options.GetPackage()
	return pkg == "routine" || pkg == "g" || pkg == "github.com/timandy/routine" || pkg == "github.com/timandy/routine/g"
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) PreHandleFile(path string, idx int, options api.CmdOptions, result *api.InjectResult) bool {
	return true
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) HandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options api.CmdOptions, result *api.InjectResult) bool {
	for _, comment := range af.Comments {
		for _, c := range comment.List {
			if r.hasTag(c) {
				return true
			}
		}
	}
	return false
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) PostHandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options api.CmdOptions, result *api.InjectResult) {
	srcDir := filepath.Dir(path)
	srcShortName := filepath.Base(path)
	destShortName := srcShortName[:len(srcShortName)-len(".go")] + "_link.go"
	destPath := filepath.Join(srcDir, destShortName)
	//存在可替换的
	if !os.IsFile(destPath) {
		return
	}
	output := (options.(*cover.CoverOptions)).Output
	if output != "" {
		//go1.19- cover only support one file
		args := options.GetArgs()
		if idx <= 0 || args[idx-1] != output {
			return
		}
		if !strings.HasSuffix(output, ".cover.go") {
			return
		}
		outputSrcDir := filepath.Dir(output)
		outputSrcShortName := filepath.Base(output)
		outputDestShortName := outputSrcShortName[:len(outputSrcShortName)-len(".cover.go")] + "_link.cover.go"
		outputDestPath := filepath.Join(outputSrcDir, outputDestShortName)
		//额外执行命令
		extraArgs := slices.Clone(args)
		extraArgs[idx-1] = outputDestPath
		extraArgs[idx] = destPath
		if options.IsDebug() || options.IsVerbose() {
			log.Infof("cover: insert counters '%v' to '%v'", destShortName, outputDestShortName)
		}
		exec.RunCmd(extraArgs)
	} else {
		//go1.20+ cover support pkgcfg option
		result.ReplaceFiles[idx] = destPath
		if options.IsDebug() || options.IsVerbose() {
			log.Infof("cover: replace source '%v' with '%v'", srcShortName, destShortName)
		}
	}
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) PostHandlePackage(options api.CmdOptions, result *api.InjectResult) {
}

// hasTag 是否有 !routinex 编译标记
func (r *RoutineXInjector) hasTag(comment *ast.Comment) bool {
	return comment != nil && strings.TrimSpace(comment.Text) == "//go:build !routinex"
}
