package exec

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func RunCmd(args []string) {
	if len(args) == 0 {
		return
	}
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	startCmd(cmd)
	statusCode := waitCmd(cmd)
	if statusCode == 0 {
		return
	}
	os.Exit(statusCode)
}

func RunCmdOutput(args []string) string {
	if len(args) == 0 {
		return ""
	}
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	out := bytes.Buffer{}
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	startCmd(cmd)
	statusCode := waitCmd(cmd)
	if statusCode == 0 {
		return strings.TrimSpace(out.String())
	}
	os.Exit(statusCode)
	return ""
}

func startCmd(cmd *exec.Cmd) {
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

func waitCmd(cmd *exec.Cmd) int {
	// wait exit and return exit code
	err := cmd.Wait()
	if err == nil {
		return 0
	}
	if status, ok := err.(*exec.ExitError); ok {
		return status.ExitCode()
	}
	return 1
}
