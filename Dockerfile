FROM golang:1.21-alpine3.18 AS builder

ENV CGO_ENABLED 0
ENV TZ=Europe/Moscow

WORKDIR /go/src/mdns

COPY . .

RUN go mod tidy 
RUN go build -o /go/src/mdns/app /go/src/mdns/cmd/main.go

FROM gruebel/upx:latest as upx
COPY --from=builder /go/src/mdns/app /app
RUN upx --best --lzma -o /app /app

FROM scratch

COPY --from=upx /mdns /app

ENV HTTP_PORT=8081
ENV DNS_TCP_PORT=53
ENV DNS_UDP_PORT=53
ENV NAME_SERVERS=1.1.1.1,1.0.0.1,8.8.8.8,8.8.4.4

EXPOSE 8081/tcp 53/tcp 53/udp

CMD ["/mdns"]