package main

import (
	"context"
	"fmt"

	searchAgent "github.com/LouYuanbo1/go-eino-agent/agents/search"
	"github.com/LouYuanbo1/go-eino-agent/tools/search/spider"
	"github.com/cloudwego/eino-ext/components/model/ollama"
)

func main() {
	ctx := context.Background()
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "qwen3:1.7b",
	})
	if err != nil {
		fmt.Printf("Error creating chat model: %v", err)
		return
	}
	searchAgent := searchAgent.NewDefaultSearchAgent(ctx, chatModel, &spider.SpiderConfig{
		Bin: "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
	})
	searchAgent.OutputMessage(ctx, "https://cloudwego.io/zh/docs/eino/quick_start/agent_llm_with_tools,讲解一下这个页面的信息", true)
}
