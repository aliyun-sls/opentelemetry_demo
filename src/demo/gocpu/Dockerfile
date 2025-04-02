FROM golang:1.23.7-alpine AS builder

WORKDIR /usr/src/app/

COPY ./tools/instgo /tmp/instgo
RUN chmod +x /tmp/instgo

RUN --mount=type=cache,target=/go/pkg/mod/ \
    #--mount=type=bind,rw,source=./src/demo/gocpu/go.sum,target=go.sum \
    --mount=type=bind,rw,source=./src/demo/gocpu/go.mod,target=go.mod \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,rw,source=./src/demo/gocpu,target=gocpu \
    cd gocpu && \
    /tmp/instgo go build -ldflags "-s -w" -o /go/bin/gocpu/ ./

FROM alpine

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/gocpu/ ./

EXPOSE 8080
ENTRYPOINT [ "./gocpu" ]