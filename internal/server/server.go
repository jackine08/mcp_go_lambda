package server

import (
	"github.com/jackine08/mcp_go_lambda/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewMCPServer creates and configures a new MCP server
func NewMCPServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-go-lambda",
		Version: "v1.0.0",
	}, nil)

	// Register all tools
	tools.RegisterCalculatorTools(server)

	return server
}
