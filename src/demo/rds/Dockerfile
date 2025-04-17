FROM golang:1.23.7-alpine AS builder

WORKDIR /usr/src/app/

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,rw,source=./src/demo/rds/go.sum,target=go.sum \
    --mount=type=bind,rw,source=./src/demo/rds/go.mod,target=go.mod \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,rw,source=./src/demo/rds,target=. \
    go mod tidy && go build -ldflags "-s -w" -o /go/bin/rds ./

FROM alpine

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/rds ./

EXPOSE 8080
ENTRYPOINT [ "./rds" ]