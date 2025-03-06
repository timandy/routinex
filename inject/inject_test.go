package inject

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/timandy/routinex/tools/consts"
	"github.com/timandy/routinex/tools/exec"
	"github.com/timandy/routinex/tools/file"
	"github.com/timandy/routinex/tools/opt"
)

const routineRepo = "https://github.com/timandy/routine.git"

var (
	appOpts      = &opt.AppOptions{Verbose: true}
	goRoot       = exec.RunCmdOutput([]string{"go", "env", "GOROOT"})
	goToolDir    = exec.RunCmdOutput([]string{"go", "env", "GOTOOLDIR"})
	routineDir   = filepath.Join(os.TempDir(), "routine")
	outDir       = filepath.Join(os.TempDir(), "out"+strconv.FormatInt(time.Now().UnixNano(), 10))
	pkgcfgFile   = filepath.Join(outDir, "pkgcfg.txt")
	coverOutFile = filepath.Join(outDir, "coveroutfiles.txt")
)

func init() {
	//克隆仓库
	if err := os.RemoveAll(routineDir); err != nil {
		panic(err)
	}
	exec.RunCmd([]string{"git", "clone", routineRepo, routineDir})
	//创建输出目录
	if err := os.MkdirAll(outDir, 0755); err != nil {
		panic(err)
	}
	//写入覆盖率测试配置
	if err := os.WriteFile(pkgcfgFile, []byte("{\"PkgPath\":\"github.com/timandy/routine/g\"}"), 0644); err != nil {
		panic(err)
	}
}

func TestCompileRuntime(t *testing.T) {
	tracker := file.NewFileTracker(&os.Stderr)
	tracker.Begin()
	defer tracker.End()
	//
	result := Execute([]string{
		filepath.Join(goToolDir, consts.CompileName),
		"-o", filepath.Join(outDir, "_pkg.a"),
		"-p", "runtime",
		"-pack", "-asmhdr",
		filepath.Join(goRoot, "src/runtime/go_asm.h"),
		filepath.Join(goRoot, "src/runtime/proc.go"),
		filepath.Join(goRoot, "src/runtime/runtime2.go"),
	}, appOpts)
	expect := []string{
		filepath.Join(goToolDir, consts.CompileName),
		"-o", filepath.Join(outDir, "_pkg.a"),
		"-p", "runtime",
		"-pack", "-asmhdr",
		filepath.Join(goRoot, "src/runtime/go_asm.h"),
		filepath.Join(outDir, "proc.go"),
		filepath.Join(outDir, "runtime2.go"),
		filepath.Join(outDir, "runtime_routine.go"),
	}
	assert.Equal(t, expect, result)
	//
	output := tracker.Value()
	lines := strings.Split(output, "\n")
	assert.Equal(t, 7, len(lines))
	assert.True(t, strings.HasSuffix(lines[0], "compile: enhance function 'runtime.goexit0' add statement 'gp.threadLocals = nil'"))
	assert.True(t, strings.HasSuffix(lines[1], "compile: enhance function 'runtime.goexit0' add statement 'gp.inheritableThreadLocals = nil'"))
	assert.True(t, strings.HasSuffix(lines[2], "compile: enhance struct 'runtime.g' add field 'threadLocals unsafe.Pointer'"))
	assert.True(t, strings.HasSuffix(lines[3], "compile: enhance struct 'runtime.g' add field 'inheritableThreadLocals unsafe.Pointer'"))
	assert.True(t, strings.HasSuffix(lines[4], "compile: create function 'runtime.getg0'"))
	assert.True(t, strings.HasSuffix(lines[5], "compile: create function 'runtime.getgp'"))
	assert.Empty(t, lines[6])
}

func TestCompileRoutine_Source(t *testing.T) {
	tracker := file.NewFileTracker(&os.Stderr)
	tracker.Begin()
	defer tracker.End()
	//
	result := Execute([]string{
		filepath.Join(goToolDir, consts.CompileName),
		"-o", filepath.Join(outDir, "_pkg.a"),
		"-p", "github.com/timandy/routine/g",
		"-pack", "-asmhdr",
		filepath.Join(goRoot, "src/runtime/go_asm.h"),
		filepath.Join(routineDir, "g/g.go"),
		filepath.Join(routineDir, "g/reflect.go"),
	}, appOpts)
	expect := []string{
		filepath.Join(goToolDir, consts.CompileName),
		"-o", filepath.Join(outDir, "_pkg.a"),
		"-p", "github.com/timandy/routine/g",
		"-pack", "-asmhdr",
		filepath.Join(goRoot, "src/runtime/go_asm.h"),
		filepath.Join(routineDir, "g/g_link.go"),
		filepath.Join(routineDir, "g/reflect.go"),
	}
	assert.Equal(t, expect, result)
	//
	output := tracker.Value()
	lines := strings.Split(output, "\n")
	assert.Equal(t, 2, len(lines))
	assert.True(t, strings.HasSuffix(lines[0], "compile: replace source 'g.go' with 'g_link.go'"))
	assert.Empty(t, lines[1])
}

