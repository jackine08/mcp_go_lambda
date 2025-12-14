package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool represents a registered MCP tool with a registration function
type Tool struct {
	Name        string
	Description string
	RegisterFn  func(*mcp.Server) // Function that knows how to register this specific tool
}

// Global registry - tools register themselves via init()
var registry []Tool

// Register adds a tool to the global registry
// Called automatically from init() functions in tool files
// The registerFn should call mcp.AddTool with the correct types
func Register(name, description string, registerFn func(*mcp.Server)) {
	registry = append(registry, Tool{
		Name:        name,
		Description: description,
		RegisterFn:  registerFn,
	})
}

// RegisterAllTools registers all tools from the global registry to the MCP server
func RegisterAllTools(server *mcp.Server) {
	for _, tool := range registry {
		tool.RegisterFn(server)
	}
}
