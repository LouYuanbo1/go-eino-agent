package githubutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mcptools "github.com/LouYuanbo1/go-eino-agent/tools/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type GithubSearchFuncWrapper struct {
	base tool.InvokableTool
}

func NewGithubSearchFuncWrapper(base tool.InvokableTool) *GithubSearchFuncWrapper {
	return &GithubSearchFuncWrapper{base: base}
}

func (w *GithubSearchFuncWrapper) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	rawOutput, err := w.base.InvokableRun(ctx, argumentsInJSON, opts...)
	if err != nil {
		return "", err // 错误也返回为正常结构
	}
	return ParseMCPGitHubSearch(rawOutput)
}

func (w *GithubSearchFuncWrapper) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return w.base.Info(ctx)
}

// ============ 1. 类型定义 ============

// GitHubSearchResponse GitHub Search API 的真实响应
type GitHubSearchResponse struct {
	TotalCount        int          `json:"total_count"`
	IncompleteResults bool         `json:"incomplete_results"`
	Items             []Repository `json:"items"`
}

// Repository 精简版仓库结构（只取常用字段，避免解析失败）
type Repository struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	HTMLURL         string    `json:"html_url"`
	Language        string    `json:"language"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
	Topics          []string  `json:"topics"`
	Private         bool      `json:"private"`
	Fork            bool      `json:"fork"`
	Archived        bool      `json:"archived"`
	// 可选：如果需要 owner 信息
	Owner struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		HTMLURL   string `json:"html_url"`
	} `json:"owner"`
}

// ============ 2. 主解析函数 ============

// ParseMCPGitHubSearch 解析 MCP 包装的 GitHub Search 响应
func ParseMCPGitHubSearch(resultJSON string) (string, error) {
	// Step 1: 解析 MCP 外层
	var mcpResp mcptools.MCPResponse
	if err := json.Unmarshal([]byte(resultJSON), &mcpResp); err != nil {
		return "", fmt.Errorf("解析 MCP 外层失败: %w", err)
	}

	// Step 2: 提取 content 中的 GitHub JSON 字符串
	var rawGitHubJSON string
	for _, item := range mcpResp.Content {
		if item.Type == "text" && strings.Contains(item.Text, `"items"`) {
			rawGitHubJSON = strings.TrimSpace(item.Text)
			break
		}
		// 兼容 resource 类型（如果有）
		if item.Type == "resource" && item.Resource != nil {
			if txt, ok := item.Resource["text"].(string); ok && strings.Contains(txt, `"items"`) {
				rawGitHubJSON = strings.TrimSpace(txt)
				break
			}
		}
	}

	if rawGitHubJSON == "" {
		return "", fmt.Errorf("未在 content 中找到有效的 GitHub Search JSON")
	}

	// Step 3: 二次解析：解析内部的 GitHub Search 响应
	var searchResp GitHubSearchResponse
	if err := json.Unmarshal([]byte(rawGitHubJSON), &searchResp); err != nil {
		return "", fmt.Errorf("解析 GitHub Search 内容失败: %w, 原始片段: %.100s...", err, rawGitHubJSON)
	}

	// Step 4: 格式化为人类可读输出
	return FormatGitHubSearchResults(&searchResp), nil
}

// ============ 3. 人性化格式化输出 ============

func FormatGitHubSearchResults(resp *GitHubSearchResponse) string {
	var sb strings.Builder

	// 📊 头部
	emoji := "🔍"
	if resp.TotalCount == 0 {
		emoji = "😕"
	}
	sb.WriteString(fmt.Sprintf("%s 找到 %d 个仓库", emoji, resp.TotalCount))
	if resp.IncompleteResults {
		sb.WriteString(" ⚠️ 结果不完整")
	}
	sb.WriteString("\n\n")

	if len(resp.Items) == 0 {
		sb.WriteString("   暂无匹配结果～\n")
		return sb.String()
	}

	// 📦 逐个格式化
	for i, repo := range resp.Items {
		sb.WriteString(formatRepoCard(repo, i+1))
		sb.WriteString("\n")
	}

	return sb.String()
}

func formatRepoCard(repo Repository, index int) string {
	var sb strings.Builder

	// 序号 + 仓库链接
	owner := repo.Owner.Login
	if owner == "" {
		// 从 FullName 解析
		parts := strings.Split(repo.FullName, "/")
		if len(parts) == 2 {
			owner = parts[0]
		}
	}
	sb.WriteString(fmt.Sprintf("%d. [%s](%s)\n", index, repo.FullName, repo.HTMLURL))

	// 描述（截断 + 换行处理）
	if repo.Description != "" {
		desc := strings.ReplaceAll(repo.Description, "\n", " ")
		if len(desc) > 100 {
			desc = desc[:97] + "..."
		}
		sb.WriteString(fmt.Sprintf("   📝 %s\n", desc))
	}

	// 核心指标行
	starStr := formatCompactNumber(repo.StargazersCount)
	sb.WriteString(fmt.Sprintf("   ⭐ %s  •  🍴 %d  •  ❗ %d\n",
		starStr, repo.ForksCount, repo.OpenIssuesCount))

	// 标签行：语言 + 状态
	var meta []string
	if repo.Language != "" {
		meta = append(meta, "🔧 "+repo.Language)
	}
	if repo.Archived {
		meta = append(meta, "🗄️ 归档")
	}
	if repo.Private {
		meta = append(meta, "🔒 私有")
	}
	if repo.Fork {
		meta = append(meta, "🔀 Fork")
	}
	if len(meta) > 0 {
		sb.WriteString(fmt.Sprintf("   %s\n", strings.Join(meta, "  ")))
	}

	// Topics（最多3个）
	if len(repo.Topics) > 0 {
		topics := repo.Topics
		if len(topics) > 3 {
			topics = topics[:3]
		}
		sb.WriteString(fmt.Sprintf("   🏷️  %s\n", strings.Join(topics, "  ")))
	}

	// 更新时间 + 作者
	updated := repo.UpdatedAt.Format("01-02")
	sb.WriteString(fmt.Sprintf("   👤 @%s  •  🕒 %s\n", owner, updated))

	return sb.String()
}

// formatCompactNumber 数字人性化：1000→1k, 1000000→1M
func formatCompactNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1e3)
	}
	return fmt.Sprintf("%d", n)
}
