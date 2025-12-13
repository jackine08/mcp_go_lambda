.PHONY: build deploy clean test init deps run test-local

# 가상환경 활성화 명령어
VENV := .venv
VENV_BIN := $(VENV)/bin
ACTIVATE := source $(VENV_BIN)/activate

# 초기 설정
init: venv deps
	@echo "초기 설정 완료"

# venv 생성
venv:
	@if [ ! -d "$(VENV)" ]; then \
		uv venv; \
		echo "venv 생성 완료"; \
	fi

# 의존성 설치 (Python + Go)
deps: venv
	$(ACTIVATE) && uv pip install -r requirements.txt
	go mod download
	go mod tidy

# 빌드
build:
	GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	zip function.zip bootstrap

# 로컬 테스트
test-local: build
	$(ACTIVATE) && sam local invoke MCPServerFunction -e events/api-gateway-event.json

# 배포
deploy: build
	$(ACTIVATE) && sam deploy --guided

# 정리
clean:
	rm -f bootstrap function.zip
	rm -rf .aws-sam

# 개발용 run (로컬)
run:
	go run main.go
