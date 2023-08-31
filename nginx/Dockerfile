# syntax=docker/dockerfile:1

FROM debian:latest
RUN apt-get update && apt-get install -y ca-certificates nginx

COPY nginx.conf /etc/nginx/

COPY app.crt app.key /etc/nginx/
COPY app_ca.crt /usr/share/ca-certificates/avito/
RUN echo avito/app_ca.crt >> /etc/ca-certificates.conf
RUN update-ca-certificates

CMD ["nginx", "-g", "daemon off;"]
