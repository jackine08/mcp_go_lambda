package auth

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackine08/mcp_go_lambda/internal/transport"
)

// Config는 인증 설정
type Config struct {
	// AuthorizationServer는 OAuth 2.1 인증 서버 URL
	AuthorizationServer string
	// RequireAuth는 인증을 필수로 할지 여부
	RequireAuth bool
	// AllowedScopes는 허용된 scope 목록
	AllowedScopes []string
	// DefaultScope는 기본 scope
	DefaultScope string
}

// Middleware는 인증 미들웨어
type Middleware struct {
	config *Config
}

// NewMiddleware는 새로운 인증 미들웨어를 생성
func NewMiddleware(config *Config) *Middleware {
	if config == nil {
		config = &Config{
			RequireAuth: false, // 기본적으로 인증 비활성화
		}
	}
	return &Middleware{config: config}
}

// ValidateToken은 Bearer 토큰을 검증
// 현재는 mock 구현 - 실제로는 OAuth 서버에 토큰을 검증해야 함
func (m *Middleware) ValidateToken(token string) (bool, error) {
	if !m.config.RequireAuth {
		// 인증이 필요 없으면 항상 통과
		return true, nil
	}

	if token == "" {
		return false, fmt.Errorf("token is empty")
	}

	// TODO: 실제 토큰 검증 로직 구현
	// - JWT 파싱 및 서명 검증
	// - 만료 시간 확인
	// - audience 확인 (RFC 8707)
	// - scope 확인

	// 현재는 mock으로 "test-token"만 허용
	if token == "test-token" || strings.HasPrefix(token, "Bearer-") {
		return true, nil
	}

	return false, fmt.Errorf("invalid token")
}

// GetWWWAuthenticateResponse는 401 Unauthorized 응답 생성
func (m *Middleware) GetWWWAuthenticateResponse(host string, stage string) map[string]interface{} {
	resourceMetadataURL := transport.GetResourceMetadataURL(host)
	wwwAuth := transport.GetWWWAuthenticateHeader(resourceMetadataURL, m.config.DefaultScope)

	return map[string]interface{}{
		"statusCode": 401,
		"headers": map[string]string{
			"Content-Type":     "application/json",
			"WWW-Authenticate": wwwAuth,
		},
		"body": mustMarshalJSON(map[string]interface{}{
			"error":             "unauthorized",
			"error_description": "Authorization required",
			"resource_metadata": resourceMetadataURL,
		}),
	}
}

// GetInsufficientScopeResponse는 403 Forbidden 응답 생성
func (m *Middleware) GetInsufficientScopeResponse(host string, requiredScope string, description string) map[string]interface{} {
	resourceMetadataURL := transport.GetResourceMetadataURL(host)
	wwwAuth := transport.GetInsufficientScopeHeader(resourceMetadataURL, requiredScope, description)

	return map[string]interface{}{
		"statusCode": 403,
		"headers": map[string]string{
			"Content-Type":     "application/json",
			"WWW-Authenticate": wwwAuth,
		},
		"body": mustMarshalJSON(map[string]interface{}{
			"error":             "insufficient_scope",
			"error_description": description,
			"required_scope":    requiredScope,
		}),
	}
}

// GetProtectedResourceMetadata는 OAuth 2.0 Protected Resource Metadata 반환
func (m *Middleware) GetProtectedResourceMetadata(host string, resourcePath string) *transport.ProtectedResourceMetadata {
	canonicalURI := fmt.Sprintf("https://%s%s", host, resourcePath)
	canonicalURI = strings.TrimSuffix(canonicalURI, "/")

	authServers := []string{m.config.AuthorizationServer}
	if m.config.AuthorizationServer == "" {
		// 기본 인증 서버 (데모용)
		authServers = []string{"https://auth.example.com"}
	}

	return &transport.ProtectedResourceMetadata{
		Resource:               canonicalURI,
		AuthorizationServers:   authServers,
		BearerMethodsSupported: []string{"header"},
		ScopesSupported:        m.config.AllowedScopes,
	}
}

// mustMarshalJSON는 JSON 마샬링을 수행하고 실패시 패닉
func mustMarshalJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return `{"error": "internal_error"}`
	}
	return string(b)
}
