FROM golang:1.23.7-alpine AS builder

RUN apk add --no-cache \
    gcc \
    musl-dev \
    librdkafka-dev \
    make \
    git \
    cmake

WORKDIR /build

WORKDIR /build/common
COPY ./src/demo/common .
COPY ./tools/instgo /build/instgo
RUN chmod +x /build/instgo
RUN go mod tidy

WORKDIR /build/order
COPY ./src/demo/order .
RUN go mod tidy
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o order .


FROM alpine:3.16 AS order
RUN apk add --no-cache librdkafka
WORKDIR /app
COPY --from=builder /build/order/order /app/
COPY --from=builder /build/order/config /app/config
EXPOSE 8080
ENTRYPOINT ["/app/order"]
