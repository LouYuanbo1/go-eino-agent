// github/mcp_parser.go
package githubutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mcptools "github.com/LouYuanbo1/go-eino-agent/tools/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type GithubGetFuncWrapper struct {
	base tool.InvokableTool
}

func NewGithubGetFuncWrapper(base tool.InvokableTool) *GithubGetFuncWrapper {
	return &GithubGetFuncWrapper{base: base}
}

func (w *GithubGetFuncWrapper) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	rawOutput, err := w.base.InvokableRun(ctx, argumentsInJSON, opts...)
	if err != nil {
		return "", err // 错误也返回为正常结构
	}
	return ParseMCPGitHubFileContents(rawOutput)
}

func (w *GithubGetFuncWrapper) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return w.base.Info(ctx)
}

type GitHubFileResource struct {
	URI      string `json:"uri"`      // e.g. "repo://cloudwego/eino/sha/xxx/contents/README.md"
	MimeType string `json:"mimeType"` // e.g. "text/plain; charset=utf-8"
	Text     string `json:"text"`     // ✅ 文件实际内容（Markdown/文本）
	// 可选扩展字段（按需添加）
	Encoding string `json:"encoding,omitempty"` // "utf-8" | "base64"
	Size     int    `json:"size,omitempty"`     // 文件大小（字节）
}

// GitHubFileResult 解析后的最终结果
type GitHubFileResult struct {
	StatusMsg string              // 状态消息（如 "successfully downloaded..."）
	File      *GitHubFileResource // 文件内容（如果有）
}

/*
// ParseMCPResponse 通用解析：支持文件/目录/状态消息
func ParseMCPResponseGitHubGetFileContents(resultJSON string) (string, error) {

	fmt.Println(resultJSON)

	var resp map[string]any
	if err := json.Unmarshal([]byte(resultJSON), &resp); err != nil {
		return "", fmt.Errorf("解析 JSON 失败: %w", err)
	}

	var statusMsg string

	contentArray, ok := resp["content"].([]any)
	if !ok {
		return "", fmt.Errorf("content 字段不是数组")
	}

	for _, item := range contentArray {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}

		typ, _ := m["type"].(string)

		switch typ {
		case "text":
			// 处理状态消息 或 目录列表
			if text, ok := m["text"].(string); ok {
				trimmed := strings.TrimSpace(text)
				if strings.HasPrefix(trimmed, "[") {
					// 目录列表（转义数组）
					return FormatDirectoryList(trimmed)
				}
				// 普通状态消息，暂存
				if statusMsg == "" {
					statusMsg = text
				}
			}

		case "resource":
			// 文件内容（最高优先级）
			if resource, ok := m["resource"].(map[string]any); ok {
				if content, ok := resource["text"].(string); ok && content != "" {
					return content, nil
				}
			}
		}
	}

	// 兜底：返回状态消息
	if statusMsg != "" {
		return statusMsg, nil
	}
	return "", fmt.Errorf("未找到可解析的内容")
}

// FormatDirectoryList 格式化目录列表（保持原样）
func FormatDirectoryList(jsonArray string) (string, error) {
	var items []map[string]any
	if err := json.Unmarshal([]byte(jsonArray), &items); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("📦 目录内容:\n")
	for _, item := range items {
		name, _ := item["name"].(string)
		typ, _ := item["type"].(string)
		size, _ := item["size"].(float64)

		icon := "📁"
		if typ == "file" {
			icon = "📄"
		}

		sizeStr := ""
		if size > 0 {
			if size < 1024 {
				sizeStr = fmt.Sprintf(" (%dB)", int64(size))
			} else {
				sizeStr = fmt.Sprintf(" (%.1fKB)", size/1024)
			}
		}

		fmt.Fprintf(&sb, "  %s %-30s%s\n", icon, name, sizeStr)
	}
	return sb.String(), nil
}
*/

