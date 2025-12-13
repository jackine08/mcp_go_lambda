# 📚 Documentation Guide

이 디렉토리에는 프로젝트의 모든 문서가 포함되어 있습니다.

## 📖 문서 목록

### 📄 [README.md](./README.md)
**프로젝트 상세 개요 및 사용 방법**

- 프로젝트 개요
- 배포된 서버 정보 및 API Endpoint
- MCP 메서드 설명
- Tool 사용 예시 (add, multiply)
- 문제 해결 팁
- 보안 고려사항

→ **언제 봐야 하나?** 프로젝트를 처음 접할 때, 사용 방법이 궁금할 때

---

### 📄 [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md)
**8개 Phase별 완료된 작업 상세 요약**

- Phase 1-8 별 구현 내용
- 최종 파일 구조
- 현재 상태 (완료 / 미구현)
- 다음 단계

→ **언제 봐야 하나?** 어떤 작업이 완료되었는지 전체적으로 파악하고 싶을 때

---

### 📄 [COMPLETION_REPORT.md](./COMPLETION_REPORT.md)
**최종 완료 보고서, 학습 성과, 코드 통계**

- 수행한 작업 요약
- 최종 코드 통계 (줄 수, 테스트 결과)
- 배포된 결과
- 파일 정리 결과
- 주요 기능 및 특징
- 학습 성과

→ **언제 봐야 하나?** 프로젝트의 최종 상태를 한눈에 파악하고 싶을 때

---

### 📄 [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md)
**AWS 배포 단계별 상세 가이드**

- 배포 전 준비사항
- 로컬 개발 및 테스트
- 빌드 및 배포 단계
- 배포 후 확인사항
- 롤백 방법

→ **언제 봐야 하나?** AWS Lambda에 직접 배포해야 할 때

---

### 📄 [DEPLOYMENT_CHECKLIST.md](./DEPLOYMENT_CHECKLIST.md)
**배포 전 확인 사항 및 문제 해결**

- 배포 전 체크리스트
- 일반적인 에러 및 해결책
- CloudWatch 로그 확인 방법
- 디버깅 팁

→ **언제 봐야 하나?** 배포 실패하거나 문제가 발생했을 때

---

### 📄 [FILE_CLEANUP_GUIDE.md](./FILE_CLEANUP_GUIDE.md)
**파일 정리 방법 및 자동화**

- 빌드 산출물 정리
- gitignore 설정 확인
- 저장소 크기 최적화
- 자동 정리 방법

→ **언제 봐야 하나?** 저장소 크기를 줄이거나 불필요한 파일을 정리하고 싶을 때

---

### 📄 [AWS_SETUP.md](./AWS_SETUP.md)
**AWS IAM 설정 및 자격증명**

- AWS IAM 사용자 생성
- 자격증명 설정 방법
- .env 파일 설정
- AWS CLI 설정

→ **언제 봐야 하나?** AWS 환경을 처음 설정할 때

---

## 🔍 문서 선택 가이드

**상황별 읽어야 할 문서:**

### 🚀 처음 시작하는 경우
1. `README.md` (프로젝트 개요)
2. `AWS_SETUP.md` (AWS 설정)
3. `DEPLOYMENT_GUIDE.md` (배포 방법)

### 🔧 코드를 수정하고 다시 배포하는 경우
1. 로컬에서 수정
2. `go test -v` 실행
3. `deploy/deploy.sh` 실행

### 🐛 배포에 실패한 경우
1. `DEPLOYMENT_CHECKLIST.md` (에러 해결)
2. CloudWatch 로그 확인
3. `DEPLOYMENT_GUIDE.md` (단계 재확인)

### 📊 프로젝트 현황을 파악하고 싶은 경우
1. `PROJECT_SUMMARY.md` (전체 작업 요약)
2. `COMPLETION_REPORT.md` (최종 상태)

### 🧹 저장소를 정리하고 싶은 경우
1. `FILE_CLEANUP_GUIDE.md` (정리 방법)
2. `deploy/Makefile`의 `make clean` 실행

---

## 📋 문서 관계도

```
README.md (시작점)
├── DEPLOYMENT_GUIDE.md (배포 방법)
│   └── DEPLOYMENT_CHECKLIST.md (문제 해결)
├── PROJECT_SUMMARY.md (작업 요약)
│   └── COMPLETION_REPORT.md (최종 상태)
├── AWS_SETUP.md (AWS 설정)
└── FILE_CLEANUP_GUIDE.md (파일 정리)
```

---

## 🎯 빠른 참조

| 질문 | 문서 |
|------|------|
| 프로젝트가 뭐예요? | README.md |
| 어떤 작업이 완료되었어요? | PROJECT_SUMMARY.md |
| 최종 상태가 어때요? | COMPLETION_REPORT.md |
| 어떻게 배포하나요? | DEPLOYMENT_GUIDE.md |
| 배포가 실패했어요. | DEPLOYMENT_CHECKLIST.md |
| AWS 설정이 필요해요. | AWS_SETUP.md |
| 저장소가 너무 커요. | FILE_CLEANUP_GUIDE.md |

---

**마지막 업데이트:** 2025년 12월
