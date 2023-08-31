# syntax=docker/dockerfile:1

FROM debian:latest
RUN apt-get update && apt-get install -y ca-certificates golang-go

WORKDIR /app
COPY source/* ./
RUN go mod download
RUN go build -o app

COPY server.crt server.key redis.crt redis.key ./
COPY redis_ca.crt /usr/share/ca-certificates/avito/
RUN echo avito/redis_ca.crt >> /etc/ca-certificates.conf
RUN update-ca-certificates

CMD ["./app"]
