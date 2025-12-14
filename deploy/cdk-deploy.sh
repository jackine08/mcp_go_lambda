#!/bin/bash

# AWS CDK 배포 스크립트
# 자동으로 Go 바이너리를 빌드하고 AWS에 배포합니다

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== AWS CDK 배포 시작 ===${NC}"

# =====================
# 1. 환경변수 확인
# =====================
echo -e "\n${YELLOW}[1/5] 환경변수 확인 중...${NC}"
if [ -f "$(dirname "$0")/../.env" ]; then
    source "$(dirname "$0")/../.env"
    echo "✓ .env 파일 로드됨"
else
    echo "⚠ .env 파일이 없습니다 (선택사항)"
fi

# =====================
# 2. Go 바이너리 빌드
# =====================
echo -e "\n${YELLOW}[2/5] Go 바이너리 빌드 중...${NC}"
cd "$(dirname "$0")/.."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap main.go server.go
echo "✓ 바이너리 빌드 완료: bootstrap"

# =====================
# 3. Python 환경 설정
# =====================
echo -e "\n${YELLOW}[3/5] Python 환경 설정 중...${NC}"
cd cdk

# Python venv 확인 또는 생성
if [ ! -d "venv" ]; then
    echo "  venv 생성 중..."
    python3 -m venv venv
fi

# venv 활성화
source venv/bin/activate
echo "✓ Python 가상환경 활성화됨"

# 의존성 설치
echo "  CDK 의존성 설치 중..."
pip install -q -r requirements.txt
echo "✓ CDK 의존성 설치 완료"

# =====================
# 4. AWS 계정 ID 설정
# =====================
echo -e "\n${YELLOW}[4/5] AWS 계정 ID 확인 중...${NC}"
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text 2>/dev/null || echo "")
if [ -z "$ACCOUNT_ID" ]; then
    echo -e "${RED}✗ AWS 계정 ID를 가져올 수 없습니다${NC}"
    echo "  app.py의 YOUR_ACCOUNT_ID를 수동으로 설정하세요"
    exit 1
fi
echo "✓ AWS 계정 ID: $ACCOUNT_ID"

# app.py의 YOUR_ACCOUNT_ID를 실제 값으로 대체
sed -i "s/YOUR_ACCOUNT_ID/$ACCOUNT_ID/g" app.py 2>/dev/null || \
    sed -i '' "s/YOUR_ACCOUNT_ID/$ACCOUNT_ID/g" app.py
echo "✓ app.py 업데이트됨"

# =====================
# 5. CDK 배포
# =====================
echo -e "\n${YELLOW}[5/5] AWS CDK 배포 중...${NC}"
cdk deploy --require-approval=never

echo -e "\n${GREEN}=== 배포 완료! ===${NC}"
echo -e "\n배포된 리소스:"
echo "  - Lambda Function: mcp-server-dev"
echo "  - API Gateway: mcp-api"
echo "  - CloudWatch Logs: /aws/lambda/mcp-server-dev"
echo "  - CloudWatch Dashboard: mcp-server-dev"

# deactivate는 옵션 (스크립트 종료 시 자동으로 deactivate됨)
