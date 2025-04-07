FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/app/

COPY ./tools/instgo /tmp/instgo
RUN chmod +x /tmp/instgo

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,rw,source=./src/demo/goApp_500/go.sum,target=go.sum \
    --mount=type=bind,rw,source=./src/demo/goApp_500/go.mod,target=go.mod \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,rw,source=./src/demo/goApp_500,target=goApp_500 \
    cd goApp && \
    /tmp/instgo go build -ldflags "-s -w" -o /go/bin/refuse/ ./

FROM alpine

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/refuse/ ./

EXPOSE 8080
ENTRYPOINT [ "./refuse" ]