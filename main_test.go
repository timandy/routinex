package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timandy/routinex/tools/consts"
	"github.com/timandy/routinex/tools/exec"
	"github.com/timandy/routinex/tools/file"
)

var goToolDir = exec.RunCmdOutput([]string{"go", "env", "GOTOOLDIR"})

func TestHelp(t *testing.T) {
	tracker := file.NewFileTracker(&os.Stderr)
	tracker.Begin()
	defer tracker.End()
	//
	args := []string{"routinex", "-h", "/demo", "-p", "ttt", "go", "version"}
	os.Args = args
	main()
	//
	output := tracker.Value()
	assert.Contains(t, output, "Usage: routinex")
}

func TestOtherCmd(t *testing.T) {
	tracker := file.NewFileTracker(&os.Stdout)
	tracker.Begin()
	defer tracker.End()
	//
	args := []string{"routinex", "-p", "/demo", "-p", "ttt", "git", "version"}
	os.Args = args
	main()
	//
	output := tracker.Value()
	assert.Contains(t, output, "git version")
}

func TestOtherCmdHelp(t *testing.T) {
	tracker := file.NewFileTracker(&os.Stdout)
	tracker.Begin()
	defer tracker.End()
	//
	args := []string{"routinex", "-p", "/demo", "-p", "ttt", "git", "-h"}
	os.Args = args
	main()
	//
	output := tracker.Value()
	assert.Contains(t, output, "usage: git")
}

func TestCompileCmdHelp(t *testing.T) {
	compilePath := filepath.Join(goToolDir, consts.CompileName)
	args := []string{"routinex", "-p", "/demo", "-p", "ttt", compilePath, "-h"}
	os.Args = args
	// expect exit 2
	// main()
}
