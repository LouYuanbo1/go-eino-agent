package tools

import (
	"context"
	"log"

	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewGitHubTools(ctx context.Context, client *client.Client, toolNames ...string) (map[string]tool.InvokableTool, error) {
	tools, err := mcpp.GetTools(ctx, &mcpp.Config{
		Cli: client,
		Meta: &mcp.Meta{
			AdditionalFields: map[string]any{
				"source": "eino-github-demo",
			},
		},
	})
	if err != nil {
		log.Fatalf("转换 MCP 工具失败: %v", err)
	}
	toolMap := make(map[string]tool.InvokableTool)
	for _, name := range toolNames {
		toolMap[name] = nil
	}
	for _, t := range tools {
		info, err := t.Info(ctx)
		if err != nil {
			log.Printf("获取工具信息失败: %v", err)
			continue
		}
		/*
			fmt.Printf("工具描述: %v\n", info.Desc)
			fmt.Printf("工具额外信息: %v\n", info.Extra)
			jsonSchema, err := info.ParamsOneOf.ToJSONSchema()
			if err != nil {
				log.Printf("转换参数失败: %v", err)
			} else if jsonSchema != nil {
				// 打印人类可读的 JSON
				schemaJSON, _ := json.MarshalIndent(jsonSchema, "", "  ")
				fmt.Printf("工具参数 Schema:\n%s\n", string(schemaJSON))
			}
		*/
		if _, ok := toolMap[info.Name]; ok {
			if invokableTool, ok := t.(tool.InvokableTool); ok {
				toolMap[info.Name] = invokableTool
			}
		}
	}
	return toolMap, nil
}
