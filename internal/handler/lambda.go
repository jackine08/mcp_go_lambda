package handler

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/jackine08/mcp_go_lambda/internal/auth"
	"github.com/jackine08/mcp_go_lambda/internal/mcp"
	"github.com/jackine08/mcp_go_lambda/internal/transport"
	"github.com/jackine08/mcp_go_lambda/internal/types"
)

const (
	// HTTP Methods
	httpMethodGET  = "GET"
	httpMethodPOST = "POST"

	// HTTP Headers
	headerContentType = "Content-Type"
	headerAllow       = "Allow"
	headerCacheControl = "Cache-Control"

	// Content Types
	contentTypeJSON = "application/json"

	// MCP Methods
	methodInitialize = "initialize"
)

// 전역 세션 스토어 (Lambda 컨테이너 재사용)
var sessionStore = transport.NewSessionStore()

// 인증 미들웨어 (환경변수로 설정 가능)
var authMiddleware = auth.NewMiddleware(&auth.Config{
	RequireAuth:         false, // 기본값: 인증 비활성화 (테스트용)
	AuthorizationServer: "https://auth.example.com",
	AllowedScopes:       []string{"mcp:read", "mcp:write"},
	DefaultScope:        "mcp:read mcp:write",
})

// HandleAPIGatewayRequest는 API Gateway 요청을 처리하는 Lambda 핸들러
func HandleAPIGatewayRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %s %s", request.HTTPMethod, request.Path)

	// 호스트 및 경로 정보 추출
	host := transport.ExtractHostFromAPIGateway(request.Headers)
	stage := request.RequestContext.Stage
	path := request.Path

	log.Printf("Host: %s, Stage: %s, Path: %s", host, stage, path)

	// 1. Protocol Version 검증
	protocolVersion := transport.NormalizeHeaderKey(request.Headers, "MCP-Protocol-Version")
	if protocolVersion != "" && !transport.ValidateProtocolVersion(protocolVersion) {
		log.Printf("Unsupported protocol version: %s", protocolVersion)
		return jsonErrorResponse(400, `{"error": "unsupported_protocol_version", "supported_versions": ["2025-06-18", "2024-11-05"]}`), nil
	}

	// 2. Well-known 경로 처리 (OAuth Protected Resource Metadata)
	if transport.IsWellKnownPath(path) {
		return handleWellKnownRequest(host, path, stage)
	}

	// 3. 인증 검증 (선택적)
	authHeader := transport.NormalizeHeaderKey(request.Headers, "Authorization")
	if authMiddleware != nil {
		token, err := transport.ParseBearerToken(authHeader)
		if err != nil && authMiddleware != nil {
			// 인증 헤더가 없거나 잘못된 경우
			log.Printf("Authentication failed: %v", err)
			respData := authMiddleware.GetWWWAuthenticateResponse(host, stage)
			return convertToAPIGatewayResponse(respData), nil
		}

		// 토큰 검증
		if token != "" {
			valid, err := authMiddleware.ValidateToken(token)
			if err != nil || !valid {
				log.Printf("Token validation failed: %v", err)
				respData := authMiddleware.GetWWWAuthenticateResponse(host, stage)
				return convertToAPIGatewayResponse(respData), nil
			}
			log.Printf("Token validated successfully")
		}
	}

	// 4. Session 관리 (선택적 - stateless 지원)
	sessionID := transport.NormalizeHeaderKey(request.Headers, "Mcp-Session-Id")
	var session *transport.Session
	if sessionID != "" {
		var exists bool
		session, exists = sessionStore.GetSession(sessionID)
		if !exists {
			// 세션이 없으면 경고만 로그, 계속 진행 (stateless 지원)
			log.Printf("Warning: Session ID provided but not found: %s. Continuing without session.", sessionID)
			session = nil
		}
	}

	// 서버 생성
	server := mcp.NewServer()

	// GET 요청 처리 (SSE 스트림용 - 현재는 간단한 상태 응답)
	if request.HTTPMethod == httpMethodGET {
		return jsonResponse(200, `{"status": "MCP Server is running", "protocol_version": "2025-06-18"}`, nil), nil
	}

	// POST 요청이 아니면 405 반환
	if request.HTTPMethod != httpMethodPOST {
		return jsonErrorResponse(405, `{"error": "method_not_allowed"}`), nil
	}

	// Body가 비어있으면 400
	if request.Body == "" {
		return jsonErrorResponse(400, `{"error": "empty_body"}`), nil
	}

	// MCP 요청 파싱
	var mcpRequest types.MCPRequest
	if err := json.Unmarshal([]byte(request.Body), &mcpRequest); err != nil {
		log.Printf("Failed to parse request: %v", err)
		errResponse := types.MCPResponse{
			JsonRPC: "2.0",
			Error: &types.MCPError{
				Code:    -32700,
				Message: "Parse error",
				Data:    err.Error(),
			},
			ID: nil,
		}
		responseBody, _ := json.Marshal(errResponse)
		return jsonResponse(400, string(responseBody), nil), nil
	}

	// MCP 요청 처리
	mcpResponse := server.Handle(ctx, mcpRequest)

	// initialize 메서드인 경우 세션 생성
	responseHeaders := map[string]string{
		headerContentType: contentTypeJSON,
	}

	if mcpRequest.Method == methodInitialize && mcpResponse.Error == nil {
		// 새 세션 생성
		newSession := sessionStore.CreateSession()
		responseHeaders["Mcp-Session-Id"] = newSession.ID
		log.Printf("Created new session: %s", newSession.ID)
	} else if sessionID != "" && session != nil {
		// 기존 세션 유지
		responseHeaders["Mcp-Session-Id"] = session.ID
	}

	// Protocol version 헤더 추가
	if protocolVersion != "" {
		responseHeaders["MCP-Protocol-Version"] = protocolVersion
	} else {
		responseHeaders["MCP-Protocol-Version"] = transport.ProtocolVersion
	}

	responseBody, _ := json.Marshal(mcpResponse)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    responseHeaders,
	}, nil
}

