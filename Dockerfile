FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /goback main.go

FROM alpine:3.18

RUN apk add --no-cache openssh && \
    mkdir -p /app/private /app/.data && \
    [ ! -f /app/private/id_rsa ] && ssh-keygen -t rsa -b 4096 -f /app/private/id_rsa -N "" || true

COPY --from=builder /goback /usr/local/bin/goback

CMD ["goback"]

