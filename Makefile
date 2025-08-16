LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=./migrations
LOCAL_MIGRATION_DSN="host=localhost port=54321 dbname=auth_service user=postgres password=postgres"

PROTO_DIR := api/user_v1
OUTPUT_DIR := pkg/user_v1

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

get-deps:
	go get -u github.com/golang/protobuf
	go get -u google.golang.org/grpc
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get github.com/joho/godotenv

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
	GOOS=linux GOARCH=amd64 go build -o auth_service ./cmd/app/main.go

copy-to-server:
	scp auth_service root@87.228.103.116:

docker-build-and-push:
	docker buildx build  --no-cache --platform linux/amd64 -t cr.selcloud.ru/cerys/auth-server:v0.0.1 .
	docker login -u token -p CRgAAAAAdQUD7n1KenY0kRQXAWPmmaddytMko6WT cr.selcloud.ru/cerys
	docker push cr.selcloud.ru/cerys/auth-server:v0.0.1

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

run-local:
	go run ./cmd/app/main.go -config-path=local.env

run-prod:
	go run ./cmd/app/main.go -config-path=prod.env