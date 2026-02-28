package duckduckgo

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/tool"
)

func NewDefaultDuckDuckGoTool(ctx context.Context) (tool.InvokableTool, error) {
	ddgConfig := &duckduckgo.Config{
		ToolName:   "duckduckgo_search",           // 工具名称，模型会通过这个名字来调用
		ToolDesc:   "搜索网络以获取最新信息，适用于新闻、文章和一般知识查询", // 工具描述，模型会根据描述决定是否使用此工具
		Timeout:    30,                            // 单次搜索超时时间（秒）
		MaxResults: 5,                             // 最多返回的搜索结果数量
		Region:     duckduckgo.RegionCN,           // 搜索区域，例如中国区 [citation:6]
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		}, // 自定义 HTTP 客户端，可选
	}

	searchTool, err := duckduckgo.NewTextSearchTool(ctx, ddgConfig)
	if err != nil {
		log.Fatalf("创建 DuckDuckGo 工具失败: %v", err)
	}
	return searchTool, nil
}

func NewDuckDuckGoTool(ctx context.Context, config *duckduckgo.Config) (tool.InvokableTool, error) {
	searchTool, err := duckduckgo.NewTextSearchTool(ctx, config)
	if err != nil {
		log.Fatalf("创建 DuckDuckGo 工具失败: %v", err)
	}
	return searchTool, nil
}
