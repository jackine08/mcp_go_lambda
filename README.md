# MCP Server on AWS Lambda (Go)

**Model Context Protocol (MCP)** 서버를 Go로 구현하여 AWS Lambda에 배포하는 프로젝트입니다.

## 🎯 프로젝트 개요

이 프로젝트는 다음을 포함합니다:

- **Go로 작성된 MCP 서버** - JSON-RPC 2.0 기반
- **AWS Lambda + API Gateway** - 서버리스 배포
- **MCP Tools 구현** - add, multiply 등의 계산 도구
- **CloudWatch 모니터링** - 로깅, 알람, 대시보드
- **자동 배포 스크립트** - 한 명령으로 배포 가능

## 📊 진행 단계

| 단계 | 내용 | 상태 |
|------|------|------|
| 1 | MCP 기본 구조 구현 | ✅ 완료 |
| 2 | 개발 환경 설정 (Go, Python venv, uv) | ✅ 완료 |
| 3 | AWS Lambda 배포 | ✅ 완료 |
| 4 | MCP Tools 추가 (add, multiply) | ✅ 완료 |
| 5 | MCP 프로토콜 호환성 수정 | ✅ 완료 |

## 🚀 배포된 서버

**API Endpoint:**
```
https://w2z3biq5nb.execute-api.ap-northeast-2.amazonaws.com/dev
```

**지원하는 메서드:**
- `initialize` - 서버 초기화
- `tools/list` - 사용 가능한 Tool 목록 조회
- `tools/call` - Tool 실행
- `resources/list` - 리소스 목록 (미구현)
- `prompts/list` - 프롬프트 목록 (미구현)

## 🛠️ 구현된 Tools

### 1. **add** - 두 숫자 더하기
```bash
curl -X POST "https://w2z3biq5nb.execute-api.ap-northeast-2.amazonaws.com/dev/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "add",
      "arguments": {"a": 10, "b": 5}
    },
    "id": 1
  }'
```

**응답:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "10 + 5 = 15"
      }
    ]
  },
  "id": 1
}
```

### 2. **multiply** - 두 숫자 곱하기
```bash
curl -X POST "https://w2z3biq5nb.execute-api.ap-northeast-2.amazonaws.com/dev/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "multiply",
      "arguments": {"a": 6, "b": 7}
    },
    "id": 2
  }'
```

## 📂 프로젝트 구조

```
mcp_go_lambda/
├── 📄 핵심 코드
│   ├── main.go              # Lambda 핸들러 (70줄)
│   ├── server.go            # MCP 서버 로직 (237줄)
│   └── server_test.go       # 단위 테스트 (6/6 PASS)
│
├── 📦 Go 의존성
│   ├── go.mod               # Go 모듈 정의
│   └── go.sum               # 의존성 체크섬
│
├── ☁️ AWS 배포
│   ├── template.yaml        # SAM 템플릿 (CloudFormation)
│   ├── samconfig.toml       # SAM 배포 설정
│   └── deploy.sh            # 자동 배포 스크립트
│
├── 🔧 개발 도구
│   ├── Makefile             # 빌드 자동화
│   ├── requirements.txt      # Python 의존성 (SAM CLI)
│   └── .env                 # AWS 자격증명 (gitignore)
│
├── 📋 테스트
│   └── events/              # 테스트용 MCP 요청 파일들
│
└── 📚 문서
    ├── README.md            # 이 파일
    ├── PROJECT_SUMMARY.md   # 완료된 작업 요약
    ├── COMPLETION_REPORT.md # 최종 완료 보고서
    ├── FILE_CLEANUP_GUIDE.md # 파일 정리 가이드
    ├── DEPLOYMENT_GUIDE.md  # 상세 배포 가이드
    ├── DEPLOYMENT_CHECKLIST.md # 배포 체크리스트
    └── AWS_SETUP.md         # AWS 설정 가이드
```

## 🏗️ 아키텍처

```
Claude / GitHub Copilot
        ↓ HTTP MCP 요청
    API Gateway (public)
        ↓
    Lambda Function (mcp-server-dev)
        ├→ MCP Handler (main.go)
        └→ MCP Server Logic (server.go)
             ├→ Tools: add, multiply
             ├→ Resources: (추후 구현)
             └→ Prompts: (추후 구현)
        ↓
    CloudWatch
        ├→ Logs: /aws/lambda/mcp-server-dev
        ├→ Dashboard: mcp-server-dev
        └→ Alarm: mcp-server-dev-errors
