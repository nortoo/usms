# Build stage
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .. .

# Build the Go application for linux/amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o usms /app/cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o usms-cli /app/cmd/cli

# Final stage
FROM alpine:latest
LABEL org.opencontainers.image.authors="heap@null.net"

ARG BUILD
LABEL VERSION=$BUILD

WORKDIR /app
COPY --from=builder /app/usms /app/usms
COPY --from=builder /app/usms-cli /app/usms-cli

EXPOSE 8080

ENTRYPOINT ["usms"]