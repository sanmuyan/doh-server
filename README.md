# DoH Server

## 使用方法

- `/dns-query`

提供基于 RFC 8484 的 DNS over HTTPS 查询，适用于所有支持 DoH 的浏览器和操作系统

**请求方法：**`GET`

**请求参数：**`dns=base64URL(dns_query)`

**请求示例：**`https://dns.example.com/dns-query?dns=uGkBAAABAAAAAAAAB2FsaWJhYmEDY29tAAABAAE`

## DNS JSON API

- `/resolve`

提供 DNS JSON API 查询

**请求方法**: `GET`

**请求示例：**`https://dns.example.com/resolve?name=www.example.com&type=1`

**返回示例**: https://developers.google.com/speed/public-dns/docs/doh/json

## 服务器配置

### 部署

```shell
git clone https://github.com/sanmuyan/doh-server
cd doh-server
docker build -t sanmuyan/doh-server:latest . -f ./build/Dockerfile
docker run --name doh-server -p 8053:8053 -d sanmuyan/doh-server:latest
```

### 启动参数

- `--cache` 开启缓存
- `--cache-ttl 60` 缓存过期时间
- `--upstream-server 8.8.8.8:53` 上游 DNS 服务器
- `--upstream-net udp` 上游 DNS 网络类型，支持 udp|tcp|tcp-tls|doh
- `--upstream-timeout 2` 上游 DNS 超时时间

### SSL 配置参考

```shell
server {
  bind  443  ssl  http2;
  ssl_certificate server.crt;
  ssl_certificate_key server.key;
  server_name  dns.example.com;

  location / {
    proxy_pass http://127.0.0.1:8053;
  }
}
```