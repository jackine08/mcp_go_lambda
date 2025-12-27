package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HelpInput defines input for the help tool
type HelpInput struct {
	// Optional: if provided, shows details for specific tool
	ToolName string `json:"tool_name,omitempty" jsonschema:"Optional: Specific tool name to get details about."`
}

// Help provides information about what this MCP server can do
func Help(ctx context.Context, req *mcp.CallToolRequest, input HelpInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	var response strings.Builder

	if input.ToolName != "" {
		// Show details for specific tool
		found := false
		for _, tool := range registry {
			if tool.Name == input.ToolName {
				response.WriteString(fmt.Sprintf("🔧 **%s**\n", tool.Name))
				response.WriteString(fmt.Sprintf("설명: %s\n", tool.Description))
				found = true
				break
			}
		}
		if !found {
			response.WriteString(fmt.Sprintf("도구 '%s'를 찾을 수 없습니다.\n", input.ToolName))
		}
	} else {
		// Show overview of all capabilities
		response.WriteString("# MCP Go Lambda 서버 기능 안내\n\n")
		response.WriteString("이 서버는 다음과 같은 기능을 제공합니다:\n\n")

		// Group tools by category
		calculatorTools := []Tool{}
		stringTools := []Tool{}
		otherTools := []Tool{}

		for _, tool := range registry {
			switch tool.Name {
			case "add", "multiply", "subtract", "divide":
				calculatorTools = append(calculatorTools, tool)
			case "to_upper", "to_lower", "reverse":
				stringTools = append(stringTools, tool)
			case "help":
				// Skip help in the list
				continue
			default:
				otherTools = append(otherTools, tool)
			}
		}

		// Calculator tools
		if len(calculatorTools) > 0 {
			response.WriteString("## 🧮 계산기 기능\n")
			for _, tool := range calculatorTools {
				response.WriteString(fmt.Sprintf("- **%s**: %s\n", tool.Name, tool.Description))
			}
			response.WriteString("\n")
		}

		// String tools
		if len(stringTools) > 0 {
			response.WriteString("## 📝 문자열 조작 기능\n")
			for _, tool := range stringTools {
				response.WriteString(fmt.Sprintf("- **%s**: %s\n", tool.Name, tool.Description))
			}
			response.WriteString("\n")
		}

		// Other tools
		if len(otherTools) > 0 {
			response.WriteString("## 🛠️ 기타 기능\n")
			for _, tool := range otherTools {
				response.WriteString(fmt.Sprintf("- **%s**: %s\n", tool.Name, tool.Description))
			}
			response.WriteString("\n")
		}

		// Count tools excluding help itself
		toolCount := 0
		for _, tool := range registry {
			if tool.Name != "help" {
				toolCount++
			}
		}

		response.WriteString(fmt.Sprintf("\n총 %d개의 도구를 사용할 수 있습니다.\n", toolCount))
		response.WriteString("\n특정 도구에 대한 자세한 정보를 보려면 tool_name 파라미터를 지정하세요.")
	}

	result := response.String()
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: result,
			},
		},
	}, map[string]interface{}{"result": result}, nil
}

// init automatically registers the help tool
func init() {
	Register("help", "이 서버가 할 수 있는 일을 설명합니다 (뭐 할 수 있어?)", Help)
}
