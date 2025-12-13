#!/bin/bash

# MCP Server on AWS Lambda - 배포 스크립트

set -e

echo "=========================================="
echo "MCP Server AWS Lambda 배포 스크립트"
echo "=========================================="
echo ""

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. AWS 자격증명 확인
echo -e "${YELLOW}1️⃣  AWS 자격증명 확인 중...${NC}"
if ! aws sts get-caller-identity &> /dev/null; then
    echo -e "${RED}❌ AWS 자격증명이 설정되지 않았습니다${NC}"
    echo ""
    echo "다음 명령어로 자격증명을 설정하세요:"
    echo "  aws configure"
    echo ""
    echo "또는 환경변수로 설정:"
    echo "  export AWS_ACCESS_KEY_ID=your_key"
    echo "  export AWS_SECRET_ACCESS_KEY=your_secret"
    echo "  export AWS_DEFAULT_REGION=ap-northeast-2"
    echo ""
    exit 1
fi

# AWS 계정 정보 출력
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
USERNAME=$(aws sts get-caller-identity --query Arn --output text | cut -d'/' -f2)
echo -e "${GREEN}✅ AWS 자격증명 확인 완료${NC}"
echo "   Account ID: $ACCOUNT_ID"
echo "   User: $USERNAME"
echo ""

# 2. 가상환경 활성화
echo -e "${YELLOW}2️⃣  Python 가상환경 활성화 중...${NC}"
if [ ! -d ".venv" ]; then
    echo -e "${RED}❌ venv 디렉토리가 없습니다${NC}"
    echo "   'make init'을 실행하여 가상환경을 생성하세요"
    exit 1
fi
source .venv/bin/activate
echo -e "${GREEN}✅ 가상환경 활성화 완료${NC}"
echo ""

# 3. SAM CLI 확인
echo -e "${YELLOW}3️⃣  SAM CLI 확인 중...${NC}"
if ! command -v sam &> /dev/null; then
    echo -e "${RED}❌ SAM CLI가 설치되지 않았습니다${NC}"
    echo "   'make deps'를 실행하여 의존성을 설치하세요"
    exit 1
fi
echo -e "${GREEN}✅ SAM CLI 준비 완료${NC}"
echo ""

# 4. Go 코드 빌드
echo -e "${YELLOW}4️⃣  Go 코드 빌드 중...${NC}"
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go server.go
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 빌드 실패${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 바이너리 빌드 완료${NC}"
echo ""

# 5. Lambda 배포 패키지 생성
echo -e "${YELLOW}5️⃣  배포 패키지 생성 중...${NC}"
zip -q -r function.zip bootstrap
echo -e "${GREEN}✅ 배포 패키지 생성 완료${NC}"
echo ""

# 6. samconfig.toml 확인
echo -e "${YELLOW}6️⃣  배포 설정 확인 중...${NC}"
if [ -f "samconfig.toml" ]; then
    echo -e "${GREEN}✅ 기존 배포 설정을 사용합니다${NC}"
    DEPLOY_OPTION=""
else
    echo -e "${YELLOW}⚠️  처음 배포입니다. 대화형 설정을 진행합니다${NC}"
    DEPLOY_OPTION="--guided"
fi
echo ""

# 7. SAM 배포
echo -e "${YELLOW}7️⃣  AWS Lambda에 배포 중...${NC}"
echo ""
sam deploy $DEPLOY_OPTION

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 배포 실패${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}=========================================="
echo "🎉 배포가 완료되었습니다!"
echo "==========================================${NC}"
echo ""

# 8. 배포 결과 출력
echo -e "${YELLOW}배포 정보:${NC}"
echo ""
aws cloudformation describe-stacks \
  --stack-name mcp-server-stack \
  --query 'Stacks[0].Outputs' \
  --region ap-northeast-2 \
  --output table 2>/dev/null || echo "⚠️  스택 정보를 가져올 수 없습니다"

echo ""
echo -e "${YELLOW}다음 단계:${NC}"
echo ""
echo "1. API 엔드포인트로 헬스체크:"
echo "   curl -X GET https://YOUR_API_ENDPOINT/dev/"
echo ""
echo "2. MCP 요청 테스트:"
echo "   curl -X POST https://YOUR_API_ENDPOINT/dev/mcp \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"jsonrpc\": \"2.0\", \"method\": \"initialize\", \"params\": {}, \"id\": 1}'"
echo ""
echo "3. CloudWatch 로그 확인:"
echo "   aws logs tail /aws/lambda/mcp-server-dev --follow --region ap-northeast-2"
echo ""
echo "자세한 배포 가이드는 DEPLOYMENT_GUIDE.md를 참고하세요"
echo ""
