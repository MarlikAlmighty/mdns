FROM golang:1.14-alpine3.12 AS builder

WORKDIR /go/src/library

COPY . .

RUN go build -o /go/src/library/app /go/src/library/cmd/main.go

FROM alpine:3.12

COPY --from=builder /go/src/library/app /

ENV PREFIX="LIBRARY"
ENV LIBRARY_MIGRATE=true
ENV LIBRARY_PATH_TO_MIGRATE="./migrations"
ENV LIBRARY_DB="postgres://user:password@localhost:5432/usecase?sslmode=disable"
ENV LIBRARY_HTTP_PORT=3000

EXPOSE 3000
CMD ["/app"]
