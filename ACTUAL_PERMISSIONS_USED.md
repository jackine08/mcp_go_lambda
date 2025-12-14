# 실제 배포 시 사용된 AWS 권한

이 문서는 실제 CDK 배포 과정에서 사용된 AWS API 호출과 권한들을 기록합니다.

## 배포 과정 분석

### 1. CDK Bootstrap 단계

CDK Bootstrap은 다음 리소스들을 생성합니다:

```
cdk-hnb659fds-deploy-role-374146078851-ap-northeast-2
cdk-hnb659fds-file-publishing-role-374146078851-ap-northeast-2
```

**사용된 권한:**
- `iam:CreateRole` - CDK 실행 역할 생성
- `iam:AttachRolePolicy` - 정책 연결
- `s3:CreateBucket` - CDK assets 버킷 생성 (cdk-hnb659fds-assets-*)
- `ssm:PutParameter` - Bootstrap 정보 저장

### 2. CDK Synthesis (템플릿 생성)

로컬에서 CloudFormation 템플릿을 생성하는 단계입니다.

**필요 권한:** 없음 (로컬 작업)

### 3. Asset Publishing (코드 업로드)

```
Publishing MCPServerFunction/Code (374146078851-ap-northeast-2-5adec38b)
Publishing mcp-server-stack Template (374146078851-ap-northeast-2-999acf45)
```

**사용된 권한:**
- `s3:PutObject` - Lambda 코드를 S3에 업로드
- `s3:GetObject` - 업로드 검증
- `s3:GetBucketLocation` - 버킷 위치 확인

**S3 버킷:**
- `cdk-hnb659fds-assets-374146078851-ap-northeast-2`

### 4. CloudFormation Stack 배포

```
mcp-server-stack: creating CloudFormation changeset...
```

**사용된 권한:**
- `cloudformation:CreateChangeSet` - 변경 사항 미리보기
- `cloudformation:DescribeChangeSet` - 변경 사항 확인
- `cloudformation:ExecuteChangeSet` - 변경 사항 적용
- `cloudformation:DescribeStacks` - 스택 상태 확인
- `cloudformation:DescribeStackEvents` - 배포 진행상황 모니터링

### 5. Lambda 함수 생성/업데이트

```
Created: mcp-server-stack-MCPServerFunctionA0453C57-jbbBGabJt4gH
```

**CloudFormation이 수행한 작업:**
- `lambda:CreateFunction` - Lambda 함수 생성
- `lambda:UpdateFunctionCode` - 코드 업데이트 (재배포 시)
- `lambda:GetFunction` - 함수 정보 조회
- `lambda:AddPermission` - API Gateway 호출 권한 추가

**Lambda 실행 역할 생성:**
- `iam:CreateRole` - Lambda 실행 역할 생성
- `iam:AttachRolePolicy` - AWSLambdaBasicExecutionRole 연결
- `iam:PassRole` - Lambda에 역할 전달

### 6. API Gateway 생성

```
Created: API Gateway REST API
Endpoint: yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com
```

**CloudFormation이 수행한 작업:**
- `apigateway:POST /restapis` - REST API 생성
- `apigateway:POST /restapis/{id}/resources` - 리소스 생성 (/mcp)
- `apigateway:PUT /restapis/{id}/resources/{id}/methods/{method}` - POST 메서드 생성
- `apigateway:PUT /restapis/{id}/resources/{id}/integrations` - Lambda 통합 설정
- `apigateway:POST /restapis/{id}/deployments` - API 배포
- `apigateway:POST /restapis/{id}/stages` - Stage 생성 (dev)

### 7. CloudWatch Logs 설정

**자동 생성:**
- `/aws/lambda/mcp-server-stack-MCPServerFunctionA0453C57-jbbBGabJt4gH`

**사용된 권한:**
- `logs:CreateLogGroup` - 로그 그룹 생성 (Lambda가 자동 생성)
- `logs:CreateLogStream` - 로그 스트림 생성
- `logs:PutLogEvents` - 로그 이벤트 기록

## 실제 배포 결과

### 생성된 리소스

1. **Lambda 함수**
   - 이름: `mcp-server-stack-MCPServerFunctionA0453C57-jbbBGabJt4gH`
   - ARN: `arn:aws:lambda:ap-northeast-2:374146078851:function:mcp-server-stack-MCPServerFunctionA0453C57-jbbBGabJt4gH`
   - 런타임: `provided.al2023` (Go custom runtime)
   - 메모리: 128 MB (기본값)
   - 타임아웃: 30초 (기본값)

2. **API Gateway**
   - REST API ID: `yhj0gi980j`
   - 엔드포인트: `https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp`
   - Stage: `dev`

