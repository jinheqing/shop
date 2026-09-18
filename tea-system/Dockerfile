# 多阶段构建
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# 运行阶段
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata wireguard-tools wget
COPY --from=builder /out/server /usr/local/bin/server
EXPOSE 8080
ENTRYPOINT ["server"]
