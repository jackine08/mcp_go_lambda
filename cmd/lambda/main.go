package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackine08/mcp_go_lambda/internal/handler"
)

func main() {
	lambda.Start(handler.HandleAPIGatewayRequest)
}
