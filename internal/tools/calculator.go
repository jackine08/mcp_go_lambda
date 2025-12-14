package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// OperationInput defines the common input for calculator operations
type OperationInput struct {
	A float64 `json:"a" jsonschema:"the first number"`
	B float64 `json:"b" jsonschema:"the second number"`
}

// Add adds two numbers
func Add(ctx context.Context, req *mcp.CallToolRequest, input OperationInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	result := input.A + input.B
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%f + %f = %f", input.A, input.B, result),
			},
		},
	}, map[string]interface{}{"result": result}, nil
}

// Multiply multiplies two numbers
func Multiply(ctx context.Context, req *mcp.CallToolRequest, input OperationInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	result := input.A * input.B
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%f × %f = %f", input.A, input.B, result),
			},
		},
	}, map[string]interface{}{"result": result}, nil
}

// Subtract subtracts two numbers
func Subtract(ctx context.Context, req *mcp.CallToolRequest, input OperationInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	result := input.A - input.B
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%f - %f = %f", input.A, input.B, result),
			},
		},
	}, map[string]interface{}{"result": result}, nil
}

// Divide divides two numbers
func Divide(ctx context.Context, req *mcp.CallToolRequest, input OperationInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	if input.B == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: division by zero",
				},
			},
			IsError: true,
		}, nil, fmt.Errorf("division by zero")
	}
	result := input.A / input.B
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%f ÷ %f = %f", input.A, input.B, result),
			},
		},
	}, map[string]interface{}{"result": result}, nil
}

// init automatically registers all calculator tools
func init() {
	Register("add", "두 개의 숫자를 더합니다", Add)
	Register("multiply", "두 개의 숫자를 곱합니다", Multiply)
	Register("subtract", "두 개의 숫자를 뺍니다", Subtract)
	Register("divide", "두 개의 숫자를 나눕니다", Divide)
}
