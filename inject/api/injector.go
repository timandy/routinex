package api

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/timandy/routinex/tools/slices"
)

type Injector interface {
	PreHandlePackage(options CmdOptions, result *InjectResult) bool

	PreHandleFile(path string, idx int, options CmdOptions, result *InjectResult) bool

	HandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options CmdOptions, result *InjectResult) bool

	PostHandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options CmdOptions, result *InjectResult)

	PostHandlePackage(options CmdOptions, result *InjectResult)
}

// ExecInjector 执行单个注入器
func ExecInjector(injector Injector, options CmdOptions, pathIdx int) {
	// define result
	result := NewInjectResult()
	// proc args after exec
	args := options.GetArgs()
	defer func() {
		for idx, path := range result.ReplaceFiles {
			args[idx] = path
		}
		args = append(args, result.ExtraFiles...)
		args = slices.Filter(args, func(str string) bool { return str != "" })
		options.SetArgs(args)
	}()
	// verify this injector can handle the package
	if !injector.PreHandlePackage(options, result) {
		return
	}
	for idx, length := pathIdx, len(args); idx < length; idx++ {
		path := args[idx]
		if !strings.HasSuffix(path, ".go") {
			continue
		}
		if !injector.PreHandleFile(path, idx, options, result) {
			continue
		}
		// parse the ast file
		fset := token.NewFileSet()
		af, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		if !injector.HandleFile(path, idx, fset, af, options, result) {
			continue
		}
		injector.PostHandleFile(path, idx, fset, af, options, result)
	}
	injector.PostHandlePackage(options, result)
}
