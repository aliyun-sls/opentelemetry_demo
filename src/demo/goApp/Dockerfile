FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/app/

COPY ./tools/instgo /tmp/instgo
COPY ./src/demo/goApp/tupian /usr/src/app/tupian
COPY ./src/demo/common /usr/src/app/common
RUN chmod +x /tmp/instgo

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,rw,source=./src/demo/goApp/go.sum,target=go.sum \
    --mount=type=bind,rw,source=./src/demo/goApp/go.mod,target=go.mod \
    --mount=type=bind,source=./src/demo/common,target=./common \
    go mod edit -replace sls-mall-go/common=./common && \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,rw,source=./src/demo/goApp,target=goApp \
    cd goApp && \
    /tmp/instgo go build -ldflags "-s -w" -o /go/bin/goapp/ ./

FROM alpine

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/goapp/ ./
COPY --from=builder /usr/src/app/tupian ./tupian

EXPOSE 8080
ENTRYPOINT [ "./goapp" ]