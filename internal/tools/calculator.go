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

// Calculator provides calculator operations as MCP tools
type Calculator struct{}

// NewCalculator creates a new Calculator instance
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Add adds two numbers
func (c *Calculator) Add(ctx context.Context, req *mcp.CallToolRequest, input OperationInput) (
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
func (c *Calculator) Multiply(ctx context.Context, req *mcp.CallToolRequest, input OperationInput) (
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
func (c *Calculator) Subtract(ctx context.Context, req *mcp.CallToolRequest, input OperationInput) (
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
func (c *Calculator) Divide(ctx context.Context, req *mcp.CallToolRequest, input OperationInput) (
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

// RegisterTools registers all calculator tools (implements ToolProvider interface)
func (c *Calculator) RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "두 개의 숫자를 더합니다",
	}, c.Add)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "multiply",
		Description: "두 개의 숫자를 곱합니다",
	}, c.Multiply)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "subtract",
		Description: "두 개의 숫자를 뺍니다",
	}, c.Subtract)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "divide",
		Description: "두 개의 숫자를 나눕니다",
	}, c.Divide)
}
