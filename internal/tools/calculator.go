package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AddInput defines the input parameters for the add tool
type AddInput struct {
	A float64 `json:"a" jsonschema:"the first number"`
	B float64 `json:"b" jsonschema:"the second number"`
}

// MultiplyInput defines the input parameters for the multiply tool
type MultiplyInput struct {
	A float64 `json:"a" jsonschema:"the first number"`
	B float64 `json:"b" jsonschema:"the second number"`
}

// SubtractInput defines the input parameters for the subtract tool
type SubtractInput struct {
	A float64 `json:"a" jsonschema:"the first number"`
	B float64 `json:"b" jsonschema:"the second number"`
}

// ResultOutput defines the output structure for calculation results
type ResultOutput struct {
	Result float64 `json:"result" jsonschema:"the calculation result"`
}

// AddTool implements the add operation
func AddTool(ctx context.Context, req *mcp.CallToolRequest, input AddInput) (
	*mcp.CallToolResult,
	ResultOutput,
	error,
) {
	result := input.A + input.B
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%f + %f = %f", input.A, input.B, result),
			},
		},
	}, ResultOutput{Result: result}, nil
}

// MultiplyTool implements the multiply operation
func MultiplyTool(ctx context.Context, req *mcp.CallToolRequest, input MultiplyInput) (
	*mcp.CallToolResult,
	ResultOutput,
	error,
) {
	result := input.A * input.B
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%f × %f = %f", input.A, input.B, result),
			},
		},
	}, ResultOutput{Result: result}, nil
}

// SubtractTool implements the subtract operation
func SubtractTool(ctx context.Context, req *mcp.CallToolRequest, input SubtractInput) (
	*mcp.CallToolResult,
	ResultOutput,
	error,
) {
	result := input.A - input.B
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%f - %f = %f", input.A, input.B, result),
			},
		},
	}, ResultOutput{Result: result}, nil
}

// RegisterCalculatorTools registers all calculator tools with the MCP server
func RegisterCalculatorTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "두 개의 숫자를 더합니다",
	}, AddTool)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "multiply",
		Description: "두 개의 숫자를 곱합니다",
	}, MultiplyTool)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "subtract",
		Description: "두 개의 숫자를 뺍니다",
	}, SubtractTool)
}
