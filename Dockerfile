# FROM golang:1.23-alpine AS builder


# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# RUN go build -o /goback main.go

FROM debian:stable-slim

WORKDIR /app

COPY . .

RUN apt-get install openssh sqlite && \
    mkdir -p /app/private /app/.data && \
    [ ! -f /app/private/id_rsa ] && ssh-keygen -t rsa -b 4096 -f /app/private/id_rsa -N "" || true

COPY goback /usr/local/bin/goback

CMD ["goback"]
