# Makefile for building Go project with optional json-iterator support

# 项目名称
PROJECT_NAME := myproject

# Go 源码目录
SRC_DIR := .

# 编译输出目录
BUILD_DIR := ./build

# JSON 库替换的标记
JSON_ITERATOR_TAG := jsoniter

# 默认 Go 编译命令
GO_BUILD := go build

# Wails 构建命令
WAILS_BUILD := wails build clean

# Go 源文件
SRC_FILES := $(SRC_DIR)/main.go

# 编译输出文件
OUTPUT_FILE := $(BUILD_DIR)/$(PROJECT_NAME)

.PHONY: all clean build build-jsoniter

# 默认任务
all: build

# 清理生成文件
clean:
		rm -rf $(BUILD_DIR)

# 普通编译
build:
		@echo "Building without json-iterator..."
		$(GO_BUILD) -o $(OUTPUT_FILE) $(SRC_FILES)

# 使用 json-iterator 编译
build-jsoniter:
		@echo "Building with json-iterator..."
		$(WAILS_BUILD) -tags=$(JSON_ITERATOR_TAG) -o $(OUTPUT_FILE)
