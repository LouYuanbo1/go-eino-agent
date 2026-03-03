package pythonAgent

import (
	"context"
	"fmt"

	"github.com/LouYuanbo1/go-eino-agent/prints"

	"github.com/LouYuanbo1/go-eino-agent/tools/pyexecutor/local"
	pyexecutor "github.com/LouYuanbo1/go-eino-agent/tools/pyexecutor/sandbox"
	"github.com/LouYuanbo1/go-eino-agent/tools/search/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/commandline/sandbox"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

type PythonAgent struct {
	Agent *adk.ChatModelAgent
}

func NewPythonAgent(ctx context.Context, config *adk.ChatModelAgentConfig) *PythonAgent {
	agent, err := adk.NewChatModelAgent(ctx, config)
	if err != nil {
		fmt.Printf("Error creating python agent: %v", err)
		return nil
	}
	return &PythonAgent{Agent: agent}
}

func NewDefaultPythonAgentLocal(ctx context.Context, model model.ToolCallingChatModel, localConfig *local.OperatorConfig) *PythonAgent {
	ddgTool, err := duckduckgo.NewDefaultDuckDuckGoTool(ctx)
	if err != nil {
		fmt.Printf("Error creating DuckDuckGo tool: %v", err)
		return nil
	}
	pythonTool, err := local.NewPythonToolLocal(ctx, localConfig)
	if err != nil {
		fmt.Printf("Error creating PyExecutor tool: %v", err)
		return nil
	}

	instruction :=
		`	
			## 角色定义
			你是一个Python运行智能体，专门负责使用Python满足用户的要求或者执行用户提供的Python代码并返回结果。

			## 核心能力
			1. 接收用户的Python代码或要求。
			2. 使用Python执行或满足用户的要求。
			3. 反思执行结果,如果出现错误,则修改代码。
			4. 输出执行结果、错误信息或必要的解释。

			## 可使用的工具
			- DuckDuckGo 搜索：用于搜寻Python相关的代码、文档、示例等。
			- Python 执行工具：用于执行用户提供的Python代码并返回结果。

			## 行为准则
			拒绝编造：不要编造Python代码的执行结果,生成代码后调用Python执行工具执行,并根据执行结果选择是修改代码还是返回结果。
			安全性：拒绝执行任何可能危害系统安全或违反法律法规的代码（如文件删除、网络攻击、恶意软件等）。如果遇到可疑代码，请警告用户并拒绝执行。

			## 输出格式模范
			以易于阅读的方式返回结果。推荐结构：
			代码：展示用户提供的代码或者按照用户要求生成的Python代码。
			执行结果：如果有输出，显示标准输出内容。
			错误信息：如果有错误，显示异常详情。
			说明：必要时附加解释或建议。
		`
	pythonAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "pythonAgent",
		Description: "一个基于大模型的python运行智能体",
		Instruction: instruction,
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					ddgTool,
					pythonTool,
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("Error creating python agent: %v", err)
		return nil
	}
	return &PythonAgent{Agent: pythonAgent}
}

func NewDefaultPythonAgentInSandbox(ctx context.Context, model model.ToolCallingChatModel, sandboxConfig *sandbox.Config) *PythonAgent {
	ddgTool, err := duckduckgo.NewDefaultDuckDuckGoTool(ctx)
	if err != nil {
		fmt.Printf("Error creating DuckDuckGo tool: %v", err)
		return nil
	}
	pythonTool, err := pyexecutor.NewPythonToolInSandbox(ctx, sandboxConfig)
	if err != nil {
		fmt.Printf("Error creating PyExecutor tool: %v", err)
		return nil
	}

	instruction :=
		`	
			## 角色定义
			你是一个Python运行智能体，专门负责使用Python满足用户的要求或者执行用户提供的Python代码并返回结果。
			Python代码一定要使用Python工具运行

			## 核心能力
			1. 接收用户的Python代码或要求。
			2. 使用Python工具执行或满足用户的要求。
			3. 输出执行结果、错误信息或必要的解释。

			## 可使用的工具
			- DuckDuckGo 搜索：用于搜寻Python相关的代码、文档、示例等。
			- Python 执行工具：使用该工具执行满足用户要求或者用户提供的Python代码。

			## 行为准则
			安全性：拒绝执行任何可能危害系统安全或违反法律法规的代码（如文件删除、网络攻击、恶意软件等）。如果遇到可疑代码，请警告用户并拒绝执行。
			错误处理：如果代码有语法错误或运行时异常，清晰展示错误类型和堆栈信息，帮助用户调试。

			## 输出格式模范
			以易于阅读的方式返回结果。推荐结构：
			代码：展示用户提供的代码或者按照用户要求生成的Python代码。
			执行结果：如果有输出，显示标准输出内容。
			错误信息：如果有错误，显示异常详情。
			说明：必要时附加解释或建议。
		`
	pythonAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "pythonAgent",
		Description: "一个基于大模型的python运行智能体",
		Instruction: instruction,
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					ddgTool,
					pythonTool,
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("Error creating python agent: %v", err)
		return nil
	}
	return &PythonAgent{Agent: pythonAgent}
}

func (a *PythonAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	return a.Agent.Run(ctx, input, options...)
}

func (a *PythonAgent) OutputMessage(ctx context.Context, input string, withReasoning bool, options ...adk.AgentRunOption) {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: a.Agent, EnableStreaming: true})
	iter := runner.Query(ctx, input, options...)
	prints.PrintMessages(iter, prints.WithReasoning(withReasoning))
}
