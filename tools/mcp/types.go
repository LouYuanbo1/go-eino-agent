package mcptools

type MCPResponse struct {
	Content []MCPContentItem `json:"content"`
}

type MCPContentItem struct {
	Type     string         `json:"type"`
	Text     string         `json:"text,omitempty"`     // ✅ 支持 text 类型
	Resource map[string]any `json:"resource,omitempty"` // ✅ 支持 resource 类型
}