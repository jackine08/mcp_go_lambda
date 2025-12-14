package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackine08/mcp_go_lambda/internal/handler"
	"github.com/jackine08/mcp_go_lambda/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Command line flags
	httpAddr := flag.String("http", "", "HTTP server address (e.g., localhost:8080)")
	flag.Parse()

	// Create MCP server with all tools registered
	mcpServer := server.NewMCPServer()

	// Check if running in Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		runLambda(mcpServer)
		return
	}

	// Check for PORT environment variable
	port := os.Getenv("PORT")
	if port == "" && *httpAddr != "" {
		port = *httpAddr
	}

	if port != "" {
		runHTTPServer(mcpServer, port)
	} else {
		runStdioServer(mcpServer)
	}
}

// runLambda starts the Lambda handler
func runLambda(mcpServer *mcp.Server) {
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return handler.HandleLambdaRequest(ctx, mcpServer, request)
	})
}

// runHTTPServer starts the HTTP server
func runHTTPServer(mcpServer *mcp.Server, port string) {
	// Create HTTP handler using go-sdk
	httpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return mcpServer
	}, nil)

	// If port is just a number, add colon
	if port[0] != ':' && port[0] >= '0' && port[0] <= '9' && len(port) < 6 {
		port = ":" + port
	}

	log.Printf("MCP server listening on %s", port)
	log.Fatal(http.ListenAndServe(port, httpHandler))
}

// runStdioServer starts the stdio server
func runStdioServer(mcpServer *mcp.Server) {
	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
