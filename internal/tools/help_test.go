package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHelp(t *testing.T) {
	ctx := context.Background()

	// Test 1: General help (no tool_name)
	t.Run("General help", func(t *testing.T) {
		result, meta, err := Help(ctx, nil, HelpInput{})
		if err != nil {
			t.Fatalf("Error calling help: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if len(result.Content) == 0 {
			t.Fatal("Expected content, got empty array")
		}

		textContent, ok := result.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatal("Expected TextContent")
		}

		if len(textContent.Text) == 0 {
			t.Fatal("Expected non-empty text")
		}

		// Check that response contains expected sections
		if !strings.Contains(textContent.Text, "계산기") {
			t.Error("Expected response to contain '계산기'")
		}

		if !strings.Contains(textContent.Text, "문자열") {
			t.Error("Expected response to contain '문자열'")
		}

		t.Logf("General help response:\n%s", textContent.Text)

		if meta == nil {
			t.Fatal("Expected meta, got nil")
		}

		if _, ok := meta["result"]; !ok {
			t.Fatal("Expected result in meta")
		}
	})

	// Test 2: Help for specific tool
	t.Run("Help for add tool", func(t *testing.T) {
		result, meta, err := Help(ctx, nil, HelpInput{ToolName: "add"})
		if err != nil {
			t.Fatalf("Error calling help: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if len(result.Content) == 0 {
			t.Fatal("Expected content, got empty array")
		}

		textContent, ok := result.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatal("Expected TextContent")
		}

		// Check that response contains tool name and description
		if !strings.Contains(textContent.Text, "add") {
			t.Error("Expected response to contain 'add'")
		}

		t.Logf("Help for 'add' tool:\n%s", textContent.Text)

		if meta == nil {
			t.Fatal("Expected meta, got nil")
		}
	})

	// Test 3: Help for non-existent tool
	t.Run("Help for non-existent tool", func(t *testing.T) {
		result, meta, err := Help(ctx, nil, HelpInput{ToolName: "nonexistent"})
		if err != nil {
			t.Fatalf("Error calling help: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if len(result.Content) == 0 {
			t.Fatal("Expected content, got empty array")
		}

		textContent, ok := result.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatal("Expected TextContent")
		}

		// Check that response indicates tool not found
		if !strings.Contains(textContent.Text, "찾을 수 없습니다") {
			t.Error("Expected response to indicate tool not found")
		}

		t.Logf("Help for non-existent tool:\n%s", textContent.Text)

		if meta == nil {
			t.Fatal("Expected meta, got nil")
		}
	})
}
