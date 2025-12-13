# MCP Server on AWS Lambda (Go)

Model Context Protocol (MCP) 서버를 Go로 구현하여 AWS Lambda에 배포하는 프로젝트입니다.

## 프로젝트 구조

```
mcp_go_lambda/
├── main.go              # MCP 서버 메인 코드
├── go.mod              # Go 모듈 정의
├── template.yaml       # AWS SAM 템플릿
├── Makefile            # 빌드 및 배포 스크립트
├── events/             # 로컬 테스트 이벤트
│   └── api-gateway-event.json
└── README.md
```

## 사전 요구사항

- Go 1.21 이상
- AWS CLI
- AWS SAM CLI
- AWS 계정 및 IAM 권한

## 개발 및 배포 가이드

### 1. 의존성 설치

```bash
make deps
```

### 2. 로컬 테스트

```bash
# 로컬 Lambda 함수 테스트
make test-local

# 또는 직접 Go 코드 실행
make run
```

### 3. 빌드

```bash
make build
```

### 4. AWS Lambda에 배포

```bash
make deploy
```

첫 배포 시에는 `--guided` 옵션이 자동으로 적용되어 대화형으로 설정을 진행합니다.

## MCP 요청 예시

```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {},
  "id": 1
}
```

## 응답 예시

```json
{
  "jsonrpc": "2.0",
  "result": {
    "message": "MCP Server is running",
    "method": "initialize"
  },
  "id": 1
}
```

## 다음 단계

- [ ] MCP 리소스 구현
- [ ] MCP 도구 구현
- [ ] MCP 프롬프트 구현
- [ ] 데이터베이스 연동
- [ ] 에러 핸들링 개선
- [ ] 로깅 및 모니터링 추가
- [ ] 단위 테스트 작성
