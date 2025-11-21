FROM  golang:1.24.1 as builder
ENV GO111MODULE=on

WORKDIR /workspace

COPY ./ ./

RUN --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go \
    go build -o doh-server ./cmd/server

FROM sanmuyan/ubuntu:22.04-base-20250427

COPY --from=builder /workspace/doh-server ./

CMD ["/opt/app/doh-server"]