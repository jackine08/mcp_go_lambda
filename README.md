# MCP Go Lambda Server

Go로 작성된 Model Context Protocol (MCP) 서버를 AWS Lambda에 배포하는 프로젝트입니다.

[Model Context Protocol](https://modelcontextprotocol.io/) 공식 [go-sdk](https://github.com/modelcontextprotocol/go-sdk)를 사용하여 구현되었습니다.

## ✨ 주요 특징

- 🚀 **공식 Go SDK 사용**: MCP 공식 go-sdk v1.1.0 기반
- 🔧 **자동 Tool 등록**: Self-registering pattern으로 파일 하나만 추가하면 자동 등록
- 📦 **AWS Lambda 배포**: API Gateway + Lambda로 HTTP Transport 지원
- 🎯 **간단한 확장**: 제네릭 기반 Register 함수로 최소한의 코드로 Tool 추가

## 프로젝트 구조

```
mcp-go-lambda/
├── cmd/
│   └── lambda/           # Lambda 엔트리포인트
│       └── main.go       # 애플리케이션 시작점
├── internal/
│   ├── handler/          # Lambda 핸들러
│   │   └── lambda.go     # API Gateway 요청 처리
│   ├── server/           # MCP 서버 팩토리
│   │   └── server.go     # 서버 생성 및 설정
│   └── tools/            # Tool 구현 (자동 등록)
│       ├── registry.go   # Tool 레지스트리 (제네릭)
│       ├── calculator.go # 계산기 툴
│       └── string.go     # 문자열 조작 툴
├── cdk/                  # AWS CDK 인프라 코드 (Python)
│   ├── app.py
│   ├── stacks/
│   │   └── mcp_lambda_stack.py
│   └── events/           # 테스트 이벤트
├── lambda_dist/          # 빌드 결과물
├── Makefile
├── go.mod
└── README.md
```

### 패키지 설명

- **cmd/lambda**: Lambda 함수의 엔트리포인트
- **internal/handler**: API Gateway 이벤트를 MCP 프로토콜로 변환
- **internal/server**: MCP 서버 생성 및 Tool 등록 관리
- **internal/tools**: Tool 구현체들 (init() 함수로 자동 등록)

## 🛠️ 구현된 Tools

### 계산기 (Calculator)
- **add**: 두 개의 숫자를 더합니다
- **multiply**: 두 개의 숫자를 곱합니다
- **subtract**: 두 개의 숫자를 뺍니다
- **divide**: 두 개의 숫자를 나눕니다

### 문자열 조작 (String Manipulation)
- **to_upper**: 텍스트를 대문자로 변환합니다
- **to_lower**: 텍스트를 소문자로 변환합니다
- **reverse**: 텍스트를 역순으로 뒤집습니다

## 🚀 새로운 Tool 추가하기

Self-registering pattern 덕분에 매우 간단합니다!

### 1. `internal/tools/` 아래에 새 파일 생성

```go
// internal/tools/mytool.go
package tools

import (
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Input 타입 정의
type MyInput struct {
	Text string `json:"text" jsonschema:"입력 텍스트"`
}

// Tool 함수 구현
func MyNewTool(ctx context.Context, req *mcp.CallToolRequest, input MyInput) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	// 로직 구현
	result := "처리된 결과: " + input.Text
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, map[string]interface{}{"result": result}, nil
}

// init()에서 자동 등록
func init() {
	Register("my_new_tool", "내 새로운 도구 설명", MyNewTool)
}
```

### 2. 빌드 및 배포

```bash
make deploy
```

끝! `server.go`나 다른 파일을 수정할 필요가 없습니다.

### 작동 원리

1. **init() 함수**: Go 패키지가 import될 때 자동 실행
2. **Register() 제네릭 함수**: Handler 함수의 타입을 자동으로 추론하여 등록
3. **RegisterAllTools()**: 서버 생성 시 한 번 호출하여 모든 Tool 등록

이 패턴은 Go 표준 라이브러리(`database/sql`, `image` 패키지)에서 사용하는 방식과 동일합니다.

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

## 🌐 API 엔드포인트

배포 후 API Gateway 엔드포인트:
```
POST https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp
```

### HTTP Transport 기능

- **InMemoryTransport**: 세션 관리를 위한 stateless Lambda 지원
- **JSON-RPC 2.0**: MCP 프로토콜 표준 준수
- **API Gateway 통합**: REST API로 외부 접근 가능

### 사용 예시

```bash
# Initialize (서버 연결 초기화)
curl -X POST "https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp" \
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

# Tools List (사용 가능한 도구 목록)
curl -X POST "https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }'

# Tool Call (도구 실행 - 덧셈 예시)
curl -X POST "https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "add",
      "arguments": {
        "a": 42,
        "b": 58
      }
    }
  }'

# 문자열 뒤집기 예시
curl -X POST "https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "reverse",
      "arguments": {
        "text": "Hello World"
      }
    }
  }'
```

## 📝 개발

### 로컬 테스트

```bash
# 테스트 이벤트로 로컬 실행
cd cdk/events
curl -X POST "http://localhost:9000/2015-03-31/functions/function/invocations" \
  -d @initialize-request.json
```

### 아키텍처 패턴

이 프로젝트는 **Self-Registering Pattern**을 사용합니다:

1. 각 Tool은 `init()` 함수에서 자동으로 등록됩니다
2. `Register()` 제네릭 함수가 타입 추론을 처리합니다
3. 중앙 레지스트리가 모든 Tool을 관리합니다
4. 새 Tool 추가 시 다른 파일 수정이 불필요합니다

이는 Go 표준 라이브러리(`database/sql`, `image`)와 동일한 패턴입니다.

## 라이선스

MIT License

## 참고

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [AWS Lambda Go](https://github.com/aws/aws-lambda-go)
- [AWS CDK Python](https://docs.aws.amazon.com/cdk/v2/guide/home.html)
