package server

import (
	"github.com/jackine08/mcp_go_lambda/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewMCPServer creates and configures a new MCP server
func NewMCPServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-go-lambda",
		Version: "v2.0.0", // Updated for 7 tools (Calculator + StringTools)
	}, nil)

	// Create tool manager and register all providers
	toolManager := tools.NewToolManager()

	// Register tool providers - 새로운 provider를 여기에 추가하면 자동으로 등록됨
	toolManager.
		Register(tools.NewCalculator()).
		Register(tools.NewStringTools())
		// .Register(새로운Provider()) // 여기에 추가하면 자동 등록!

	// Register all tools from all providers
	toolManager.RegisterAll(server)

	return server
}
