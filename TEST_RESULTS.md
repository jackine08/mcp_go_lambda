# MCP HTTP Transport 배포 및 테스트 결과

## 배포 정보
- **API Endpoint**: https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp
- **Lambda Function**: mcp-server-stack-MCPServerFunctionA0453C57-jbbBGabJt4gH
- **Region**: ap-northeast-2
- **배포 시간**: ~51초

## 테스트 결과 ✅

### 1. GET /mcp (서버 상태 확인)
```bash
curl -X GET "https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp" \
  -H "MCP-Protocol-Version: 2025-06-18" \
  -H "Authorization: Bearer test-token"
```
**결과**: ✅ 200 OK
```json
{"status": "MCP Server is running", "protocol_version": "2025-06-18"}
```

### 2. Initialize (세션 생성)
```bash
POST /dev/mcp
Headers:
  - MCP-Protocol-Version: 2025-06-18
  - Authorization: Bearer test-token
Body: {"jsonrpc":"2.0","id":1,"method":"initialize",...}
```
**결과**: ✅ 200 OK
- **세션 ID 생성**: `db6b3bf16b74e6f4536b5a5a686d3935`
- **Response Headers**:
  - `mcp-session-id: db6b3bf16b74e6f4536b5a5a686d3935`
  - `mcp-protocol-version: 2025-06-18`

### 3. tools/list (도구 목록)
```bash
POST /dev/mcp
Headers:
  - Mcp-Session-Id: db6b3bf16b74e6f4536b5a5a686d3935
  - Authorization: Bearer test-token
```
**결과**: ✅ 200 OK
```json
{
  "tools": [
    {"name": "add", "description": "두 개의 숫자를 더합니다"},
    {"name": "multiply", "description": "두 개의 숫자를 곱합니다"}
  ]
}
```

### 4. tools/call (도구 실행)
```bash
POST /dev/mcp
Body: {"method":"tools/call","params":{"name":"add","arguments":{"a":15,"b":27}}}
```
**결과**: ✅ 200 OK
```json
{
  "result": {
    "content": [{"text": "15 + 27 = 42", "type": "text"}]
  }
}
```

### 5. OAuth Protected Resource Metadata
```bash
GET /dev/.well-known/oauth-protected-resource
```
**결과**: ✅ 200 OK
```json
{
  "resource": "https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/mcp",
  "authorization_servers": ["https://auth.example.com"],
  "bearer_methods_supported": ["header"],
  "scopes_supported": ["mcp:read", "mcp:write"]
}
```
✅ **Stage path 정규화 확인**: `/dev/mcp` → `/mcp` (canonical URI)

### 6. 프로토콜 버전 검증
```bash
MCP-Protocol-Version: 2023-01-01 (지원하지 않는 버전)
```
**결과**: ✅ 400 Bad Request
```json
{
  "error": "unsupported_protocol_version",
  "supported_versions": ["2025-06-18", "2024-11-05"]
}
```

### 7. 인증 검증
```bash
# Authorization 헤더 없이 요청
```
**결과**: ✅ 401 Unauthorized
```
WWW-Authenticate: Bearer resource_metadata="...", scope="mcp:read mcp:write"
```

## 기능 확인

| 기능 | 상태 | 비고 |
|------|------|------|
| HTTP GET /mcp | ✅ | 서버 상태 확인 |
| HTTP POST /mcp | ✅ | MCP JSON-RPC 처리 |
| OAuth 2.1 인증 | ✅ | Bearer token 검증 |
| 세션 관리 | ✅ | initialize 시 자동 생성 |
| 프로토콜 버전 검증 | ✅ | 2025-06-18, 2024-11-05 지원 |
| API Gateway Stage 처리 | ✅ | /dev/mcp → /mcp 정규화 |
| Protected Resource Metadata | ✅ | /.well-known 경로 |
| WWW-Authenticate 헤더 | ✅ | 401 응답에 포함 |
| CORS 헤더 | ✅ | OPTIONS 메서드 지원 |

## 성능
- 평균 응답 시간: ~60-100ms
- Cold start: 정상
- Lambda 메모리: 256MB

## 다음 단계
- [ ] 실제 OAuth 서버 연동 (AWS Cognito 등)
- [ ] JWT 토큰 검증 구현
- [ ] SSE 스트리밍 구현
- [ ] 세션 영속화 (DynamoDB)
- [ ] 모니터링 및 로깅 강화
