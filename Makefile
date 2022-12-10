GOIMPORTS_VERSION = v0.3.0
GOLANDCI_LINT_VERSION = v1.50.1
GOTESTSUM_VERSION = v1.8.2
LOCAL_PACKAGES = github.com/kapitanov/habrabot
PROJECT_NAME = habrabot

all: clean build test lint

help:
	@echo "Usage:"
	@echo "\tmake build"
	@echo "\tmake clean"
	@echo "\tmake lint"
	@echo "\tmake test"
	@echo "\tmake docker"
	@echo "\tmake fmt"
	@echo "\tmake run"

build:
	@go build -v

clean:
	@go clean ./...

lint:
	@which golangci-lint > /dev/null || echo "Installing golangci-lint $(GOLANDCI_LINT_VERSION)" && \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANDCI_LINT_VERSION)
	@golangci-lint run --config .golangci.yml --timeout 10m --new-from-rev=$$(git merge-base HEAD master) ./...

test: build
	@which gotestsum > /dev/null || echo "Installing gotestsum $(GOTESTSUM_VERSION)" && \
		go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)
	@gotestsum ./...

docker:
	@docker build -t $(PROJECT_NAME):latest .

fmt:
	@which goimports > /dev/null || echo "Installing goimports $(GOIMPORTS_VERSION)" && \
		go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	@goimports -local $(LOCAL_PACKAGES) -w -format-only $$(find . -type f -name '*.go')

format: fmt # alias for fmt

run:
	@test -f .env || (echo "Error: missing .env file" && exit 1)
	@go build
	@./$(PROJECT_NAME) -env ./.env

