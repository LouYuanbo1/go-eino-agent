package githubutil

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	mcptools "github.com/LouYuanbo1/go-eino-agent/tools/mcp"
)

// ============ 1. 结构体定义 ============

// GitHubDirEntry GitHub Contents API 目录项（支持目录/文件）
type GitHubDirEntry struct {
	Type        string `json:"type"` // "dir" | "file"
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	SHA         string `json:"sha"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url,omitempty"` // 仅 file 类型有
}

// GitHubFileResource resource 字段的强类型定义（针对文件内容）
type GitHubFileResource struct {
	URI      string `json:"uri"`                // e.g. "repo://cloudwego/eino/sha/xxx/contents/README.md"
	MimeType string `json:"mimeType"`           // e.g. "text/plain; charset=utf-8"
	Text     string `json:"text"`               // ✅ 文件实际内容
	Encoding string `json:"encoding,omitempty"` // "utf-8" | "base64"
	Size     int    `json:"size,omitempty"`     // 文件大小（字节）
}

// GitHubFileResult 文件内容解析结果
type GitHubFileResult struct {
	StatusMsg string
	File      *GitHubFileResource
}

// ============ 2. 通用解析入口（自动路由） ============

// ParseMCPGitHubContents 通用解析入口：自动识别目录/文件/搜索/状态响应
// ParseMCPGitHubContents 通用解析入口（修复版）
func ParseMCPGitHubContents(resultJSON string) (string, error) {
	trimmed := strings.TrimSpace(resultJSON)
	if trimmed == "" {
		return "", fmt.Errorf("空响应")
	}

	// 🔍 类型探测 + 路由分发
	// Case 1: 原始目录列表数组 [{"type":"dir"...}]
	if strings.HasPrefix(trimmed, "[") {
		return ParseGitHubDirectoryList(trimmed)
	}

	// Case 2: MCP 包装响应 {"content":[...]}
	var mcpResp mcptools.MCPResponse
	if err := json.Unmarshal([]byte(trimmed), &mcpResp); err != nil {
		// Case 3: 原始 GitHub Search 响应
		if strings.Contains(trimmed, `"total_count"`) && strings.Contains(trimmed, `"items"`) {
			return ParseGitHubSearchResults(trimmed)
		}
		return "", fmt.Errorf("无法识别的响应格式: %w", err)
	}

	// 🎯 两阶段遍历：先找高优先级的 resource，再兜底 text
	// ========== Phase 1: 优先查找 resource（文件内容） ==========
	for _, item := range mcpResp.Content {
		if item.Type == "resource" && item.Resource != nil {
			var fileRes GitHubFileResource
			resourceBytes, err := json.Marshal(item.Resource)
			if err != nil {
				continue
			}
			if err := json.Unmarshal(resourceBytes, &fileRes); err != nil {
				continue
			}
			// 处理 base64 编码
			if fileRes.Encoding == "base64" {
				if decoded, err := base64.StdEncoding.DecodeString(fileRes.Text); err == nil {
					fileRes.Text = string(decoded)
					fileRes.Encoding = "utf-8"
				}
			}
			if strings.TrimSpace(fileRes.Text) != "" {
				return FormatGitHubFileResult(&GitHubFileResult{File: &fileRes}), nil
			}
		}
	}

	// ========== Phase 2: 如果没有 resource，再处理 text ==========
	var statusMsg string
	for _, item := range mcpResp.Content {
		if item.Type == "text" {
			text := strings.TrimSpace(item.Text)
			if text == "" {
				continue
			}

			// 🔹 目录列表（转义数组字符串）
			if strings.HasPrefix(text, "[") && strings.Contains(text, `"type"`) {
				return ParseGitHubDirectoryList(text)
			}

			// 🔹 GitHub Search 响应（转义对象字符串）
			if strings.Contains(text, `"total_count"`) && strings.Contains(text, `"items"`) {
				return ParseGitHubSearchResults(text)
			}

			// 🔹 普通状态消息（暂存，作为兜底）
			if statusMsg == "" {
				statusMsg = text
			}
		}
	}

	// 兜底：返回状态消息（如果没有找到 resource）
	if statusMsg != "" {
		return formatStatusMessage(statusMsg), nil
	}

	return "", fmt.Errorf("未找到可解析的有效内容")
}

