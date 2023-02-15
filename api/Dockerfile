# Build stage
FROM golang:1.20.0-alpine3.17 AS builder

WORKDIR /app

COPY . .
RUN go build -v -o main main.go

# Run stage
FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 3000
CMD ["/app/main"]