3. **IAM 역할**
   - Lambda 실행 역할: `mcp-server-stack-MCPServerFunctionServiceRole*`
   - 정책: `AWSLambdaBasicExecutionRole` (CloudWatch Logs 쓰기 권한)

4. **CloudWatch Log Group**
   - `/aws/lambda/mcp-server-stack-MCPServerFunctionA0453C57-jbbBGabJt4gH`

5. **S3 Bucket** (CDK Assets)
   - `cdk-hnb659fds-assets-374146078851-ap-northeast-2`

## 배포 과정에서 발견된 권한 메시지

```
current credentials could not be used to assume 'arn:aws:iam::374146078851:role/cdk-hnb659fds-deploy-role-374146078851-ap-northeast-2', but are for the right account. Proceeding anyway.
```

이 메시지는 다음을 의미합니다:
- 현재 credentials가 CDK에서 생성한 역할을 assume할 수 없음
- 하지만 같은 계정이므로 직접 권한으로 진행
- **결론**: 현재 사용자가 직접 모든 권한을 가지고 있어야 함 (역할 위임 없이)

## 권한 요약

### 필수 직접 권한 (역할 위임 없이 배포하는 경우)

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "cloudformation:*"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "arn:aws:s3:::cdk-*",
        "arn:aws:s3:::cdk-*/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "lambda:CreateFunction",
        "lambda:UpdateFunctionCode",
        "lambda:UpdateFunctionConfiguration",
        "lambda:GetFunction",
        "lambda:AddPermission",
        "lambda:RemovePermission",
        "lambda:DeleteFunction",
        "lambda:TagResource"
      ],
      "Resource": "arn:aws:lambda:*:*:function:mcp-server-stack-*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "apigateway:*"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "iam:CreateRole",
        "iam:GetRole",
        "iam:PassRole",
        "iam:AttachRolePolicy",
        "iam:DetachRolePolicy",
        "iam:DeleteRole",
        "iam:PutRolePolicy",
        "iam:DeleteRolePolicy",
        "iam:GetRolePolicy",
        "iam:TagRole"
      ],
      "Resource": [
        "arn:aws:iam::*:role/mcp-server-stack-*",
        "arn:aws:iam::*:role/cdk-*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:DescribeLogGroups",
        "logs:DeleteLogGroup",
        "logs:PutRetentionPolicy"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter",
        "ssm:PutParameter"
      ],
      "Resource": "arn:aws:ssm:*:*:parameter/cdk-bootstrap/*"
    }
  ]
}
```

### Lambda 실행 시 필요한 권한 (자동 생성됨)

Lambda 함수 자체가 실행될 때 필요한 권한 (IAM 역할에 자동 연결):

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:log-group:/aws/lambda/mcp-server-stack-*:*"
    }
  ]
}
```

이 권한은 CDK가 자동으로 생성하고 Lambda 실행 역할에 연결합니다.

## 테스트 결과

### 배포 성공
```bash
✅ mcp-server-stack
✨ Deployment time: 32.85s
```

### API 테스트 성공
```bash
# tools/list 호출 - 7개 tool 반환
curl -X POST https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp \
  -d '{"jsonrpc": "2.0", "method": "tools/list", "params": {}, "id": 1}'

# tools/call 호출 - add tool 실행
curl -X POST https://yhj0gi980j.execute-api.ap-northeast-2.amazonaws.com/dev/mcp \
  -d '{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "add", "arguments": {"a": 42, "b": 58}}, "id": 2}'
```

모두 정상 작동!

## 권장 사항

1. **CloudTrail 활성화**: API 호출 추적을 위해 CloudTrail 활성화 권장
2. **최소 권한**: 프로덕션에서는 위의 권한에서 불필요한 부분 제거
3. **리소스 태그**: Cost Allocation을 위해 모든 리소스에 태그 추가
4. **모니터링**: CloudWatch Alarms 설정 (Lambda 에러율, API Gateway 5xx 등)

## 추가 로깅 설정

더 자세한 로그를 확인하려면:

```bash
# CloudWatch Logs Insights 쿼리
aws logs start-query \
  --log-group-name /aws/lambda/mcp-server-stack-MCPServerFunctionA0453C57-jbbBGabJt4gH \
  --start-time $(date -d '1 hour ago' +%s) \
  --end-time $(date +%s) \
  --query-string 'fields @timestamp, @message | sort @timestamp desc | limit 20'
```

또는 AWS Console에서:
1. CloudWatch → Log groups
2. `/aws/lambda/mcp-server-stack-MCPServerFunctionA0453C57-jbbBGabJt4gH` 선택
3. 최근 로그 스트림 확인
