# syntax=docker/dockerfile:1

FROM debian:latest
RUN apt-get update && apt-get install -y ca-certificates redis-server

WORKDIR /data
COPY redis.conf /etc/redis/
COPY server.crt server.key ca.crt /etc/redis/

CMD ["redis-server", "/etc/redis/redis.conf"]
