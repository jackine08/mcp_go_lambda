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

	// All tools are automatically registered via init() functions in tool packages
	// 새로운 tool을 추가하려면 internal/tools/ 아래에 파일만 생성하고
	// init() 함수에서 Register() 호출하면 자동 등록됨
	tools.RegisterAllTools(server)

	return server
}
