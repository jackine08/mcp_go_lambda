# 📋 프로젝트 요약: MCP 서버 AWS Lambda 배포

> **마지막 업데이트:** 2024년 현재
> **프로젝트 상태:** ✅ 완료 및 배포됨

## 📌 개요

Go로 작성한 Model Context Protocol (MCP) 서버를 AWS Lambda에 배포했습니다. 
Claude Desktop 클라이언트와 GitHub Copilot에서 직접 사용 가능합니다.

---

## 🎯 수행한 작업 (완료된 항목)

### 1️⃣ **개발 환경 설정** ✅

```bash
# 설치된 도구
✅ Go 1.22
✅ Python 3.13 + venv + uv
✅ AWS CLI
✅ AWS SAM CLI (aws-sam-cli==1.120.0)

# 의존성
✅ github.com/aws/aws-lambda-go v1.41.0
```

### 2️⃣ **MCP 서버 핵심 구현** ✅

#### `main.go` (108줄)
- API Gateway 이벤트 수신
- JSON-RPC 요청 파싱
- Server.Handle() 호출
- HTTP 응답 변환

#### `server.go` (227줄)
- **MCPRequest/MCPResponse 구조체** 정의
- **Handle() 메서드** - 라우팅 로직
- **initialize** - 서버 초기화
- **tools/list** - Tool 목록 반환
- **tools/call** - Tool 실행
- **resources/list** - 리소스 목록 (구조체 정의됨)
- **prompts/list** - 프롬프트 목록 (구조체 정의됨)

#### `server_test.go` 
- **6개 test case 모두 PASS** ✅
  - TestHandleInitialize
  - TestHandleToolsList
  - TestHandleToolCallAdd
  - TestHandleToolCallMultiply
  - TestHandleResourcesList
  - TestHandleMethodNotFound

### 3️⃣ **AWS Lambda 배포** ✅

#### 빌드 프로세스
```bash
# Makefile에 정의된 빌드 명령
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap main.go server.go
```

**핵심 요소:**
- `CGO_ENABLED=0` - 정적 링크 (GLIBC 호환성)
- `GOOS=linux GOARCH=amd64` - Linux 64비트
- `bootstrap` - Lambda가 인식하는 실행 파일명

#### AWS 리소스 (template.yaml)
```yaml
Resources:
  ✅ MCPServerFunction        # Lambda 함수
  ✅ MCPApi                   # API Gateway (REST)
  ✅ MCPServerLogGroup        # CloudWatch 로그
  ✅ MCPServerDashboard       # 모니터링 대시보드
  ✅ MCPServerErrorAlarm      # 에러 알람

Parameters:
  ✅ EnvironmentName          # dev/staging/prod
```

#### 배포 설정
```
Region: ap-northeast-2 (Seoul)
Stack: mcp-server-stack
Function: mcp-server-dev
S3: s3://aws-sam-cli-managed-default-samclisourcebucket-...
```

### 4️⃣ **Tool 구현** ✅

#### `add` Tool
```go
입력: {"a": 10, "b": 5}
출력: "10 + 5 = 15"
응답: {"result": {"content": [{"type": "text", "text": "10 + 5 = 15"}]}}
```

#### `multiply` Tool
```go
입력: {"a": 6, "b": 7}
출력: "6 × 7 = 42"
응답: {"result": {"content": [{"type": "text", "text": "6 × 7 = 42"}]}}
```

**MCP 프로토콜 호환성:**
- ✅ `content` 배열 형식
- ✅ `type: "text"` 타입 지정
- ✅ JSON-RPC 2.0 응답

### 5️⃣ **자동화 스크립트** ✅

#### `deploy.sh`
```bash
✅ .env 파일 검증
✅ Python venv 활성화
✅ Go 빌드 (CGO_ENABLED=0)
✅ SAM 배포
✅ 에러 처리 및 로깅
```

#### `Makefile`
```bash
make init    # 초기 설정
make build   # Go 컴파일
make deploy  # AWS 배포
make clean   # 산출물 정리
make test    # 단위 테스트
```

### 6️⃣ **배포 및 통합** ✅

#### API Gateway 설정
```
Method: POST
Resource: /mcp
Endpoint: https://w2z3biq5nb.execute-api.ap-northeast-2.amazonaws.com/dev/mcp
```

#### 인증 설정
```yaml
✅ DefaultAuthorizer: NONE  # 공개 접근 허용
```

#### CloudWatch 모니터링
```
✅ 로그그룹: /aws/lambda/mcp-server-dev
✅ 메트릭: Invocations, Duration, Errors
✅ 대시보드: mcp-server-dev
✅ 알람: 에러 발생 시 알림
```

### 7️⃣ **Claude/GitHub Copilot 통합** ✅

#### 문제 해결
| 문제 | 원인 | 해결책 |
|------|------|--------|
| GLIBC_2.34 에러 | 동적 링크 바이너리 | CGO_ENABLED=0 추가 |
| 403 Missing Auth Token | API Gateway 인증 활성화 | DefaultAuthorizer: NONE 설정 |
| TypeError: content not iterable | 응답 형식 오류 | content 배열 형식 적용 |

