# Stage 1: Build
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application from cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o jwt-server ./cmd/server/main.go

# Stage 2
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/jwt-server .
COPY .env .

EXPOSE ${JWT_PORT:-8080}

# Run the application
CMD ["./jwt-server"]