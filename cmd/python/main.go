package main

import (
	"context"
	"fmt"
	"os"

	pythonAgent "github.com/LouYuanbo1/go-eino-agent/agents/python"
	"github.com/LouYuanbo1/go-eino-agent/tools/pyexecutor/local"

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
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v", err)
		return
	}
	pythonAgent := pythonAgent.NewDefaultPythonAgentLocal(ctx, chatModel, &local.OperatorConfig{
		TaskIDFormat:   local.IDFormatTime,
		RootDir:        wd,
		WorkDir:        wd,
		ExecutablePath: "D:\\ANACONDA\\envs\\python-llm\\python.exe",
		//ExecutablePath: "your_python_executable_path",
	})
	fmt.Printf("work dir: %s\n", wd)
	pythonAgent.OutputMessage(ctx, "使用Python生成一个简易的线性回归模型", true)
}
