FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/app/

COPY ./tools/instgo /tmp/instgo
COPY ./src/demo/common common
RUN go mod tidy
RUN go mod download
RUN chmod +x /tmp/instgo

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=./src/demo/product/go.sum,target=go.sum \
    --mount=type=bind,source=./src/demo/product/go.mod,target=go.mod \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,rw,source=./src/demo/product,target=. \
    /tmp/instgo go build -ldflags "-s -w" -o /go/bin/product/ ./

FROM alpine

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/product/ ./

EXPOSE 8080
ENTRYPOINT [ "./product" ]
