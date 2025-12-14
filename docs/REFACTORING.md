# 프로젝트 리팩토링 가이드

## 변경 사항

### 구조 변경
기존의 flat한 구조에서 Go 표준 레이아웃으로 변경했습니다.

#### Before (Flat Structure)
```
mcp-go-lambda/
├── main.go          # Lambda handler + 엔트리포인트
├── server.go        # MCP 서버 + 타입 정의
└── server_test.go   # 테스트
```

#### After (Standard Go Layout)
```
mcp-go-lambda/
├── cmd/lambda/main.go              # 엔트리포인트만
├── internal/
│   ├── types/types.go              # 타입 정의 분리
│   ├── mcp/server.go               # MCP 비즈니스 로직
│   └── handler/lambda.go           # Lambda 핸들러 분리
└── old_structure/                  # 백업
```

## 장점

### 1. 관심사의 분리 (Separation of Concerns)
- **types**: 순수한 데이터 타입만 정의
- **mcp**: MCP 프로토콜 비즈니스 로직만
- **handler**: Lambda 통합 로직만
- **cmd**: 앱 시작 로직만

### 2. 테스트 용이성
각 패키지를 독립적으로 테스트 가능:
```bash
go test ./internal/types      # 타입 테스트
go test ./internal/mcp        # MCP 로직 테스트
go test ./internal/handler    # 핸들러 테스트
```

### 3. 재사용성
- `internal/mcp`를 다른 엔트리포인트에서도 사용 가능
- 예: CLI 도구, HTTP 서버, Lambda 등 다양한 런타임에서 같은 MCP 로직 사용

### 4. Go 표준 규칙 준수
- `cmd/`: 애플리케이션 엔트리포인트
- `internal/`: 외부에서 import 불가능한 내부 패키지
- 패키지별로 명확한 역할 분담

## 마이그레이션 가이드

### 의존성 방향
```
cmd/lambda/main.go
    ↓ import
internal/handler/lambda.go
    ↓ import
internal/mcp/server.go
    ↓ import
internal/types/types.go (최하위, 의존성 없음)
```

### Import 경로
모든 import는 모듈명을 기준으로:
```go
import (
    "github.com/jackine08/mcp_go_lambda/internal/types"
    "github.com/jackine08/mcp_go_lambda/internal/mcp"
    "github.com/jackine08/mcp_go_lambda/internal/handler"
)
```

### 빌드 명령어 변경
**Before:**
```bash
go build -o lambda_dist/bootstrap main.go server.go
```

**After:**
```bash
go build -o lambda_dist/bootstrap cmd/lambda/main.go
# 또는
make build
```

## 추가 개선 가능 항목

### 1. 설정 관리
```
internal/
└── config/
    └── config.go    # 환경변수, 설정 로딩
```

### 2. 로깅
```
internal/
└── logging/
    └── logger.go    # 구조화된 로깅
```

### 3. OAuth 기능 (별도 브랜치)
```
internal/
└── oauth/
    ├── dcr.go       # Dynamic Client Registration
    ├── metadata.go  # RFC 8414 endpoints
    └── pkce.go      # PKCE validation
```

### 4. 도구 플러그인 시스템
```
internal/
└── tools/
    ├── registry.go  # Tool 등록/관리
    ├── add.go       # Add tool
    ├── multiply.go  # Multiply tool
    └── custom.go    # Custom tools
```

## 백업된 파일

원본 파일은 `old_structure/` 디렉토리에 백업되어 있습니다:
- `old_structure/main.go`
- `old_structure/server.go`
- `old_structure/server_test.go`

필요시 해당 파일들을 참고할 수 있습니다.

## 배포 변경 사항

CDK 배포는 변경사항 없음:
- `lambda_dist/bootstrap` 파일을 사용하므로 빌드 결과물 경로 동일
- API Gateway 엔드포인트 동일
- Lambda 설정 동일

## 테스트 결과

리팩토링 후 모든 기능 정상 작동 확인:
- ✅ Initialize: 서버 초기화
- ✅ Tools/list: 도구 목록 조회
- ✅ Tools/call: Add 도구 (15 + 25 = 40)
- ✅ API Gateway 통합
- ✅ Lambda 실행
