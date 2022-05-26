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

ENV HTTP_PORT=8081
ENV NAME_SERVERS="1.1.1.1:53,1.0.0.1:53,8.8.8.8:53,8.8.4.4:53"
ENV DNS_HOST="127.0.0.1"

EXPOSE 8081/tcp 53/tcp 53/udp

CMD ["/mdns"]