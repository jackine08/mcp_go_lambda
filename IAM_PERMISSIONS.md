# AWS IAM 권한 요구사항

이 프로젝트를 배포하고 운영하는데 필요한 AWS IAM 권한 목록입니다.

## 최소 필수 권한 (Minimum Required Permissions)

### 1. AWS CDK 배포 권한

CDK는 CloudFormation을 사용하므로 다음 권한이 필요합니다:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "cloudformation:CreateStack",
        "cloudformation:UpdateStack",
        "cloudformation:DeleteStack",
        "cloudformation:DescribeStacks",
        "cloudformation:DescribeStackEvents",
        "cloudformation:DescribeStackResources",
        "cloudformation:GetTemplate",
        "cloudformation:ValidateTemplate",
        "cloudformation:CreateChangeSet",
        "cloudformation:DescribeChangeSet",
        "cloudformation:ExecuteChangeSet",
        "cloudformation:DeleteChangeSet",
        "cloudformation:ListStacks"
      ],
      "Resource": "*"
    }
  ]
}
```

### 2. Lambda 함수 관리

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "lambda:CreateFunction",
        "lambda:UpdateFunctionCode",
        "lambda:UpdateFunctionConfiguration",
        "lambda:DeleteFunction",
        "lambda:GetFunction",
        "lambda:GetFunctionConfiguration",
        "lambda:ListFunctions",
        "lambda:PublishVersion",
        "lambda:AddPermission",
        "lambda:RemovePermission",
        "lambda:InvokeFunction"
      ],
      "Resource": "arn:aws:lambda:*:*:function:mcp-server-stack-*"
    }
  ]
}
```

### 3. API Gateway 관리

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "apigateway:GET",
        "apigateway:POST",
        "apigateway:PUT",
        "apigateway:DELETE",
        "apigateway:PATCH"
      ],
      "Resource": "arn:aws:apigateway:*::/*"
    }
  ]
}
```

### 4. IAM 역할 및 정책 관리

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iam:CreateRole",
        "iam:DeleteRole",
        "iam:GetRole",
        "iam:PassRole",
        "iam:AttachRolePolicy",
        "iam:DetachRolePolicy",
        "iam:PutRolePolicy",
        "iam:DeleteRolePolicy",
        "iam:GetRolePolicy",
        "iam:TagRole",
        "iam:UntagRole"
      ],
      "Resource": [
        "arn:aws:iam::*:role/mcp-server-stack-*",
        "arn:aws:iam::*:role/cdk-*"
      ]
    }
  ]
}
```

### 5. S3 (CDK Assets 저장소)

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:CreateBucket",
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket",
        "s3:GetBucketLocation",
        "s3:GetBucketPolicy",
        "s3:PutBucketPolicy"
      ],
      "Resource": [
        "arn:aws:s3:::cdk-*",
        "arn:aws:s3:::cdk-*/*"
      ]
    }
  ]
}
```

### 6. CloudWatch Logs

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams",
        "logs:DeleteLogGroup"
      ],
      "Resource": "arn:aws:logs:*:*:log-group:/aws/lambda/mcp-server-stack-*"
    }
  ]
}
```

### 7. SSM Parameter Store (CDK Context)

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter",
        "ssm:PutParameter",
        "ssm:DeleteParameter"
      ],
      "Resource": "arn:aws:ssm:*:*:parameter/cdk-bootstrap/*"
    }
  ]
}
```

## 통합 권한 정책 (All-in-One Policy)

위의 모든 권한을 하나로 합친 정책:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "CloudFormationAccess",
      "Effect": "Allow",
      "Action": [
        "cloudformation:*"
      ],
      "Resource": "*"
    },
    {
      "Sid": "LambdaAccess",
      "Effect": "Allow",
      "Action": [
        "lambda:*"
      ],
      "Resource": "arn:aws:lambda:*:*:function:mcp-server-stack-*"
    },
    {
      "Sid": "APIGatewayAccess",
      "Effect": "Allow",
      "Action": [
        "apigateway:*"
      ],
      "Resource": "arn:aws:apigateway:*::/*"
    },
    {
      "Sid": "IAMRoleAccess",
      "Effect": "Allow",
      "Action": [
        "iam:CreateRole",
        "iam:DeleteRole",
        "iam:GetRole",
        "iam:PassRole",
        "iam:AttachRolePolicy",
        "iam:DetachRolePolicy",
        "iam:PutRolePolicy",
        "iam:DeleteRolePolicy",
        "iam:GetRolePolicy",
        "iam:TagRole",
        "iam:UntagRole"
      ],
      "Resource": [
        "arn:aws:iam::*:role/mcp-server-stack-*",
        "arn:aws:iam::*:role/cdk-*"
      ]
    },
    {
      "Sid": "S3CDKAssetsAccess",
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
      "Sid": "CloudWatchLogsAccess",
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams",
        "logs:DeleteLogGroup"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    },
    {
      "Sid": "SSMParameterAccess",
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter",
        "ssm:PutParameter",
        "ssm:DeleteParameter"
      ],
      "Resource": "arn:aws:ssm:*:*:parameter/cdk-bootstrap/*"
    },
    {
      "Sid": "STSAccess",
      "Effect": "Allow",
      "Action": [
        "sts:GetCallerIdentity"
      ],
      "Resource": "*"
    }
  ]
}
```

