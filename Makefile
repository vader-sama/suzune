APP_NAME := suzune
PKG := ./...
BUILD_DIR := ./bin
RELEASE_DIR := ./release
MAIN_FILE := ./cmd/suzune
GO_CMD := go
TARGETS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: all build run clean test lint tidy deps release

all: build

build:
	@echo "🔨 Building..."
	@start=$$(date +%s); \
	$(GO_CMD) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE); \
	end=$$(date +%s); \
	duration=$$((end - start)); \
	echo "✅ Build complete! Binary is at $(BUILD_DIR)/$(APP_NAME) (took $${duration}s)"

release: clean
	@echo "🚀 Building release binaries..."
	@mkdir -p $(RELEASE_DIR)
	@for target in $(TARGETS); do \
		OS=$$(echo $$target | cut -d/ -f1); \
        ARCH=$$(echo $$target | cut -d/ -f2); \
		BIN_NAME=$(APP_NAME); \
		if [ "$$OS" = "windows" ]; then BIN_NAME=$${BIN_NAME}.exe; fi; \
		echo "🔨 Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH $(GO_CMD) build -o $(RELEASE_DIR)/$$BIN_NAME $(MAIN_FILE); \
	done
	@echo "✅ Release build complete! Binaries are in $(RELEASE_DIR)"

run: build
	@echo "🚀 Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

lint:
	@echo "🔍 Linting code..."
	@golangci-lint run
	@echo "✅ Linting complete. No major issues found!"

tidy:
	@echo "🧹 Tidying modules..."
	@$(GO_CMD) mod tidy
	@$(GO_CMD) mod verify
	@echo "✅ Modules tidy and verified."

deps:
	@echo "📦 Downloading dependencies..."
	@$(GO_CMD) mod download
	@echo "✅ Dependencies downloaded."

clean:
	@echo "🧺 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(RELEASE_DIR)
	@echo "✅ Clean complete. binaries removed."
