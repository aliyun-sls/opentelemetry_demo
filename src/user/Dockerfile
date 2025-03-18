FROM golang:1.18-alpine3.16 AS builder

WORKDIR /demo
ENV GOPROXY https://goproxy.cn

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o demo .

FROM alpine:3.16 AS product

WORKDIR /app
COPY --from=builder /demo/demo /app/
EXPOSE 8080

ENTRYPOINT ["/app/demo"]