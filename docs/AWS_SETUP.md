# SAM 배포 설정 예시
# AWS 자격증명 설정 후 'make deploy' 명령으로 배포할 수 있습니다

## AWS 자격증명 설정 방법:

### 1. AWS CLI 자격증명 설정 (권장)
```bash
aws configure
# AWS Access Key ID: [입력]
# AWS Secret Access Key: [입력]
# Default region: ap-northeast-2 (또는 원하는 리전)
# Default output format: json
```

### 2. 또는 환경변수로 설정
```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=ap-northeast-2
```

### 3. AWS SSO로 설정
```bash
aws sso login --profile your-profile
```

## 배포 명령어
```bash
make deploy
```

## 배포 후 확인
```bash
# Lambda 함수 목록 조회
aws lambda list-functions --region ap-northeast-2

# 특정 함수의 설정 확인
aws lambda get-function --function-name MCPServerFunction --region ap-northeast-2

# API Gateway 엔드포인트 확인
aws apigateway get-rest-apis --region ap-northeast-2
```

## samconfig.toml
처음 배포 시 'sam deploy --guided'를 실행하면 설정이 자동으로 저장됩니다.
