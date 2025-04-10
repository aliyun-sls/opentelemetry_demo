# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0


FROM golang:1.24-alpine AS builder

WORKDIR /usr/src/app/

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=./src/demo_bench/go.sum,target=go.sum \
    --mount=type=bind,source=./src/demo_bench/go.mod,target=go.mod \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,rw,source=./src/demo_bench,target=. \
    go build -ldflags "-s -w" -o /go/bin/demo_bench/ ./

FROM alpine

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/demo_bench/ ./
COPY people.json ./

ENTRYPOINT [ "./demo_bench" ]
