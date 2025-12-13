# AWS 실제 배포 단계별 가이드

## 1단계: AWS 자격증명 설정

### 옵션 A: AWS CLI configure 명령어 사용 (가장 간단)

```bash
aws configure
```

다음 정보를 입력합니다:
- **AWS Access Key ID**: [IAM에서 발급받은 Access Key]
- **AWS Secret Access Key**: [IAM에서 발급받은 Secret Key]
- **Default region name**: ap-northeast-2 (서울)
- **Default output format**: json

### 옵션 B: 환경변수 사용

```bash
export AWS_ACCESS_KEY_ID=your_access_key_id
export AWS_SECRET_ACCESS_KEY=your_secret_access_key
export AWS_DEFAULT_REGION=ap-northeast-2
```

### 옵션 C: AWS SSO 사용 (조직 내 권장)

```bash
aws sso login --profile your-profile
```

## 2단계: AWS IAM 권한 확인

배포하는 IAM 사용자가 다음 권한을 가져야 합니다:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "lambda:*",
        "apigateway:*",
        "cloudformation:*",
        "iam:PassRole",
        "logs:*",
        "cloudwatch:*",
        "s3:*"
      ],
      "Resource": "*"
    }
  ]
}
```

또는 기본 권한으로 충분합니다:
- `AdministratorAccess` (개발 환경)
- `PowerUserAccess` + `IAMPassRolePolicy` (권장)

## 3단계: S3 버킷 생성 (SAM 배포용)

SAM이 배포 아티팩트를 저장할 S3 버킷이 필요합니다:

```bash
# 버킷 생성
aws s3 mb s3://mcp-server-deployment-artifacts-$(date +%s) --region ap-northeast-2

# 또는 직접 선택한 이름으로 생성
aws s3 mb s3://YOUR_BUCKET_NAME --region ap-northeast-2
```

## 4단계: 실제 배포 실행

```bash
# 터미널에서 프로젝트 디렉토리로 이동
cd /home/syoh/workspace/mcp_go_lambda

# 가상환경 활성화
source .venv/bin/activate

# 빌드
make build

# 배포 (첫 배포 시)
sam deploy --guided

# 배포 설정 저장 후 다시 배포할 때는
sam deploy
```

### sam deploy --guided 시 입력값 예시:

```
Stack Name [sam-app]: mcp-server-stack
Region [ap-northeast-2]: ap-northeast-2
Confirm changes before deploy [y/N]: y
Allow SAM CLI IAM role creation [Y/n]: Y
Save parameters to samconfig.toml [Y/n]: Y
```

## 5단계: 배포 결과 확인

```bash
# 배포된 스택 정보 확인
aws cloudformation describe-stacks --stack-name mcp-server-stack --region ap-northeast-2

# Lambda 함수 확인
aws lambda list-functions --region ap-northeast-2

# API Gateway 엔드포인트 확인
aws cloudformation describe-stacks \
  --stack-name mcp-server-stack \
  --query 'Stacks[0].Outputs' \
  --region ap-northeast-2
```

## 6단계: API 테스트

배포 후 받은 API 엔드포인트로 테스트:

```bash
# API 엔드포인트 저장
API_ENDPOINT="https://XXXXX.execute-api.ap-northeast-2.amazonaws.com/dev"

# 헬스체크
curl -X GET "$API_ENDPOINT/"

# MCP initialize 요청
curl -X POST "$API_ENDPOINT/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {"name": "test", "version": "1.0"}
    },
    "id": 1
  }'

# tools/list 요청
curl -X POST "$API_ENDPOINT/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "params": {},
    "id": 2
  }'
```

## 7단계: 로그 확인

```bash
# CloudWatch 로그 실시간 확인
aws logs tail /aws/lambda/mcp-server-dev --follow --region ap-northeast-2

# 또는 로그 그룹 나열
aws logs describe-log-groups --region ap-northeast-2
```

## 8단계: 비용 최적화 (선택사항)

```bash
# CloudWatch 대시보드에서 메트릭 모니터링
# Lambda 메모리 크기 최적화
# API Gateway 캐싱 설정 (필요시)

# 테스트 후 비용 절감을 위해 dev 환경 스택 삭제
aws cloudformation delete-stack --stack-name mcp-server-stack --region ap-northeast-2
```

## 문제 해결

### 배포 실패 시

```bash
# SAM 빌드 로그 확인
sam build --debug

# CloudFormation 이벤트 확인
aws cloudformation describe-stack-events --stack-name mcp-server-stack --region ap-northeast-2

# Lambda 함수 로그 확인
aws logs tail /aws/lambda/mcp-server-dev --follow
```

### IAM 권한 부족

```bash
# 사용자의 현재 권한 확인
aws iam list-user-policies --user-name YOUR_USERNAME --region ap-northeast-2

# 정책 추가 필요시
aws iam attach-user-policy \
  --user-name YOUR_USERNAME \
  --policy-arn arn:aws:iam::aws:policy/AdministratorAccess
```

## 배포 후 관리

배포 후 필요한 관리 작업:
1. CloudWatch 알람 설정 (이미 template.yaml에 포함됨)
2. 로그 보관 기간 설정 (현재 14일)
3. 자동 스케일링 정책 (필요시)
4. 비용 모니터링

자세한 내용은 AWS 공식 문서 참고:
- https://docs.aws.amazon.com/lambda/latest/dg/
- https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/
