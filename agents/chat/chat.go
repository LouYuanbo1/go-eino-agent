package chatAgent

import (
	"context"
	"fmt"

	"github.com/LouYuanbo1/go-eino-agent/prints"
	"github.com/LouYuanbo1/go-eino-agent/tools/search/duckduckgo"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

type ChatAgent struct {
	Agent *adk.ChatModelAgent
}

func NewChatAgent(ctx context.Context, config *adk.ChatModelAgentConfig) *ChatAgent {
	agent, err := adk.NewChatModelAgent(ctx, config)
	if err != nil {
		fmt.Printf("Error creating chat model agent: %v", err)
		return nil
	}
	return &ChatAgent{Agent: agent}
}

func NewDefaultChatAgent(ctx context.Context, model model.ToolCallingChatModel) *ChatAgent {
	ddgTool, err := duckduckgo.NewDefaultDuckDuckGoTool(ctx)
	if err != nil {
		fmt.Printf("Error creating DuckDuckGo tool: %v", err)
		return nil
	}

	instruction :=
		`你是一个专业、友好且乐于助人的AI聊天助手。请遵循以下原则与用户交流：
		1. 【角色定位】保持中立、客观，不提供医疗、法律等专业领域的诊断或建议
		2. 【回答风格】语言简洁清晰，逻辑严谨，必要时分点陈述；避免过度冗长
		3. 【知识边界】对于不确定的信息，坦诚说明"我不确定"，并建议用户查阅权威来源
		4. 【安全合规】拒绝生成违法、有害、歧视性或侵犯隐私的内容
		5. 【上下文理解】充分利用对话历史，保持回复的连贯性和一致性
		6. 【多轮交互】当用户需求模糊时，主动提问以澄清意图，而非猜测作答
		7. 【语言适配】默认使用用户提问的语言回复；技术类问题优先使用中文解释

		## 可使用的工具
		- DuckDuckGo 搜索：用于一般信息查询、最新动态、专业知识

		如果用户提出代码相关需求，请：
		- 优先使用用户指定的编程语言
		- 提供可运行的代码示例 + 简要注释
		- 说明关键逻辑和潜在注意事项

		现在，请开始与用户对话。`
	return NewChatAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "chatAgent",
		Description: "一个基于大模型的聊天智能体",
		Instruction: instruction,
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{ddgTool},
			},
		},
	})
}

func (a *ChatAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	return a.Agent.Run(ctx, input, options...)
}

func (a *ChatAgent) OutputMessage(ctx context.Context, input string, withReasoning bool, options ...adk.AgentRunOption) {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: a.Agent, EnableStreaming: true})
	iter := runner.Query(ctx, input, options...)
	prints.PrintMessages(iter, prints.WithReasoning(withReasoning))
}
