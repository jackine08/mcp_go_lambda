package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GreetingInput defines input for greeting operations
type GreetingInput struct {
	Name     string `json:"name" jsonschema:"the name to greet (optional)"`
	Language string `json:"language" jsonschema:"language for greeting: 'ko' for Korean, 'en' for English (default: 'ko')"`
}

// SayHello responds to greetings in Korean or English
func SayHello(ctx context.Context, req *mcp.CallToolRequest, input GreetingInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	name := input.Name
	if name == "" {
		name = "친구" // "friend" in Korean
	}

	var greeting string
	language := input.Language
	if language == "" {
		language = "ko"
	}

	switch language {
	case "ko":
		greeting = fmt.Sprintf("안녕하세요, %s님! 👋", name)
	case "en":
		greeting = fmt.Sprintf("Hello, %s! 👋", name)
	default:
		greeting = fmt.Sprintf("ㅎㅇ %s! 👋", name)
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: greeting,
				},
			},
		}, map[string]interface{}{
			"greeting": greeting,
			"name":     name,
			"language": language,
		}, nil
}

// init automatically registers greeting tool
func init() {
	Register("say_hello", "인사를 건넵니다. 한국어(ko) 또는 영어(en)로 인사할 수 있습니다", SayHello)
}
