package main

import (
	"context"
	"encoding/json"
)

// MCPRequest는 MCP 클라이언트로부터 받는 요청 구조
type MCPRequest struct {
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id"`
}

// MCPResponse는 MCP 서버에서 응답하는 구조
type MCPResponse struct {
	JsonRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

// MCPError는 MCP 에러 응답
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool은 MCP Tool 정의
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// Server는 MCP 서버 구조체
type Server struct {
	// 필요한 필드는 나중에 추가
}

// NewServer는 새로운 MCP 서버를 생성합니다
func NewServer() *Server {
	return &Server{}
}

// Handle은 MCP 요청을 처리합니다
func (s *Server) Handle(ctx context.Context, request MCPRequest) MCPResponse {
	response := MCPResponse{
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
		response.Error = &MCPError{
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
	return map[string]interface{}{
		"resources": []interface{}{},
	}
}

// handleToolsList는 tools/list 메서드를 처리합니다
func (s *Server) handleToolsList(ctx context.Context, params interface{}) interface{} {
	tools := []Tool{
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
			"error": "Invalid arguments: a and b must be numbers",
		}
	}

	result := a + b
	return map[string]interface{}{
		"result": result,
	}
}

// toolMultiply는 두 숫자를 곱합니다
func (s *Server) toolMultiply(args map[string]interface{}) interface{} {
	a, ok1 := args["a"].(float64)
	b, ok2 := args["b"].(float64)

	if !ok1 || !ok2 {
		return map[string]interface{}{
			"error": "Invalid arguments: a and b must be numbers",
		}
	}

	result := a * b
	return map[string]interface{}{
		"result": result,
	}
}

// handlePromptsList는 prompts/list 메서드를 처리합니다
func (s *Server) handlePromptsList(ctx context.Context, params interface{}) interface{} {
	return map[string]interface{}{
		"prompts": []interface{}{},
	}
}
