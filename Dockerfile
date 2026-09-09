# ==============================================================================
# Multi-Stage Hardened Dockerfile: A Bank FinTech Wallet Service
# Architecture: Distroless Non-Root (PCI-DSS & ISO 27001 Compliant)
# ==============================================================================

# Stage 1: Build & Compile
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install security certificates
RUN apk add --no-cache ca-certificates tzdata

# Cache dependency layers
COPY go.mod ./
RUN go mod download || true

# Copy source code
COPY . .

# Compile static binary with zero CGO and stripped debugging symbols
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=1.0.0" \
    -o /app/wallet-service ./cmd/wallet-service

# Create non-root unprivileged user and group
RUN adduser -D -u 10001 -g "" appuser

# Stage 2: Final Minimal Scratch / Distroless Runtime (<15MB)
FROM scratch

# Copy TLS certificates and timezone database
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy unprivileged user
COPY --from=builder /etc/passwd /etc/passwd

# Copy compiled binary
COPY --from=builder /app/wallet-service /wallet-service

# Enforce non-root execution
USER 10001:10001

EXPOSE 8080

ENTRYPOINT ["/wallet-service"]
