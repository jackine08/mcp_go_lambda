# 🏗️ AWS CDK (Cloud Development Kit)

AWS CDK로 인프라를 관리합니다. Python으로 서버리스 리소스를 정의하고 배포합니다.

## 📁 파일 구조

```
cdk/
├── app.py                      # CDK 애플리케이션 진입점
├── stacks/
│   ├── __init__.py
│   └── mcp_lambda_stack.py    # Lambda, API Gateway, CloudWatch 스택
├── cdk.json                    # CDK 설정
├── requirements.txt            # Python 의존성
└── venv/                       # Python 가상환경 (생성됨)
```

## 🚀 빠른 시작

### 1️⃣ 초기 설정

```bash
cd cdk

# Python 가상환경 생성
python3 -m venv venv
source venv/bin/activate

# 의존성 설치
pip install -r requirements.txt
```

### 2️⃣ CDK 스택 확인

```bash
# CloudFormation 템플릿 생성 (검증만)
cdk synth

# 배포할 변경사항 미리보기
cdk diff
```

### 3️⃣ CDK 배포

```bash
# 자동으로 Go 빌드 + CDK 배포
cd ..
./deploy/cdk-deploy.sh

# 또는 수동 배포
cd cdk
cdk deploy --require-approval=never
```

### 4️⃣ 배포 확인

```bash
# 배포된 스택 확인
aws cloudformation describe-stacks --stack-name mcp-server-stack

# CloudFormation 이벤트 확인
aws cloudformation describe-stack-events --stack-name mcp-server-stack

# Lambda 함수 확인
aws lambda list-functions --region ap-northeast-2

# API Gateway 확인
aws apigateway get-rest-apis --region ap-northeast-2
```

## 📝 app.py 설명

```python
from stacks.mcp_lambda_stack import MCPLambdaStack

app = cdk.App()

MCPLambdaStack(
    app,
    "mcp-server-stack",
    env=cdk.Environment(
        account="YOUR_ACCOUNT_ID",  # AWS 계정 ID
        region="ap-northeast-2",
    ),
)

app.synth()
```

**핵심 요소:**
- `app`: CDK 애플리케이션 인스턴스
- `MCPLambdaStack`: 정의된 스택 클래스
- `env`: AWS 계정 ID와 리전 지정

## 🔧 stacks/mcp_lambda_stack.py 설명

### 생성되는 리소스

#### 1. Lambda 함수
```python
lambda_fn = lambda_.Function(
    self,
    "MCPServerFunction",
    runtime=lambda_.Runtime.GO_1_X,
    handler="bootstrap",
    code=lambda_.Code.from_asset(".."),
    timeout=Duration.seconds(30),
    memory_size=256,
)
```

#### 2. CloudWatch 로그 그룹
```python
log_group = logs.LogGroup(
    self,
    "MCPServerLogGroup",
    log_group_name="/aws/lambda/mcp-server-dev",
    retention=logs.RetentionDays.ONE_WEEK,
)
```

#### 3. API Gateway (REST API)
```python
api = apigw.RestApi(
    self,
    "MCPApi",
    rest_api_name="mcp-api",
)

mcp_resource = api.root.add_resource("mcp")
mcp_resource.add_method("POST", apigw.LambdaIntegration(lambda_fn))
```

#### 4. CloudWatch 대시보드
```python
dashboard = cloudwatch.Dashboard(
    self,
    "MCPServerDashboard",
    dashboard_name="mcp-server-dev",
)

dashboard.add_widgets(
    cloudwatch.GraphWidget(...),  # Lambda Invocations
    cloudwatch.GraphWidget(...),  # Lambda Duration
    cloudwatch.GraphWidget(...),  # Lambda Errors
)
```

#### 5. CloudWatch 알람
```python
error_alarm = cloudwatch.Alarm(
    self,
    "MCPServerErrorAlarm",
    metric=lambda_.Function.metric_errors(statistic="Sum"),
    threshold=1,
    evaluation_periods=1,
)
```

#### 6. Outputs
```python
core.CfnOutput(
    self,
    "APIEndpoint",
    value=api.url_for_path("/mcp"),
    description="MCP API Endpoint",
)
```

## 💻 주요 CDK 명령어

```bash
# 스택 합성 (CloudFormation 템플릿 생성)
cdk synth

# 배포할 변경사항 미리보기
cdk diff

# AWS에 배포
cdk deploy

# 대화형 배포 (확인 필요)
cdk deploy --require-approval=always

# CloudFormation 템플릿 출력
cdk synth -q

# 스택 삭제 (AWS 리소스 정리)
cdk destroy

# CDK 설명 보기
cdk list
```

## 🔄 SAM vs CDK 비교

| 항목 | SAM | CDK |
|------|-----|-----|
| 문법 | YAML | Python 코드 |
| 학습곡선 | 낮음 | 중간 |
| 확장성 | 낮음 | 높음 |
| 재사용성 | 낮음 | 높음 |
| 로컬 테스트 | sam local start-api | cdk synth |
| 복잡한 로직 | 어려움 | 쉬움 |

## 📚 CDK 구조 확장

향후 새로운 스택 추가 시:

```python
# cdk/stacks/database_stack.py
from aws_cdk import aws_dynamodb as dynamodb

class DatabaseStack(core.Stack):
    def __init__(self, scope, id, **kwargs):
        super().__init__(scope, id, **kwargs)
        
        table = dynamodb.Table(
            self, "MCPTable",
            partition_key=dynamodb.Attribute(
                name="id",
                type=dynamodb.AttributeType.STRING
            )
        )
```

```python
# cdk/app.py
from stacks.database_stack import DatabaseStack

# 기존 스택
MCPLambdaStack(app, "mcp-server-stack", ...)

# 새로운 스택 추가 (1줄!)
DatabaseStack(app, "database-stack", ...)
```

## 🚨 주의사항

### 1. AWS 계정 ID 설정
`app.py`의 `YOUR_ACCOUNT_ID`를 실제 AWS 계정 ID로 변경하거나 `cdk-deploy.sh`가 자동으로 설정합니다.

```bash
aws sts get-caller-identity --query Account --output text
```

### 2. 리전 설정
기본값: `ap-northeast-2` (서울)
변경하려면 `app.py`의 `region` 값 수정

### 3. CloudFormation 스택 이름
`app.py`에서 `"mcp-server-stack"` = CloudFormation 스택 이름

### 4. 상태 관리
- SAM의 `samconfig.toml` ❌ (더 이상 필요 없음)
- CDK는 CloudFormation으로 관리 ✅

## 📊 CDK vs SAM 마이그레이션

**변경사항:**
- ❌ `template.yaml` (SAM) - 삭제 또는 보관
- ❌ `samconfig.toml` - 삭제 또는 보관
- ✅ `cdk/` 디렉토리 - 새로 생성
- ✅ `deploy/cdk-deploy.sh` - 새로운 배포 스크립트

**배포 방식:**
```bash
# 이전 (SAM)
./deploy/deploy.sh

# 현재 (CDK)
./deploy/cdk-deploy.sh
```

## 🔗 참고 자료

- [AWS CDK Python API Reference](https://docs.aws.amazon.com/cdk/api/v2/python/)
- [AWS CDK Workshop](https://cdkworkshop.com/)
- [AWS CDK Best Practices](https://docs.aws.amazon.com/cdk/v2/guide/best-practices.html)

---

**마지막 업데이트:** 2025년 12월
