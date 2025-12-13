# 📂 파일 정리 가이드

> 빌드 산출물과 임시 파일을 정리하는 방법

## 🗂️ 현재 구조와 정리 전략

### 1. **gitignore에 이미 포함된 파일들** (이미 정리됨)
```bash
bootstrap          # 컴파일된 바이너리 (9.6MB)
function.zip       # 배포 패키지 (5.3MB)
.venv/             # Python 가상환경
__pycache__/       # Python 캐시
.aws-sam/          # SAM 빌드 산출물
.env               # AWS 자격증명
```

### 2. **필수 파일 (유지)** ✅
```
main.go                    # Lambda 핸들러
server.go                  # MCP 서버 로직
server_test.go             # 단위 테스트
go.mod, go.sum             # Go 의존성
template.yaml              # SAM 배포 템플릿
samconfig.toml             # SAM 설정
Makefile                   # 빌드 자동화
deploy.sh                  # 배포 스크립트
requirements.txt           # Python 의존성
.gitignore, .envrc         # 환경 설정
README.md                  # 프로젝트 개요
```

### 3. **문서 파일 (유지)** 📚
```
PROJECT_SUMMARY.md         # 프로젝트 요약 (NEW)
DEPLOYMENT_GUIDE.md        # 배포 가이드
DEPLOYMENT_CHECKLIST.md    # 배포 체크리스트
AWS_SETUP.md               # AWS 설정 가이드
FILE_CLEANUP_GUIDE.md      # 이 파일 (NEW)
```

### 4. **테스트 파일 (선택)** 📋
```
events/
  ├── initialize-request.json
  ├── tools-list-request.json
  ├── resources-list-request.json
  └── prompts-list-request.json
```

---

## 🧹 정리 절차

### 방법 A: 로컬 개발용 정리 (권장)

```bash
# 1. 로컬 빌드 산출물 제거
rm -f bootstrap function.zip

# 2. 캐시 파일 제거
rm -rf __pycache__
rm -rf .aws-sam

# 3. 불필요한 파일 제거
rm -f *.log *.tmp

# 4. Makefile의 clean 타겟 사용
make clean
```

**결과:**
- 소스 코드는 유지
- 배포 파일은 필요할 때 다시 생성
- 저장소 크기 감소 (전체 -15MB)

### 방법 B: Git 저장소 정리

```bash
# 1. 변경사항 확인
git status

# 2. gitignore된 파일 확인
git check-ignore -v bootstrap function.zip .venv

# 3. 변경사항 커밋
git add -A
git commit -m "chore: 정리 및 문서 업데이트"

# 4. (선택) 원격에 푸시
git push origin main
```

### 방법 C: 완전 정리 (리셋)

```bash
# 경고: 로컬 변경사항 모두 제거됨
git clean -fdx

# 또는 특정 패턴만 제거
rm -rf bootstrap function.zip .venv __pycache__ .aws-sam
```

---

## 📋 정리 체크리스트

- [ ] `bootstrap` 파일 제거 (또는 `make clean`)
- [ ] `function.zip` 파일 제거
- [ ] `__pycache__` 디렉토리 제거
- [ ] `.aws-sam` 디렉토리 제거
- [ ] `.env` 파일 확인 (gitignore 포함 확인)
- [ ] 불필요한 로그 파일 제거
- [ ] `README.md` 업데이트 확인 ✅
- [ ] `PROJECT_SUMMARY.md` 생성 확인 ✅
- [ ] `git status` 확인
- [ ] `git commit` 및 `git push`

---

## 🚀 배포 후 재빌드

정리 후 다시 배포해야 할 경우:

```bash
# 1. 빌드
make build

# 2. 배포
./deploy.sh

# 또는 단계별
make clean
make build
make deploy
```

---

## 💾 저장소 상태 관리

### 개발 중
```
git add    : 소스 코드만 (go, yaml, sh, md)
git ignore : 빌드 산출물, 자격증명, 캐시
```

### 배포 전
```bash
# 최종 확인
git status
git log --oneline -5

# 깔끔한 상태 확인
git check-ignore -v *
```

### 배포 후
```bash
# CI/CD 자동화 대비
# (현재는 수동이지만 GitHub Actions 추가 고려)
```

---

## 📊 파일 크기 비교

### 정리 전
```
bootstrap     9.6 MB  ❌ 제거 가능
function.zip  5.3 MB  ❌ 제거 가능
.venv/        500 MB  ❌ gitignore (로컬만)
__pycache__   2 MB    ❌ 제거 가능
.aws-sam/     10 MB   ❌ gitignore (로컬만)
소스 코드     200 KB  ✅ 유지
문서          100 KB  ✅ 유지
---
총계          ~525 MB (로컬) / ~100 KB (git)
```

### 정리 후
```
소스 코드     200 KB  ✅
문서          100 KB  ✅
---
총계          ~300 KB (매우 경량)
```

---

## 🔄 자동 정리 (Makefile)

`Makefile`에 이미 정의된 clean 타겟:

```bash
make clean
```

내용:
```makefile
clean:
	rm -f bootstrap
	rm -f function.zip
	rm -f .aws-sam
	rm -rf __pycache__
```

---

## ✅ 정리 완료 후 검증

```bash
# 1. 파일 목록 확인
ls -la | grep -E "bootstrap|function.zip|__pycache__|.aws-sam"
# (아무것도 출력 안 됨 = 성공)

# 2. Git 상태 확인
git status
# (clean working tree = 성공)

# 3. 소스 코드 확인
ls -la *.go
# (main.go, server.go, server_test.go 존재 = 성공)

# 4. 배포 설정 확인
ls -la template.yaml samconfig.toml Makefile deploy.sh
# (모두 존재 = 성공)
```

---

## 📚 참고

- **Makefile**: `make clean` 목표 사용
- **.gitignore**: 이미 올바르게 설정됨
- **자동 배포**: `deploy.sh`가 필요할 때마다 새로 생성

---

## 🎯 결론

정리의 이점:
- ✅ 저장소 크기 감소 (500MB → 300KB)
- ✅ 배포 산출물 자동 생성
- ✅ 깔끔한 프로젝트 구조
- ✅ 협업 시 불필요한 파일 제외

**권장:** `make clean` 정기적 실행 및 배포 전 자동 빌드
