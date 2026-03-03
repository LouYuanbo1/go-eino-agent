# go-eino-agent

基于 CloudWeGo Eino 框架实现的智能体仓库，提供开箱即用的智能体组件，用于快速构建和组装AI应用。

## 项目结构

```
go-eino-agent/
├── agents/                 # 智能体实现
│   ├── chat/               # 聊天智能体
│   ├── model.go            # 智能体模型定义
│   ├── python/             # Python执行智能体
│   ├── retriever/          # 检索智能体
│   └── search/             # 搜索智能体
├── cmd/                    # 命令行工具
│   ├── chat/               # 聊天智能体命令行
│   ├── python/             # Python智能体命令行
│   ├── retriever/          # 检索智能体命令行
│   └── search/             # 搜索智能体命令行
├── config/                 # 配置管理
│   ├── config.go           # 配置处理
│   ├── config.yaml         # 配置文件
│   └── config_example.yaml # 配置示例
├── prints/                 # 输出工具
│   ├── options.go          # 输出选项
│   └── print.go            # 输出实现
└── tools/                  # 工具实现
    ├── pyexecutor/         # Python执行工具
    │   ├── local/          # 本地执行模式
    │   ├── params/         # 参数定义
    │   └── sandbox/        # 沙箱执行模式
    ├── retriever/          # 检索工具
    │   ├── elasticsearch/  # Elasticsearch实现
    │   ├── model.go        # 检索模型定义
    │   ├── redisstack/     # Redis Stack实现
    │   └── tool.go         # 检索工具接口
    └── search/             # 搜索工具
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

4. **RetrieverAgent**
   - 支持从Elasticsearch和Redis Stack检索信息
   - 适用于知识库查询和信息检索场景

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

4. **检索工具**
   - 支持Elasticsearch和Redis Stack
   - 适用于知识库查询和信息检索

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

## Python执行工具详解

### 核心功能

Python执行工具是go-eino-agent的核心工具之一，提供了两种执行模式：

1. **本地执行模式**
   - 直接在本地环境执行Python代码
   - 生成临时Python文件并执行
   - 执行成功时保存源码为.py文件
   - 支持自定义工作目录和Python解释器路径

2. **沙箱执行模式**
   - 使用Docker容器隔离执行环境
   - 提供更安全的执行环境
   - 防止恶意代码对系统造成损害

### 实现原理

#### 本地执行模式

```go
// 核心执行流程
1. 生成唯一ID（基于UUID或时间戳）
2. 创建临时Python文件
3. 写入用户提供的代码
4. 执行Python代码
5. 收集标准输出和错误
6. 执行成功时保存源码
7. 返回执行结果
```

#### 沙箱执行模式

```go
// 核心执行流程
1. 创建Docker沙箱环境
2. 在沙箱中执行Python代码
3. 收集执行结果
4. 清理沙箱环境
5. 返回执行结果
```

### 使用示例

#### 本地执行模式

```go
import (
    "context"
    "github.com/LouYuanbo1/go-eino-agent/agents/python"
    "github.com/LouYuanbo1/go-eino-agent/tools/pyexecutor/local"
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
    localConfig := &local.OperatorConfig{
        WorkDir:          "./temp",         // 工作目录
        ExecutablePath:   "python",         // Python解释器路径
        FileName:         "python_script",  // 文件名前缀
        TaskIDFormat:     local.IDFormatTime, // 任务ID格式
    }
    
    // 创建默认Python智能体
    agent := pythonAgent.NewDefaultPythonAgentLocal(ctx, model, localConfig)
    
    // 执行Python代码
    agent.OutputMessage(ctx, "编写一个函数计算斐波那契数列的第10个数", false)
}
```

#### 沙箱执行模式

```go
import (
    "context"
    "github.com/LouYuanbo1/go-eino-agent/agents/python"
    "github.com/cloudwego/eino-ext/components/model/ollama"
    "github.com/cloudwego/eino-ext/components/tool/commandline/sandbox"
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
    
    // 配置沙箱环境
    sandboxConfig := &sandbox.Config{
        Image: "python:3.10-slim", // Docker镜像
        Timeout: 30,               // 执行超时时间（秒）
    }
    
    // 创建默认Python智能体（沙箱模式）
    agent := pythonAgent.NewDefaultPythonAgentInSandbox(ctx, model, sandboxConfig)
    
    // 执行Python代码
    agent.OutputMessage(ctx, "编写一个函数计算斐波那契数列的第10个数", false)
}
```

### 安全特性

1. **本地执行模式**
   - 临时文件自动清理
   - 错误处理机制
   - 支持自定义工作目录，便于隔离文件操作

2. **沙箱执行模式**
   - Docker容器隔离
   - 资源限制
   - 网络访问控制

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
    "github.com/LouYuanbo1/go-eino-agent/tools/pyexecutor/local"
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
    localConfig := &local.OperatorConfig{
        WorkDir:          "./temp",
        ExecutablePath:   "python",
        FileName:         "python_script",
        TaskIDFormat:     local.IDFormatTime,
    }
    
    // 创建默认Python智能体
    agent := pythonAgent.NewDefaultPythonAgentLocal(ctx, model, localConfig)
    
    // 执行Python代码
    agent.OutputMessage(ctx,"使用Python生成一个简易的线性回归模型", false)
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

#### RetrieverAgent 示例

```go
import (
    "context"
    "github.com/LouYuanbo1/go-eino-agent/agents/retriever"
    "github.com/LouYuanbo1/go-eino-agent/tools/retriever/elasticsearch"
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
    
    // 配置Elasticsearch
    esConfig := &elasticsearch.Config{
        Address: "http://localhost:9200",
        Index:   "documents",
    }
    
    // 创建默认检索智能体
    agent := retrieverAgent.NewDefaultRetrieverAgent(ctx, model, esConfig)
    
    // 执行检索
    agent.OutputMessage(ctx, "检索关于CloudWeGo Eino框架的信息", false)
}
```

## 配置管理

项目使用YAML配置文件管理工具参数，默认配置文件为`config/config.yaml`。

### 配置示例

```yaml
# 搜索工具配置
search:
  duckduckgo:
    region: "cn-en"
    max_results: 5

