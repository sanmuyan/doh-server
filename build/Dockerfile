FROM  golang:1.22.1 as builder
ENV GO111MODULE=on

WORKDIR /workspace

COPY ./ ./

RUN --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go \
    go build -o doh-server ./cmd/server

FROM ubuntu:22.04
ENV LANG=en_US.utf8

WORKDIR /workspace

RUN apt-get update && \
    apt-get install -y ca-certificates curl net-tools tcpdump dnsutils

COPY --from=builder /workspace/doh-server /workspace

CMD ["/bin/bash", "-c", "/workspace/doh-server"]