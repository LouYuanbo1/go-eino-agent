# go-eino-agent

基于 CloudWeGo Eino 框架实现的智能体仓库，提供开箱即用的智能体组件，用于快速构建和组装AI应用。

## 项目结构

```
go-eino-agent/
├── agents/                 # 智能体实现
│   ├── base/               # 基础接口定义
│   ├── chat/               # 聊天智能体
│   ├── python/             # Python执行智能体
│   └── search/             # 搜索智能体
├── cmd/                    # 命令行工具
│   ├── chat/               # 聊天智能体命令行
│   ├── python/             # Python智能体命令行
│   └── search/             # 搜索智能体命令行
├── prints/                 # 输出工具
└── tools/                  # 工具实现
    ├── pyexecutor/         # Python执行工具
    └── search/             # 搜索相关工具
        ├── duckduckgo/     # DuckDuckGo搜索
        └── spider/         # 网络爬虫
```

## 核心功能

### 智能体类型

1. **ChatAgent**
   - 基于大模型的聊天智能体
   - 集成DuckDuckGo搜索工具
   - 支持流式输出
   - 适用于一般聊天和信息查询场景

2. **PythonAgent**
   - 支持执行Python代码
   - 提供本地执行和沙箱执行两种模式
   - 集成DuckDuckGo搜索工具
   - 适用于代码执行、数据分析等场景

3. **SearchAgent**
   - 集成DuckDuckGo搜索和网络爬虫工具
   - 支持深度信息获取和整合
   - 适用于需要详细信息检索的场景

### 工具模块

1. **DuckDuckGo搜索**
   - 提供网络搜索能力
   - 支持自定义搜索区域和结果数量
   - 适用于获取最新信息和一般知识查询

2. **Python执行**
   - 支持本地执行Python代码
   - 支持沙箱环境执行（安全隔离）
   - 适用于代码测试、数据分析等场景

3. **网络爬虫**
   - 支持爬取动态JavaScript网页
   - 提取网页核心内容
   - 适用于深度信息获取

## 技术亮点

1. **模块化设计**
   - 清晰的组件划分
   - 易于扩展和定制
   - 支持快速组装复杂智能体

2. **工具集成**
   - 丰富的内置工具
   - 标准化的工具接口
   - 支持自定义工具扩展

3. **安全执行**
   - Python代码沙箱执行
   - 超时控制
   - 错误处理机制

4. **流式输出**
   - 实时响应
   - 良好的用户体验

5. **基于CloudWeGo Eino**
   - 利用Eino框架的强大能力
   - 与大模型无缝集成
   - 支持多种模型后端

## 快速开始

### 安装

```bash
go get github.com/LouYuanbo1/go-eino-agent
```

### 示例代码

#### ChatAgent 示例

```go
import (
    "context"
    "github.com/LouYuanbo1/go-eino-agent/agents/chat"
    "github.com/cloudwego/eino-ext/components/model/ollama"
)

func main() {
    ctx := context.Background()
    
    // 初始化模型
    model, err := ollama.NewOllamaChatModel(ctx, &ollama.Config{
        BaseURL: "http://localhost:11434",
        Model:   "llama3",
    })
    if err != nil {
        panic(err)
    }
    
    // 创建默认聊天智能体
    agent := chatAgent.NewDefaultChatAgent(ctx, model)
    
    // 执行查询
    agent.OutputMessage(ctx, "你好，请介绍一下你自己", false)
}
```

#### PythonAgent 示例

```go
import (
    "context"
    "github.com/LouYuanbo1/go-eino-agent/agents/python"
    "github.com/LouYuanbo1/go-eino-agent/tools/pyexecutor"
    "github.com/cloudwego/eino-ext/components/model/ollama"
)

func main() {
    ctx := context.Background()
    
    // 初始化模型
    model, err := ollama.NewOllamaChatModel(ctx, &ollama.Config{
        BaseURL: "http://localhost:11434",
        Model:   "llama3",
    })
    if err != nil {
        panic(err)
    }
    
    // 配置本地Python执行环境
    localConfig := &pyexecutor.LocalOperatorConfig{
        WorkDir:          "./temp",
        ExecutablePath:   "python",
        TaskIDFormat:     pyexecutor.IDFormatTime,
    }
    
    // 创建默认Python智能体
    agent := pythonAgent.NewDefaultPythonAgentLocal(ctx, model, localConfig)
    
    // 执行Python代码
    agent.OutputMessage(ctx, "编写一个函数计算斐波那契数列的第10个数", false)
}
```

#### SearchAgent 示例

```go
import (
    "context"
    "github.com/LouYuanbo1/go-eino-agent/agents/search"
    "github.com/LouYuanbo1/go-eino-agent/tools/search/spider"
    "github.com/cloudwego/eino-ext/components/model/ollama"
)

func main() {
    ctx := context.Background()
    
    // 初始化模型
    model, err := ollama.NewOllamaChatModel(ctx, &ollama.Config{
        BaseURL: "http://localhost:11434",
        Model:   "llama3",
    })
    if err != nil {
        panic(err)
    }
    
    // 配置爬虫
    spiderConfig := &spider.SpiderConfig{
        // 浏览器可执行文件路径,例如:
        Bin: "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
    }
    
    // 创建默认搜索智能体
    agent := searchAgent.NewDefaultSearchAgent(ctx, model, spiderConfig)
    
    // 执行搜索
    agent.OutputMessage(ctx, "https://cloudwego.io/zh/docs/eino/quick_start/agent_llm_with_tools,讲解一下这个页面的信息", false)
}
```

## 依赖

- [github.com/cloudwego/eino](https://github.com/cloudwego/eino) - 核心框架
- [github.com/cloudwego/eino-ext](https://github.com/cloudwego/eino-ext) - 扩展组件
- [github.com/go-rod/rod](https://github.com/go-rod/rod) - 网络爬虫
- [github.com/go-shiori/go-readability](https://github.com/go-shiori/go-readability) - 网页内容提取

## 注意事项

- 此仓库与CloudWeGo Eino框架完全耦合
- 使用Python执行工具时，请注意安全风险
- 网络爬虫功能需要安装浏览器（如Chrome/Chromium）
- 部分功能可能需要配置API密钥或环境变量

## 许可证

[MIT License](LICENSE)