// ParseMCPGitHubFileContents 解析 MCP 包装的 GitHub 文件内容响应
func ParseMCPGitHubFileContents(resultJSON string) (string, error) {
	// Step 1: 绑定 MCP 外层（类型安全 ✅）
	var mcpResp mcptools.MCPResponse
	if err := json.Unmarshal([]byte(resultJSON), &mcpResp); err != nil {
		return "", fmt.Errorf("解析 MCP 外层失败: %w", err)
	}

	result := &GitHubFileResult{}

	// Step 2: 遍历 content 数组，提取不同类型内容
	for _, item := range mcpResp.Content {
		switch item.Type {
		case "text":
			// 状态消息（如下载成功提示）
			if strings.TrimSpace(item.Text) != "" {
				result.StatusMsg = item.Text
			}

		case "resource":
			// 文件内容：将 map[string]any 绑定到强类型结构
			if item.Resource != nil {
				var fileRes GitHubFileResource
				// 🔑 关键：把 resource 的 map 重新 marshal + unmarshal 到目标结构
				// 这是 Go 处理 "map → struct" 的标准做法
				resourceBytes, err := json.Marshal(item.Resource)
				if err != nil {
					continue // 跳过解析失败的项
				}
				if err := json.Unmarshal(resourceBytes, &fileRes); err != nil {
					continue
				}
				// ✅ 优先返回有内容的文件
				if strings.TrimSpace(fileRes.Text) != "" {
					result.File = &fileRes
				}
			}
		}
	}

	// Step 3: 校验结果
	if result.File == nil && result.StatusMsg == "" {
		return "", fmt.Errorf("未找到有效内容：既无文件资源，也无状态消息")
	}

	return FormatGitHubFileResult(result), nil
}

// FormatGitHubFileResult 将解析结果格式化为易读文本
func FormatGitHubFileResult(result *GitHubFileResult) string {
	var sb strings.Builder

	// 📌 状态消息（如果有）
	if result.StatusMsg != "" {
		//sb.WriteString(fmt.Sprintf("✅ %s\n\n", result.StatusMsg))
	}

	// 📄 文件内容
	if result.File != nil {
		// 提取文件名（从 URI 解析）
		filename := extractFilenameFromURI(result.File.URI)

		sb.WriteString(fmt.Sprintf("📄 %s\n", filename))

		// MIME 类型 + 大小（如果有）
		var meta []string
		if result.File.MimeType != "" {
			meta = append(meta, result.File.MimeType)
		}
		if result.File.Size > 0 {
			meta = append(meta, formatFileSize(result.File.Size))
		}
		if len(meta) > 0 {
			sb.WriteString(fmt.Sprintf("   %s\n", strings.Join(meta, "  •  ")))
		}
		sb.WriteString("\n")

		/*
			// 内容预览（前 30 行，避免刷屏）
			preview := limitLines(result.File.Text, 30)
			sb.WriteString(preview)

			// 如果内容有截断，添加提示
			if countLines(result.File.Text) > 30 {
				sb.WriteString("\n   ...（内容已截断，共 ")
				sb.WriteString(fmt.Sprintf("%d", countLines(result.File.Text)))
				sb.WriteString(" 行）\n")
			}
		*/
		for line := range strings.SplitSeq(result.File.Text, "\n") {
			sb.WriteString(line + "\n")
		}
	}

	return sb.String()
}

// ============ 辅助函数 ============

// extractFilenameFromURI 从 GitHub URI 提取文件名
// e.g. "repo://cloudwego/eino/sha/xxx/contents/README.md" → "README.md"
func extractFilenameFromURI(uri string) string {
	parts := strings.Split(uri, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

// limitLines 限制文本显示行数
/*
func limitLines(text string, maxLines int) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= maxLines {
		return text
	}
	return strings.Join(lines[:maxLines], "\n") + "\n"
}

// countLines 统计文本行数
func countLines(text string) int {
	return len(strings.Split(strings.TrimSpace(text), "\n"))
}
*/

// formatFileSize 人性化文件大小
func formatFileSize(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}