```

## 🔄 워크플로우

### 로컬 개발
```bash
# 1. 코드 수정
vim server.go

# 2. 테스트 실행
go test -v

# 3. 로컬 빌드 확인
go build -o bootstrap main.go server.go
```

### 배포
```bash
# 1. 자동 배포 스크립트 실행
./deploy.sh

# 또는 수동 배포
make build
sam deploy --stack-name mcp-server-stack --no-confirm-changeset --resolve-s3 --capabilities CAPABILITY_IAM
```

### 테스트
```bash
# 1. Initialize
curl -X POST "https://w2z3biq5nb.execute-api.ap-northeast-2.amazonaws.com/dev/mcp" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "initialize", "params": {}, "id": 1}'

# 2. Tool 호출
curl -X POST "https://w2z3biq5nb.execute-api.ap-northeast-2.amazonaws.com/dev/mcp" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "add", "arguments": {"a": 10, "b": 5}}, "id": 2}'

# 3. 로그 확인
aws logs tail /aws/lambda/mcp-server-dev --follow
```

## 📋 주요 파일 설명

### `main.go` (Lambda 핸들러)
- API Gateway 요청을 받음
- JSON-RPC 요청 파싱
- Server로 요청 전달
- MCP 응답을 HTTP 응답으로 변환

### `server.go` (MCP 서버)
- MCPRequest/MCPResponse 타입 정의
- Handle 메서드로 라우팅
- initialize, tools/list, tools/call 등 구현
- add, multiply Tool 로직

### `template.yaml` (SAM 템플릿)
- Lambda 함수 정의
- API Gateway 설정
- IAM Role, CloudWatch 설정
- 파라미터 (Environment: dev/staging/prod)

### `Makefile` (빌드 자동화)
- `make init`: 초기 설정 (venv, 의존성)
- `make build`: Go 컴파일 및 zip 생성
- `make deploy`: AWS Lambda 배포
- `make clean`: 산출물 정리

## 🔐 보안 고려사항

현재 상태:
- ❌ API 인증 없음 (public)
- ✅ HTTPS only (API Gateway 제공)
- ⚠️ Tool이 제한적 (계산만 가능)

추후 개선:
- [ ] API Key 또는 OAuth 추가
- [ ] Tool 입력 검증 강화
- [ ] Rate limiting 추가
- [ ] CloudWatch 알람 확대

## 🔧 커스터마이징

### 새로운 Tool 추가
`server.go`에서:
```go
case "new_tool":
    return s.handleNewTool(toolParams.Arguments)

func (s *Server) handleNewTool(args map[string]interface{}) interface{} {
    // 구현
}
```

### Tool 스키마 변경
`handleToolsList()`에서 Tool 정의 수정

## 📊 비용 (월간)

AWS 프리 티어 (충분함):
- Lambda: 100만 요청 무료
- API Gateway: 100만 호출 무료
- CloudWatch: 5GB 로그 무료

## 🚨 문제 해결

### "Missing Authentication Token" 에러
→ `template.yaml`에서 `Auth: DefaultAuthorizer: NONE` 확인

### "GLIBC 버전" 에러
→ `Makefile`에서 `CGO_ENABLED=0` 확인

### "TypeError: o.content is not iterable"
→ MCP 응답이 `content` 배열 형식이어야 함

## 📚 참고 자료

- [Model Context Protocol 문서](https://modelcontextprotocol.io/)
- [AWS Lambda 문서](https://docs.aws.amazon.com/lambda/)
- [AWS SAM 문서](https://docs.aws.amazon.com/serverless-application-model/)
- [Go AWS SDK](https://github.com/aws/aws-sdk-go)

## 📝 라이선스

MIT License

## 👤 작성자

jackine08

---

## 🎓 학습 목표 달성

- ✅ Go로 MCP 서버 구현
- ✅ AWS Lambda에 배포
- ✅ CloudFormation으로 인프라 자동화
- ✅ CI/CD 파이프라인 구성 (deploy.sh)
- ✅ CloudWatch 모니터링 설정
- ✅ JSON-RPC 2.0 프로토콜 이해
- ✅ API Gateway 통합
