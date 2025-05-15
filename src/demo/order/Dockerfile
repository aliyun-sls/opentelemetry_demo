FROM golang:1.23.7-alpine AS builder


WORKDIR /usr/src/app/
COPY ./src/demo/order .
COPY ./tools/instgo /tmp/instgo
RUN go mod tidy
RUN chmod +x /tmp/instgo
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o order .

FROM alpine:3.16 AS order

WORKDIR /usr/src/app/
COPY --from=builder /usr/src/app/order /usr/src/app/
COPY --from=builder /usr/src/app/tupian ./tupian
EXPOSE 8080
ENTRYPOINT ["/usr/src/app/"]