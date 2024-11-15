FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /goback main.go

FROM alpine:latest

RUN apk add openssh

RUN mkdir -p /app/private /app/data

RUN ssh-keygen -t rsa -b 4096 -f /app/private/id_rsa 

COPY --from=builder /goback /usr/local/bin/goback

CMD ["goback"]

