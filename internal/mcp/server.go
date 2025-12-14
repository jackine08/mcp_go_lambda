package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackine08/mcp_go_lambda/internal/types"
)

// Server는 MCP 서버 구조체
type Server struct {
	// 필요한 필드는 나중에 추가
}

// NewServer는 새로운 MCP 서버를 생성합니다
func NewServer() *Server {
	return &Server{}
}

// Handle은 MCP 요청을 처리합니다
func (s *Server) Handle(ctx context.Context, request types.MCPRequest) types.MCPResponse {
	response := types.MCPResponse{
		JsonRPC: "2.0",
		ID:      request.ID,
	}

	switch request.Method {
	case "initialize":
		response.Result = s.handleInitialize(ctx, request.Params)
	case "resources/list":
		response.Result = s.handleResourcesList(ctx, request.Params)
	case "tools/list":
		response.Result = s.handleToolsList(ctx, request.Params)
	case "tools/call":
		response.Result = s.handleToolCall(ctx, request.Params)
	case "prompts/list":
		response.Result = s.handlePromptsList(ctx, request.Params)
	default:
		response.Error = &types.MCPError{
			Code:    -32601,
			Message: "Method not found",
		}
	}

	return response
}

// handleInitialize는 initialize 메서드를 처리합니다
func (s *Server) handleInitialize(ctx context.Context, params interface{}) interface{} {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"resources": map[string]interface{}{
				"listChanged": true,
			},
			"tools": map[string]interface{}{
				"listChanged": true,
			},
			"prompts": map[string]interface{}{
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "Go MCP Server",
			"version": "0.2.0",
		},
	}
}

// handleResourcesList는 resources/list 메서드를 처리합니다
func (s *Server) handleResourcesList(ctx context.Context, params interface{}) interface{} {
	resources := []map[string]interface{}{
		{
			"uri":         "file:///project/README.md",
			"name":        "프로젝트 README",
			"description": "프로젝트 설명 문서",
			"mimeType":    "text/markdown",
		},
		{
			"uri":         "https://api.example.com/status",
			"name":        "API 상태",
			"description": "서버 상태 확인 엔드포인트",
			"mimeType":    "application/json",
		},
		{
			"uri":         "file:///logs/app.log",
			"name":        "애플리케이션 로그",
			"description": "최근 애플리케이션 로그 파일",
			"mimeType":    "text/plain",
		},
	}

	return map[string]interface{}{
		"resources": resources,
	}
}

// handleToolsList는 tools/list 메서드를 처리합니다
func (s *Server) handleToolsList(ctx context.Context, params interface{}) interface{} {
	tools := []types.Tool{
		{
			Name:        "add",
			Description: "두 개의 숫자를 더합니다",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "number",
						"description": "첫 번째 숫자",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "두 번째 숫자",
					},
				},
				"required": []string{"a", "b"},
			},
		},
		{
			Name:        "multiply",
			Description: "두 개의 숫자를 곱합니다",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "number",
						"description": "첫 번째 숫자",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "두 번째 숫자",
					},
				},
				"required": []string{"a", "b"},
			},
		},
	}

	return map[string]interface{}{
		"tools": tools,
	}
}

// ToolCallParams는 tool call 요청 파라미터
type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// handleToolCall은 tools/call 메서드를 처리합니다
func (s *Server) handleToolCall(ctx context.Context, params interface{}) interface{} {
	var toolParams ToolCallParams

	// params를 JSON으로 변환했다가 다시 파싱
	paramsJSON, _ := json.Marshal(params)
	json.Unmarshal(paramsJSON, &toolParams)

	switch toolParams.Name {
	case "add":
		return s.toolAdd(toolParams.Arguments)
	case "multiply":
		return s.toolMultiply(toolParams.Arguments)
	default:
		return map[string]interface{}{
			"error": "Unknown tool: " + toolParams.Name,
		}
	}
}

// toolAdd는 두 숫자를 더합니다
func (s *Server) toolAdd(args map[string]interface{}) interface{} {
	a, ok1 := args["a"].(float64)
	b, ok2 := args["b"].(float64)

	if !ok1 || !ok2 {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "Error: Invalid arguments - a and b must be numbers",
				},
			},
		}
	}

	result := a + b
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("%.0f + %.0f = %.0f", a, b, result),
			},
		},
	}
}

// toolMultiply는 두 숫자를 곱합니다
func (s *Server) toolMultiply(args map[string]interface{}) interface{} {
	a, ok1 := args["a"].(float64)
	b, ok2 := args["b"].(float64)

	if !ok1 || !ok2 {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "Error: Invalid arguments - a and b must be numbers",
				},
			},
		}
	}

	result := a * b
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("%.0f × %.0f = %.0f", a, b, result),
			},
		},
	}
}

// handlePromptsList는 prompts/list 메서드를 처리합니다
func (s *Server) handlePromptsList(ctx context.Context, params interface{}) interface{} {
	prompts := []map[string]interface{}{
		{
			"name":        "code-review",
			"description": "코드 리뷰를 위한 프롬프트",
			"arguments": []map[string]interface{}{
				{
					"name":        "language",
					"description": "프로그래밍 언어",
					"required":    true,
				},
				{
					"name":        "code",
					"description": "리뷰할 코드",
					"required":    true,
				},
			},
		},
		{
			"name":        "bug-analysis",
			"description": "버그 분석 프롬프트",
			"arguments": []map[string]interface{}{
				{
					"name":        "error_message",
					"description": "에러 메시지",
					"required":    true,
				},
				{
					"name":        "context",
					"description": "에러 발생 컨텍스트",
					"required":    false,
				},
			},
		},
		{
			"name":        "documentation",
			"description": "문서 생성 프롬프트",
			"arguments": []map[string]interface{}{
				{
					"name":        "function_name",
					"description": "문서화할 함수명",
					"required":    true,
				},
			},
		},
	}

	return map[string]interface{}{
		"prompts": prompts,
	}
}
