FROM golang:1.23.7-alpine AS builder

WORKDIR /usr/src/app/

COPY ./src/chaosBlade/yaml /usr/src/app/yaml

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,rw,source=./src/chaosBlade/go.sum,target=go.sum \
    --mount=type=bind,rw,source=./src/chaosBlade/go.mod,target=go.mod \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,rw,source=./src/chaosBlade,target=. \
    go build -ldflags "-s -w" -o /go/bin/chaosBlade ./

FROM alpine

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/chaosBlade ./
COPY --from=builder /usr/src/app/yaml ./yaml

EXPOSE 8080
ENTRYPOINT [ "./chaosBlade" ]