PROTO_DIR := proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
BUILD_DIR := build
BINARY := $(BUILD_DIR)/kvstore

.PHONY: all proto build clean deps

all: deps proto build

deps:
	@echo "==> Installing Go dependencies..."
	@go mod tidy
	@echo "==> Installing Protobuf plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto: $(PROTO_FILES)
	@echo "==> Generating protobuf code..."
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		$(PROTO_FILES)

build:
	@echo "==> Building binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) ./main.go

clean:
	@echo "==> Cleaning..."
	@rm -rf $(BUILD_DIR)
	@find . -name "*.pb.go" -delete
