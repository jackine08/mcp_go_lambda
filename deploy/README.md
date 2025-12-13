# 🚀 Deployment & Build

배포 및 빌드 자동화 도구들이 저장되는 디렉토리입니다.

## 📁 파일 구조

```
deploy/
├── Makefile             # 빌드 자동화
├── deploy.sh            # 배포 자동화 스크립트
└── events/              # 테스트 이벤트 샘플
    ├── api-gateway-event.json
    ├── initialize-request.json
    ├── tools-list-request.json
    ├── resources-list-request.json
    └── prompts-list-request.json
```

## 📄 파일 설명

### Makefile
**빌드 자동화 및 타겟화**

#### 주요 타겟

```bash
make init
```
**초기 설정** - Python venv, 의존성 설치
- Python 가상환경 생성 (.venv/)
- SAM CLI 설치 (requirements.txt)
- Go 모듈 다운로드

```bash
make build
```
**Go 빌드** - 정적 바이너리 생성
- CGO_ENABLED=0으로 정적 링크
- GOOS=linux GOARCH=amd64로 Linux 64비트 컴파일
- 결과: `bootstrap` 바이너리 (9.6MB)

```bash
make deploy
```
**AWS 배포** - Lambda 함수 배포
- SAM build 실행
- SAM deploy 실행
- 결과: CloudFormation Stack 업데이트

```bash
make clean
```
**정리** - 빌드 산출물 제거
- bootstrap 파일 제거
- function.zip 제거
- .aws-sam 디렉토리 제거

#### 사용 예시

```bash
# 초기 설정 (처음 한 번만)
make init

# 코드 수정 후 배포
make build
make deploy

# 또는 한 번에
make clean && make build && make deploy
```

---

### deploy.sh
**완전 자동화된 배포 스크립트**

#### 역할
1. `.env` 파일 검증 (AWS 자격증명 확인)
2. Python venv 활성화
3. Go 빌드 실행
4. SAM 배포 실행
5. 에러 처리 및 로깅

#### 사용 방법
```bash
./deploy.sh
```

#### 스크립트의 동작

```bash
# 1. 환경 변수 로드
source ../.env

# 2. venv 활성화
source ../.venv/bin/activate

# 3. Go 빌드
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap \
  main.go server.go

# 4. SAM 배포
sam deploy --stack-name mcp-server-stack \
  --no-confirm-changeset \
  --resolve-s3 \
  --capabilities CAPABILITY_IAM
```

#### 사전 요구사항
- `.env` 파일 존재 (AWS 자격증명)
- AWS CLI 설치
- SAM CLI 설치
- Go 1.22 이상

#### 에러 처리
```bash
set -e  # 에러 발생 시 즉시 종료
```

---

### events/ 디렉토리
**MCP 요청 샘플 데이터**

#### 파일 설명

**initialize-request.json**
```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {},
  "id": 1
}
```
→ 서버 초기화 요청

**tools-list-request.json**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "params": {},
  "id": 2
}
```
→ 사용 가능한 Tool 목록 조회

**resources-list-request.json**
```json
{
  "jsonrpc": "2.0",
  "method": "resources/list",
  "params": {},
  "id": 3
}
```
→ 리소스 목록 조회

**prompts-list-request.json**
```json
{
  "jsonrpc": "2.0",
  "method": "prompts/list",
  "params": {},
  "id": 4
}
```
→ 프롬프트 목록 조회

**api-gateway-event.json**
```json
{
  "httpMethod": "POST",
  "resource": "/mcp",
  "body": "{...json-rpc-request...}"
}
```
→ API Gateway 형식의 이벤트

#### 사용 방법

로컬 테스트:
```bash
sam local invoke MCPServerFunction -e events/initialize-request.json
```

---

## 🔄 배포 워크플로우

### Step 1: 로컬 개발 & 테스트
```bash
# 코드 수정
vi ../main.go

# 테스트 실행
cd ..
go test -v

# 돌아오기
cd deploy
```

### Step 2: 로컬 빌드
```bash
make build
```

결과:
- `bootstrap` 바이너리 생성 (9.6MB)
- `function.zip` 생성 (5.3MB)

### Step 3: 배포
```bash
./deploy.sh
```

또는:
```bash
make deploy
```

결과:
- CloudFormation Stack 업데이트
- Lambda 함수 배포
- API Gateway 활성화

### Step 4: 검증
```bash
# API 테스트
curl -X POST https://w2z3biq5nb.execute-api.ap-northeast-2.amazonaws.com/dev/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "initialize", "params": {}, "id": 1}'

# 로그 확인
aws logs tail /aws/lambda/mcp-server-dev --follow
```

---

## ✅ 배포 체크리스트

배포 전:
- [ ] 코드 수정 완료
- [ ] `go test -v` 통과
- [ ] `.env` 파일 확인
- [ ] AWS 자격증명 확인

배포 중:
- [ ] `make build` 성공
- [ ] 빌드 파일 생성됨 (bootstrap, function.zip)
- [ ] `./deploy.sh` 성공

배포 후:
- [ ] CloudFormation Stack 업데이트 확인
- [ ] Lambda 함수 배포 확인
- [ ] API Gateway 엔드포인트 응답 확인
- [ ] CloudWatch 로그 생성됨

---

## 🐛 일반적인 문제

### "Permission denied" 에러
```bash
chmod +x deploy.sh
```

### ".env: No such file or directory"
```bash
# 루트 디렉토리의 .env 생성
cd ..
cat > .env << 'EOF'
export aws_access_key=YOUR_KEY
export aws_secret_key=YOUR_SECRET
EOF
```

### "GLIBC version not found" 에러
```bash
# Makefile에서 CGO_ENABLED=0 확인
cat Makefile | grep CGO_ENABLED
```

### "sam: command not found"
```bash
# SAM CLI 설치
make init
```

---

## 📊 빌드 산출물

### bootstrap (9.6MB)
- Go 컴파일된 바이너리
- 정적 링크 (GLIBC 없음)
- Lambda가 실행할 파일

### function.zip (5.3MB)
- bootstrap을 압축한 배포 패키지
- CloudFormation으로 S3 업로드

### .aws-sam/ (10MB)
- SAM 빌드 캐시
- 다음 배포 시 활용

---

## 🧹 정리 및 최적화

### 빌드 산출물 제거
```bash
make clean
```

### 저장소 최적화
```bash
# 모든 불필요한 파일 제거
make clean
rm -rf ../.aws-sam
rm -rf ../__pycache__
```

### 저장소 크기 비교
```bash
# 정리 전: ~525MB
# 정리 후: ~300KB
# 절감: 99.9%
```

---

## 📚 참고 자료

- [AWS SAM CLI 문서](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html)
- [Go 빌드 옵션](https://golang.org/cmd/go/)
- [Lambda 배포 패키지](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)

---

**마지막 업데이트:** 2025년 12월
