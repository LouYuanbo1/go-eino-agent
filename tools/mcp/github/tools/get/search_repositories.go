package get

import (
	"context"

	githubutil "github.com/LouYuanbo1/go-eino-agent/tools/mcp/github/utils"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type GithubSearchFuncWrapper struct {
	base tool.InvokableTool
}

func NewGithubSearchFuncWrapper(base tool.InvokableTool) *GithubSearchFuncWrapper {
	return &GithubSearchFuncWrapper{base: base}
}

func (w *GithubSearchFuncWrapper) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	rawOutput, err := w.base.InvokableRun(ctx, argumentsInJSON, opts...)
	if err != nil {
		return "", err // 错误也返回为正常结构
	}
	return githubutil.ParseMCPGitHubSearch(rawOutput)
}

func (w *GithubSearchFuncWrapper) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return w.base.Info(ctx)
}
