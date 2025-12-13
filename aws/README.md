# ☁️ AWS Infrastructure

AWS Lambda, API Gateway, CloudFormation 관련 파일들이 저장되는 디렉토리입니다.

## 📁 파일 구조

```
aws/
├── template.yaml        # CloudFormation/SAM 템플릿
└── samconfig.toml       # SAM 배포 설정
```

## 📄 파일 설명

### template.yaml
**AWS SAM (Serverless Application Model) CloudFormation 템플릿**

#### 주요 내용
- **Lambda 함수 정의** (MCPServerFunction)
  - Runtime: go1.x
  - Handler: bootstrap
  - Environment: 환경변수 설정
  - CloudWatch 로그 활성화

- **API Gateway 설정** (MCPApi)
  - Type: REST API
  - Resource: /mcp
  - Method: POST
  - 인증: 비활성화 (public access)

- **CloudWatch 리소스**
  - LogGroup: /aws/lambda/mcp-server-dev
  - Dashboard: 모니터링 대시보드
  - Alarm: 에러 감지 알람

- **Outputs**
  - MCPApiEndpoint: API Gateway URL
  - FunctionArn: Lambda 함수 ARN
  - LogGroupName: CloudWatch 로그그룹명

#### 핵심 설정
```yaml
Auth:
  DefaultAuthorizer: NONE        # 인증 비활성화 (Claude Desktop 접근 가능)

Environment:
  Variables:
    ENVIRONMENT: dev             # 환경 구분
```

#### 수정 시 주의사항
- `DefaultAuthorizer` 변경 시 API 인증 활성화됨
- `Runtime`은 반드시 `go1.x`로 유지 (Lambda AL2 runtime)
- `Handler`는 반드시 `bootstrap`으로 유지 (Go 바이너리명)

---

### samconfig.toml
**SAM CLI 배포 설정 파일**

#### 주요 내용
```toml
region = "ap-northeast-2"           # AWS 리전 (Seoul)
s3_prefix = "aws-sam-cli-managed-default-samclisourcebucket-*"
confirm_changeset = false           # Changeset 자동 승인
capabilities = "CAPABILITY_IAM"     # IAM 리소스 자동 생성
```

#### 역할
- 배포 시 `sam deploy --guided` 프롬프트 없음
- 기존 설정을 기억하여 다시 입력할 필요 없음
- `deploy/deploy.sh`에서 자동으로 적용됨

#### 수정 시 주의사항
- `region` 변경 시 AWS 리소스가 다른 리전에 생성됨
- S3 버킷명은 자동 관리 (수정 불필요)
- `confirm_changeset = false`는 자동 배포를 위해 유지

---

## 🚀 배포 프로세스

### 1️⃣ template.yaml 검증
```bash
sam validate --template aws/template.yaml
```

### 2️⃣ 로컬 빌드
```bash
cd deploy
make build
```

### 3️⃣ SAM 배포
```bash
sam deploy --template aws/template.yaml --config-file aws/samconfig.toml
```

또는 자동화된 방법:
```bash
cd deploy
./deploy.sh
```

---

## 📊 생성되는 AWS 리소스

| 리소스 | 이름 | 설명 |
|--------|------|------|
| Lambda | mcp-server-dev | MCP 서버 함수 |
| API Gateway | MCPApi | REST API 엔드포인트 |
| CloudWatch Logs | /aws/lambda/mcp-server-dev | 함수 실행 로그 |
| CloudWatch Dashboard | mcp-server-dev | 모니터링 대시보드 |
| CloudWatch Alarm | mcp-server-dev-errors | 에러 발생 시 알람 |
| CloudFormation Stack | mcp-server-stack | 전체 인프라 관리 |

---

## 🔗 API Gateway 정보

### Endpoint
```
https://w2z3biq5nb.execute-api.ap-northeast-2.amazonaws.com/dev/mcp
```

### Method
```
POST /mcp
```

### Request Format (JSON-RPC 2.0)
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "add",
    "arguments": {"a": 10, "b": 5}
  },
  "id": 1
}
```

### Response Format
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

---

## 🛠️ CloudFormation Stack 관리

### Stack 상태 확인
```bash
aws cloudformation describe-stacks \
  --stack-name mcp-server-stack \
  --region ap-northeast-2
```

### Stack 삭제
```bash
aws cloudformation delete-stack \
  --stack-name mcp-server-stack \
  --region ap-northeast-2
```

### 변경사항 미리보기
```bash
sam deploy \
  --no-confirm-changeset \
  --no-execute-changeset \
  --region ap-northeast-2
```

---

## 🔍 문제 해결

### "Missing Authentication Token" 에러
**원인:** API Gateway에 인증이 활성화되어 있음

**해결:**
```yaml
# template.yaml에서 확인
Auth:
  DefaultAuthorizer: NONE
```

### "GLIBC" 호환성 에러
**원인:** Go 바이너리가 동적 링크됨

**해결:** deploy/Makefile에서
```makefile
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap
```

---

## 📚 참고 자료

- [AWS SAM 문서](https://docs.aws.amazon.com/serverless-application-model/)
- [CloudFormation 템플릿 작성](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/)
- [Lambda Go 런타임](https://docs.aws.amazon.com/lambda/latest/dg/lambda-go-how-to-create-deployment-package.html)
- [API Gateway REST API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api.html)

---

**마지막 업데이트:** 2025년 12월
