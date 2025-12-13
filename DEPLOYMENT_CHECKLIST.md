# AWS 배포 체크리스트

## ✅ 배포 전 필수 확인 사항

- [ ] **AWS 계정 생성**
  - AWS 회원가입: https://aws.amazon.com/
  - 결제 정보 등록

- [ ] **IAM 사용자 생성 및 액세스 키 발급**
  - AWS Management Console → IAM → Users → Create user
  - 권한 할당: AdministratorAccess 또는 PowerUserAccess + IAMPassRole
  - 액세스 키 생성 및 저장

- [ ] **AWS CLI 자격증명 설정**
  ```bash
  aws configure
  # 또는
  export AWS_ACCESS_KEY_ID=your_key
  export AWS_SECRET_ACCESS_KEY=your_secret
  export AWS_DEFAULT_REGION=ap-northeast-2
  ```

- [ ] **자격증명 확인**
  ```bash
  aws sts get-caller-identity
  ```

## 🚀 배포 단계

### 1단계: 로컬 빌드 테스트
```bash
cd /home/syoh/workspace/mcp_go_lambda
go test -v                    # 테스트 실행
go build -o bootstrap main.go server.go  # 빌드 확인
```

### 2단계: AWS 배포 실행
```bash
# 가상환경 활성화
source .venv/bin/activate

# 방법 1: deploy.sh 스크립트 사용 (권장)
./deploy.sh

# 방법 2: Makefile 사용
make deploy

# 방법 3: 직접 SAM 사용
sam deploy --guided
```

### 3단계: 배포 결과 확인
```bash
# 스택 정보 확인
aws cloudformation describe-stacks \
  --stack-name mcp-server-stack \
  --region ap-northeast-2

# API 엔드포인트 출력
aws cloudformation describe-stacks \
  --stack-name mcp-server-stack \
  --query 'Stacks[0].Outputs' \
  --region ap-northeast-2 \
  --output table
```

### 4단계: API 테스트
```bash
# API 엔드포인트를 변수로 저장 (위의 단계 3에서 얻은 값 사용)
export API_ENDPOINT="https://XXXXX.execute-api.ap-northeast-2.amazonaws.com/dev"

# 헬스체크
curl -X GET "$API_ENDPOINT/"

# Initialize 메서드 테스트
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

# Tools list 메서드 테스트
curl -X POST "$API_ENDPOINT/mcp" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "params": {},
    "id": 2
  }'
```

### 5단계: 로그 확인
```bash
# 실시간 로그 확인
aws logs tail /aws/lambda/mcp-server-dev --follow --region ap-northeast-2

# 최근 N분 로그 확인
aws logs tail /aws/lambda/mcp-server-dev --since 5m --region ap-northeast-2

# 특정 기간의 로그 조회
aws logs filter-log-events \
  --log-group-name /aws/lambda/mcp-server-dev \
  --region ap-northeast-2
```

## 🔧 배포 후 운영

### 모니터링
```bash
# CloudWatch 메트릭 확인
aws cloudwatch get-metric-statistics \
  --namespace AWS/Lambda \
  --metric-name Invocations \
  --dimensions Name=FunctionName,Value=mcp-server-dev \
  --start-time 2025-12-13T00:00:00Z \
  --end-time 2025-12-13T23:59:59Z \
  --period 3600 \
  --statistics Sum,Average \
  --region ap-northeast-2
```

### 함수 업데이트
```bash
# 코드 변경 후 재배포
make clean build deploy

# 또는 SAM으로 직접
sam build && sam deploy
```

### 리소스 정리
```bash
# 테스트 후 비용 절감을 위해 스택 삭제
aws cloudformation delete-stack --stack-name mcp-server-stack --region ap-northeast-2

# 삭제 확인
aws cloudformation describe-stacks --stack-name mcp-server-stack --region ap-northeast-2
```

## 🚨 문제 해결

### AWS CLI 자격증명 오류
```
Unable to locate credentials
```
→ `aws configure`를 실행하거나 환경변수 설정

### IAM 권한 부족
```
User is not authorized to perform: iam:PassRole
```
→ IAM 사용자에게 "iam:PassRole" 권한 추가 필요

### SAM 빌드 실패
```bash
sam build --debug  # 상세 로그 확인
```

### Lambda 함수 에러
```bash
# CloudWatch 로그 확인
aws logs tail /aws/lambda/mcp-server-dev --follow

# 또는 Lambda 콘솔에서 "Logs" 탭 확인
```

## 📊 비용 추정

AWS 프리 티어 (월간):
- Lambda: 100만 건의 요청 무료
- API Gateway: 100만 건의 호출 무료
- CloudWatch: 5GB 로그 저장 무료

초과 시 비용:
- Lambda: 백만 건당 $0.20
- API Gateway: 백만 건당 $3.50
- CloudWatch: GB당 $0.50

## 📚 유용한 링크

- [AWS Lambda 문서](https://docs.aws.amazon.com/lambda/)
- [AWS SAM 문서](https://docs.aws.amazon.com/serverless-application-model/)
- [AWS CLI 참고](https://docs.aws.amazon.com/cli/)
- [Model Context Protocol](https://modelcontextprotocol.io/)
