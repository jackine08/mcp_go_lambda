package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HandleLambdaRequest processes each API Gateway request with a temporary session
func HandleLambdaRequest(ctx context.Context, server *mcp.Server, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse JSON-RPC request
	var jsonRPCRequest map[string]interface{}
	if err := json.Unmarshal([]byte(request.Body), &jsonRPCRequest); err != nil {
		return errorResponse(400, fmt.Sprintf("Invalid JSON: %v", err)), nil
	}

	method, _ := jsonRPCRequest["method"].(string)

	// Create temporary in-memory session for this request
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	// Connect server
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		return errorResponse(500, fmt.Sprintf("Failed to connect server: %v", err)), nil
	}
	defer serverSession.Close()

	// Create client
	client := mcp.NewClient(&mcp.Implementation{Name: "lambda-client", Version: "1.0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		return errorResponse(500, fmt.Sprintf("Failed to connect client: %v", err)), nil
	}
	defer clientSession.Close()

	var response interface{}

	// Handle MCP methods
	switch method {
	case "initialize":
		// Return server capabilities
		response = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      jsonRPCRequest["id"],
			"result": map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "mcp-go-lambda",
					"version": "v1.0.0",
				},
			},
		}

	case "tools/list":
		result, err := clientSession.ListTools(ctx, nil)
		if err != nil {
			return errorResponse(500, err.Error()), nil
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
			return errorResponse(500, err.Error()), nil
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
		return errorResponse(500, fmt.Sprintf("Marshal failed: %v", err)), nil
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

// errorResponse creates a standardized error response
func errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       fmt.Sprintf(`{"error":"%s"}`, message),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}
