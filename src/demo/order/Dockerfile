FROM golang:1.23.7-alpine AS builder

WORKDIR /build

WORKDIR /build/common
COPY ./src/demo/common .
COPY ./tools/instgo /build/instgo
RUN chmod +x /build/instgo
RUN go mod tidy

WORKDIR /build/order
COPY ./src/demo/order .
RUN go mod tidy
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux /build/instgo go build -ldflags "-s -w" -o order .


FROM alpine:3.16 AS order

WORKDIR /app
COPY --from=builder /build/order/order /app/
COPY --from=builder /build/order/config /app/config
EXPOSE 8080
ENTRYPOINT ["/app/order"]
