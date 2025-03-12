package injector

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/timandy/routinex/inject/api"
	"github.com/timandy/routinex/tools/log"
	"github.com/timandy/routinex/tools/os"
)

type RoutineXInjector struct {
}

func NewRoutineXInjector() api.Injector {
	return &RoutineXInjector{}
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) PreHandlePackage(options api.CmdOptions, result *api.InjectResult) bool {
	pkg := options.GetPackage()
	return pkg == "github.com/timandy/routine" || pkg == "github.com/timandy/routine/g"
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
	destShortName := ""
	if strings.HasSuffix(srcShortName, ".cover.go") {
		//插桩文件
		destShortName = srcShortName[:len(srcShortName)-len(".cover.go")] + "_link.cover.go"
	} else if strings.HasSuffix(srcShortName, "_test.go") {
		//测试文件
		destShortName = srcShortName[:len(srcShortName)-len("_test.go")] + "_link_test.go"
	} else {
		//源码文件
		destShortName = srcShortName[:len(srcShortName)-len(".go")] + "_link.go"
	}
	destPath := filepath.Join(srcDir, destShortName)
	//存在可替换的
	if os.IsFile(destPath) {
		result.ReplaceFiles[idx] = destPath
		if options.IsDebug() || options.IsVerbose() {
			log.Infof("compile: replace source '%v' with '%v'", srcShortName, destShortName)
		}
		return
	}
	//不存在可替换的
	result.ReplaceFiles[idx] = ""
	if options.IsDebug() || options.IsVerbose() {
		log.Infof("compile: remove source '%v'", srcShortName)
	}
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) PostHandlePackage(options api.CmdOptions, result *api.InjectResult) {
}

// hasTag 是否有 !routinex 编译标记
func (r *RoutineXInjector) hasTag(comment *ast.Comment) bool {
	return comment != nil && strings.TrimSpace(comment.Text) == "//go:build !routinex"
}
