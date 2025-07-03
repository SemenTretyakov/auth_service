FROM golang:1.23-alpine AS builder

COPY . /github.com/SemenTretyakov/auth_service/grpc/source
WORKDIR /github.com/SemenTretyakov/auth_service/grpc/source

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/SemenTretyakov/auth_service/grpc/source/bin/crud_server .

CMD ["./crud_server"]