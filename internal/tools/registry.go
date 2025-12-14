package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolProvider는 MCP tool을 제공하는 인터페이스
type ToolProvider interface {
	// RegisterTools는 이 provider의 모든 tool을 서버에 등록
	RegisterTools(server *mcp.Server)
}

// ToolManager는 여러 tool provider를 관리하고 자동으로 등록
type ToolManager struct {
	providers []ToolProvider
}

// NewToolManager는 새로운 ToolManager를 생성
func NewToolManager() *ToolManager {
	return &ToolManager{
		providers: make([]ToolProvider, 0),
	}
}

// Register는 tool provider를 등록 (체이닝 가능)
func (tm *ToolManager) Register(provider ToolProvider) *ToolManager {
	tm.providers = append(tm.providers, provider)
	return tm
}

// RegisterAll은 등록된 모든 provider의 tool을 MCP 서버에 등록
func (tm *ToolManager) RegisterAll(server *mcp.Server) {
	for _, provider := range tm.providers {
		provider.RegisterTools(server)
	}
}
