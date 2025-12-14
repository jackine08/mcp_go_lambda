# MCP Go Lambda Server

Go로 작성된 Model Context Protocol (MCP) 서버를 AWS Lambda에 배포하는 프로젝트입니다.

## 프로젝트 구조

```
mcp-go-lambda/
├── cmd/
│   └── lambda/           # Lambda 엔트리포인트
│       └── main.go       # 애플리케이션 시작점
├── internal/             # 내부 패키지 (외부에서 import 불가)
│   ├── handler/          # Lambda 핸들러
│   │   └── lambda.go     # API Gateway 요청 처리
│   ├── mcp/              # MCP 서버 구현
│   │   └── server.go     # MCP 프로토콜 로직
│   └── types/            # 공통 타입 정의
│       └── types.go      # Request/Response 타입
├── cdk/                  # AWS CDK 인프라 코드 (Python)
│   ├── app.py            # CDK 앱
│   └── stacks/
│       └── mcp_lambda_stack.py  # Lambda + API Gateway 스택
├── docs/                 # 프로젝트 문서
├── deploy/               # 배포 스크립트 및 이벤트 샘플
├── lambda_dist/          # 빌드 결과물 (bootstrap 바이너리)
├── old_structure/        # 기존 flat 구조 백업
│   ├── main.go
│   ├── server.go
│   └── server_test.go
├── Makefile              # 빌드/테스트/배포 명령어
├── go.mod                # Go 모듈 정의
├── go.sum                # 의존성 체크섬
└── README.md
```

### 패키지 설명

- **cmd/lambda**: Lambda 함수의 엔트리포인트. `main()` 함수만 포함
- **internal/types**: MCP 프로토콜의 타입 정의 (Request, Response, Error, Tool 등)
- **internal/mcp**: MCP 서버 비즈니스 로직 (초기화, 도구 목록, 도구 실행 등)
- **internal/handler**: AWS Lambda와 MCP 서버를 연결하는 핸들러 계층
- **cdk**: AWS 인프라를 코드로 정의 (Infrastructure as Code)

## 주요 구성 요소

### Internal Packages

- **cmd/lambda**: Lambda 함수의 엔트리포인트
- **internal/handler**: API Gateway 요청을 처리하는 Lambda 핸들러
- **internal/mcp**: MCP 프로토콜 서버 구현

## ✨ VS Code/GitHub Copilot 연동

이 MCP 서버는 **VS Code의 GitHub Copilot**에서 바로 사용할 수 있습니다!

### 빠른 시작

1. **VS Code에서 mcp.json 생성** (`~/.vscode/mcp.json` 또는 `.vscode/mcp.json`):

```json
{
  "servers": {
    "mcp-go-lambda": {
      "type": "http",
      "url": "https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp",
      "headers": {
        "Authorization": "Bearer test-token",
        "MCP-Protocol-Version": "2025-06-18"
      }
    }
  }
}
```

2. **VS Code에서 사용**:
   - `Ctrl+Alt+I`로 Chat View 열기
   - Tools 버튼 클릭
   - `add`, `multiply` 도구 선택
   - "15와 27을 더해줘" 같은 프롬프트 입력

자세한 내용은 [VS Code 설정 가이드](docs/VSCODE_SETUP.md)를 참고하세요.
- **internal/types**: 공통 타입 정의 (Request, Response, Error 등)

### MCP 기능

현재 구현된 MCP 메서드:
- `initialize`: 서버 초기화
- `tools/list`: 사용 가능한 도구 목록 조회
- `tools/call`: 도구 실행
- `resources/list`: 리소스 목록 조회
- `prompts/list`: 프롬프트 목록 조회

### Tools

- **add**: 두 숫자 더하기
- **multiply**: 두 숫자 곱하기

## 빌드 및 배포

### 1. AWS Credentials 설정

`.env` 파일에 AWS credentials 추가:
```bash
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
AWS_DEFAULT_REGION=ap-northeast-2
```

### 2. Makefile 사용 (권장)

```bash
# 빌드만
make build

# 빌드 + 배포
make deploy

# 스택 삭제
make destroy

# 테스트
make test

# 모든 단계 (clean, deps, test, build)
make all
```

### 3. 수동 빌드 (필요시)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lambda_dist/bootstrap cmd/lambda/main.go
```

## 테스트

```bash
# 단위 테스트 실행
go test ./...

# 특정 패키지 테스트
go test ./internal/handler

# 커버리지 포함
go test -cover ./...
```

## API 엔드포인트

배포 후 API Gateway 엔드포인트:
```
POST https://<api-id>.execute-api.<region>.amazonaws.com/dev/mcp
GET  https://<api-id>.execute-api.<region>.amazonaws.com/dev/mcp
GET  https://<api-id>.execute-api.<region>.amazonaws.com/dev/.well-known/oauth-protected-resource
```

### HTTP Transport 기능 (2025-06-18)

- **Streamable HTTP Transport**: POST/GET 메서드 지원
- **OAuth 2.1 인증**: Bearer Token 기반 (선택적)
- **세션 관리**: Mcp-Session-Id 헤더
- **프로토콜 버전**: MCP-Protocol-Version 헤더 지원
- **API Gateway Stage 자동 처리**: /dev/mcp → /mcp 정규화

자세한 내용은 [HTTP Transport Guide](docs/HTTP_TRANSPORT_GUIDE.md) 참고

### 사용 예시

```bash
# Initialize
curl -X POST "https://<endpoint>/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }'

# Tools List
curl -X POST "https://<endpoint>/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'

# Tool Call
curl -X POST "https://<endpoint>/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "add",
      "arguments": {
        "a": 10,
        "b": 20
      }
    }
  }'
```

## GitHub Copilot 연동

VS Code의 `mcp.json` 파일에 추가:

```json
{
  "mcpServers": {
    "aws-lambda-mcp": {
      "url": "https://<api-id>.execute-api.<region>.amazonaws.com/dev/mcp",
      "transport": {
        "type": "http"
      },
      "name": "AWS Lambda MCP Server",
      "description": "Go MCP Server running on AWS Lambda"
    }
  }
}
```

## 개발

### 로컬 테스트

```bash
# 로컬에서 테스트 이벤트 사용
cd deploy
sam local invoke MCPServerFunction -e events/initialize-request.json
```

### 새로운 Tool 추가

1. `internal/mcp/server.go`의 `handleToolsList`에 도구 정의 추가
2. `handleToolsCall`에 도구 로직 구현
3. 테스트 작성

## 라이선스

MIT License

## 참고

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [AWS Lambda Go](https://github.com/aws/aws-lambda-go)
- [AWS CDK Python](https://docs.aws.amazon.com/cdk/v2/guide/home.html)
