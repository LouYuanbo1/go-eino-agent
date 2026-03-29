package githubutil

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// CreateLocalGitHubMCPClient 通过 Docker 运行本地 GitHub MCP Server，并通过 stdio 连接
func CreateLocalGitHubMCPClient(ctx context.Context, token string) (*client.Client, error) {
	// 使用 Docker 运行官方镜像，通过 stdio 通信
	cli, err := client.NewStdioMCPClient(
		"docker", // command
		[]string{ // env（第二个参数）
			"GITHUB_PERSONAL_ACCESS_TOKEN=" + token,
		},
		// args...（从第三个参数开始，可变参数）
		"run", "-i", "--rm",
		"-e", "GITHUB_PERSONAL_ACCESS_TOKEN",
		"ghcr.io/github/github-mcp-server",
		"stdio", // 👈 显式指定启动模式（从 --help 可知支持）
	)
	if err != nil {
		return nil, fmt.Errorf("创建 stdio 客户端失败: %w", err)
	}

	// 启动客户端
	if err := cli.Start(ctx); err != nil {
		return nil, fmt.Errorf("启动客户端失败: %w", err)
	}

	// MCP 协议初始化握手
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "eino-github-demo",
		Version: "1.0.0",
	}
	if _, err := cli.Initialize(ctx, initReq); err != nil {
		return nil, fmt.Errorf("初始化握手失败: %w", err)
	}

	log.Println("本地 GitHub MCP 客户端连接成功")
	return cli, nil
}
