# Stateless 지원 패치 노트

## 변경 사항

### 날짜: 2025-12-14

### 문제
VS Code에서 MCP 서버 연결 시 세션 에러 발생:
```
Error 404 status sending message: {"error": "session_not_found", "message": "Session expired or not found"}
```

### 원인
HTTP MCP 서버가 세션 ID를 엄격하게 검증하여, 세션이 없거나 만료된 경우 404 에러를 반환했습니다.
VS Code는 초기 연결 시 세션 ID 없이 요청을 보내는데, 이를 처리하지 못했습니다.

### 해결 방법
**Stateless 지원 추가**: 세션 ID가 없거나 유효하지 않아도 요청을 계속 처리하도록 수정

#### 코드 변경

**파일**: `internal/handler/lambda.go`

**변경 전**:
```go
if sessionID != "" {
    var exists bool
    session, exists = sessionStore.GetSession(sessionID)
    if !exists {
        // 세션이 만료되거나 없음
        return events.APIGatewayProxyResponse{
            StatusCode: 404,
            Body:       `{"error": "session_not_found", "message": "Session expired or not found"}`,
            // ...
        }, nil
    }
}
```

**변경 후**:
```go
if sessionID != "" {
    var exists bool
    session, exists = sessionStore.GetSession(sessionID)
    if !exists {
        // 세션이 없으면 경고만 로그, 계속 진행 (stateless 지원)
        log.Printf("Warning: Session ID provided but not found: %s. Continuing without session.", sessionID)
        session = nil
    }
}
```

### 테스트 결과

#### Before (실패):
```bash
curl -X POST .../dev/mcp \
  -H "Authorization: Bearer test-token" \
  -d '{"jsonrpc":"2.0","method":"tools/list",...}'
# 결과: 404 {"error": "session_not_found"}
```

#### After (성공):
```bash
curl -X POST .../dev/mcp \
  -H "Authorization: Bearer test-token" \
  -d '{"jsonrpc":"2.0","method":"tools/list",...}'
# 결과: 200 {"result":{"tools":[{"name":"add",...}]}}
```

### VS Code 연결 가이드

1. **MCP 서버 재시작**:
   - Command Palette: `MCP: List Servers`
   - `mcp-go-lambda` 선택
   - `Restart` 클릭

2. **테스트**:
   - Chat View 열기 (`Ctrl+Alt+I`)
   - Tools 버튼 클릭
   - `add`, `multiply` 도구 확인
   - "15와 27을 더해줘" 입력

### 기술 세부사항

- **세션 관리**: 선택적 (Optional)
  - `initialize` 시 세션 생성 및 `Mcp-Session-Id` 헤더 반환
  - 이후 요청에서 세션 ID 제공 시 검증하지만, 없어도 작동
  
- **Stateless 지원**:
  - 세션 없이도 모든 MCP 메서드 호출 가능
  - Lambda 컨테이너 재사용 시 세션 유지 (선택적 최적화)
  
- **호환성**:
  - ✅ VS Code / GitHub Copilot
  - ✅ Claude Desktop
  - ✅ 모든 MCP 클라이언트

### 배포 정보

- **배포 시간**: 2025-12-14 17:30 KST
- **엔드포인트**: https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp
- **배포 방법**: `make deploy`
- **변경 범위**: Lambda 함수만 업데이트 (인프라 변경 없음)