// ============ 3. 目录列表解析 ============

// ParseGitHubDirectoryList 解析目录列表（支持原始数组 或 MCP 包装的字符串）
func ParseGitHubDirectoryList(jsonInput string) (string, error) {
	trimmed := strings.TrimSpace(jsonInput)

	// 兼容 MCP 包装：{"content":[{"type":"text","text":"[...]"}}]
	if strings.HasPrefix(trimmed, `{"content"`) {
		var mcp mcptools.MCPResponse
		if err := json.Unmarshal([]byte(trimmed), &mcp); err == nil {
			for _, item := range mcp.Content {
				if item.Type == "text" && strings.HasPrefix(strings.TrimSpace(item.Text), "[") {
					trimmed = strings.TrimSpace(item.Text)
					break
				}
			}
		}
	}

	// 解析目录数组
	var entries []GitHubDirEntry
	if err := json.Unmarshal([]byte(trimmed), &entries); err != nil {
		return "", fmt.Errorf("解析目录列表失败: %w", err)
	}

	return FormatGitHubDirectoryList(entries), nil
}

// FormatGitHubDirectoryList 格式化目录列表输出
func FormatGitHubDirectoryList(entries []GitHubDirEntry) string {
	if len(entries) == 0 {
		return "📁 空目录"
	}

	// 分类：目录在前，文件在后（按名称排序）
	var dirs, files []GitHubDirEntry
	for _, e := range entries {
		if e.Type == "dir" {
			dirs = append(dirs, e)
		} else {
			files = append(files, e)
		}
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name < dirs[j].Name })
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📁 共 %d 项 (%d 目录, %d 文件)\n\n", len(entries), len(dirs), len(files)))

	// 📁 目录
	for _, d := range dirs {
		sb.WriteString(fmt.Sprintf("  📁 %s/\n", d.Name))
	}

	// 📄 文件
	for _, f := range files {
		sizeStr := ""
		if f.Size > 0 {
			sizeStr = fmt.Sprintf(" (%s)", formatFileSize(int(f.Size)))
		}
		sb.WriteString(fmt.Sprintf("  📄 %s%s\n", f.Name, sizeStr))
	}

	return sb.String()
}

// ============ 4. 文件内容解析 ============

// ParseMCPGitHubFileContents 解析 MCP 包装的文件内容响应（兼容旧调用）
func ParseMCPGitHubFileContents(resultJSON string) (string, error) {
	return ParseMCPGitHubContents(resultJSON) // 委托给通用入口
}

// FormatGitHubFileResult 格式化文件内容输出
func FormatGitHubFileResult(result *GitHubFileResult) string {
	var sb strings.Builder

	// 📌 状态消息（如果有）
	if result.StatusMsg != "" {
		sb.WriteString(fmt.Sprintf("✅ %s\n\n", result.StatusMsg))
	}

	// 📄 文件内容
	if result.File != nil {
		filename := extractFilenameFromURI(result.File.URI)
		sb.WriteString(fmt.Sprintf("📄 %s\n", filename))

		// MIME 类型 + 大小
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

		// 完整内容输出（如需截断可取消注释 limitLines）
		for line := range strings.SplitSeq(result.File.Text, "\n") {
			sb.WriteString(line + "\n")
		}
	}

	return sb.String()
}

// ============ 5. GitHub Search 解析 ============

// ParseGitHubSearchResults 解析搜索结果（兼容旧调用）
func ParseGitHubSearchResults(resultJSON string) (string, error) {
	var resp GitHubSearchResponse
	if err := json.Unmarshal([]byte(resultJSON), &resp); err != nil {
		return "", fmt.Errorf("解析 Search JSON 失败: %w", err)
	}
	return FormatGitHubSearchResults(&resp), nil
}

// FormatGitHubSearchResults 格式化搜索结果输出

// ============ 6. 辅助函数 ============

func formatStatusMessage(msg string) string {
	return fmt.Sprintf("✅ %s\n", msg)
}

func extractFilenameFromURI(uri string) string {
	parts := strings.Split(uri, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

func formatFileSize(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}
