FROM golang:1.25-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o pr-service ./cmd/pr_service/main.go

FROM alpine:latest

RUN apk add --no-cache postgresql-client curl && \
    curl -L https://github.com/pressly/goose/releases/download/v3.19.1/goose_linux_x86_64 -o /usr/local/bin/goose && \
    chmod +x /usr/local/bin/goose

COPY --from=builder /app/pr-service /root/

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["/root/pr-service"]
