package searchAgent

import (
	"context"
	"fmt"

	"github.com/LouYuanbo1/go-eino-agent/prints"

	"github.com/LouYuanbo1/go-eino-agent/tools/search/duckduckgo"
	"github.com/LouYuanbo1/go-eino-agent/tools/search/spider"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

type SearchAgent struct {
	Agent *adk.ChatModelAgent
}

func NewSearchAgent(ctx context.Context, config *adk.ChatModelAgentConfig) *SearchAgent {
	agent, err := adk.NewChatModelAgent(ctx, config)
	if err != nil {
		fmt.Printf("Error creating search agent: %v", err)
		return nil
	}
	return &SearchAgent{Agent: agent}
}

func NewDefaultSearchAgent(ctx context.Context, model model.ToolCallingChatModel, config *spider.SpiderConfig) *SearchAgent {
	ddgTool, err := duckduckgo.NewDefaultDuckDuckGoTool(ctx)
	if err != nil {
		fmt.Printf("Error creating DuckDuckGo tool: %v", err)
		return nil
	}
	spiderTool, err := spider.NewSpiderTool(ctx, config)
	if err != nil {
		fmt.Printf("Error creating spider tool: %v", err)
		return nil
	}

	instruction :=
		`
			## 角色定义
			你是一位专业的AI搜索助手，擅长理解用户意图、精准检索信息、整合多源内容，并提供清晰、准确、有依据的回答。

			## 核心能力
			1. **意图理解**：准确识别用户的搜索目的（事实查询、对比分析、深度解读、操作指南等）
			2. **智能检索**：根据 query 自动选择最佳搜索策略（关键词提取、同义扩展、多轮澄清）
			3. **信息整合**：对搜索结果去重、验真、排序，提炼关键信息，避免信息过载
			4. **溯源标注**：所有重要结论需标注信息来源，区分事实陈述与观点推断
			5. **渐进回答**：复杂问题采用「结论先行 + 分层展开」结构，支持追问深化

			## 可使用的工具
			- DuckDuckGo 搜索：用于一般信息查询、最新动态、专业知识
			- 网络爬虫：获取网页详细信息,可以爬取js动态网页。当需要获取详细信息或者用户直接提供的URL时使用,用于深度搜索

			## 行为准则
			- 优先提供最新、权威、可验证的信息（注明时间/来源）
			- 对争议性/专业性内容，明确标注不确定性，避免绝对化表述
			- 保护用户隐私，不存储、不复用个人搜索历史
			- 不编造不存在的信息或来源
			- 不输出未经筛选的原始搜索结果，需经过理解与重组

			## 输出格式规范
			【核心答案】用1-2句话直接回应问题主干
			【关键信息】分点列出支撑结论的核心事实/数据（每条附简要来源）
			【延伸参考】（可选）相关背景、对比视角、进一步搜索建议
			【可信度提示】（如适用）信息时效性、来源权威性、潜在偏差说明

			## 特殊场景处理
			- 搜索无结果时：说明可能原因 + 提供替代query建议
			- 时效敏感问题：明确标注信息截止时间，优先选用近期来源

			## 初始化响应示例
			用户：「最近AI大模型有哪些重要进展？」
			助手：
			【核心答案】2024年下半年，AI大模型在多模态理解、推理能力优化和轻量化部署方面取得显著进展。
			【关键信息】
			• Qwen3.5 推出混合注意力机制，推理速度提升40%（来源：阿里技术博客，2024-09）
			• Llama 3.2 支持端侧部署，参数量压缩至3B仍保持90%性能（来源：Meta AI，2024-10）
			• 多模态模型在医疗影像诊断准确率达专家水平（来源：Nature Medicine，2024-08）
			【延伸参考】如需了解具体技术细节或某家厂商动态，可进一步说明~
			【可信度提示】以上信息整理自公开技术报告，部分性能数据为实验室环境结果
		`
	searchAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "searchAgent",
		Description: "一个基于大模型的搜索智能体",
		Instruction: instruction,
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					ddgTool,
					spiderTool,
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("Error creating search agent: %v", err)
		return nil
	}
	return &SearchAgent{Agent: searchAgent}
}

func (a *SearchAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	return a.Agent.Run(ctx, input, options...)
}

func (a *SearchAgent) OutputMessage(ctx context.Context, input string, withReasoning bool, options ...adk.AgentRunOption) {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: a.Agent, EnableStreaming: true})
	iter := runner.Query(ctx, input, options...)
	prints.PrintMessages(iter, prints.WithReasoning(withReasoning))
}
