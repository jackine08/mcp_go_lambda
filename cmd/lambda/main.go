package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

// SubtractInput defines the input parameters for the subtract tool
type SubtractInput struct {
	A float64 `json:"a" jsonschema:"the first number"`
	B float64 `json:"b" jsonschema:"the second number"`
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

func main() {
	// Command line flags
	httpAddr := flag.String("http", "", "HTTP server address (e.g., localhost:8080)")
	flag.Parse()

	// Create MCP server with implementation details
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-go-lambda",
		Version: "v1.0.0",
	}, nil)

	// Add the add tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "두 개의 숫자를 더합니다",
	}, AddTool)

	// Add the multiply tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "multiply",
		Description: "두 개의 숫자를 곱합니다",
	}, MultiplyTool)

	// Add the subtract tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "subtract",
		Description: "두 개의 숫자를 뺍니다",
	}, SubtractTool)

	// Create HTTP handler using go-sdk
	httpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	// Check if running in Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda mode - use stateless handler
		lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			return handleLambdaRequest(ctx, server, request)
		})
		return
	}

	// Check for PORT environment variable
	port := os.Getenv("PORT")
	if port == "" && *httpAddr != "" {
		port = *httpAddr
	}

	if port != "" {
		// If port is just a number, add colon
		if port[0] != ':' && port[0] >= '0' && port[0] <= '9' && len(port) < 6 {
			port = ":" + port
		}

		log.Printf("MCP server listening on %s", port)
		log.Fatal(http.ListenAndServe(port, httpHandler))
	} else {
		// Stdio mode - use StdioTransport
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}
}

// handleLambdaRequest processes each request with a temporary session
func handleLambdaRequest(ctx context.Context, server *mcp.Server, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse JSON-RPC request
	var jsonRPCRequest map[string]interface{}
	if err := json.Unmarshal([]byte(request.Body), &jsonRPCRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error":"Invalid JSON: %v"}`, err),
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	method, _ := jsonRPCRequest["method"].(string)

	// Create temporary in-memory session for this request
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	// Connect server
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error":"Failed to connect server: %v"}`, err),
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer serverSession.Close()

	// Create client
	client := mcp.NewClient(&mcp.Implementation{Name: "lambda-client", Version: "1.0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error":"Failed to connect client: %v"}`, err),
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer clientSession.Close()

	var response interface{}

	// Handle MCP methods
	switch method {
	case "tools/list":
		result, err := clientSession.ListTools(ctx, nil)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf(`{"error":"%v"}`, err),
				Headers: map[string]string{
					"Content-Type":                "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			}, nil
		}
		response = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      jsonRPCRequest["id"],
			"result":  result,
		}

	case "tools/call":
		params, _ := jsonRPCRequest["params"].(map[string]interface{})
		toolName, _ := params["name"].(string)
		args, _ := params["arguments"].(map[string]interface{})

		result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name:      toolName,
			Arguments: args,
		})
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf(`{"error":"%v"}`, err),
				Headers: map[string]string{
					"Content-Type":                "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			}, nil
		}
		response = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      jsonRPCRequest["id"],
			"result":  result,
		}

	case "notifications/initialized":
		// No response needed
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, nil

	default:
		// Return empty result for unsupported methods
		response = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      jsonRPCRequest["id"],
			"result":  map[string]interface{}{},
		}
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error":"Marshal failed: %v"}`, err),
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBytes),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}
