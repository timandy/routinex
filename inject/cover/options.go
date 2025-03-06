package cover

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/timandy/routinex/tools/json"
	"github.com/timandy/routinex/tools/os"
	"github.com/timandy/routinex/tools/slices"
)

type CoverOptions struct {
	Output  string   `name:"output" shorthand:"o" usage:"file for output; default: stdout"`
	PkgCfg  string   `name:"pkgcfg" shorthand:"pkgcfg" usage:"enable full-package instrumentation mode using params from specified config file"`
	Package string   // package
	Debug   bool     // debug mode enabled or not
	Verbose bool     // verbose mode enabled or not
	Args    []string // remain args exclude the options of current program
}

func (c *CoverOptions) IsDebug() bool {
	return c.Debug
}

func (c *CoverOptions) IsVerbose() bool {
	return c.Verbose
}

func (c *CoverOptions) GetArgs() []string {
	return c.Args
}

func (c *CoverOptions) SetArgs(args []string) {
	c.Args = args
}

func (c *CoverOptions) GetPackage() string {
	return c.Package
}

func (c *CoverOptions) GetWorkDir() string {
	return ""
}

func (c *CoverOptions) IsValid(execName string) bool {
	cmd := filepath.Base(execName)
	if ext := filepath.Ext(cmd); ext != "" {
		cmd = strings.TrimSuffix(cmd, ext)
	}
	return cmd == "cover" && (c.PkgCfg != "" || c.Output != "")
}

func (c *CoverOptions) Clone() *CoverOptions {
	return &CoverOptions{
		Output:  c.Output,
		PkgCfg:  c.PkgCfg,
		Package: c.Package,
		Debug:   c.Debug,
		Verbose: c.Verbose,
		Args:    slices.Clone(c.Args),
	}
}

func (c *CoverOptions) ReadConfig(args []string) {
	c.tryReadPkgCfg()
	if c.GetPackage() != "" {
		return
	}
	c.tryParseSource(args)
}

func (c *CoverOptions) tryReadPkgCfg() {
	if !os.IsFile(c.PkgCfg) {
		return
	}
	bytes := os.ReadFile(c.PkgCfg)
	m := json.Unmarshal[map[string]any](bytes)
	if m == nil {
		return
	}
	pkg, ok := m["PkgPath"]
	if !ok {
		return
	}
	c.Package = pkg.(string)
}

func (c *CoverOptions) tryParseSource(args []string) {
	if len(args) != 1 {
		return
	}
	path := args[0]
	if !os.IsFile(path) {
		return
	}
	if !strings.HasSuffix(path, ".go") {
		return
	}
	// parse the ast file, stop parsing after package clause
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.PackageClauseOnly)
	if err != nil {
		panic(err)
	}
	c.Package = node.Name.Name
}
