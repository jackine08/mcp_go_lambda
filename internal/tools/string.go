package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// StringInput defines input for string operations
type StringInput struct {
	Text string `json:"text" jsonschema:"the input text"`
}

// ReverseInput defines input for reverse operation
type ReverseInput struct {
	Text string `json:"text" jsonschema:"the text to reverse"`
}

// CaseInput defines input for case conversion
type CaseInput struct {
	Text string `json:"text" jsonschema:"the text to convert"`
}

// StringTools provides string manipulation operations
type StringTools struct{}

// NewStringTools creates a new StringTools instance
func NewStringTools() *StringTools {
	return &StringTools{}
}

// ToUpper converts text to uppercase
func (s *StringTools) ToUpper(ctx context.Context, req *mcp.CallToolRequest, input CaseInput) (
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
func (s *StringTools) ToLower(ctx context.Context, req *mcp.CallToolRequest, input CaseInput) (
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
func (s *StringTools) Reverse(ctx context.Context, req *mcp.CallToolRequest, input ReverseInput) (
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

// RegisterTools registers all string tools (implements ToolProvider interface)
func (s *StringTools) RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "to_upper",
		Description: "텍스트를 대문자로 변환합니다",
	}, s.ToUpper)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "to_lower",
		Description: "텍스트를 소문자로 변환합니다",
	}, s.ToLower)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "reverse",
		Description: "텍스트를 역순으로 뒤집습니다",
	}, s.Reverse)
}
