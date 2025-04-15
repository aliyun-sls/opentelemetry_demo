FROM golang:1.20-alpine3.16 AS builder

WORKDIR /build

ENV GOPROXY https://goproxy.cn


FROM builder AS builder_chaosblade

WORKDIR /build/chaosblade
COPY chaosblade .
RUN go mod tidy -compat=1.20
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o run .

FROM alpine:3.16 AS chaosblade

WORKDIR /app
COPY --from=builder_chaosblade /build/chaosblade/run /app/run
COPY yaml /app/chaosblade/yaml
#RUN ls -l /app/chaosblade



CMD ["/app/run"]


