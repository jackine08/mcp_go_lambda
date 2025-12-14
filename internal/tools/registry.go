package tools

import (
	"context"

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
// The handler function is automatically wrapped and registered with proper types
func Register[In any](name, description string, handler func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, map[string]interface{}, error)) {
	registry = append(registry, Tool{
		Name:        name,
		Description: description,
		RegisterFn: func(server *mcp.Server) {
			mcp.AddTool(server, &mcp.Tool{
				Name:        name,
				Description: description,
			}, handler)
		},
	})
}

// RegisterAllTools registers all tools from the global registry to the MCP server
func RegisterAllTools(server *mcp.Server) {
	for _, tool := range registry {
		tool.RegisterFn(server)
	}
}
