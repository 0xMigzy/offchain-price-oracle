# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker ./cmd/worker

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/worker /usr/local/bin/worker

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/worker"]
