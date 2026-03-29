package main

import (
	"context"
	"fmt"

	githubAgent "github.com/LouYuanbo1/go-eino-agent/agents/github"
	"github.com/LouYuanbo1/go-eino-agent/config"
	githubutil "github.com/LouYuanbo1/go-eino-agent/tools/mcp/github/utils"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
)

func main() {
	ctx := context.Background()
	config, err := config.InitConfig()
	if err != nil {
		fmt.Printf("Error initializing config: %v", err)
		return
	}
	/*
		chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
			BaseURL: "http://localhost:11434",
			Model:   "qwen3.5:2b",
		})
	*/
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey: config.Deepseek.APIKey,
		Model:  "deepseek-reasoner",
	})
	if err != nil {
		fmt.Printf("Error creating chat model: %v", err)
		return
	}
	cli, err := githubutil.CreateLocalGitHubMCPClient(ctx, config.Github.Token)
	if err != nil {
		fmt.Printf("Error creating GitHub MCP client: %v", err)
		return
	}
	githubAgent := githubAgent.NewDefaultGitHubAgent(ctx, chatModel, cli, config.Github.ToolNames...)
	githubAgent.OutputMessage(ctx, "详细讲解一下字节的eino框架的功能,用中文回答", true)
}
