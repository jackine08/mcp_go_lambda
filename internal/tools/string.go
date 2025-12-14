package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TextInput defines input for string operations
type TextInput struct {
	Text string `json:"text" jsonschema:"the input text"`
}

// ToUpper converts text to uppercase
func ToUpper(ctx context.Context, req *mcp.CallToolRequest, input TextInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	result := strings.ToUpper(input.Text)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("변환 결과: %s", result),
			},
		},
	}, map[string]interface{}{"result": result}, nil
}

// ToLower converts text to lowercase
func ToLower(ctx context.Context, req *mcp.CallToolRequest, input TextInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	result := strings.ToLower(input.Text)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("변환 결과: %s", result),
			},
		},
	}, map[string]interface{}{"result": result}, nil
}

// Reverse reverses a string
func Reverse(ctx context.Context, req *mcp.CallToolRequest, input TextInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	runes := []rune(input.Text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	result := string(runes)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("원본: %s → 역순: %s", input.Text, result),
			},
		},
	}, map[string]interface{}{"result": result}, nil
}

// init automatically registers all string tools
func init() {
	Register("to_upper", "텍스트를 대문자로 변환합니다", func(server *mcp.Server) {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "to_upper",
			Description: "텍스트를 대문자로 변환합니다",
		}, ToUpper)
	})

	Register("to_lower", "텍스트를 소문자로 변환합니다", func(server *mcp.Server) {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "to_lower",
			Description: "텍스트를 소문자로 변환합니다",
		}, ToLower)
	})

	Register("reverse", "텍스트를 역순으로 뒤집습니다", func(server *mcp.Server) {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "reverse",
			Description: "텍스트를 역순으로 뒤집습니다",
		}, Reverse)
	})
}
