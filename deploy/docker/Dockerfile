# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy source
COPY . .

# Build
RUN go build -o forge .

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/forge .

# Copy default config
COPY deploy/config /app/config

# Create data and log directories
RUN mkdir -p /app/data /app/logs

EXPOSE 8080

ENTRYPOINT ["./forge"]
CMD ["--config", "/app/config/runtime.yaml"]
