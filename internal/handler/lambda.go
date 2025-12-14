package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HandleLambdaRequest processes each API Gateway request with a temporary session
func HandleLambdaRequest(ctx context.Context, server *mcp.Server, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// JSON 포맷으로 구조화된 로그 출력 (CloudWatch에서 파싱 가능)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // DEBUG 레벨까지 모두 출력
	}))

	startTime := time.Now()

	logger.Info("Lambda request received",
		"path", request.Path,
		"httpMethod", request.HTTPMethod,
		"bodyLength", len(request.Body),
	)

	// Parse JSON-RPC request
	var jsonRPCRequest map[string]interface{}
	if err := json.Unmarshal([]byte(request.Body), &jsonRPCRequest); err != nil {
		logger.Error("Failed to parse JSON-RPC request",
			"error", err,
			"body", request.Body,
		)
		return errorResponse(400, fmt.Sprintf("Invalid JSON: %v", err)), nil
	}

	method, _ := jsonRPCRequest["method"].(string)
	requestID, _ := jsonRPCRequest["id"]

	logger.Info("JSON-RPC request parsed",
		"method", method,
		"requestId", requestID,
		"params", jsonRPCRequest["params"],
	)

	// Create temporary in-memory session for this request
	logger.Debug("Creating in-memory transport session")
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	// Connect server
	logger.Debug("Connecting MCP server")
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		logger.Error("Failed to connect server", "error", err)
		return errorResponse(500, fmt.Sprintf("Failed to connect server: %v", err)), nil
	}
	defer serverSession.Close()

	// Create client
	logger.Debug("Creating MCP client")
	client := mcp.NewClient(&mcp.Implementation{Name: "lambda-client", Version: "1.0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		logger.Error("Failed to connect client", "error", err)
		return errorResponse(500, fmt.Sprintf("Failed to connect client: %v", err)), nil
	}
	defer clientSession.Close()

	logger.Debug("MCP session established successfully")

	var response interface{}

	// Handle MCP methods
	logger.Info("Processing MCP method", "method", method)

	switch method {
	case "initialize":
		logger.Debug("Handling initialize request")
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
		logger.Info("Initialize request completed")

	case "tools/list":
		logger.Debug("Calling clientSession.ListTools")
		result, err := clientSession.ListTools(ctx, nil)
		if err != nil {
			logger.Error("ListTools failed", "error", err)
			return errorResponse(500, err.Error()), nil
		}
		logger.Info("Tools list retrieved", "toolCount", len(result.Tools))
		response = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      jsonRPCRequest["id"],
			"result":  result,
		}

	case "tools/call":
		params, _ := jsonRPCRequest["params"].(map[string]interface{})
		toolName, _ := params["name"].(string)
		args, _ := params["arguments"].(map[string]interface{})

		logger.Info("Calling tool",
			"toolName", toolName,
			"arguments", args,
		)

		result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name:      toolName,
			Arguments: args,
		})
		if err != nil {
			logger.Error("Tool call failed",
				"toolName", toolName,
				"error", err,
			)
			return errorResponse(500, err.Error()), nil
		}

		logger.Info("Tool call completed",
			"toolName", toolName,
			"resultContentCount", len(result.Content),
		)

		response = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      jsonRPCRequest["id"],
			"result":  result,
		}

	case "notifications/initialized":
		logger.Debug("Handling initialized notification")
		// No response needed
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, nil

	default:
		logger.Warn("Unsupported method", "method", method)
		// Return empty result for unsupported methods
		response = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      jsonRPCRequest["id"],
			"result":  map[string]interface{}{},
		}
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal response", "error", err)
		return errorResponse(500, fmt.Sprintf("Marshal failed: %v", err)), nil
	}

	duration := time.Since(startTime)
	logger.Info("Request completed successfully",
		"method", method,
		"duration", duration.String(),
		"responseSize", len(responseBytes),
	)

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
