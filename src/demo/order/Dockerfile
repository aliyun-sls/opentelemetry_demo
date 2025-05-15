FROM golang:1.20-alpine AS builder

# 安装必要的依赖，包括 librdkafka
RUN apk add --no-cache gcc musl-dev librdkafka-dev pkgconf

WORKDIR /build

WORKDIR /build/common
COPY ./src/demo/common .
COPY ./tools/instgo /build/instgo
RUN chmod +x /build/instgo
RUN go mod tidy

WORKDIR /build/order
COPY ./src/demo/order .
COPY ./src/demo/order/go.mod .
COPY ./src/demo/order/go.sum .
RUN go mod download
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o order .


FROM alpine:3.16 AS order
RUN apk add --no-cache librdkafka
WORKDIR /app
COPY --from=builder /build/order/order /app/
COPY --from=builder /build/order/config /app/config
EXPOSE 8080
ENTRYPOINT ["/app/order"]