# Python执行工具配置
pyexecutor:
  local:
    work_dir: "./temp"
    executable_path: "python"
    file_name: "python_script"
    task_id_format: "time"

# 爬虫配置
spider:
  bin: "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
  headless: true

# 检索工具配置
retriever:
  elasticsearch:
    address: "http://localhost:9200"
    index: "documents"
  redisstack:
    address: "localhost:6379"
    password: ""
    db: 0
```

## 依赖

- [github.com/cloudwego/eino](https://github.com/cloudwego/eino) - 核心框架
- [github.com/cloudwego/eino-ext](https://github.com/cloudwego/eino-ext) - 扩展组件
- [github.com/go-rod/rod](https://github.com/go-rod/rod) - 网络爬虫
- [github.com/go-shiori/go-readability](https://github.com/go-shiori/go-readability) - 网页内容提取
- [github.com/google/uuid](https://github.com/google/uuid) - UUID生成
- [github.com/redis/go-redis/v9](https://github.com/redis/go-redis/v9) - Redis客户端
- [github.com/elastic/go-elasticsearch/v8](https://github.com/elastic/go-elasticsearch/v8) - Elasticsearch客户端

## 注意事项

- 此仓库与CloudWeGo Eino框架完全耦合
- 使用Python执行工具时，请注意安全风险
- 网络爬虫功能需要安装浏览器（如Chrome/Chromium）
- 沙箱执行模式需要安装Docker
- 部分功能可能需要配置API密钥或环境变量

## 许可证

[MIT License](LICENSE)
