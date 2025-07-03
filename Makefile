LOCAL_BIN:=$(CURDIR)/bin
PROTO_DIR := api/user_v1
OUTPUT_DIR := pkg/user_v1

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u github.com/golang/protobuf
	go get -u google.golang.org/grpc
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-user-api

generate-user-api:
	@echo "Generating User API protobuf files..."
	@mkdir -p $(OUTPUT_DIR)
	protoc \
		--proto_path=$(PROTO_DIR) \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
		--go_out=$(OUTPUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUTPUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/user.proto

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux ./cmd/grpc_server/main.go

copy-to-server:
	scp service_linux root@87.228.103.116:

docker-build-and-push:
	docker buildx build  --no-cache --platform linux/amd64 -t cr.selcloud.ru/cerys/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAAdQUD7n1KenY0kRQXAWPmmaddytMko6WT cr.selcloud.ru/cerys
	docker push cr.selcloud.ru/cerys/test-server:v0.0.1