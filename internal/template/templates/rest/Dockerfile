# Build stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

# Install git for downloading dependencies
RUN apk add --no-cache git tzdata

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notesapp ./cmd/api

# Final stage
FROM scratch

WORKDIR /app

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/notesapp .

# Set timezone env variable
ENV TZ=UTC

# Use numeric user:group
USER 1000:1000

# Expose port
EXPOSE 8080

# Command to run
ENTRYPOINT ["/app/notesapp"]