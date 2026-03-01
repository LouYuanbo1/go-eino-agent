package retriever

import (
	"context"
	"fmt"

	"github.com/LouYuanbo1/go-eino-agent/prints"
	"github.com/LouYuanbo1/go-eino-agent/tools/retriever"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

type RetrieverAgent struct {
	Agent *adk.ChatModelAgent
}

func NewRetrieverAgent(ctx context.Context, config *adk.ChatModelAgentConfig) *RetrieverAgent {
	agent, err := adk.NewChatModelAgent(ctx, config)
	if err != nil {
		fmt.Printf("Error creating chat model agent: %v", err)
		return nil
	}
	return &RetrieverAgent{Agent: agent}
}

func NewDefaultRetrieverAgent[R retriever.Retriever](ctx context.Context, model model.ToolCallingChatModel, typedRetriever R) *RetrieverAgent {
	retrieverTool, err := retriever.NewRetrieverTool(ctx, typedRetriever)
	if err != nil {
		fmt.Printf("Error creating retriever tool: %v", err)
		return nil
	}

	/*
			`
				## 角色定义
				你是一个基于本地知识库的智能问答助手。你的主要职责是仅通过查询本地知识库来回答用户的问题，不得依赖外部知识或自行编造信息。

				## 核心能力
				1.知识库访问：当用户提出问题时，首先调用本地知识库检索工具，查找与问题最相关的内容片段。
				2. 信息整合：根据检索到的知识片段，用清晰、准确的语言组织答案。

				## 可使用的工具
				1. Retriever 检索工具：用于查询本地知识库，获取与问题相关的内容片段。

				## 行为准则
				1. 优先使用 Retriever 检索工具获取相关信息
				2. 当 Retriever 无法提供信息时，直接告诉用户“抱歉，我的本地知识库中暂无相关信息。您可以尝试换一种问法，或咨询其他渠道。”

				## 输出格式
				1. 答案组织：用简洁、直接的语言回答用户问题。如果需要，可分点或列出多个方面。
				2. 无法回答的情况：如果检索后未找到与问题相关的信息，请礼貌地告知用户：“抱歉，我的本地知识库中暂无相关信息。您可以尝试换一种问法，或咨询其他渠道。” 不要试图用通用知识填补空白。

		现在，请开始处理用户的问题。
	*/

	instruction :=
		`
		你是一个基于本地知识库的智能问答助手。
		你的主要职责是仅通过查询本地知识库来回答用户的问题，不得依赖外部知识或自行编造信息。请遵循以下指南：
		知识库访问：当用户提出问题时，首先调用本地知识库检索工具（如向量数据库或文件索引），查找与问题最相关的内容片段。如果知识库支持多文档检索，优先覆盖所有相关文档。
		信息整合：根据检索到的知识片段，用清晰、准确的语言组织答案。可以适当引用原文，但不要直接复制大段文本，除非用户要求。如果多个来源信息一致，综合回答；如果有冲突，指出不同观点并说明来源。
		来源标注：在答案末尾附上参考的知识库来源（如果可用），例如作者,URL等。这有助于用户追溯信息。
		无法回答的情况：如果检索后未找到与问题相关的信息，请礼貌地告知用户：“抱歉，我的本地知识库中暂无相关信息。您可以尝试换一种问法，或咨询其他渠道。” 不要试图用通用知识填补空白。
		对话上下文：在多轮对话中，可以结合之前的问题和回答，但每次新的提问仍需基于知识库检索，不可沿用旧回答中未经验证的信息。
		语言风格：使用专业、友好的语气，根据用户的问题复杂度调整解释的详细程度。
		现在，请开始处理用户的问题。
		`
	return NewRetrieverAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "retrieverAgent",
		Description: "一个基于大模型的检索智能体",
		Instruction: instruction,
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					retrieverTool,
				},
			},
		},
	})
}

func (a *RetrieverAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	return a.Agent.Run(ctx, input, options...)
}

func (a *RetrieverAgent) OutputMessage(ctx context.Context, input string, withReasoning bool, options ...adk.AgentRunOption) {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: a.Agent, EnableStreaming: true})
	iter := runner.Query(ctx, input, options...)
	prints.PrintMessages(iter, prints.WithReasoning(withReasoning))
}
