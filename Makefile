.PHONY: build test clean deploy destroy

# .env 파일 로드
include .env
export

# 변수
BINARY_NAME=bootstrap
BUILD_DIR=lambda_dist
CMD_DIR=cmd/lambda
REGION=$(AWS_DEFAULT_REGION)

# 빌드
build:
	@echo "Building Lambda function..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 테스트
test:
	@echo "Running tests..."
	go test -v -cover ./...

# 테스트 (커버리지 리포트)
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# 클린
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -rf cdk/cdk.out
	@echo "Clean complete"

# CDK 배포
deploy: build
	@echo "Deploying with AWS CDK..."
	@echo "Region: $(AWS_DEFAULT_REGION)"
	cd cdk && \
	bash -c "source ../.venv/bin/activate && cdk deploy --require-approval never"

# CDK 제거
destroy:
	@echo "Destroying CDK stack..."
	cd cdk && \
	bash -c "source ../.venv/bin/activate && cdk destroy --force"

# CDK Synth
synth:
	@echo "Synthesizing CDK stack..."
	cd cdk && \
	bash -c "source ../.venv/bin/activate && cdk synth"

# 의존성 설치
deps:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "Installing Python dependencies..."
	pip install -r requirements.txt

# 로컬 테스트 (SAM Local)
local-test: build
	@echo "Testing locally with SAM..."
	cd deploy && \
	sam local invoke MCPServerFunction -e events/initialize-request.json

# 포맷팅
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete"

# Lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# 모든 빌드 및 테스트
all: clean deps fmt test build

# 도움말
help:
	@echo "Available targets:"
	@echo "  build          - Build Lambda function"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Remove build artifacts"
	@echo "  deploy         - Build and deploy to AWS"
	@echo "  destroy        - Destroy AWS resources"
	@echo "  synth          - Synthesize CDK stack"
	@echo "  deps           - Install dependencies"
	@echo "  local-test     - Test locally with SAM"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  all            - Clean, deps, format, test, and build"
	@echo "  help           - Show this help message"