// handleWellKnownRequest는 .well-known 경로 요청을 처리
func handleWellKnownRequest(host, path, stage string) (events.APIGatewayProxyResponse, error) {
	if strings.Contains(path, "/.well-known/oauth-protected-resource") {
		// OAuth Protected Resource Metadata 반환
		resourcePath := "/mcp"
		if stage != "" {
			// stage를 제외한 경로
			resourcePath = transport.ExtractResourcePath(path)
			if resourcePath == "/" || resourcePath == "" {
				resourcePath = "/mcp"
			}
		}

		metadata := authMiddleware.GetProtectedResourceMetadata(host, resourcePath)
		body, _ := json.Marshal(metadata)

		return jsonResponse(200, string(body), map[string]string{
			headerCacheControl: "public, max-age=3600",
		}), nil
	}

	return jsonErrorResponse(404, `{"error": "not_found"}`), nil
}

// convertToAPIGatewayResponse는 map을 APIGatewayProxyResponse로 변환
func convertToAPIGatewayResponse(data map[string]interface{}) events.APIGatewayProxyResponse {
	statusCode := 500
	if sc, ok := data["statusCode"].(int); ok {
		statusCode = sc
	}

	headers := make(map[string]string)
	if h, ok := data["headers"].(map[string]string); ok {
		headers = h
	}

	body := ""
	if b, ok := data["body"].(string); ok {
		body = b
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
	}
}

// jsonErrorResponse는 JSON 형식의 에러 응답을 생성
func jsonErrorResponse(statusCode int, errorMsg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       errorMsg,
		Headers: map[string]string{
			headerContentType: contentTypeJSON,
		},
	}
}

// jsonResponse는 JSON 형식의 성공 응답을 생성
func jsonResponse(statusCode int, body string, extraHeaders map[string]string) events.APIGatewayProxyResponse {
	headers := map[string]string{
		headerContentType: contentTypeJSON,
	}
	for k, v := range extraHeaders {
		headers[k] = v
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
	}
}
