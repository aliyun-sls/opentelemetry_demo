FROM golang:1.23-alpine3.20 AS builder

WORKDIR /demo
ENV GOPROXY https://goproxy.cn

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o demo .

FROM alpine:3.20 AS product

WORKDIR /app
COPY --from=builder /demo/demo /app/
EXPOSE 8080

ENTRYPOINT ["/app/demo"]