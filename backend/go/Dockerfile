# Build stage: compile Go binaries inside the container
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build deps (e.g. for CGO if needed later)
RUN apk add --no-cache git ca-certificates

# Copy module files first for better layer caching
COPY go.mod go.sum* ./
RUN go mod download && go mod verify

# Copy source and build both binaries (go build will update go.sum if missing)
COPY . .
RUN go build -o /bin/migrate ./cmd/migrate && \
    go build -o /bin/api ./cmd/api

# Run stage: minimal image to run the binaries
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

# Non-root user for security
RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder /bin/migrate /bin/migrate
COPY --from=builder /bin/api /bin/api

# Default command: run the API server.
# To run migrations instead, override the command, e.g.:
#   docker run --rm <image> /bin/migrate up
#   docker compose run --rm backend /bin/migrate up
CMD ["/bin/api"]
