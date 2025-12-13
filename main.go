package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HandleAPIGatewayRequest는 API Gateway 요청을 처리하는 Lambda 핸들러
func HandleAPIGatewayRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %s %s", request.HTTPMethod, request.Path)

	// 서버 생성
	server := NewServer()

	// 요청이 빈 경우 (GET /)
	if request.Body == "" {
		response := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"status": "MCP Server is running"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return response, nil
	}

	// MCP 요청 파싱
	var mcpRequest MCPRequest
	err := json.Unmarshal([]byte(request.Body), &mcpRequest)
	if err != nil {
		log.Printf("Failed to parse request: %v", err)
		errResponse := MCPResponse{
			JsonRPC: "2.0",
			Error: &MCPError{
				Code:    -32700,
				Message: "Parse error",
				Data:    err.Error(),
			},
			ID: nil,
		}
		responseBody, _ := json.Marshal(errResponse)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseBody),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// MCP 요청 처리
	mcpResponse := server.Handle(ctx, mcpRequest)
	responseBody, _ := json.Marshal(mcpResponse)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleAPIGatewayRequest)
}