## Lambda 실행 역할 (Lambda Execution Role)

CDK가 자동으로 생성하지만, 참고용으로 Lambda 함수 자체에 필요한 권한:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:log-group:/aws/lambda/mcp-server-stack-*:*"
    }
  ]
}
```

이 역할은 CDK가 자동으로 생성하고 Lambda 함수에 연결합니다.

## 권한 설정 방법

### 방법 1: IAM 사용자에 정책 직접 연결

1. AWS Console → IAM → Users → 사용자 선택
2. "Add permissions" → "Attach policies directly"
3. "Create policy" 버튼 클릭
4. 위의 JSON 정책을 붙여넣기
5. 정책 이름: `MCPLambdaDeploymentPolicy`
6. 사용자에게 연결

### 방법 2: IAM 그룹 생성 후 사용자 추가

```bash
# AWS CLI로 그룹 생성
aws iam create-group --group-name MCPDeployers

# 정책 파일 생성 (policy.json)
# 위의 통합 권한 정책을 파일로 저장

# 그룹에 정책 연결
aws iam put-group-policy \
  --group-name MCPDeployers \
  --policy-name MCPDeploymentPolicy \
  --policy-document file://policy.json

# 사용자를 그룹에 추가
aws iam add-user-to-group \
  --group-name MCPDeployers \
  --user-name your-username
```

### 방법 3: Managed Policy 사용 (개발/테스트용)

간단한 개발/테스트 환경이라면 AWS Managed Policy 사용:

```bash
aws iam attach-user-policy \
  --user-name your-username \
  --policy-arn arn:aws:iam::aws:policy/AdministratorAccess
```

⚠️ **주의**: `AdministratorAccess`는 모든 권한을 부여하므로 프로덕션 환경에서는 권장하지 않습니다.

## CDK Bootstrap 권한

CDK를 처음 사용하는 경우 bootstrap이 필요합니다:

```bash
cdk bootstrap aws://ACCOUNT-ID/REGION
```

Bootstrap에 필요한 추가 권한:
- CloudFormation 스택 생성
- S3 버킷 생성 (cdk-bootstrap-*)
- IAM 역할 생성 (cdk-*)
- ECR 저장소 생성 (필요시)

## 권한 검증

배포 전에 권한을 검증하려면:

```bash
# 현재 사용자 확인
aws sts get-caller-identity

# CloudFormation 권한 확인
aws cloudformation describe-stacks --stack-name mcp-server-stack

# Lambda 권한 확인
aws lambda list-functions

# API Gateway 권한 확인
aws apigateway get-rest-apis
```

## 트러블슈팅

### AccessDenied 오류

```
User: arn:aws:iam::xxx:user/xxx is not authorized to perform: xxx
```

위 오류가 발생하면 해당 액션을 통합 권한 정책에 추가하세요.

### CDK Bootstrap 실패

```bash
# Bootstrap 재실행
cdk bootstrap aws://$(aws sts get-caller-identity --query Account --output text)/ap-northeast-2
```

### IAM PassRole 오류

Lambda에 역할을 연결할 때 `iam:PassRole` 권한이 필요합니다. 통합 정책에 포함되어 있습니다.

## 보안 권장사항

1. **최소 권한 원칙**: 프로덕션에서는 통합 정책 대신 필요한 권한만 부여
2. **MFA 활성화**: IAM 사용자에 MFA 설정 권장
3. **액세스 키 로테이션**: 정기적으로 액세스 키 갱신
4. **CloudTrail 활성화**: 모든 API 호출 로깅
5. **리소스 태그 활용**: Cost Allocation과 권한 관리에 활용

## 참고 자료

- [AWS CDK Permissions](https://docs.aws.amazon.com/cdk/v2/guide/permissions.html)
- [Lambda Execution Role](https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html)
- [API Gateway Permissions](https://docs.aws.amazon.com/apigateway/latest/developerguide/permissions.html)
