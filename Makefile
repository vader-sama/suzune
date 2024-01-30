BINARY=suzune
BUILD_DIR=bin
MAIN_PATH=cmd/suzune/main.go
COLOR_RED=\033[0;31m
COLOR_GREEN=\033[0;32m
COLOR_YELLOW=\033[0;33m
COLOR_RESET=\033[0m

all: clean build_mac build_linux build_windows
	@echo "Complied succesfully!"

build_mac:
	@echo "Compiling for [${COLOR_GREEN}MacOS${COLOR_RESET}] (${COLOR_YELLOW}x86-64${COLOR_RESET})"
	env GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY}_macos_x86-64 ${MAIN_PATH}

	@echo "Compiling for [${COLOR_GREEN}MacOS${COLOR_RESET}] (${COLOR_YELLOW}ARM64${COLOR_RESET})"
	env GOOS=darwin GOARCH=arm64 go build -o ${BUILD_DIR}/${BINARY}_macos_arm64 ${MAIN_PATH}
 
build_linux: 
	@echo "Compiling for [${COLOR_GREEN}Linux${COLOR_RESET}] (${COLOR_YELLOW}x86-64${COLOR_RESET})"
	env GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY}_linux_x86-64 ${MAIN_PATH}

	@echo "Compiling for [${COLOR_GREEN}Linux${COLOR_RESET}] (${COLOR_YELLOW}ARM64${COLOR_RESET})"
	env GOOS=linux GOARCH=arm64 go build -o ${BUILD_DIR}/${BINARY}_linux_arm64 ${MAIN_PATH}
 
build_windows: 
	@echo "Compiling for [${COLOR_GREEN}Windows${COLOR_RESET}] (${COLOR_YELLOW}x86-64${COLOR_RESET})"
	env GOOS=windows GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY}_widnows_x86-64 ${MAIN_PATH}

	@echo "Compiling for [${COLOR_GREEN}Windows${COLOR_RESET}] (${COLOR_YELLOW}ARM64${COLOR_RESET})"
	env GOOS=windows GOARCH=arm64 go build -o ${BUILD_DIR}/${BINARY}_windows_arm64 ${MAIN_PATH}

clean:
	@echo "Cleaning ${COLOR_RED}${BUILD_DIR}${COLOR_RESET}"
	rm -rf bin/*

.PHONY: all clean build_mac build_linux build_windows