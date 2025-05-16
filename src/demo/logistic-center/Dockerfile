FROM golang:1.23.7 AS builder

RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

WORKDIR /build/common
COPY ./src/demo/common .
COPY ./tools/instgo /build/instgo
RUN chmod +x /build/instgo
RUN go mod tidy

WORKDIR /build/logistic-center
COPY ./src/demo/logistic-center .
RUN go mod tidy
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux /build/instgo go build -ldflags "-s -w" -o logistic-center .


FROM debian:stable-slim AS logistic-center
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /build/logistic-center/logistic-center /app/
COPY --from=builder /build/logistic-center/config /app/config
EXPOSE 8080
ENTRYPOINT ["/app/logistic-center"]
