# Builder Stage
FROM golang:1.21.1-alpine3.18 as builder

WORKDIR /app

COPY . .

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o api-fiber-gorm main.go
RUN go get github.com/golang-jwt/jwt
RUN go get github.com/go-redis/redis/v8
RUN go install github.com/cosmtrek/air@v1.49.0
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Development Stage
FROM golang:1.21.1-alpine3.18

WORKDIR /app

COPY --from=builder /app/api-fiber-gorm /app/api-fiber-gorm
COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv
COPY . /app

EXPOSE 8080

CMD ["air", "-c", "/app/air.toml"]
