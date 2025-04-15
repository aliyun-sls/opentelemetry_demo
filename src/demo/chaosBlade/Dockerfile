FROM golang:1.23.7-alpine AS builder

WORKDIR /build

ENV GOPROXY https://goproxy.cn


FROM builder AS builder_chaosblade

WORKDIR /build/chaosblade
COPY ./src/demo/chaosBlade .
RUN go mod tidy -compat=1.20
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o run .

FROM alpine AS chaosblade

WORKDIR /app
COPY --from=builder_chaosblade /build/chaosblade/run /app/run

CMD ["/app/run"]