func TestCompileRoutine_Test(t *testing.T) {
	remainArgs := Execute([]string{
		filepath.Join(goToolDir, consts.CoverName),
		"-mode", "atomic",
		"-var", "GoCover_0_653661653238656532343338",
		"-o", filepath.Join(outDir, "g.cover.go"),
		filepath.Join(routineDir, "g/g.go"),
	}, appOpts)
	exec.RunCmd(remainArgs)
	//
	tracker := file.NewFileTracker(&os.Stderr)
	tracker.Begin()
	defer tracker.End()
	//
	result := Execute([]string{
		filepath.Join(goToolDir, consts.CompileName),
		"-o", filepath.Join(outDir, "_pkg.a"),
		"-p", "github.com/timandy/routine/g",
		"-pack",
		filepath.Join(goRoot, "src/runtime/go_asm.h"),
		filepath.Join(outDir, "g.cover.go"),
		filepath.Join(routineDir, "g/g_test.go"),
	}, appOpts)
	expect := []string{
		filepath.Join(goToolDir, consts.CompileName),
		"-o", filepath.Join(outDir, "_pkg.a"),
		"-p", "github.com/timandy/routine/g",
		"-pack",
		filepath.Join(goRoot, "src/runtime/go_asm.h"),
		filepath.Join(outDir, "g_link.cover.go"),
		filepath.Join(routineDir, "g/g_test.go"),
	}
	assert.Equal(t, expect, result)
	//
	output := tracker.Value()
	lines := strings.Split(output, "\n")
	assert.Equal(t, 2, len(lines))
	assert.True(t, strings.HasSuffix(lines[0], "compile: replace source 'g.cover.go' with 'g_link.cover.go'"))
	assert.Empty(t, lines[1])
}

func TestCoverRoutine_Go1_18(t *testing.T) {
	tracker := file.NewFileTracker(&os.Stderr)
	tracker.Begin()
	defer tracker.End()
	//
	result := Execute([]string{
		filepath.Join(goToolDir, consts.CoverName),
		"-mode", "atomic",
		"-var", "GoCover_0_653661653238656532343338",
		"-o", filepath.Join(outDir, "g.cover.go"),
		filepath.Join(routineDir, "g/g.go"),
	}, appOpts)
	expect := []string{
		filepath.Join(goToolDir, consts.CoverName),
		"-mode", "atomic",
		"-var", "GoCover_0_653661653238656532343338",
		"-o", filepath.Join(outDir, "g.cover.go"),
		filepath.Join(routineDir, "g/g.go"),
	}
	assert.Equal(t, expect, result)
	//
	output := tracker.Value()
	lines := strings.Split(output, "\n")
	assert.Equal(t, 2, len(lines))
	assert.True(t, strings.HasSuffix(lines[0], "cover: insert counters 'g_link.go' to 'g_link.cover.go'"))
	assert.Empty(t, lines[1])
}

func TestCoverRoutine_Go1_20(t *testing.T) {
	tracker := file.NewFileTracker(&os.Stderr)
	tracker.Begin()
	defer tracker.End()
	//
	result := Execute([]string{
		filepath.Join(goToolDir, consts.CoverName),
		"-pkgcfg", pkgcfgFile,
		"-mode", "atomic",
		"-var", "goCover_e6ae28ee2438_",
		"-outfilelist", coverOutFile,
		filepath.Join(routineDir, "g/g.go"), filepath.Join(routineDir, "g/reflect.go"),
	}, appOpts)
	expect := []string{
		filepath.Join(goToolDir, consts.CoverName),
		"-pkgcfg", pkgcfgFile,
		"-mode", "atomic",
		"-var", "goCover_e6ae28ee2438_",
		"-outfilelist", coverOutFile,
		filepath.Join(routineDir, "g/g_link.go"), filepath.Join(routineDir, "g/reflect.go"),
	}
	assert.Equal(t, expect, result)
	//
	output := tracker.Value()
	lines := strings.Split(output, "\n")
	assert.Equal(t, 2, len(lines))
	assert.True(t, strings.HasSuffix(lines[0], "cover: replace source 'g.go' with 'g_link.go'"))
	assert.Empty(t, lines[1])
}