#### 통합 결과
```bash
✅ Claude Desktop에서 MCP 연결 성공
✅ GitHub Copilot에서 Tool 호출 가능
✅ "10 + 5 = ?"에 답변 능력 확인
```

---

## 📊 최종 파일 구조

```
mcp_go_lambda/
├── main.go                 # Lambda 핸들러 (108줄)
├── server.go               # MCP 서버 (227줄)
├── server_test.go          # 단위 테스트 (모두 PASS)
├── go.mod, go.sum
├── template.yaml           # SAM 배포 템플릿
├── samconfig.toml          # SAM 설정
├── Makefile                # 빌드 자동화
├── deploy.sh               # 배포 스크립트
├── requirements.txt        # Python 의존성
├── .env                    # AWS 자격증명 (gitignore)
├── .gitignore, .envrc
├── README.md               # 프로젝트 개요
├── events/                 # 테스트 이벤트 파일
├── .venv/                  # Python 가상환경 (gitignore)
├── bootstrap               # 빌드 산출물 (gitignore)
└── function.zip            # 배포 패키지 (gitignore)
```

---

## 🚀 현재 상태

### ✅ 완료된 기능
- MCP 서버 핵심 기능 5개 메서드
- 2개의 Tool 구현 및 테스트
- AWS Lambda 배포
- CloudWatch 모니터링
- Claude/Copilot 통합
- 자동 배포 파이프라인

### 📋 미구현 기능
- Resources 기능 (구조 정의, 실제 구현 필요)
- Prompts 기능 (구조 정의, 실제 구현 필요)
- 추가 Tool들 (DB 쿼리, API 호출 등)

### 📈 다음 단계 (선택사항)
- [ ] 데이터베이스 연동
- [ ] 외부 API 호출 도구 추가
- [ ] 이미지 처리 도구 추가
- [ ] 파일 시스템 접근 도구
- [ ] API 인증 추가 (OAuth, API Key)
- [ ] Rate limiting 구현
- [ ] 더 상세한 로깅

---

## 💡 주요 학습 내용

### Go 언어
- Lambda handler 작성
- JSON 파싱 및 생성
- 구조체 정의 및 메서드

### AWS
- Lambda 함수 개발
- API Gateway 설정
- CloudFormation/SAM 배포
- CloudWatch 모니터링

### DevOps
- 정적 바이너리 컴파일
- 자동 배포 스크립트
- 인프라 코드화 (IaC)

### 프로토콜
- JSON-RPC 2.0 이해
- MCP (Model Context Protocol) 구현
- Claude/Copilot 통합

---

## 📞 실제 사용 예시

### Tool 호출 (curl)
```bash
API_ENDPOINT="https://w2z3biq5nb.execute-api.ap-northeast-2.amazonaws.com/dev"

# Add Tool 호출
curl -X POST "$API_ENDPOINT/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "add",
      "arguments": {"a": 15, "b": 27}
    },
    "id": 1
  }'

# 응답: {"jsonrpc": "2.0", "result": {"content": [{"type": "text", "text": "15 + 27 = 42"}]}, "id": 1}
```

### Tool 호출 (Claude Desktop)
```
Claude에 질문: "10과 5를 더하면 몇이야?"

Claude가 MCP Tool 호출:
  → tools/call with name: "add", arguments: {"a": 10, "b": 5}
  → Lambda 응답: "10 + 5 = 15"

Claude 답변: "10과 5를 더하면 15입니다."
```

---

## 🎯 성공 지표

| 항목 | 목표 | 달성 |
|------|------|------|
| Go MCP 서버 구현 | ✅ | ✅ |
| AWS Lambda 배포 | ✅ | ✅ |
| CloudFormation IaC | ✅ | ✅ |
| 자동 배포 파이프라인 | ✅ | ✅ |
| 단위 테스트 | ✅ | ✅ 6/6 |
| Tool 구현 | ✅ | ✅ 2개 |
| Claude 통합 | ✅ | ✅ |
| 모니터링 | ✅ | ✅ |

---

## 📝 결론

이 프로젝트는 **Go로 작성한 MCP 서버를 AWS Lambda에 배포하여 Claude Desktop 및 GitHub Copilot과 통합하는 완전한 예제**입니다.

**주요 성과:**
- 프로덕션 레벨의 정적 바이너리 배포
- 자동화된 배포 파이프라인
- 모니터링 및 로깅 인프라
- Claude/Copilot과의 완벽한 통합
- 확장 가능한 아키텍처

**코드 품질:**
- 모든 테스트 통과 ✅
- 에러 처리 적절 ✅
- 문서화 완료 ✅
- 자동화 스크립트 준비 ✅

이제 필요한 Tool을 추가하여 자신의 사용 사례에 맞게 확장할 수 있습니다.
