# MCP HTTP Transport 구현 가이드

## 개요
이 프로젝트는 MCP (Model Context Protocol) 서버를 AWS Lambda에 배포하며, 최신 Streamable HTTP Transport와 OAuth 2.1 인증을 지원합니다.

## 주요 기능

### 1. Streamable HTTP Transport (2025-06-18)
- **POST /mcp**: JSON-RPC 요청 처리
- **GET /mcp**: SSE 스트림 지원 (현재는 상태 확인용)
- **프로토콜 버전 검증**: `MCP-Protocol-Version` 헤더 지원

### 2. OAuth 2.1 인증 (선택적)
- Bearer Token 인증
- WWW-Authenticate 헤더를 통한 인증 요구
- Protected Resource Metadata 엔드포인트

### 3. 세션 관리
- 세션 ID 자동 생성 (initialize 시)
- `Mcp-Session-Id` 헤더를 통한 세션 추적
- 메모리 기반 세션 스토어 (Lambda 컨테이너 재사용)

### 4. API Gateway Stage 처리
- Stage path 자동 정규화 (`/dev/mcp` → `/mcp`)
- Canonical URI 생성

## API 엔드포인트

### MCP 서버
```
POST /mcp
GET  /mcp
```

### OAuth 메타데이터
```
GET /.well-known/oauth-protected-resource
```

## 요청 헤더

### 필수 헤더
- `Content-Type: application/json`
- `MCP-Protocol-Version: 2025-06-18` (권장)

### 선택적 헤더
- `Authorization: Bearer <token>` (인증 활성화 시)
- `Mcp-Session-Id: <session-id>` (세션 유지 시)

## 응답 헤더

- `MCP-Protocol-Version`: 서버가 사용하는 프로토콜 버전
- `Mcp-Session-Id`: 세션 ID (initialize 시 생성)
- `WWW-Authenticate`: 인증 실패 시 (401/403)

## 인증 설정

현재 인증은 **기본적으로 비활성화**되어 있습니다. 활성화하려면:

1. 환경변수 설정:
```bash
export MCP_REQUIRE_AUTH=true
export MCP_AUTH_SERVER=https://your-auth-server.com
```

2. 또는 코드에서 직접 수정:
```go
// internal/handler/lambda.go
var authMiddleware = auth.NewMiddleware(&auth.Config{
    RequireAuth:         true,  // 인증 활성화
    AuthorizationServer: "https://auth.example.com",
    AllowedScopes:      []string{"mcp:read", "mcp:write"},
})
```

## 테스트 이벤트

### 1. Initialize (세션 생성)
```bash
# deploy/events/initialize-with-auth.json
POST /dev/mcp
MCP-Protocol-Version: 2025-06-18
```

### 2. Tools List (세션 사용)
```bash
# deploy/events/tools-list-with-auth.json
POST /dev/mcp
Authorization: Bearer test-token
Mcp-Session-Id: abc123def456
```

### 3. OAuth Metadata 조회
```bash
# deploy/events/well-known-metadata.json
GET /dev/.well-known/oauth-protected-resource
```

### 4. 서버 상태 확인
```bash
# deploy/events/get-mcp-status.json
GET /dev/mcp
```

## 프로토콜 버전 호환성

- **2025-06-18**: 최신 Streamable HTTP Transport
- **2024-11-05**: 하위 호환 (Fallback)
- 버전 헤더 없음: 2024-11-05로 간주

## API Gateway Stage 처리

API Gateway는 stage를 path에 자동으로 추가합니다:
- 실제 요청: `https://xxx.amazonaws.com/dev/mcp`
- Canonical URI: `https://xxx.amazonaws.com/mcp`

서버는 자동으로 stage를 제거하고 정규화된 경로를 생성합니다.

## 에러 코드

- `400 Bad Request`: 잘못된 요청, 지원하지 않는 프로토콜 버전
- `401 Unauthorized`: 인증 필요
- `403 Forbidden`: 권한 부족
- `404 Not Found`: 세션 만료 또는 리소스 없음
- `405 Method Not Allowed`: 지원하지 않는 HTTP 메서드

## 배포

```bash
# 빌드 및 배포
make deploy

# 또는 개별 단계
make build
cd cdk && cdk deploy
```

## 로컬 테스트

```bash
# 빌드
make build

# SAM으로 로컬 테스트 (선택)
sam local invoke -e deploy/events/initialize-with-auth.json
```

## 보안 고려사항

1. **토큰 검증**: 현재는 mock 구현, 프로덕션에서는 실제 JWT 검증 필요
2. **CORS**: 프로덕션에서는 특정 도메인으로 제한
3. **HTTPS**: API Gateway는 자동으로 HTTPS 제공
4. **Session Expiry**: 1시간 후 자동 만료

## 다음 단계

- [ ] JWT 토큰 검증 구현
- [ ] SSE 스트리밍 구현
- [ ] 프로덕션용 인증 서버 연동
- [ ] 세션 영속화 (DynamoDB 등)
- [ ] Rate limiting
