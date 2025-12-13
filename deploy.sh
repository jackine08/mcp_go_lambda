#!/bin/bash

# MCP Server on AWS Lambda - 배포 스크립트

set -e

echo "=========================================="
echo "MCP Server AWS Lambda 배포 스크립트"
echo "=========================================="
echo ""

# 1. 가상환경 활성화
echo "1️⃣  Python 가상환경 활성화 중..."
source .venv/bin/activate
echo "✅ 가상환경 활성화 완료"
echo ""

# 2. AWS 자격증명 확인
echo "2️⃣  AWS 자격증명 확인 중..."
if ! command -v sam &> /dev/null; then
    echo "❌ SAM CLI가 설치되지 않았습니다"
    echo "   'make deps'를 실행하여 의존성을 설치하세요"
    exit 1
fi
echo "✅ SAM CLI 준비 완료"
echo ""

# 3. Go 코드 빌드
echo "3️⃣  Go 코드 빌드 중..."
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go server.go
echo "✅ 바이너리 빌드 완료"
echo ""

# 4. Lambda 배포 패키지 생성
echo "4️⃣  배포 패키지 생성 중..."
zip -q -r function.zip bootstrap
echo "✅ 배포 패키지 생성 완료"
echo ""

# 5. SAM 배포
echo "5️⃣  AWS Lambda에 배포 중..."
echo "    첫 배포인 경우 --guided 옵션으로 설정을 입력합니다"
echo "    리전은 'ap-northeast-2' (서울)을 권장합니다"
echo ""

sam deploy --guided

echo ""
echo "=========================================="
echo "🎉 배포가 완료되었습니다!"
echo "=========================================="
echo ""
echo "배포 후 다음 명령어로 API 엔드포인트를 확인하세요:"
echo "  aws cloudformation describe-stacks --stack-name MCPServerStack --region ap-northeast-2"
echo ""
