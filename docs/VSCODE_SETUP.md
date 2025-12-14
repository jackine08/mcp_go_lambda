# VS Code에서 MCP 서버 사용하기

## 개요

현재 배포된 HTTP MCP 서버를 VS Code의 GitHub Copilot과 함께 사용할 수 있습니다.

## VS Code MCP 설정

### 1. mcp.json 파일 생성

VS Code의 MCP 서버 설정 파일을 생성합니다:

**위치**: 
- 사용자 전역: `~/.vscode/mcp.json` (macOS/Linux) 또는 `%APPDATA%\Code\User\mcp.json` (Windows)
- 워크스페이스: `.vscode/mcp.json` (프로젝트 루트)

**설정 내용**:

```json
{
  "servers": {
    "mcp-go-lambda": {
      "type": "http",
      "url": "https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp",
      "headers": {
        "Authorization": "Bearer ${input:mcp-go-lambda-token}",
        "MCP-Protocol-Version": "2025-06-18"
      }
    }
  },
  "inputs": [
    {
      "type": "promptString",
      "id": "mcp-go-lambda-token",
      "description": "MCP Server API Token (test-token for demo)",
      "password": true
    }
  ]
}
```

### 2. VS Code에서 MCP 서버 시작

1. **Command Palette** 열기: `Ctrl+Shift+P` (Windows/Linux) 또는 `Cmd+Shift+P` (macOS)
2. **`MCP: List Servers`** 명령 실행
3. **`mcp-go-lambda`** 서버 선택
4. **Start** 선택
5. 토큰 입력 프롬프트가 나타나면 `test-token` 입력

### 3. Chat에서 MCP 도구 사용

1. **Chat View** 열기: `Ctrl+Alt+I`
2. **Tools** 버튼 클릭 (채팅 입력창 상단)
3. **mcp-go-lambda** 서버의 도구 확인:
   - `add`: 두 숫자 더하기
   - `multiply`: 두 숫자 곱하기

### 4. 도구 사용 예시

#### 자동 도구 호출 (Agent 모드):
```
15와 27을 더해줘
```

#### 명시적 도구 호출:
```
#add를 사용해서 15와 27을 더해줘
```

## 주요 기능

### 지원되는 전송 방식
- ✅ **HTTP Transport**: 현재 서버가 사용하는 방식
- ✅ **SSE (Server-Sent Events)**: VS Code가 자동으로 fallback

### 지원되는 기능
- ✅ **Tools**: `add`, `multiply` 도구
- ✅ **Authentication**: Bearer Token 인증
- ✅ **Session Management**: 자동 세션 관리
- ✅ **Protocol Version Validation**: 버전 호환성 검증

## 트러블슈팅

### MCP 서버가 시작되지 않는 경우

1. **Output 로그 확인**:
   - Command Palette > `MCP: List Servers`
   - 서버 선택 > `Show Output`

2. **인증 토큰 확인**:
   - 현재 mock 인증은 `test-token` 또는 `Bearer-*` 패턴 허용
   - 실제 환경에서는 JWT 토큰 필요

3. **네트워크 연결 확인**:
   ```bash
   curl -X POST https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp \
     -H "Authorization: Bearer test-token" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"vscode","version":"1.0.0"}}}'
   ```

### 도구가 보이지 않는 경우

1. **캐시 초기화**:
   - Command Palette > `MCP: Reset Cached Tools`

2. **서버 재시작**:
   - Command Palette > `MCP: List Servers`
   - 서버 선택 > `Restart`

## 고급 설정

### 환경 변수 사용

```json
{
  "servers": {
    "mcp-go-lambda": {
      "type": "http",
      "url": "${env:MCP_SERVER_URL}",
      "headers": {
        "Authorization": "Bearer ${env:MCP_API_TOKEN}"
      }
    }
  }
}
```

### 워크스페이스별 설정

프로젝트 루트에 `.vscode/mcp.json` 파일을 생성하면 워크스페이스별로 다른 MCP 서버를 사용할 수 있습니다.

## 참고 자료

- [VS Code MCP 공식 문서](https://code.visualstudio.com/docs/copilot/customization/mcp-servers)
- [Model Context Protocol 스펙](https://modelcontextprotocol.io/)
- [프로젝트 README](../README.md)
- [HTTP Transport 가이드](./HTTP_TRANSPORT_GUIDE.md)
