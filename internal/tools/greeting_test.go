package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSayHello(t *testing.T) {
	tests := []struct {
		name     string
		input    GreetingInput
		wantText string
	}{
		{
			name: "Korean greeting with name",
			input: GreetingInput{
				Name:     "철수",
				Language: "ko",
			},
			wantText: "안녕하세요, 철수님! 👋",
		},
		{
			name: "English greeting with name",
			input: GreetingInput{
				Name:     "John",
				Language: "en",
			},
			wantText: "Hello, John! 👋",
		},
		{
			name: "Default (Korean) greeting without name",
			input: GreetingInput{
				Language: "ko",
			},
			wantText: "안녕하세요, 친구님! 👋",
		},
		{
			name:     "Default language and name",
			input:    GreetingInput{},
			wantText: "안녕하세요, 친구님! 👋",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, metadata, err := SayHello(context.Background(), &mcp.CallToolRequest{}, tt.input)

			if err != nil {
				t.Errorf("SayHello() error = %v", err)
				return
			}

			if result == nil {
				t.Error("SayHello() result is nil")
				return
			}

			if len(result.Content) == 0 {
				t.Error("SayHello() result has no content")
				return
			}

			textContent, ok := result.Content[0].(*mcp.TextContent)
			if !ok {
				t.Error("SayHello() result content is not TextContent")
				return
			}

			if textContent.Text != tt.wantText {
				t.Errorf("SayHello() text = %v, want %v", textContent.Text, tt.wantText)
			}

			// Check metadata
			if metadata["greeting"] != tt.wantText {
				t.Errorf("SayHello() metadata[greeting] = %v, want %v", metadata["greeting"], tt.wantText)
			}
		})
	}
}
