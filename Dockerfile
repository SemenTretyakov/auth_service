FROM golang:1.23-alpine AS builder

COPY . /github.com/SemenTretyakov/auth_service/app
WORKDIR /github.com/SemenTretyakov/auth_service/app

RUN go mod download
RUN go build -o ./bin/auth_server cmd/app/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/SemenTretyakov/auth_service/app/bin/auth_server .

CMD ["./auth_server"]