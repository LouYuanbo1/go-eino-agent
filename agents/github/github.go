package searchAgent

import (
	"context"
	"fmt"
	"strings"

	"github.com/LouYuanbo1/go-eino-agent/prints"
	"github.com/LouYuanbo1/go-eino-agent/tools/mcp/github/tools"
	"github.com/LouYuanbo1/go-eino-agent/tools/mcp/github/tools/get"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/mark3labs/mcp-go/client"
)

type GitHubAgent struct {
	Agent *adk.ChatModelAgent
}

func NewGitHubAgent(ctx context.Context, config *adk.ChatModelAgentConfig) *GitHubAgent {
	agent, err := adk.NewChatModelAgent(ctx, config)
	if err != nil {
		fmt.Printf("Error creating search agent: %v", err)
		return nil
	}
	return &GitHubAgent{Agent: agent}
}

func NewDefaultGitHubAgent(ctx context.Context, model model.ToolCallingChatModel, client *client.Client, toolNames ...string) *GitHubAgent {
	toolMap, err := tools.NewGitHubTools(ctx, client, toolNames...)
	if err != nil {
		fmt.Printf("Error creating GitHub tools: %v", err)
		return nil
	}
	tools := make([]tool.BaseTool, 0)
	for toolName, tool := range toolMap {
		if strings.HasPrefix(toolName, "get") {
			wrapper := get.NewGithubGetFuncWrapper(tool)
			tools = append(tools, wrapper)
		} else if strings.HasPrefix(toolName, "search") {
			wrapper := get.NewGithubSearchFuncWrapper(tool)
			tools = append(tools, wrapper)
		} else {
			tools = append(tools, tool)
		}
	}
	if _, ok := toolMap["get_file_contents"]; ok {
		overviewTool, err := get.NewOverviewTool(ctx, toolMap["get_file_contents"])
		if err != nil {
			fmt.Printf("Error creating overview tool: %v", err)
			return nil
		}
		tools = append(tools, overviewTool)
	}
	
	instruction :=
		`
			##角色定义
			你是一位专业的 GitHub 智能体（GitHub Agent），擅长理解开发意图、通过 MCP 协议操作 GitHub 资源、分析代码上下文，并提供安全、高效、可追溯的解决方案。
			##核心能力
			意图识别：准确识别用户任务类型（代码查询、Bug 修复、功能开发、仓库管理、CI/CD 配置等）
			工具调度：根据任务需求自动选择最佳 GitHub MCP 工具（读取文件、搜索代码、创建 Issue、提交 PR、触发 Action 等）
			上下文分析：深入理解仓库结构、依赖关系及代码逻辑，避免破坏性变更
			操作溯源：所有代码变更或操作需关联具体文件路径、Commit ID 或 Issue 链接，确保可追溯
			渐进执行：复杂任务采用「计划确认 -> 分步执行 -> 结果验证」流程，支持中途干预
			##可使用的工具
			GitHub MCP 工具集：
			读取类：read_file, search_code, list_issues, get_repo_info（用于获取上下文）
			写入类：create_file, update_file, create_issue, create_pull_request（需用户确认或符合安全策略）
			自动化类： trigger_workflow, check_ci_status（用于验证变更）
			注：具体可用工具以实际加载的 MCP Server 能力为准，操作前需确认权限
			##行为准则
			安全第一：严禁硬编码敏感信息（Token/Key），不执行未经确认的强制推送（force push）或删除操作
			最小权限：仅申请完成任务所需的最小仓库权限，操作前说明影响范围
			变更透明：代码修改需提供 Diff 对比或变更说明，重要变更建议走 PR 流程而非直接 commit
			错误处理：遇到 API 限制、权限不足或冲突时，明确告知原因并提供解决建议
			不编造代码：不生成无法验证存在的 API 或库函数，不确定处需标注“需进一步验证”
			##输出格式规范
			【任务目标】用 1 句话明确当前任务的核心目的
			【执行方案】分步列出计划调用的工具及操作逻辑
			【操作详情】
			• 文件变更：路径 + 变更摘要（或代码块）
			• 资源创建：Issue/PR 链接 或 编号
			• 状态反馈：成功/失败/待确认
			【后续建议】（可选）测试建议、关联任务、优化方向
			【风险提示】（如适用）潜在冲突、依赖影响、权限说明
			##特殊场景处理
			权限不足时：说明缺失的具体权限 + 引导用户检查 Token Scope 或仓库设置
			代码冲突时：指出冲突文件 + 提供合并建议（手动解决或重试策略）
			工具不可用时：说明 MCP 工具加载状态 + 提供替代手动操作指南
			时效敏感任务：涉及 CI/CD 状态时，明确标注检查时间点，建议刷新验证
			##初始化响应示例
			用户：「详细讲解一下字节的eino框架的功能」
			助手：
			【任务目标】使用 GitHub MCP 工具获取字节eino框架的详细信息
			【执行方案】
			调用 search_repositories 获取字节eino框架的仓库信息
			调用 get_file_contents 获取字节eino框架的README.md和代码内容
			总结内容,将主要信息总结并回复给用户
			
		`
	githubAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "githubAgent",
		Description: "一个基于大模型的GitHub智能体",
		Instruction: instruction,
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
	})
	if err != nil {
		fmt.Printf("Error creating search agent: %v", err)
		return nil
	}
	return &GitHubAgent{Agent: githubAgent}
}

func (a *GitHubAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	return a.Agent.Run(ctx, input, options...)
}

func (a *GitHubAgent) OutputMessage(ctx context.Context, input string, withReasoning bool, withStreaming bool, options ...adk.AgentRunOption) {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: a.Agent, EnableStreaming: true})
	iter := runner.Query(ctx, input, options...)
	prints.PrintMessages(iter, prints.WithReasoning(withReasoning), prints.WithStreaming(withStreaming))
}
