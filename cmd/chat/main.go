package main

import (
	"context"
	"fmt"

	chatAgent "github.com/LouYuanbo1/go-eino-agent/agents/chat"
	"github.com/cloudwego/eino-ext/components/model/ollama"
)

func main() {
	ctx := context.Background()
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "qwen3.5:2b",
	})
	if err != nil {
		fmt.Printf("Error creating chat model: %v", err)
		return
	}
	chatAgent := chatAgent.NewDefaultChatAgent(ctx, chatModel)
	chatAgent.OutputMessage(ctx, "你好,帮我搜索一下eino是什么", false)
}
