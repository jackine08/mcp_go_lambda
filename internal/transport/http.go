package transport

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ProtocolVersion은 지원하는 MCP 프로토콜 버전
const (
	ProtocolVersion         = "2025-06-18"
	FallbackProtocolVersion = "2024-11-05"
)

// Session은 MCP 세션 정보
type Session struct {
	ID        string
	CreatedAt time.Time
	LastUsed  time.Time
}

// SessionStore는 세션 저장소
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewSessionStore는 새로운 세션 저장소를 생성
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
	}
}

// CreateSession은 새로운 세션을 생성
func (s *SessionStore) CreateSession() *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := generateSessionID()
	session := &Session{
		ID:        sessionID,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}
	s.sessions[sessionID] = session
	return session
}

// GetSession은 세션 ID로 세션을 조회
func (s *SessionStore) GetSession(sessionID string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if exists {
		session.LastUsed = time.Now()
	}
	return session, exists
}

// DeleteSession은 세션을 삭제
func (s *SessionStore) DeleteSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}

// CleanExpiredSessions는 만료된 세션을 정리 (1시간 이상 사용하지 않은 세션)
func (s *SessionStore) CleanExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	expireTime := time.Now().Add(-1 * time.Hour)
	for id, session := range s.sessions {
		if session.LastUsed.Before(expireTime) {
			delete(s.sessions, id)
		}
	}
}

// generateSessionID는 암호학적으로 안전한 세션 ID를 생성
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ProtectedResourceMetadata는 OAuth 2.0 Protected Resource Metadata
type ProtectedResourceMetadata struct {
	Resource               string   `json:"resource"`
	AuthorizationServers   []string `json:"authorization_servers"`
	BearerMethodsSupported []string `json:"bearer_methods_supported,omitempty"`
	ScopesSupported        []string `json:"scopes_supported,omitempty"`
}

// GetCanonicalURI는 API Gateway stage를 제거한 canonical URI를 생성
// 예: /dev/mcp -> /mcp
//
//	/prod/mcp -> /mcp
func GetCanonicalURI(host, path, stage string) string {
	// stage가 path의 prefix인 경우 제거
	if stage != "" && strings.HasPrefix(path, "/"+stage+"/") {
		path = strings.TrimPrefix(path, "/"+stage)
	} else if stage != "" && path == "/"+stage {
		path = "/"
	}

	// URL 구성 (trailing slash 제거)
	uri := fmt.Sprintf("https://%s%s", host, path)
	uri = strings.TrimSuffix(uri, "/")
	return uri
}

// ValidateProtocolVersion은 프로토콜 버전을 검증
func ValidateProtocolVersion(version string) bool {
	if version == "" {
		// 버전이 없으면 fallback 버전으로 간주 (하위 호환성)
		return true
	}
	return version == ProtocolVersion || version == FallbackProtocolVersion
}

// GetResourceMetadataURL은 resource metadata URL을 생성
func GetResourceMetadataURL(host string) string {
	return fmt.Sprintf("https://%s/.well-known/oauth-protected-resource", host)
}

// ParseBearerToken은 Authorization 헤더에서 Bearer 토큰을 추출
func ParseBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid authorization header format")
	}

	if strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("authorization scheme is not Bearer")
	}

	return parts[1], nil
}

// GetWWWAuthenticateHeader는 401 응답용 WWW-Authenticate 헤더를 생성
func GetWWWAuthenticateHeader(resourceMetadataURL string, scope string) string {
	header := fmt.Sprintf(`Bearer resource_metadata="%s"`, resourceMetadataURL)
	if scope != "" {
		header += fmt.Sprintf(`, scope="%s"`, scope)
	}
	return header
}

// GetInsufficientScopeHeader는 403 응답용 WWW-Authenticate 헤더를 생성
func GetInsufficientScopeHeader(resourceMetadataURL string, scope string, description string) string {
	header := fmt.Sprintf(`Bearer error="insufficient_scope", resource_metadata="%s", scope="%s"`,
		resourceMetadataURL, scope)
	if description != "" {
		header += fmt.Sprintf(`, error_description="%s"`, escapeQuotes(description))
	}
	return header
}

// escapeQuotes는 문자열의 따옴표를 이스케이프
func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

// ExtractHostFromAPIGateway는 API Gateway 요청에서 호스트를 추출
func ExtractHostFromAPIGateway(headers map[string]string) string {
	// Host 헤더 확인
	if host := headers["Host"]; host != "" {
		return host
	}
	if host := headers["host"]; host != "" {
		return host
	}

	// X-Forwarded-Host 헤더 확인
	if host := headers["X-Forwarded-Host"]; host != "" {
		return host
	}
	if host := headers["x-forwarded-host"]; host != "" {
		return host
	}

	return "localhost"
}

// NormalizeHeaderKey는 헤더 키를 정규화 (API Gateway는 소문자로 변환할 수 있음)
func NormalizeHeaderKey(headers map[string]string, key string) string {
	// 원본 키 확인
	if val, ok := headers[key]; ok {
		return val
	}

	// 소문자 키 확인
	lowerKey := strings.ToLower(key)
	if val, ok := headers[lowerKey]; ok {
		return val
	}

	// Title case 확인
	titleKey := strings.Title(strings.ToLower(key))
	if val, ok := headers[titleKey]; ok {
		return val
	}

	return ""
}

// IsWellKnownPath는 well-known 경로인지 확인
func IsWellKnownPath(path string) bool {
	return strings.Contains(path, "/.well-known/")
}

// ExtractResourcePath는 well-known 경로에서 리소스 경로를 추출
func ExtractResourcePath(path string) string {
	// /.well-known/oauth-protected-resource 또는
	// /.well-known/oauth-protected-resource/mcp 형태
	if strings.HasPrefix(path, "/.well-known/oauth-protected-resource") {
		if path == "/.well-known/oauth-protected-resource" {
			return "/"
		}
		// 경로 추출
		suffix := strings.TrimPrefix(path, "/.well-known/oauth-protected-resource")
		return suffix
	}
	return path
}
