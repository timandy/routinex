package injector

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/timandy/routinex/inject/api"
	"github.com/timandy/routinex/tools/astutil"
	"github.com/timandy/routinex/tools/log"
	"github.com/timandy/routinex/tools/os"
	"github.com/timandy/routinex/tools/stringutil"
)

type RuntimeInjector struct {
}

func NewRuntimeInjector() api.Injector {
	return &RuntimeInjector{}
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInjector) PreHandlePackage(options api.CmdOptions, result *api.InjectResult) bool {
	pkg := options.GetPackage()
	return pkg == "runtime"
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInjector) PreHandleFile(path string, idx int, options api.CmdOptions, result *api.InjectResult) bool {
	return strings.HasSuffix(path, "runtime2.go") || strings.HasSuffix(path, "proc.go")
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInjector) HandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options api.CmdOptions, result *api.InjectResult) bool {
	handled := false
	ast.Inspect(af, func(node ast.Node) bool {
		if r.handleNode(node, options) {
			handled = true
			return false
		}
		return true
	})
	return handled
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInjector) PostHandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options api.CmdOptions, result *api.InjectResult) {
	srcShortName := filepath.Base(path)
	destPath := filepath.Join(options.GetWorkDir(), srcShortName)
	astutil.SaveAs(destPath, fset, af)
	result.ReplaceFiles[idx] = destPath
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInjector) PostHandlePackage(options api.CmdOptions, result *api.InjectResult) {
	code := stringutil.ExecuteTemplate(`package runtime

import _ "unsafe"

//go:nosplit
//go:linkname getg0
func getg0() interface{} {
	return g0
}

//go:nosplit
//go:linkname getgp
func getgp() *g {
	return getg()
}
`, nil)
	// save file
	destShortName := "runtime_routine.go"
	destPath := filepath.Join(options.GetWorkDir(), destShortName)
	os.WriteFile(destPath, code)
	result.ExtraFiles = append(result.ExtraFiles, destPath)
	if options.IsDebug() || options.IsVerbose() {
		log.Info("compile: create function 'runtime.getg0'")
		log.Info("compile: create function 'runtime.getgp'")
	}
}

func (r *RuntimeInjector) handleNode(node ast.Node, options api.CmdOptions) bool {
	switch n := node.(type) {
	case *ast.TypeSpec:
		ident := n.Name
		if ident == nil || ident.Name != "g" {
			return false
		}
		st, isSt := n.Type.(*ast.StructType)
		if !isSt {
			return false
		}
		fields := st.Fields
		if fields == nil {
			return false
		}
		fieldList := fields.List
		if len(fieldList) == 0 {
			return false
		}
		threadLocalsField := astutil.CreateField("threadLocals", "unsafe.Pointer")
		inheritableThreadLocalsField := astutil.CreateField("inheritableThreadLocals", "unsafe.Pointer")
		fields.List = append(fieldList, threadLocalsField, inheritableThreadLocalsField)
		if options.IsDebug() || options.IsVerbose() {
			log.Info("compile: enhance struct 'runtime.g' add field 'threadLocals unsafe.Pointer'")
			log.Info("compile: enhance struct 'runtime.g' add field 'inheritableThreadLocals unsafe.Pointer'")
		}
		return true
	case *ast.FuncDecl:
		// check name
		ident := n.Name
		if ident == nil || ident.Name != "goexit0" {
			return false
		}
		// check type not nil
		funcType := n.Type
		if funcType == nil {
			return false
		}
		// check no results
		results := funcType.Results
		if results != nil && len(results.List) > 0 {
			return false
		}
		// check only one params
		params := funcType.Params
		if params == nil || len(params.List) != 1 {
			return false
		}
		paramField0 := params.List[0]
		if paramField0 == nil || len(paramField0.Names) != 1 {
			return false
		}
		gp := paramField0.Names[0]
		// check body not empty
		body := n.Body
		if body == nil || len(body.List) == 0 {
			return false
		}
		// add set nil statements
		threadLocalsStmt := astutil.CreateAssignNilStmt(gp, "threadLocals")
		inheritableThreadLocalsStmt := astutil.CreateAssignNilStmt(gp, "inheritableThreadLocals")
		body.List = append([]ast.Stmt{threadLocalsStmt, inheritableThreadLocalsStmt}, body.List...)
		if options.IsDebug() || options.IsVerbose() {
			log.Info("compile: enhance function 'runtime.goexit0' add statement 'gp.threadLocals = nil'")
			log.Info("compile: enhance function 'runtime.goexit0' add statement 'gp.inheritableThreadLocals = nil'")
		}
		return true
	default:
		return false
	}
}
