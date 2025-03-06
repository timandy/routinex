package api

import "github.com/timandy/routinex/tools/opt"

type Cmd interface {
	// Resolve 解析参数
	Resolve(args []string, app *opt.AppOptions)
	// IsValid 是否有效
	IsValid() bool
	// Execute 执行注入器, 返回修改后的参数
	Execute() []string
}
