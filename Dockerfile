FROM golang:alpine AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0

COPY . /code

RUN set -ex \
    && cd /code \
    && go build -ldflags "-s -w" -o wol-forwarder .

FROM alpine:latest

COPY --from=builder /code/wol-forwarder /usr/sbin/wol-forwarder

ENV WOL_ADDR="0.0.0.0"
ENV WOL_PORT="1999"
ENV WOL_BADDR="255.255.255.255"
ENV WOL_BPORT="9"

ENTRYPOINT ["/usr/sbin/wol-forwarder"]
