package get

import (
	"context"
	"fmt"

	githubutil "github.com/LouYuanbo1/go-eino-agent/tools/mcp/github/utils"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type OverviewParams struct {
	Repo  string `json:"repo" jsonschema:"description=要获取概述的仓库名称"`
	Owner string `json:"owner" jsonschema:"description=仓库的所有者名称"`
}

func OverviewFunc(ctx context.Context, getTool tool.InvokableTool) func(ctx context.Context, params *OverviewParams) (string, error) {
	return func(ctx context.Context, params *OverviewParams) (string, error) {
		output, err := getTool.InvokableRun(ctx, fmt.Sprintf(`{"repo": "%s", "owner": "%s","path":"README.md"}`, params.Repo, params.Owner))
		if err != nil {
			return "", err
		}
		return githubutil.ParseMCPGitHubContents(output)
	}
}

func NewOverviewTool(ctx context.Context, getTool tool.InvokableTool) (tool.InvokableTool, error) {
	retrieverTool, err := utils.InferTool(
		"overview", // tool name
		`Overview is used to search the overview (e.g. README.md) of a repository; 
		This tool is used when a quick overview and summary of information is needed (e.g., "What does this project do?" "What are its main functions?").`,
		OverviewFunc(ctx, getTool))
	if err != nil {
		return nil, err
	}
	return retrieverTool, nil
}
