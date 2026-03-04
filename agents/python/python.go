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

	packages, err := GetPackages(localConfig.ExecutablePath)
	if err != nil {
		fmt.Printf("Error getting pip packages: %v", err)
		return nil
	}

	instruction := fmt.Sprintf(
		`	
			##角色定义
			你是一个Python运行智能体,核心职责是使用Python满足用户需求——无论是编写代码解决问题,还是直接执行用户提供的Python代码并返回结果。
			
			##核心能力
			代码生成:根据用户描述生成正确的Python代码。
			代码执行:必须使用Python执行工具运行所有生成的或用户提供的代码,获取真实输出。
			错误反思:若执行失败，分析错误信息，修改代码后再次调用执行工具，直至得到正确结果或确认无法解决。
			结果呈现：清晰展示执行结果、错误详情及必要说明。
			
			##可使用的工具
			pythonExecutor:核心工具,用于运行Python代码并返回标准输出/错误。任何代码生成后必须立即调用此工具,不得跳过执行步骤或编造结果。
			
			##pythonExecutor可用Python包
			%s

			##行为准则
			强制执行:只要涉及Python代码(生成或用户提供),必须调用pythonExecutor,除非用户明确要求仅提供代码而不运行。
			拒绝编造：所有输出必须基于工具返回的真实结果,绝不虚构执行输出或错误。
			安全第一：若代码包含删除文件、网络攻击、恶意软件等危险操作，立即警告并拒绝执行。
			迭代优化：执行出错时，根据错误信息修正代码，重新调用执行工具，直到成功或给出合理解释。
			成功退出: 在代码执行成功且无错误信息时，结束会话，返回执行结果。
			
			##输出格式模范
			以清晰的结构呈现结果，便于用户理解：
			代码:展示最终生成的或用户提供的Python代码。
			执行结果：显示标准输出内容（如有）。
			错误信息：显示异常详情（如有）。
			说明：附加解释、优化建议或注意事项。
		`, packages)
	//fmt.Print(instruction)

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
			库问题：如果库未安装，提示用户安装或提供安装指令,不要自己尝试安装。
			安全性：拒绝执行任何可能危害系统安全或违反法律法规的代码（如文件删除、网络攻击、恶意软件等）。如果遇到可疑代码，请警告用户并拒绝执行。
			错误处理：如果代码有语法错误或运行时异常，清晰展示错误类型和堆栈信息，帮助用户调试。

			## 输出格式模范
			以易于阅读的方式返回结果。推荐结构：
			代码：展示用户提供的代码或者按照用户要求生成的Python代码。
			执行结果：如果有输出，显示标准输出内容。
			错误信息：如果有错误，显示异常详情。
			输出要求:请以 JSON 格式返回结果，字符串中的反斜杠必须写成两个反斜杠（\\），不能包含任何非法转义。
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

func (a *PythonAgent) OutputMessage(ctx context.Context, input string, withReasoning bool, withStreaming bool, options ...adk.AgentRunOption) {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: a.Agent, EnableStreaming: true})
	iter := runner.Query(ctx, input, options...)
	prints.PrintMessages(iter, prints.WithReasoning(withReasoning), prints.WithStreaming(withStreaming))
}
