package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/timandy/routinex/tools/trace"
)

var traceId = trace.NewTraceId()

func Info(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, "--> ["+traceId+"] "+msg)
}

func Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	Info(msg)
}

func PrintArg(prefix string, arg string) {
	_, _ = fmt.Fprintln(os.Stderr, "==> ["+traceId+"] ["+prefix+"] "+arg)
}

func PrintArgs(prefix string, args []string) {
	builder := strings.Builder{}
	if len(args) > 0 {
		builder.WriteString("\"")
		builder.WriteString(args[0])
		builder.WriteString("\"")
	}
	for _, arg := range args[1:] {
		builder.WriteString(" \"")
		builder.WriteString(arg)
		builder.WriteString("\"")
	}
	PrintArg(prefix, builder.String())
}
