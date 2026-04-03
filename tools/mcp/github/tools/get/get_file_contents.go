// github/mcp_parser.go
package get

import (
	"context"

	githubutil "github.com/LouYuanbo1/go-eino-agent/tools/mcp/github/utils"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ============ 7. Wrapper 集成 ============

type GithubGetFuncWrapper struct {
	base tool.InvokableTool
}

func NewGithubGetFuncWrapper(base tool.InvokableTool) *GithubGetFuncWrapper {
	return &GithubGetFuncWrapper{base: base}
}

func (w *GithubGetFuncWrapper) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	rawOutput, err := w.base.InvokableRun(ctx, argumentsInJSON, opts...)
	if err != nil {
		return "", err
	}
	// ✅ 使用通用解析入口，自动处理目录/文件/搜索
	return githubutil.ParseMCPGitHubContents(rawOutput)
}

func (w *GithubGetFuncWrapper) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return w.base.Info(ctx)
}
