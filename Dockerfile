FROM golang:1.18-alpine AS builder

ENV CGO_ENABLED 0
ENV TZ=Europe/Moscow

WORKDIR /go/src/mdns

COPY . .

RUN go mod tidy && go build -o /go/src/mdns/app /go/src/mdns/cmd/main.go

FROM gruebel/upx:latest as upx
COPY --from=builder /go/src/mdns/app /app
RUN upx --best --lzma -o /mdns /app

FROM scratch

COPY --from=upx /mdns /mdns

ENV REDIS_URL="redis://localhost:6379"
ENV REDIS_KEY="DUMP"
ENV HTTP_PORT=8081
ENV UDP_PORT=5353
ENV ACME_URL="https://acme-staging-v02.api.letsencrypt.org/directory"
ENV DOMAIN="example.com"
ENV IPV4="0.0.0.0"

EXPOSE 8081/tcp 5353/udp

CMD ["/mdns"]